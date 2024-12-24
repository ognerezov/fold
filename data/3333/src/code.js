export function authQueryResponse(document){
    let params = new URL(document.location.toString()).searchParams;
    console.log(params)
    return {
        error : params.get("error"),
        code: params.get("code")
    }
}