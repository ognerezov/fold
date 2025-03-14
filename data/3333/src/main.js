import {deleteCookie, getCookie} from "./cookies.js"
import {authHeaders, send, hide, show, linkify} from "./util.js";
import project from "./project.js"
const iss = project.name

const START_LOADING = "start_loading";
export const startLoadingEvent = new Event(START_LOADING);
const END_LOADING = "end_loading";
export const endLoadingEvent = new Event(END_LOADING);
export const ERROR_RECEIVED = "error_received"
export const ERROR_CLEARED  = "error_cleared"
export const clearErrorEvent = new Event(ERROR_CLEARED);
export const USER_LOGIN = "user_logged_in"
export const USER_LOGOUT = "user_logged_out"
export const USER_RECEIVED = "user_received"
export const logoutEvent = new Event(USER_LOGOUT);


export function logout(){
    window.dispatchEvent(logoutEvent)
    deleteCookie(`${iss}Token`)
}


export async function listenForUser(element, unAuthElements, authElements, userDisplay, logoutBtn){
    logoutBtn.onclick = logout
    for(const element of authElements){
        hide(element)
    }
    element.addEventListener(
        USER_LOGOUT,
        ()=>{
            for(const element of authElements){
                hide(element)
            }
            for(const element of unAuthElements){
                show(element)
            }
        }
    );
    element.addEventListener(
        USER_RECEIVED,
        (event)=>{
            for(const element of unAuthElements){
                hide(element)
            }
            for(const element of authElements){
                show(element)
            }
            userDisplay.innerText = event.detail
        }
    );

    const { json: user, status } = await me()
    if (status !== 200){
        if (status === 401) {
            logout()
        }else {
            window.dispatchEvent(logoutEvent)
        }
        return
    }
    window.dispatchEvent(new CustomEvent(USER_RECEIVED, {detail : user.id}))
}


export function listenForErrors(element, terminalElements, textOutput, closeButton){
    closeButton.onclick = ()=> window.dispatchEvent(clearErrorEvent)
    element.addEventListener(
        ERROR_RECEIVED,
        (event)=>{
            const errorMessage = event?.detail || "Failed to process request"
            show(textOutput)
            textOutput.innerHTML = linkify(errorMessage)
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
        return Promise.resolve({status: 401, json : null});
    }

    return await send("/me", null, "GET", authHeaders(token) )
}


 export function redirectToControl(){
    window.location.replace("/dashboard")
 }