import { getCookie } from "./cookies.js"
import {authHeaders, send, hide, show} from "./util.js";
const iss = "fold"

const START_LOADING = "start_loading";
export const startLoadingEvent = new Event(START_LOADING);
const END_LOADING = "end_loading";
export const endLoadingEvent = new Event(END_LOADING);
export const ERROR_RECEIVED = "error_received"
export const ERROR_CLEARED  = "error_cleared"
export const clearErrorEvent = new Event(ERROR_CLEARED);

export function listenForErrors(element, terminalElements, textOutput, closeButton){
    closeButton.onclick = ()=> window.dispatchEvent(clearErrorEvent)
    element.addEventListener(
        ERROR_RECEIVED,
        (event)=>{
            console.log(event)
            const errorMessage = event?.detail || "Failed to process request"
            show(textOutput)
            textOutput.innerText = errorMessage
            for(const el of terminalElements){
                show(el)
            }
        }
    );
    element.addEventListener(
        ERROR_CLEARED,
        ()=>{
            for(const el of terminalElements){
                hide(el)
            }
        }
    );
}

export function listenLoading(element, onLoad, onEndLoading){
    element.addEventListener(
        START_LOADING,
         onLoad
    );
    element.addEventListener(
        END_LOADING,
         onEndLoading
    );
}

export function getStoredToken(){
    const cookieKey = `${iss}Token`
    return  getCookie(cookieKey)
}

export async function me(){
    const token = getStoredToken()

    if (!token){
        return Promise.resolve(null);
    }

    return await send("/me", null, "GET", authHeaders(token) )
}


 export function redirectToControl(){
    window.location.replace("/control")
 }