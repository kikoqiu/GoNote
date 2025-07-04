# Markdown File Management Web Service - Technical Documentation

This document details a Go-based web service for managing multi-user Markdown files, providing features like version control, attachment management, and full-text search.

## 1. Program Structure

The project uses a single-file structure (`main.go`) but is logically divided into multiple modules for better understanding and maintenance.

### 1.1. Logical Module Division

The code within `main.go` is logically separated by comment blocks:

-   **`config.go`**: Handles program configuration. It loads from `config.json` and allows overrides via command-line arguments.
-   **`auth.go`**: Implements simple file-based Basic Authentication. User credentials are stored in `users.txt`.
-   **`backup.go`**: Implements scheduled backup tasks based on `cron` expressions, responsible for packaging the markdown directory and managing backup files.
-   **`store.go`**: Implements an in-memory file cache (`InMemoryStore`) to speed up file reading and searching.
-   **`file_monitor.go`**: Uses the `fsnotify` library to monitor file system changes and update the in-memory cache in real-time.
-   **`versioning.go`**: Implements incremental version control for files based on `bbolt` (BoltDB).
-   **`search.go`**: Provides real-time full-text search functionality based on the in-memory cache.
-   **`handlers.go`**: Contains all HTTP request handlers, forming the core of the API logic.
-   **`utils.go`**: Provides auxiliary utility functions, such as SHA1 calculation, certificate generation, etc.
-   **`main.go (entry point)`**: The program's entry point, responsible for initialization, setting up routes, and starting the service.

### 1.2. File System Layout

When the program runs, it automatically creates and uses the following directory structure:

```
.
├── config.json          # Main configuration file
├── users.txt            # User authentication file
├── visit.log            # Access log
├── www/                 # Static web files root directory
│   └── index.html       # Default homepage
├── markdown/            # Markdown files root directory
│   └── [username]/      # Independent directory for each user
│       ├── _resource/   # (Optional) User-level shared resource directory
│       ├── notes/
│       │   ├── doc1.md
│       │   └── doc1.md.attach/  # Attachment directory for doc1.md
│       │       └── image.png
│       └── .extra/            # Special directory for internal system use
│           ├── versions.db    # Version history database
│           └── .recycle/      # Recycle bin
│               └── [sha1]/
├── backup/              # Directory for automatic backup files
│   └── markdown-2023-10-27T15-04-05.zip
└── main.go              # Main program file
```

## 2. Core Business Logic

### 2.1. Authentication and Authorization

-   Uses **HTTP Basic Authentication**.
-   User credentials are loaded from `users.txt` at startup. The format is `username password`, one user per line.
-   All requests to `/api/` must be authenticated.
-   Each user can only access and operate on files within their `markdown/[username]/` directory.

### 2.2. File Caching and Monitoring

-   **At Startup**: The program completely scans the `markdown` directory, loading the content and SHA1 hash of all `.md` files into memory to build the `InMemoryStore` cache.
-   **At Runtime**: A background `goroutine` uses `fsnotify` to monitor the `markdown` directory. Any external file creation, modification, or deletion is captured and synchronized with the in-memory cache in real-time to ensure data consistency.
-   **API Operations**: All write operations via the API (create, modify, delete, rename) also update the in-memory cache.

### 2.3. Version Control

-   **Trigger**: Triggered when an **existing** `.md` file is modified.
-   **Database**: Each user has an independent `versions.db` (BoltDB) file in their `.extra/` directory.
-   **Storage Strategy**:
    -   **Differential Backup (Patch)**: By default, the system calculates the difference (patch) between the new and old file content and stores this patch.
    -   **Full Backup**: Every **50** differential backups, the system automatically performs a full backup, storing the complete file content in the version repository. This avoids applying too many patches when restoring a historical version, thus improving recovery efficiency.
    -   The version chain is linked by the file's SHA1 hash.

### 2.4. Attachment Management

-   **Storage Location**:
    -   **Local Attachments**: Attachments associated with a specific Markdown file are stored in a folder named `[markdown_filename].md.attach/` in the same directory.
    -   **Shared Attachments**: Users can create shared directories (e.g., `_resource/`) and reference these resources in any Markdown file using a relative path (`../../_resource/image.png`).
-   **Access Permissions**:
    -   **Read**: Extremely flexible. Allows reading **any** file in the user's directory via the API, as long as the correct relative path is provided.
    -   **Write/Delete**: Strictly limited. The API only allows creating and deleting files in the local attachment directory (`.attach/`) to prevent accidental modification of shared resources or other files.

### 2.5. Search

-   **Mechanism**: Search operations are performed entirely on the in-memory `InMemoryStore` cache, making them extremely fast.
-   **Modes**: Supports multi-keyword search and regular expression search.
-   **Results**: Returns the paths of matching files and the context lines where the matches occurred.

### 2.6. Automatic Backup and Cleanup

-   **Configuration**: Configured via the `backup` object in `config.json`. You can enable/disable, set the backup directory, CRON expression, and retention days.
-   **Workflow**:
    1.  If `backup.enabled` is `true`, a CRON scheduler is initialized at service startup.
    2.  The `performBackup` function is triggered periodically according to the `backup.cron` expression, compressing the entire `markdown_dir` into a timestamped zip file and storing it in `backup.dir`.
    3.  The scheduler also triggers the `performBackupCleanup` function at a fixed time daily (1 AM) to check for and delete any old backup files that exceed the `backup.retention_days` limit.

## 3. API Parameter Conventions and Call Examples

All API root paths are `/api`.

---

### 3.1. Directory Operations (`/api/dir`)

-   **Endpoint**: `/api/dir`
-   **Method**: `POST`
-   **Body (JSON)**:
    -   `action`: (string) "create", "delete", "rename"
    -   `path`: (string) Target directory path, relative to the user's root directory.
    -   `new_path`: (string, optional) The new path for a rename operation.
-   **Example (Create Directory)**:
    ```bash
    curl -u "user:pass" -X POST -H "Content-Type: application/json" \
         -d '{"action": "create", "path": "new_project/assets"}' \
         https://localhost:8080/api/dir
    ```
-   **Success Response (JSON)**:
    ```json
    {
      "status": "success"
    }
    ```

---

### 3.2. List (`/api/list`)

-   **Endpoint**: `/api/list`
-   **Method**: `GET`
-   **Query Parameters**:
    -   `path`: (string, optional) The directory path to list. If empty, lists the user's root directory.
    -   `recursive`: (bool, optional) Whether to list all content recursively, defaults to `false`.
-   **Example**:
    ```bash
    curl -u "user:pass" "https://localhost:8080/api/list?path=notes"
    ```
-   **Success Response (JSON)**:
    Returns an array of `TreeItem` objects.
    ```json
    [
      {
        "name": "Doc",
        "is_dir": true,
        "size": 4096,
        "mod_time": "2023-10-27T08:10:25.123Z",
        "children": null
      },
      {
        "name": "example.md",
        "is_dir": false,
        "size": 1234,
        "mod_time": "2023-10-27T09:15:00.456Z",
        "attach_count": 1
      }
    ]
    ```
    **`TreeItem` Object Property Description**:
    | Property | Type | Description |
    | :--- | :--- | :--- |
    | `name` | string | File or directory name. |
    | `is_dir` | bool | Whether it is a directory. |
    | `size` | int64 | File size (bytes). |
    | `mod_time` | string | Last modification time (RFC3339 format). |
    | `attach_count` | int | (Files only) Number of associated attachments. |
    | `children` | array | (Only if `recursive=true`) Array of child `TreeItem` objects, otherwise `null`. |

---

### 3.3. File Read/Write (`/api/file`)

#### Create/Modify File

-   **Method**: `POST`
-   **Body (JSON)**:
    -   `path`: (string) File path, **must** end with `.md`.
    -   `content`: (string) File content.
    -   `comment`: (string, optional) Version comment.
-   **Example**:
    ```bash
    curl -u "user:pass" -X POST -H "Content-Type: application/json" \
         -d '{"path": "notes/idea.md", "content": "# My Great Idea"}' \
         https://localhost:8080/api/file
    ```
-   **Success Response (JSON)**:
    ```json
    {
      "status": "success",
      "sha1": "bf8b4533d759389c9684b3b1904791550c4b31a8"
    }
    // Or when the file content has not changed
    {
      "status": "no change"
    }
    ```

#### Read File

-   **Method**: `GET`
-   **Query Parameters**:
    -   `path`: (string) The path of the file to read.
-   **Example**:
    ```bash
    curl -u "user:pass" "https://localhost:8080/api/file?path=notes/idea.md"
    ```
-   **Success Response (JSON)**:
    Returns a `Document` object.
    ```json
    {
      "Path": "notes/idea.md",
      "SHA1": "bf8b4533d759389c9684b3b1904791550c4b31a8",
      "Content": "# My Great Idea"
    }
    ```

#### Rename/Delete File

-   **Method**: `PATCH`
-   **Body (JSON)**:
    -   `action`: (string) "rename", "delete"
    -   `path`: (string) Target file path.
    -   `new_path`: (string, optional) The new path for a rename operation.
-   **Example (Rename)**:
    ```bash
    curl -u "user:pass" -X PATCH -H "Content-Type: application/json" \
         -d '{"action": "rename", "path": "notes/idea.md", "new_path": "notes/final_idea.md"}' \
         https://localhost:8080/api/file
    ```
-   **Success Response (JSON)**:
    ```json
    {
      "status": "success"
    }
    ```

---

### 3.4. Attachment Operations (`/api/attach/*`)

#### Upload Attachment

-   **Endpoint**: `/api/attach/upload`
-   **Method**: `POST`
-   **Body**: `multipart/form-data`
    -   `path`: (string) The path of the associated Markdown file.
    -   `attachment`: (file) The file to upload.
-   **Example**:
    ```bash
    curl -u "user:pass" -X POST \
         -F "path=notes/doc.md" \
         -F "attachment=@/path/to/local/image.png" \
         https://localhost:8080/api/attach/upload
    ```
-   **Success Response (JSON)**:
    ```json
    {
        "status": "success", 
        "mdPath": "notes/doc.md", 
        "attachPath": "doc.md.attach/image.png"
    }
    ```

#### List Attachments

-   **Endpoint**: `/api/attach/list`
-   **Method**: `GET`
-   **Query Parameters**:
    -   `path`: (string) The path of the associated Markdown file.
-   **Example**:
    ```bash
    curl -u "user:pass" "https://localhost:8080/api/attach/list?path=notes/doc.md"
    ```
-   **Success Response (JSON)**:
    ```json
    {
      "mdPath": "notes/doc.md",
      "attachments": [
        {
          "name": "image.png",
          "attachPath": "doc.md.attach/image.png",
          "size": 102400,
          "mod_time": "2023-10-27T10:30:00.789Z"
        }
      ]
    }
    ```
    **`AttachmentInfo` Object Property Description**:
    | Property | Type | Description |
    | :--- | :--- | :--- |
    | `name` | string | Attachment file name. |
    | `attachPath` | string | Relative path that can be used for downloading. |
    | `size` | int64 | File size (bytes). |
    | `mod_time` | string | Last modification time (RFC3339 format). |

#### Get/Download Attachment

-   **Endpoint**: `/api/attach/get/[path]`
-   **Method**: `GET`
-   **Description**: Returns the file stream directly. The path is relative to the user's root directory.
-   **Example (Local Attachment)**: `.../api/attach/get/notes/doc.md.attach/image.png`
-   **Example (Shared Attachment)**: `.../api/attach/get/_resource/shared.jpg`

#### Delete Attachment

-   **Endpoint**: `/api/attach/delete`
-   **Method**: `POST`
-   **Body (JSON)**:
    -   `mdPath`: (string) The path of the associated Markdown file.
    -   `attachPath`: (string) The relative path of the attachment to delete (**must** be within the `.attach/` directory).
-   **Example**:
    ```bash
    curl -u "user:pass" -X POST -H "Content-Type: application/json" \
         -d '{"mdPath": "notes/doc.md", "attachPath": "doc.md.attach/image.png"}' \
         https://localhost:8080/api/attach/delete
    ```
-   **Success Response (JSON)**:
    ```json
    {
      "status": "success"
    }
    ```

---

### 3.5. Other APIs

-   **Search (`/api/search`)**: `GET` request, parameters `q` (search term) and `regex` (bool).
    - **Success Response (JSON)**: `[{"path": "file/path.md", "context": ["12: matching line content..."]}]`
-   **History (`/api/history`)**: `GET` request, parameter `path` (file path).
    - **Success Response (JSON)**: Returns an array of `VersionRecord` objects.
-   **Version (`/api/version`)**: `GET` request, parameters `path` (file path) and `id` (version ID).
    - **Success Response (JSON)**: `{"content": "Content of the specific version"}`

## 4. Function Descriptions

### `main()`
The main entry point of the program. It is responsible for calling initialization functions, setting up Chi routes, applying middleware (logging, CORS, authentication), and starting the HTTP or HTTPS server. It now also calls `StartBackupScheduler` to start the backup scheduler.

### `LoadConfig()`
-   **Function**: Initializes the program configuration.
-   **Logic**:
    1.  Sets default configurations.
    2.  Tries to read the configuration from `config.json` and override the defaults.
    3.  If `config.json` does not exist, or if new fields (like `backup`) are missing, it creates/updates it with default or existing values.
    4.  Parses command-line arguments to override existing configurations again.

### `StartBackupScheduler()`
-   **Function**: Initializes and starts the background backup tasks.
-   **Logic**:
    1.  Checks if `AppConfig.Backup.Enabled` is `true`.
    2.  If enabled, creates a `cron` instance.
    3.  Adds two scheduled tasks: one to execute `performBackup` according to `AppConfig.Backup.Cron`, and another to execute `performBackupCleanup` daily at a fixed time.
    4.  Starts the scheduler in a background goroutine.

### `performBackup()`
-   **Function**: Performs a full backup.
-   **Logic**: Packages all content under the `AppConfig.MarkdownDir` directory into a timestamped `.zip` file and stores it in the `AppConfig.Backup.Dir` directory.

### `performBackupCleanup()`
-   **Function**: Cleans up expired backup files.
-   **Logic**: Iterates through the backup directory and deletes all backup files created earlier than the `AppConfig.Backup.RetentionDays` period.

### `LoadUsers()`
-   **Function**: Loads user authentication information.
-   **Logic**: Reads `username password` pairs line by line from the `users.txt` file and stores them in the `userCredentials` map.

### `AuthMiddleware()`
-   **Function**: A Chi middleware for protecting API routes.
-   **Logic`:
    1.  Parses the Basic Auth information from the HTTP request header.
    2.  Compares it with the in-memory `userCredentials`.
    3.  If authentication is successful, it stores the username in the request's `context` and passes it to the next handler.
    4.  If authentication fails, it returns a `401 Unauthorized`.

### `InMemoryStore.Scan()`
-   **Function**: Scans the file system to build the initial in-memory file cache.
-   **Logic**: Recursively traverses each user's directory, reads the content of all `.md` files, calculates their SHA1, and stores them in the `store.docs` map.

### `WatchMarkdownDir()`
-   **Function**: Starts a background goroutine to monitor the file system.
-   **Logic**: Uses `fsnotify` to listen for file events (`Create`, `Write`, `Remove`, `Rename`) and calls `store.UpdateDoc` or `store.DeleteDoc` accordingly to update the in-memory cache.

### `NewVersionManager()`
-   **Function**: Creates a version manager instance for a specified user.
-   **Logic**: Opens or creates the `versions.db` (BoltDB) file in the user's directory and ensures the bucket for storing backups exists.

### `VersionManager.CreateBackup()`
-   **Function**: Creates a version record for a file modification.
-   **Logic`:
    1.  Determines whether to perform a full or differential backup.
    2.  If it's a differential backup, it uses the `diffmatchpatch` library to generate a patch text.
    3.  Packages the version information (SHA1, patch, time, etc.) into a `VersionRecord` struct.
    4.  Serializes it to JSON and stores it in BoltDB.

### `getSafeAttachmentPath(r, mdPath, attachPath)`
-   **Function**: **Core security function**. Parses and validates attachment paths.
-   **Parameters**:
    -   `r`: `*http.Request`, used to get user information.
    -   `mdPath`: The path of the associated Markdown file (relative to the user's root directory).
    -   `attachPath`: The path of the attachment relative to `mdPath`.
-   **Logic**:
    1.  Gets the user's absolute root directory path (`userBasePath`).
    2.  Gets the absolute path of the directory containing `mdPath` (`mdDirAbsPath`).
    3.  Joins `mdDirAbsPath` and `attachPath` to get a preliminary absolute path.
    4.  Uses `filepath.Clean()` to resolve `..` etc., resulting in the final absolute path `cleanedPath`.
    5.  **Security Check**: Verifies that `cleanedPath` is still under `userBasePath`. If not, it indicates a directory traversal attack, and an error is returned immediately.
    6.  Returns the safe, real absolute path on the file system.

### `handle*` Functions
All functions starting with `handle` are HTTP handlers for the Chi router. They are responsible for:
1.  Parsing request parameters (URL query, JSON body, form data).
2.  Calling core business logic functions (e.g., `getSafeAttachmentPath`, `vm.CreateBackup`).
3.  Calling `respondJSON` or `respondError` to return a standard formatted response to the client.

```