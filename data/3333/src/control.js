import settings from "./fold.js"
import {authHeaders, hide, parseForm, send, show, withJsonContent} from "./util.js";
import {
    clearErrorEvent,
    endLoadingEvent,
    ERROR_RECEIVED,
    getStoredToken,
    startLoadingEvent,
    USER_RECEIVED
} from "./main.js";

let fileInputs = {
}

export function actionForm(label, params, formParent, hideOnForm, submit){
    hide(hideOnForm)
    fileInputs = {}
    const form = document.createElement("form")
    form.onsubmit = submit
    form.className = "m-frame fill-x start-column g-2 with-shadow wave"

    const title = document.createElement("div")
    title.className = "font-bold text-xl"
    title.innerText = label
    form.appendChild(title)

    for(const param of params){
        const label = document.createElement("label")
        label.className = "label-text fill-x"
        label.for = param.name
        label.innerText = param.schema.description
        form.appendChild(label)

        const input = document.createElement("input")
        input.className = "input"
        input.id = param.name
        if(['json'].includes(param.schema.type)){
           input.type = "file"
           input.accept = "." + param.schema.type
           input.onchange = onFileInput
        }else {
            input.placeholder = param.schema.description
            input.type = ["integer", "float"].includes(param.schema.type) ? "number" : "text"
        }
        input.name = param.name
        form.appendChild(input)
    }

    const footer = document.createElement("div")
    footer.className = "end-row pt_2"
    form.appendChild(footer)

    const cancel = document.createElement("button")
    cancel.innerText = "Cancel"
    cancel.type = "reset"
    cancel.className = "btn btn-wave m_05"
    cancel.onclick = ()=>{
        show(hideOnForm)
        window.dispatchEvent(clearErrorEvent)
        form.remove()
    }
    footer.appendChild(cancel)
    const ok = document.createElement("button")
    ok.innerText = "Submit"
    ok.type = "submit"
    ok.className = "btn btn-deep m_05"
    ok.onclick = () => {
        window.dispatchEvent(clearErrorEvent)
    }
    footer.appendChild(ok)


    formParent.appendChild(form)
}

export function actionButton(key, s, token, formParent, hideOnForm, className = "btn btn-deep m_05"){
    const button = document.createElement('button')
    const params = settings[key].parameters
    const showForm = params && params.length > 0
    const schema = {};
    (params || []).forEach((param)=>{
        schema[param.name] = param.schema.type
    })
    button.id = key
    button.innerText = s.label
    button.className = className
    button.onclick = showForm ? ()=>  actionForm(s.label, params, formParent, hideOnForm,
            async (event)=> {
                event.preventDefault();
                let data = parseForm(event, schema)
                data = {...data, ...fileInputs}
                console.log(data)

                await useControlEndpoint(key, token, data)
            }) :
        ()=>  useControlEndpoint(key, token)
    return button
}

export function actionButtons(element, formParent, hideOnForm, pub = true){
    const token = pub ? null : getStoredToken();
    const keys = Object
                        .keys(settings)
                        .filter((key)=> Boolean(settings[key].public) === pub)
    const res = []
    keys.forEach((key)=>{
        const button = actionButton(key, settings[key], token, formParent, hideOnForm)
        element.appendChild(button)
        res.push(button)
    })
    
    return res
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
    const res = await send(config.path,
        method === "GET" ? params : undefined,
        config.method || "GET",
        headers, body )
    if (res?.status !== 200){
        console.log(res)
        const msg = `Error: ${res?.status}. Message: ${res?.json?.message}`
        window.dispatchEvent(new CustomEvent(ERROR_RECEIVED, {detail : msg}))
    }
    window.dispatchEvent(endLoadingEvent)
    return res || { error : "no response"}
}

function onFileInput(event) {
    const file = event.target.files[0]; // Get the selected file
    console.log(event)
    if (file && file.type === "application/json") { // Ensure it's a JSON file
        const reader = new FileReader(); // Create a FileReader

        reader.onload = function(e) {
            const content = e.target.result; // Get the file content
            try {
                const json = JSON.parse(content); // Parse the JSON
                console.log(json)
                fileInputs[event.target.id] = json
            } catch (error) {
                console.error("Error parsing JSON:", error);
                alert("Invalid JSON file. Please upload a valid JSON file.");
            }
        };

        reader.readAsText(file); // Read the file as text
    } else {
        alert("Please upload a valid JSON file.");
    }
}