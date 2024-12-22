import {get, toQueryString} from "./util.js";
import google from "../providers/google.js";

export function authQueryResponse(document){
    let params = new URL(document.location.toString()).searchParams;
    console.log(params)
    return {
        error : params.get("error"),
        code: params.get("code")
    }
}

export function googleAuthLink(){
    const query = {
        client_id: google?.web?.client_id,
        response_type: "code",
        state: "state_parameter_passthrough_value",
        scope:"https://www.googleapis.com/auth/userinfo.profile",
        redirect_uri: getRedirectUri(google?.web?.redirect_uris),
        prompt: "consent",
        include_granted_scopes: true
    }
    const url = google?.web?.auth_uri;
    return url + toQueryString(query)
}

function getRedirectUri(allowedUris){
    const host = window.location.host;
    const allUris = allowedUris || []
    const uris = allUris
        .filter((uri)=> uri && uri.includes(host))
    if (uris.length === 0){
        throw new Error(`Host ${host} not found in uris list ${allUris.join(", ")}`)
    }
    // take matching uri with min symbols
    uris.sort((u1, u2) => u1.length - u2.length)
    console.log(uris)
    return uris[0]
}