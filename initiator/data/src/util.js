export async function send(url, query, method = "GET", headers, body = undefined) {
    try {
        const response = await fetch(url + toQueryString(query), {
            method : method,
            headers : headers || {},
            body : body
        });
        const json =  await response.json();

        return {
            json : json,
            status : response.status
        }

    } catch (error) {
        console.log(error)
        return {
            error : error.message,
            status: 500
        }
    }
}

export async function echo(data){
    return await send("/echo", "", "POST", {
        "Content-Type": "application/json"
    }, JSON.stringify(data))
}

export const JSON_HEADERS =  Object.freeze({
    "Content-Type": "application/json"
})

export function toQueryString(query) {
    if (!query || !Object.keys(query).length) {
        return "";
    }
    return `?${new URLSearchParams(query).toString()}`;
}

export function withJsonContent(obj){
    return {... JSON_HEADERS, ...obj}
}

export function authHeaders(token){
    return withJsonContent({
        "Authorization" : `Bearer ${token}`
    })
}

export function hide(element){
    if(Array.isArray(element)){
        for(const el of element){
            el.classList.add("hidden");
        }
        return
    }
    element.classList.add("hidden");
}

export function show(element){
    if(Array.isArray(element)){
        for(const el of element){
            el.classList.remove("hidden");
        }
        return
    }
    element.classList.remove("hidden");
}

export function parseForm(form, schema){
    if (!schema){
        return schema
    }

    const result = cloneObject(schema)

    for(const el of form.target){
        if (!el || !el.id || result[el.id] === undefined){
            continue
        }
        result[el.id] = el.value
    }

    return result
}

export function cloneObject(obj){
    return JSON.parse(JSON.stringify(obj))
}

export function linkify(text) {
    const urlRegex = /(\b(https?|ftp|file):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/ig;

    const res= text.replace(urlRegex, function(url) {
        return '<a class="paper-link pointer" href="' + url + '" target="_blank">' + url + '</a>';
    });

    return `<div class="fill-x msg-box">${res}</div>`
}