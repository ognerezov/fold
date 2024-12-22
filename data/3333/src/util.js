export async function get(url, query, headers) {
    try {
        const response = await fetch(url + toQueryString(query), {
            headers : headers || {}
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
        return {
            error : error.message,
            status: 500
        }
    }
}

export function toQueryString(query) {
    if (!query || !Object.keys(query).length) {
        return "";
    }
    return `?${new URLSearchParams(query).toString()}`;
}