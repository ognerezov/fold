import {JSON_HEADERS, send} from "./util.js";
import {saveTokenResponse} from "./cookies.js"
import {redirectToControl} from "./main.js";
import {getRedirect} from "./login.js";

export function authQueryResponse(document){
    let params = new URL(document.location.toString()).searchParams;
    return {
        error : params.get("error"),
        code: params.get("code")
    }
}

export async function processCode(document){
    const request = authQueryResponse(document)

    if (!request.code){
        return Promise.resolve(request)
    }

    request.redirect_uri = getRedirect()
    const res = await send("/auth", null, "POST", JSON_HEADERS, JSON.stringify(request))
    console.log(res)
    if (res.status === 200) {
        saveTokenResponse(res)
        redirectToControl()
    }
    return res.json
}