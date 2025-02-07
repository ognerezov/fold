import { getCookie } from "./cookies.js"
import {authHeaders, send} from "./util.js";
const iss = "fold"

const START_LOADING = "start_loading";
export const startLoadingEvent = new Event(START_LOADING);
const END_LOADING = "end_loading";
export const endLoadingEvent = new Event(END_LOADING);

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


 export function redirectToLogin(){
    window.location.replace("/login")
 }