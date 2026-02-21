# VC Badge Action

[日本語版 (Japanese)](README_ja.md)

This GitHub Action issues **SD-JWT Verifiable Credentials (VC)** to developers' DIDs based on their GitHub activities (e.g., merging PRs).

## Features

- **SD-JWT VC (RFC 9901)** compliant credential issuance
- **Selective Disclosure**: Prove contributions while revealing only necessary data
- **Ed25519 Signing**: High-performance and secure signatures
- **Web Verification Portal**: Built-in portal to verify and view credentials in browser
- **did:webvh** compatible (upcoming)

## Usage

Add the following to your `.github/workflows/badge.yml`:

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

### Inputs

| Name | Description | Required | Default |
| :--- | :--- | :---: | :--- |
| `issuer_did` | DID of the issuer | Yes | - |
| `issuer_private_key` | Issuer's private key (Hex encoded Ed25519) | Yes | - |
| `subject_did_regex` | Regex to extract subject's DID from PR body | No | `DID: (did:.*)` |

## Verification Portal

Once issued, the SD-JWT can be verified using the built-in portal at:
`https://YOUR_GITHUB_ID.github.io/vc-badge-action/`

## License

MIT
