# 📘 Fold Server – Command Line Reference

`fold` is a CLI tool that allows you to record API responses, sanitize data, serve mock APIs, and manage local/cloud migrations.

---

## 📦 Installation

Fold Server is a zero-dependency single-binary application. No runtime or package manager is required.

### 🧱 Download the Binary

Visit the [Releases page](https://github.com/ognerezov/fold/releases) and download the latest binary for your platform:

| OS      | File Name Example        |
|---------|--------------------------|
| macOS   | `fold-darwin-amd64` or `fold-darwin-arm64` |
| Linux   | `fold-linux-amd64`       |
| Windows  | `fold-windows-amd64.exe` |

---

### 💻 Setup Instructions

#### macOS / Linux

```bash
# Download and give execute permissions
chmod +x fold-linux-amd64

# Optional: Move to a directory in your PATH
sudo mv fold-linux-amd64 /usr/local/bin/fold
```

Now you can run it using:

```bash
fold -h
```

#### Windows

1. Download the `fold-windows-amd64.exe` file.
2. (Optional) Rename it to `fold.exe` and add its folder to your system PATH:
  - Open **System Properties** → **Environment Variables** → **System PATH** → Add folder.
3. Open Command Prompt and run:

```cmd
fold.exe -h
```

---

## 🔧 Basic Syntax

```bash
fold [options]
```

You can use the tool to:

- Record API responses
- Initialize a new project
- Migrate data
- Start a local server
- Serve from a single file or a folder

---

## 🧩 Global Options

| Option                  | Alias | Description                                      | Default       |
|-------------------------|-------|--------------------------------------------------|---------------|
| `--dir`                 | `-d`  | Working directory for the project                | `./`          |
| `--file`                | `-f`  | Serve or initialize using a single file          |               |
| `--port`                | `-p`  | Port to serve on                                 | `3333`        |
| `--api`                 | `-a`  | Base path for the API routes                     |               |
| `--name`                | `-n`  | Project/app name                                 | `fold`        |
| `--description`         |       | Project description                              |               |
| `--origin`              | `-o`  | CORS allowed origin                              | `http://localhost:<port>` |
| `--cache`               |       | Enable/disable file caching                      | `true`        |
| `--reg`                 |       | Allow user registration                          | `true`        |
| `--version`             |       | Project version                                  | `1.0.0`       |

---

## 🚀 Actions

- All major actions (`--init`, `--record`, `--migrate`) are **mutually exclusive**, meaning only one should be executed per run.
- ✅ **Exception**: `--init` can be combined with `--record` to initialize a project and record your first API call in one command.
- You can also run multiple `--record` commands independently later to add more data and endpoints.
- `--init` command requires empty --dir to exist

### 1. `--init` / `-i`

Initialize a new project folder.

| Option              | Alias | Description                              |
|---------------------|-------|------------------------------------------|
| `--template`        | `-t`  | Template to use                          |
| `--credentials`     | `-c`  | Google Drive credentials (if using cloud)|
| `--dr`              |       | Google Drive folder ID (destination)     |

### 2. `--record` / `-r`

Record an API call using a description JSON file.

```bash
fold -record ./my-request.json -dir ./mocks
```

See record file description [below](#-record-file-format)

### 3. `--migrate` / `-m`

Migrate data between sources.

| Option              | Alias | Description                                  |
|---------------------|-------|----------------------------------------------|
| `--source`          | `-s`  | Source type: `dir` or `drive`                |
| `--destination`     | `-ds` | Destination type: `dir` or `drive`           |

---

## 🌐 Serve Mode

Starts a local server based on directory or file mode.

```bash
fold -d ./mocks
```

Or with a single file:

```bash
fold -f ./mocks/data.json
```

Backend endpoints and security settings are configured by a [File structure](#-server-file-structure)

---

## 🆘 Help

```bash
fold -h
```

Prints all available commands and flags.

---

## ✅ Examples

### Start with a folder

```bash
fold -d ./mocks -p 3333
```

### Start with a single file

```bash
fold -f ./data.json
```

### Initialize project with Drive

```bash
fold -i -dr FOLDER_ID -c credentials.json
```

### Create empty project

```bash
fold --init --template blanc --dir examples/blanc
```

### Migrate from Drive to local folder

```bash
fold -m -s drive -ds dir -dr FOLDER_ID -d ./localdata
```
---

# 📥 Record File Format
The `--record` flag uses a structured JSON file that defines how requests are made, sanitized, and authenticated.
```javascript
{
  "invocations": [ //list of requests
    {
      "url": "https://places.googleapis.com/v1/places:searchNearby", //request url
      "method": "POST", //request method
      "data": { //request body
        "includedTypes": ["restaurant"],
        "maxResultCount": 20,
        "locationRestriction": {
          "circle": {
            "center": {
              "latitude": 39.45434,
              "longitude": -0.33025
            },
            "radius": 1500.0
          }
        }
      },
      "headers": { //request headers
        "X-Goog-FieldMask": "*"
      },
      "securitySchemes": { //request authentication algorithm
        "apiKey": { //schema name
          "type": "ApiKey", //type : ApiKey | https | http 
          "in": "headers", //pass credentials via Headers | Query
          "name": "X-Goog-Api-Key", //name of Header or query parameter
          "scheme": "" //Scheme to be passed along with a token for example 'Bearer'
        }
      },
      "sanitize": { //How response data should be sanitized 
        "displayName": { //Field name to process
          "method": "randomize", //Replacement method randomize | erase(default)
          "combine": 2, //How random data should be combined. 1 is default - means no data combination
          "values": [ //Values used for randomization
            ["John", "Jane", "Martin", "Patricia"],
            ["Smith", "Doe", "Jones", "Brown"]
          ],
          "parents": ["authorAttributions"] //if provided the field would be processed only if parents found in full path to data
        },
        "nationalPhoneNumber": {
          // no method - values would be erased
        },
        "internationalPhoneNumber": {
          "method": "randomize",
          "values": [
            ["+34 961 00 00 01", "+34 961 00 00 02", "+34 961 00 00 03"]
          ]
        }
      }
    }
  ],
  "credentials": {
    // Credentials could be placed here with schema names as values. For example "apiKey": "your-api-key"
    // If some securityScheme doesn't have credentials provided record process would stop and prompt for input
  }
}
```

---
# 📁 Server File Structure
Folder structure provides all backend service settings.


## ⚙️Core Components

### 1. `project.json` (Required)
```json
{
"name": "project name",
"description": "",
"version": "1.0.0",
"allow_origin": "http://localhost:3333",
"auth_providers": null,
"jwt_secret": "auto generated secret",
"guest_password": "auto generated password"
}
```
JWT secret and guest role password are autogenerated on --init action run. A user could change them.
Server restart required to load project.json changes. 

### 📁 2. Port Folders (e.g., `3333`)
- Numerically named folders represent different ports the service will listen to
- All files inside are served at paths matching their relative location
- Special routing rules:
    - Extensions like `.json` and `.csv` are stripped from routes
    - `index` files are served at the parent path (e.g., `/user/index.csv` → `/user`)

## Reserved Directories

| Directory | Purpose |
|-----------|---------|
| `security` | Backend security configuration |
| `user` | User database and authentication |
| `www` | External API settings |
| `providers` | External authentication services |
| `src` | Frontend files (automatically updated by backend) |

### Security Folder (`security/`)
- `audience.csv`: Defines protected resources
- `authorities.csv`: Specifies user roles
- `rules.csv`: Authorization rules mapping roles to resources

### User Folder (`user/`)
- `index.csv`: User records
- `roles.csv`: User role assignments
- `oauth.csv`: External service authentication records

## Custom API Structure

- All non-reserved folders can be used to create custom APIs
- Follows the same routing rules as port folders
- Supports nested directory structures

## File Type Handling

| Extension | Behavior |
|-----------|----------|
| `.csv` | Transformed into tables with CRUD operations |
| `.json` | Transformed into NoSQL collections with CRUD operations |
| `.fold` | Represents built-in actions |
| `.drive` | Links to external Google Drive folders |
| Others | Served as-is (`.js`, `.css`, `.html` are cached by default) |

## Example API Endpoints

1. `pub/employees.json` → `/pub/employees` (NoSQL collection)
2. `pub/project.csv` → `/pub/project` (Database table)
3. `user/index.csv` → `/user` (User records table)
4. `security/rules.csv` → `/security/rules` (Authorization rules)

## Frontend Files (`src/`)
- Contains client-side JavaScript, CSS, and other assets
- Backend automatically updates:
    - `fold.js`: Service-specific functionality
    - `project.js`: Project configuration

## File structure examples

use `fold --init` or `fold --init --template public` commands to create basic file structures

---

## 🐳 Running with Docker

You can run Fold Server using the official Docker image:

```yaml
services:
  fold:
    image: ognerezov/fold:latest
    environment:
      - INPUT_DIR=/service/dataPath
    ports:
      - '3333:3333'
    volumes:
      - /local-folder:/service:rw
```

- Pass arguments via environmental variables `INPUT_<V>` maps to `--<v>`. For example: `INPUT_DIR` maps to `--dir` flag and should point to data folder inside the container.
- Mount local data and use `/<mounted>/<dataPath>` as internal project directory.
- Export ports your service uses.

Run directly:

```bash
docker run -it --rm \
  -e INPUT_DIR=/service/data \
  -v $(pwd)/mocks:/service \
  -p 3333:3333 \
  ognerezov/fold:latest
```

---

## 🧪 GitHub Actions Integration

Use Fold Server in your CI pipeline with the [ognerezov/fold-test-node](https://github.com/ognerezov/fold-test-node) GitHub Action:

```yaml
jobs:
  integration_tests:
    runs-on: ubuntu-latest
    name: A job to test with mock server
    steps:
      - name: Checkout Repository
        uses: actions/checkout@v3

      - name: Test
        id: fold
        uses: ognerezov/fold-test-node@0.3
        with:
          dir: '/github/workspace/mock'
          work_dir: '/github/workspace'
          test: 'npm run test'
          run: 'npm run dev'  # Optional: run your frontend app if needed
```

### 🔧 Action Parameters

| Key       | Description                                  |
|-----------|----------------------------------------------|
| `dir`     | Path to folder with mock data                |
| `work_dir`| Directory where test command will execute    |
| `test`    | Command to run tests                         |
| `run`     | (Optional) Command to run frontend app       |

# License

Fold Server is licensed under the [GPLv3](http://choosealicense.com/licenses/gpl-3.0) license for all open source applications.
A commercial license is required for all commercial applications (including sites, themes and apps you plan to sell).

# Contacts
If you have any questions regarding technical details, support, or licensing for our software, feel free to reach out — we're here to help!

📩 **Contact:** [ognerezov@foldserver.net](mailto:ognerezov@foldserver.net) for prompt and friendly assistance with:

- 🛠️ Technical guidance or troubleshooting
- 🤝 Integration and usage support
- 📄 Licensing terms and commercial use cases

We look forward to hearing from you!

(c) Sergey Okhotnikov 2025