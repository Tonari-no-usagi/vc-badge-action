package sdjwt

import (
	"crypto/ed25519"
	"strings"
	"testing"
)

func TestSDJWTFlow(t *testing.T) {
	_, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}

	issuer := &Issuer{IssuerDID: "did:webvh:example.com"}

	publicClaims := map[string]any{
		"sub": "did:example:holder",
	}
	privateClaims := map[string]any{
		"given_name":  "John",
		"family_name": "Doe",
	}

	claims, disclosures, err := issuer.CreateSDClaims(publicClaims, privateClaims)
	if err != nil {
		t.Fatalf("CreateSDClaims failed: %v", err)
	}

	// _sd クレームが存在することを確認
	if _, ok := claims["_sd"]; !ok {
		t.Error("_sd claim missing")
	}

	sdJWT, err := issuer.SignAndCombine(claims, disclosures, priv)
	if err != nil {
		t.Fatalf("SignAndCombine failed: %v", err)
	}

	// 形式チェック: JWT~D1~D2~...~
	parts := strings.Split(sdJWT, "~")
	if len(parts) != 4 { // JWT, Disclosure1, Disclosure2, Empty
		t.Errorf("expected 4 parts, got %d", len(parts))
	}

	t.Logf("Generated SD-JWT: %s", sdJWT)
}
