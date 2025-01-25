import {JSON_HEADERS, send} from "./util.js";
import {createCookie} from "./cookies.js"

export function authQueryResponse(document){
    let params = new URL(document.location.toString()).searchParams;
    console.log(params)
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
    // TODO this should go from settings
    request.redirect_uri = "http://localhost:3333"
    const res = await send("/auth", null, "POST", JSON_HEADERS, JSON.stringify(request))

    const key = `${res?.json?.iss}Token`
    if(res?.json?.token){
        createCookie(key, res?.json?.token, 100);
    }
    return res.json
}