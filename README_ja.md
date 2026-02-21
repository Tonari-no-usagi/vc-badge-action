# VC Badge Action

[English Version](README.md)

GitHub上の活動（PRマージなど）をトリガーに、開発者のDIDに対して **SD-JWT形式のVerifiable Credential (VC)** を発行するGitHub Actionです。

## 特徴

- **SD-JWT VC (RFC 9901)** 準拠の証明書発行
- **選択的開示 (Selective Disclosure)**: 最小限のデータ露出で実績を証明
- **Ed25519** による高速かつセキュアな署名
- **Web検証ポータル**: ブラウザ上で誰でも実績を検証可能なポータルサイトを同梱
- **did:webvh** 対応（予定）

## 使用方法

`.github/workflows/badge.yml` に以下のように設定します：

```yaml
name: Issue Contribution Badge
on:
  pull_request:
    types: [closed]
    branches: [main]

jobs:
  issue_badge:
    if: github.event.pull_request.merged == true
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Issue VC Badge
        uses: your-org/vc-badge-action@v1
        with:
          issuer_did: "did:webvh:example.com:scid-123"
          issuer_private_key: ${{ secrets.VC_ISSUER_KEY }}
          subject_did_regex: "DID: (did:.*)"
```

### 入力パラメータ

| 名前 | 説明 | 必須 | デフォルト |
| :--- | :--- | :---: | :--- |
| `issuer_did` | 発行者のDID | Yes | - |
| `issuer_private_key` | 発行者の秘密鍵 (Hex形式のEd25519) | Yes | - |
| `subject_did_regex` | PR本文から保持者のDIDを抽出する正規表現 | No | `DID: (did:.*)` |

## 検証ポータル

発行された SD-JWT は、同梱のポータルサイトで検証可能です：
`https://tonari-no-usagi.github.io/vc-badge-action/`

## ライセンス

MIT
