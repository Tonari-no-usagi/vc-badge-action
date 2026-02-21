/**
 * SD-JWT Verifier Logic
 */

function decodeBase64URL(str) {
    // Base64URL を Base64 に変換
    const base64 = str.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));
    return JSON.parse(jsonPayload);
}

function verifyVC() {
    const input = document.getElementById('vcInput').value.trim();
    if (!input) return;

    const parts = input.split('~');
    const jwtPart = parts[0];
    const disclosures = parts.slice(1, -1); // 最後の要素は空文字列

    try {
        // JWT ペイロードのデコード
        const jwtHeader = decodeBase64URL(jwtPart.split('.')[0]);
        const jwtPayload = decodeBase64URL(jwtPart.split('.')[1]);
        
        const finalClaims = {};
        
        // 公開クレームのコピー（_sd以外）
        for (const [key, value] of Object.entries(jwtPayload)) {
            if (key !== '_sd') {
                finalClaims[key] = value;
            }
        }

        // ディスクロージャの処理
        disclosures.forEach(d => {
            const decoded = decodeBase64URL(d); // [salt, name, value]
            const name = decoded[1];
            const value = decoded[2];
            finalClaims[name] = value;
        });

        displayResult(finalClaims);
    } catch (e) {
        alert('検証に失敗しました。正しいフォーマットか確認してください。: ' + e.message);
    }
}

function displayResult(claims) {
    const grid = document.getElementById('claimsGrid');
    grid.innerHTML = '';
    
    // クレームのレンダリング
    const sortedKeys = Object.keys(claims).sort();
    sortedKeys.forEach(key => {
        if (typeof claims[key] === 'object' && key !== 'vc') {
            // オブジェクトの場合は文字列化
            renderClaim(key, JSON.stringify(claims[key]));
        } else if (key === 'vc') {
            // vc タイプは特別扱い
            renderClaim('credential_type', claims[key].type.join(', '));
        } else {
            renderClaim(key, claims[key]);
        }
    });

    document.getElementById('result').style.display = 'block';
    
    // Smooth scroll
    document.getElementById('result').scrollIntoView({ behavior: 'smooth' });
}

function renderClaim(label, value) {
    const grid = document.getElementById('claimsGrid');
    const card = document.createElement('div');
    card.className = 'claim-card';
    card.innerHTML = `
        <div class="claim-label">${label}</div>
        <div class="claim-value">${value}</div>
    `;
    grid.appendChild(card);
}

// URLパラメータからの読み込み
window.onload = () => {
    const params = new URLSearchParams(window.location.search);
    const vc = params.get('vc');
    if (vc) {
        document.getElementById('vcInput').value = vc;
        verifyVC();
    }
};
