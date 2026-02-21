package sdjwt

import (
	"crypto/ed25519"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// SignAndCombine は JWT に署名し、ディスクロージャを結合して SD-JWT 形式（JWT~D1~D2~...~）を生成します。
func (i *Issuer) SignAndCombine(claims map[string]any, disclosures []string, privKey ed25519.PrivateKey) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims(claims))

	// 発行者DIDなどの標準クレームを設定
	claims["iss"] = i.IssuerDID

	signedJWT, err := token.SignedString(privKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	// SD-JWT 形式: <JWT>~<Disclosure 1>~<Disclosure 2>~...~
	var sb strings.Builder
	sb.WriteString(signedJWT)
	for _, d := range disclosures {
		sb.WriteString("~")
		sb.WriteString(d)
	}
	sb.WriteString("~") // 最後にチルダが必要

	return sb.String(), nil
}
