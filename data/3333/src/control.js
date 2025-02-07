import settings from "./fold.js"
import {authHeaders, send, withJsonContent} from "./util.js";
import {endLoadingEvent, getStoredToken, startLoadingEvent} from "./main.js";

export function actionButton(key, s, token = null, className = "btn btn-deep m_05"){
    const button = document.createElement('button')
    button.id = key
    button.innerText = s.label
    button.className = className
    button.onclick = ()=>useControlEndpoint(key, token)
    return button
}

export function actionButtons(element, pub = true){
    const token = pub ? null : getStoredToken();
    const keys = Object
                        .keys(settings)
                        .filter((key)=> Boolean(settings[key].public) === pub)
    keys.forEach((key)=>{
        element.appendChild(actionButton(key, settings[key], token))
    })
}

export async function useControlEndpoint(id, token = null, params ){
    const config = settings[id]
    console.log(config)
    if (!config){
        return Promise.resolve({ error : "config not found"})
    }

    const method = config.method || "GET"
    const body = method === "GET" || !params ? undefined : JSON.stringify(params)
    const headers = token ? authHeaders(token) : withJsonContent({})
    window.dispatchEvent(startLoadingEvent)
    const res = await send(config.path, params, config.method || "GET", headers, body )
    window.dispatchEvent(endLoadingEvent)
    return res || { error : "no response"}
}