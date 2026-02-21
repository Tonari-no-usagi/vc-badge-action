package sdjwt

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	"github.com/google/uuid"
)

// Disclosure は SD-JWT の選択的開示要素を表します。
type Disclosure struct {
	Salt      string `json:"salt"`
	ClaimName string `json:"claim_name"`
	Value     any    `json:"value"`
}

// NewDisclosure は新しいディスクロージャを作成します。
func NewDisclosure(name string, value any) (*Disclosure, error) {
	salt := uuid.New().String()
	return &Disclosure{
		Salt:      salt,
		ClaimName: name,
		Value:     value,
	}, nil
}

// Encode はディスクロージャを Base64URL エンコードします。
func (d *Disclosure) Encode() (string, error) {
	// [salt, name, value] の配列として直列化
	arr := []any{d.Salt, d.ClaimName, d.Value}
	data, err := json.Marshal(arr)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

// Hash はエンコードされたディスクロージャのハッシュを計算します。
func (d *Disclosure) Hash() (string, error) {
	encoded, err := d.Encode()
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256([]byte(encoded))
	return base64.RawURLEncoding.EncodeToString(hash[:]), nil
}

// Issuer は SD-JWT の発行者を表します。
type Issuer struct {
	IssuerDID string
}

// CreateSDClaims は公開クレームと隠蔽するクレームから _sd フィールドを含むマップを生成します。
func (i *Issuer) CreateSDClaims(publicClaims map[string]any, privateClaims map[string]any) (map[string]any, []string, error) {
	sdClaims := make(map[string]any)
	for k, v := range publicClaims {
		sdClaims[k] = v
	}

	var disclosures []string
	var sdHashes []string

	for name, value := range privateClaims {
		d, err := NewDisclosure(name, value)
		if err != nil {
			return nil, nil, err
		}

		encoded, err := d.Encode()
		if err != nil {
			return nil, nil, err
		}
		disclosures = append(disclosures, encoded)

		hash, err := d.Hash()
		if err != nil {
			return nil, nil, err
		}
		sdHashes = append(sdHashes, hash)
	}

	if len(sdHashes) > 0 {
		sdClaims["_sd"] = sdHashes
	}

	return sdClaims, disclosures, nil
}
