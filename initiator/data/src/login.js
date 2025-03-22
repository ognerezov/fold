import {hide, JSON_HEADERS, send, show, toQueryString} from "./util.js";
import google from "../providers/google.js";
import {saveTokenResponse} from "./cookies.js";
import {endLoadingEvent, startLoadingEvent, ERROR_RECEIVED} from "./main.js";
import project from "./project.js"

export function showSupportedAuthOptions(elements){
    if (!elements){
        return
    }
    const supportedAuth = project.auth_providers || []
    const res = []
    Object.keys(elements).forEach((key)=>{
        for(const element of elements[key]) {
            if (supportedAuth.includes(key)) {
                show(element)
                res.push(element)
                console.log("show", element)
            } else {
                hide(element)
                console.log("hide", element)
            }
        }
    })

    return res;
}

// https://www.googleapis.com/auth/drive https://www.googleapis.com/auth/spreadsheets
export function googleAuthLink(){
    if(google.disabled){
        return ""
    }
    const query = {
        client_id: google?.web?.client_id,
        response_type: "code",
        state: "state_parameter_passthrough_value",
        scope:"https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/cloud-platform",
        redirect_uri: getRedirect(),
        prompt: "consent",
        include_granted_scopes: true,
        access_type: "offline"
    }
    const url = google?.web?.auth_uri;
    return url + toQueryString(query)
}

export function getRedirect(){
    return getRedirectUri(google?.web?.redirect_uris)
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
    return uris[0]
}

export async function login(data){
    console.log(data)
    window.dispatchEvent(startLoadingEvent)
    const res = await send("/login", null, "POST", JSON_HEADERS, JSON.stringify(data))
    window.dispatchEvent(endLoadingEvent)
    if (res.status === 200) {
        saveTokenResponse(res)
    } else {
        const message = res.status === 401 ? "Wrong login or password" : res.status === 400 ? "Bad request" : "Server error"
        window.dispatchEvent(new CustomEvent(ERROR_RECEIVED, { detail : message }))
    }
    return res
}