export async function send(url, query, method = "GET", headers, body = undefined) {
    try {
        const response = await fetch(url + toQueryString(query), {
            method : method,
            headers : headers || {},
            body : body
        });
        if (!response.ok) {
            return {
                status: response.status
            }
        }
        const json =  await response.json();

        return {
            json : json,
            status : 200
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