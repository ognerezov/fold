import settings from "./fold.js"
import {send} from "./util.js";

export function actionButton(key, s, className = "btn btn-deep m_05"){
    const button = document.createElement('button')
    button.id = key
    button.innerText = s.label
    button.className = className
    button.onclick = ()=>usePublicControl(key)
    return button
}

export function publicActionButtons(element){
    const keys = Object
                        .keys(settings)
                        .filter((key)=> settings[key].public)
    keys.forEach((key)=>{
        element.appendChild(actionButton(key, settings[key]))
    })
}

export async function usePublicControl(id, params ){
    const config = settings[id]
    console.log(config)
    if (!config){
        return Promise.resolve({ error : "config not found"})
    }

    if(!config.public){
        return Promise.resolve({ error : "config is private"})
    }
    const method = config.method || "GET"
    const body = method === "GET" || !params ? undefined : JSON.stringify(params)

    const res = await send(config.path, params, config.method || "GET", body )
    return res || { error : "no response"}
}