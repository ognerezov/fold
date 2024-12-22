import { get } from "./util";
import { SETTINGS } from "./settings";

export function authQueryResponse(document){
    let params = new URL(document.location.toString()).searchParams;
    console.log(params)
    return {
        error : params.get("error"),
        code: params.get("code")
    }
}

export async function initiateGoogleLogin(){
    const query = {
        client_id: SETTINGS.GoogleProject,
        response_type: "code",
        state: "state_parameter_passthrough_value",
        scope:"https://www.googleapis.com/auth/userinfo.profile",
        redirect_uri: SETTINGS.GoogleRedirectUrl,
        prompt: "consent",
        include_granted_scopes: true
    }

    await get("https://accounts.google.com/o/oauth2/v2/auth", query)
}