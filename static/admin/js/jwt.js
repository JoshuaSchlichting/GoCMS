function parseJwt(token) {
    let base64Url = token.split('.')[1];
    let base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    let jsonPayload = decodeURIComponent(window.atob(base64).split('').map(function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));

    return JSON.parse(jsonPayload);
}

function validateJwt(token) {
    let jwt = parseJwt(token);
    let now = new Date();
    let exp = new Date(jwt.exp * 1000);
    return now < exp;
}

function validateJwtSignature(token) {
    let jwt = parseJwt(token);
    let signature = jwt.signature;
    let jwtWithoutSignature = token.split('.')[0] + '.' + token.split('.')[1];
    let calculatedSignature = CryptoJS.HmacSHA256(jwtWithoutSignature).toString(CryptoJS.enc.Base64);
    return signature === calculatedSignature;
}

function setJwtToLocalStorage() {
    // get jwt from cookie
    let jwt = getCookie('jwt');
    if (jwt !== undefined) {
        // validate jwt
        if (validateJwt(jwt) && validateJwtSignature(jwt)) {
            // set jwt to local storage
            localStorage.setItem('jwt', jwt);
        }
    }
}
