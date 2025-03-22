import project from "./project.js"
const iss = project.name

window.unescape = window.unescape || window.decodeURI;
export function createCookie(name, value, days) {
    let expires;
    if (days) {
        const date = new Date();
        date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
        expires = "; expires=" + date.toGMTString();
    }
    else {
        expires = "";
    }
    document.cookie = name + "=" + value + expires + "; path=/";
}

export function saveTokenResponse(res){
    const key = `${iss}Token`
    if(res?.json?.token){
        createCookie(key, res?.json?.token, 100);
    }
}

export function getCookie(key) {
    if (document.cookie.length > 0) {
        let c_start = document.cookie.indexOf(key + "=");
        if (c_start !== -1) {
            c_start = c_start + key.length + 1;
            let c_end = document.cookie.indexOf(";", c_start);
            if (c_end === -1) {
                c_end = document.cookie.length;
            }
            return window.unescape(document.cookie.substring(c_start, c_end));
        }
    }
    return "";
}

export function deleteCookie(name) {
    document.cookie = name +'=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;';
}