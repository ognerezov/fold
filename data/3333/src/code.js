import {JSON_HEADERS, send} from "./util.js";

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
    console.log(request)
    if (!request.code){
        return Promise.resolve(request)
    }
    request.redirect_uri = "http://localhost:3333"
    return send("/auth", null, "POST", JSON_HEADERS, JSON.stringify(request))
}