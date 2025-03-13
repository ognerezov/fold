# How to Obtain Google Credentials JSON for a New User

## 1. Create a New Project in Google Cloud Console
- Go to the [Google Cloud Console](https://console.cloud.google.com/).
- Click on the project dropdown at the top of the page and select **New Project**.
- Enter a project name and click **Create**.

## 2. Enable the Required API
- In the Google Cloud Console, navigate to **APIs & Services** > **Library**.
- Search for the API you need (e.g., Google Drive API, Google Sheets API, etc.).
- Click on the API and then click **Enable**.

## 3. Create Credentials
- Go to **APIs & Services** > **Credentials**.
- Click **Create Credentials** and select **OAuth client ID**.
- If prompted, configure the consent screen:
    - Choose **External** for user type.
    - Fill in the required fields (app name, support email, etc.).
    - Add any necessary scopes (e.g., `https://www.googleapis.com/auth/drive`).
    - Add test users (your email or the new user's email).
    - Save and continue.
- Under **Application type**, select **Web application**.
- Add authorized redirect URIs (e.g., `http://localhost:8080` for testing).
- Click **Create**.

## 4. Download the JSON File
- After creating the OAuth client ID, you’ll see it listed in the Credentials page.
- Click the download icon (⭳) next to the client ID to download the JSON file.

## 5. Use the Credentials
- The downloaded JSON file contains your `client_id`, `client_secret`, and other details needed for authentication.
- Use this file in your application to authenticate and access Google APIs.

---

### Notes:
- If you’re in a hurry, you can skip some optional steps like adding logos or detailed app information in the consent screen.
- For testing purposes, you can use `http://localhost` as the redirect URI.