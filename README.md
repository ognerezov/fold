# Installation
Download fold binary for your platform or use ognerezov/fold docker image

#  File Structure
Folder structure provides all backend service settings.


## Core Components

### 1. `project.json` (Required)
- Contains project configuration:
    - Project name
    - Secrets
    - Allowed origins
    - Other global settings

### 2. Port Folders (e.g., `3333`)
- Numerically named folders represent different ports the service will listen on
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

# Command Line Options

## Server Options
| Option | Shorthand | Type | Description                      | Default |
|--------|-----------|------|----------------------------------|---------|
| `--api` | `-a` | string | Server base path. v1 for example | |
| `--port` | `-p` | int | Port to listen on                | `3333` |
| `--origin` | `-o` | string | Allow origin                     | |
| `--reg` | | bool | User registration allowed        | `true` |

## File & Directory Options
| Option | Shorthand | Type | Description                                                        | Default |
|--------|-----------|------|--------------------------------------------------------------------|---------|
| `--credentials` | `-c` | string | Credentials file. Required for serving from Google Drive           | |
| `--dir` | `-d` | string | Working directory. Required for serving from local or cloud folder | |
| `--file` | `-f` | string | Serve or init with single file                                     | |

## Project Options
| Option | Shorthand | Type | Description                               | Default |
|--------|-----------|------|-------------------------------------------|---------|
| `--name` | `-n` | string | Application name                          | `"fold"` |
| `--description` | | string | Project description                       | |
| `--version` | | string | Project version                           | `"1.0.0"` |
| `--template` | `-t` | string | Project template. Used with --init action | `"default"` |

## Migration Options
| Option | Shorthand | Type | Description | Default |
|--------|-----------|------|-------------|---------|
| `--source` | `-s` | string | Data migration source type | `"dir"` |
| `--destination` | `-ds` | string | Data migration destination type | `"drive"` |
| `--drive` | `-dr` | string | Google Drive folder ID | |

## Actions
| Option | Shorthand | Type | Description |
|--------|-----------|------|-------------|
| `--init` | `-i` | bool | Init new project folder |
| `--migrate` | `-m` | bool | Migrate data |
| `--cache` | | bool | Cache files requests | `true` |

## Help
| Option | Shorthand | Type | Description |
|--------|-----------|------|-------------|
| `--help` | `-h` | bool | Show help |


## License

Fold Server is licensed under the [GPLv3](http://choosealicense.com/licenses/gpl-3.0) license for all open source applications.
A commercial license is required for all commercial applications (including sites, themes and apps you plan to sell).