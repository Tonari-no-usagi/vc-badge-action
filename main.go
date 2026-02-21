package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/your-org/vc-badge-action/pkg/github"
	"github.com/your-org/vc-badge-action/pkg/sdjwt"
)

func main() {
	// 入力パラメータの取得
	issuerDID := os.Getenv("INPUT_ISSUER_DID")
	issuerPrivKeyHex := os.Getenv("INPUT_ISSUER_PRIVATE_KEY")
	subjectRegexStr := os.Getenv("INPUT_SUBJECT_DID_REGEX")
	eventPath := os.Getenv("GITHUB_EVENT_PATH")

	if issuerDID == "" || issuerPrivKeyHex == "" || eventPath == "" {
		log.Fatal("Missing required inputs")
	}

	// 秘密鍵のデコード
	privKeyBytes, err := hex.DecodeString(issuerPrivKeyHex)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}
	privKey := ed25519.PrivateKey(privKeyBytes)

	// イベントの読み込み
	event, err := github.LoadEvent(eventPath)
	if err != nil {
		log.Fatalf("Failed to load event: %v", err)
	}

	// 保持者DIDの抽出
	subjectDID := ""
	if subjectRegexStr != "" {
		re, err := regexp.Compile(subjectRegexStr)
		if err == nil {
			matches := re.FindStringSubmatch(event.PullRequest.Body)
			if len(matches) > 1 {
				subjectDID = matches[1]
			}
		}
	}

	if subjectDID == "" {
		log.Println("Subject DID not found in PR body, skipping VC issuance")
		return
	}

	// VCクレームの構築
	issuer := &sdjwt.Issuer{IssuerDID: issuerDID}
	publicClaims := map[string]any{
		"sub": subjectDID,
		"vc": map[string]any{
			"type": []string{"VerifiableCredential", "ContributionBadge"},
		},
	}
	privateClaims := map[string]any{
		"repository":        event.PullRequest.Base.Repo.FullName,
		"contribution_type": "Pull Request Merge",
		"pull_request":      event.PullRequest.Number,
		"github_user":       event.PullRequest.User.Login,
	}

	claims, disclosures, err := issuer.CreateSDClaims(publicClaims, privateClaims)
	if err != nil {
		log.Fatalf("Failed to create SD claims: %v", err)
	}

	sdJWT, err := issuer.SignAndCombine(claims, disclosures, privKey)
	if err != nil {
		log.Fatalf("Failed to sign and combine SD-JWT: %v", err)
	}

	// 結果の出力（GitHub Actionsのステップ出力やファイル保存など）
	fmt.Printf("::set-output name=sd_jwt_vc::%s\n", sdJWT)

	// Badge Directoryへの保存
	badgeDir := fmt.Sprintf("badges/%s", event.PullRequest.User.Login)
	if err := os.MkdirAll(badgeDir, 0755); err != nil {
		log.Fatalf("Failed to create badge directory: %v", err)
	}

	hashPath := fmt.Sprintf("%s/contribution.jwt", badgeDir)
	if err := os.WriteFile(hashPath, []byte(sdJWT), 0644); err != nil {
		log.Fatalf("Failed to write badge file: %v", err)
	}

	fmt.Printf("Badge generated and saved to %s\n", hashPath)
	fmt.Println("--------------------------------------------------")
	fmt.Println("Verification Link (Paste your SD-JWT here):")
	fmt.Println("https://your-org.github.io/vc-badge-action/verification/")
	fmt.Println("--------------------------------------------------")
}
