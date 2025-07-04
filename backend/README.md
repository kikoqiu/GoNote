# Markdown 文件管理 Web 服务 - 技术文档

本文档详细说明了一个基于 Go 语言的 Web 服务，该服务用于管理多用户的 Markdown 文件，并提供版本控制、附件管理和全文搜索等功能。

## 1. 程序结构

项目采用单文件结构（`main.go`），但在逻辑上划分为多个模块，以便于理解和维护。

### 1.1. 逻辑模块划分

代码在 `main.go` 文件内通过注释块进行逻辑上的区分：

-   **`config.go`**: 负责处理程序的配置。从 `config.json` 文件加载，并允许通过命令行参数覆盖。
-   **`auth.go`**: 实现简单的基于文件的 Basic Authentication。用户凭证存储在 `users.txt` 中。
-   **`store.go`**: 实现一个内存中的文件缓存 (`InMemoryStore`)，用于加速文件读取和搜索。
-   **`file_monitor.go`**: 使用 `fsnotify` 库监控文件系统的变更，并实时更新内存缓存。
-   **`versioning.go`**: 基于 `bbolt` (BoltDB) 实现文件的增量版本控制。
-   **`search.go`**: 提供基于内存缓存的实时全文搜索功能。
-   **`handlers.go`**: 包含所有 HTTP 请求的处理函数 (Handlers)，是 API 逻辑的核心。
-   **`utils.go`**: 提供一些辅助工具函数，如 SHA1 计算、证书生成等。
-   **`main.go (entry point)`**: 程序的入口，负责初始化、设置路由和启动服务。

### 1.2. 文件系统布局

程序运行时，会自动创建并使用以下目录结构：

```
.
├── config.json          # 主配置文件
├── users.txt            # 用户认证文件
├── visit.log            # 访问日志
├── www/                 # 静态 Web 文件根目录
│   └── index.html       # 默认首页
├── markdown/            # Markdown 文件存储根目录
│   └── [username]/      # 各用户的独立目录
│       ├── _resource/   # (可选) 用户级的共享资源目录
│       ├── notes/
│       │   ├── doc1.md
│       │   └── doc1.md.attach/  # doc1.md 的附件目录
│       │       └── image.png
│       └── .extra/            # 系统内部使用的特殊目录
│           ├── versions.db    # 版本历史数据库
│           └── .recycle/      # 回收站
│               └── [sha1]/
└── main.go              # 程序主文件
```

## 2. 核心业务逻辑

### 2.1. 认证与授权

-   采用 **HTTP Basic Authentication**。
-   用户凭证在启动时从 `users.txt` 加载。文件格式为 `username password`，每行一个用户。
-   所有对 `/api/` 的请求都必须通过认证。
-   每个用户只能访问和操作其在 `markdown/[username]/` 目录下的文件。

### 2.2. 文件缓存与监控

-   **启动时**: 程序会完整扫描 `markdown` 目录，将所有 `.md` 文件的内容和 SHA1 哈希值加载到内存中，构建 `InMemoryStore` 缓存。
-   **运行时**: 一个后台 `goroutine` 使用 `fsnotify` 监控 `markdown` 目录。任何外部对文件的创建、修改、删除操作都会被捕获，并实时同步到内存缓存中，确保数据的一致性。
-   **API 操作**: 所有通过 API 对文件的写操作（创建、修改、删除、重命名）也会同步更新内存缓存。

### 2.3. 版本控制

-   **触发**: 当一个**已存在**的 `.md` 文件被修改时触发。
-   **数据库**: 每个用户在自己的 `.extra/` 目录下都有一个独立的 `versions.db` (BoltDB) 文件。
-   **存储策略**:
    -   **差量备份 (Patch)**: 默认情况下，系统会计算新旧文件内容的差异（patch），并存储这个 patch。
    -   **全量备份 (Full)**: 每隔 **50** 次差量备份，系统会自动进行一次全量备份，即将文件的完整内容存入版本库。这可以避免恢复历史版本时需要应用过多的 patch，从而提高恢复效率。
    -   版本链通过文件的 SHA1 哈希值关联。

### 2.4. 附件管理

-   **存储位置**:
    -   **本地附件**: 与特定 Markdown 文件关联的附件存储在同目录下，名为 `[markdown_filename].md.attach/` 的文件夹内。
    -   **共享附件**: 用户可以创建共享目录（如 `_resource/`），并通过相对路径 (`../../_resource/image.png`) 在任何 Markdown 文件中引用这些资源。
-   **访问权限**:
    -   **读取**: 极其灵活。允许通过 API 读取用户目录下的**任何**文件，只要提供了正确的相对路径。
    -   **写入/删除**: 严格受限。API 只允许在本地附件目录 (`.attach/`) 中创建和删除文件，防止对共享资源或其他文件的意外修改。

### 2.5. 搜索

-   **机制**: 搜索操作完全基于内存中的 `InMemoryStore` 缓存进行，速度极快。
-   **模式**: 支持多关键字搜索和正则表达式搜索。
-   **结果**: 返回匹配文件的路径以及匹配内容所在的上下文行。

## 3. API 参数约定与调用示例

所有 API 的根路径为 `/api`。

---

### 3.1. 目录操作 (`/api/dir`)

-   **Endpoint**: `/api/dir`
-   **Method**: `POST`
-   **Body (JSON)**:
    -   `action`: (string) "create", "delete", "rename"
    -   `path`: (string) 目标目录路径，相对于用户根目录。
    -   `new_path`: (string, optional) 重命名时的新路径。
-   **示例 (创建目录)**:
    ```bash
    curl -u "user:pass" -X POST -H "Content-Type: application/json" \
         -d '{"action": "create", "path": "new_project/assets"}' \
         https://localhost:8080/api/dir
    ```

---

### 3.2. 列表 (`/api/list`)

-   **Endpoint**: `/api/list`
-   **Method**: `GET`
-   **Query Parameters**:
    -   `path`: (string, optional) 要列出的目录路径。如果为空，则列出用户根目录。
-   **示例**:
    ```bash
    curl -u "user:pass" "https://localhost:8080/api/list?path=notes"
    ```

---

### 3.3. 文件读写 (`/api/file`)

#### 创建/修改文件

-   **Method**: `POST`
-   **Body (JSON)**:
    -   `path`: (string) 文件路径，**必须**以 `.md` 结尾。
    -   `content`: (string) 文件内容。
    -   `comment`: (string, optional) 版本备注。
-   **示例**:
    ```bash
    curl -u "user:pass" -X POST -H "Content-Type: application/json" \
         -d '{"path": "notes/idea.md", "content": "# My Great Idea"}' \
         https://localhost:8080/api/file
    ```

#### 读取文件

-   **Method**: `GET`
-   **Query Parameters**:
    -   `path`: (string) 要读取的文件路径。
-   **示例**:
    ```bash
    curl -u "user:pass" "https://localhost:8080/api/file?path=notes/idea.md"
    ```

#### 重命名/删除文件

-   **Method**: `PATCH`
-   **Body (JSON)**:
    -   `action`: (string) "rename", "delete"
    -   `path`: (string) 目标文件路径。
    -   `new_path`: (string, optional) 重命名时的新路径。
-   **示例 (重命名)**:
    ```bash
    curl -u "user:pass" -X PATCH -H "Content-Type: application/json" \
         -d '{"action": "rename", "path": "notes/idea.md", "new_path": "notes/final_idea.md"}' \
         https://localhost:8080/api/file
    ```

---

### 3.4. 附件操作 (`/api/attach/*`)

#### 上传附件

-   **Endpoint**: `/api/attach/upload`
-   **Method**: `POST`
-   **Body**: `multipart/form-data`
    -   `path`: (string) 关联的 Markdown 文件路径。
    -   `attachment`: (file) 要上传的文件。
-   **示例**:
    ```bash
    curl -u "user:pass" -X POST \
         -F "path=notes/doc.md" \
         -F "attachment=@/path/to/local/image.png" \
         https://localhost:8080/api/attach/upload
    ```
-   **成功响应**: `{"status": "success", "mdPath": "notes/doc.md", "attachPath": "doc.md.attach/image.png"}`

#### 列出附件

-   **Endpoint**: `/api/attach/list`
-   **Method**: `GET`
-   **Query Parameters**:
    -   `path`: (string) 关联的 Markdown 文件路径。
-   **示例**:
    ```bash
    curl -u "user:pass" "https://localhost:8080/api/attach/list?path=notes/doc.md"
    ```

#### 获取/下载附件

-   **Endpoint**: `/api/attach/get/[path]`
-   **Method**: `GET`
-   **Query Parameters**:
-   **示例 (本地附件)**:
    ```bash
    curl -u "user:pass" "https://localhost:8080/api/attach/get/notes/doc.md/../doc.md.attach/image.png"
    ```
-   **示例 (共享附件)**:
    ```bash
    curl -u "user:pass" "https://localhost:8080/api/attach/get/notes/doc.md/../../_resource/shared.jpg"
    ```

#### 删除附件

-   **Endpoint**: `/api/attach/delete`
-   **Method**: `POST`
-   **Body (JSON)**:
    -   `mdPath`: (string) 关联的 Markdown 文件路径。
    -   `attachPath`: (string) 要删除的附件的相对路径 (**必须**在 `.attach/` 目录内)。
-   **示例**:
    ```bash
    curl -u "user:pass" -X POST -H "Content-Type: application/json" \
         -d '{"mdPath": "notes/doc.md", "attachPath": "doc.md.attach/image.png"}' \
         https://localhost:8080/api/attach/delete
    ```

---

### 3.5. 其他 API

-   **搜索 (`/api/search`)**: `GET` 请求，参数 `q` (搜索词) 和 `regex` (bool)。
-   **历史 (`/api/history`)**: `GET` 请求，参数 `path` (文件路径)。
-   **版本 (`/api/version`)**: `GET` 请求，参数 `path` (文件路径) 和 `id` (版本ID)。

## 4. 函数功能说明

### `main()`
程序主入口。负责调用初始化函数、设置 Chi 路由、应用中间件（日志、CORS、认证）并启动 HTTP 或 HTTPS 服务器。

### `LoadConfig()`
-   **功能**: 初始化程序配置。
-   **逻辑**:
    1.  设置默认配置。
    2.  尝试从 `config.json`读取配置并覆盖默认值。
    3.  如果 `config.json` 不存在，则用默认值创建它。
    4.  解析命令行参数，再次覆盖现有配置。

### `LoadUsers()`
-   **功能**: 加载用户认证信息。
-   **逻辑**: 从 `users.txt` 文件逐行读取 `username password` 对，并存入 `userCredentials` map 中。

### `AuthMiddleware()`
-   **功能**: Chi 中间件，用于保护 API 路由。
-   **逻辑**:
    1.  解析 HTTP 请求头中的 Basic Auth 信息。
    2.  与内存中的 `userCredentials` 进行比对。
    3.  验证通过，则将用户名存入请求的 `context` 中，并传递给下一个 handler。
    4.  验证失败，则返回 `401 Unauthorized`。

### `InMemoryStore.Scan()`
-   **功能**: 扫描文件系统，构建初始的内存文件缓存。
-   **逻辑**: 递归遍历每个用户的目录，读取所有 `.md` 文件内容，计算 SHA1，并存入 `store.docs` map。

### `WatchMarkdownDir()`
-   **功能**: 启动一个后台 goroutine 来监控文件系统。
-   **逻辑**: 使用 `fsnotify` 监听文件事件 (`Create`, `Write`, `Remove`, `Rename`)，并相应地调用 `store.UpdateDoc` 或 `store.DeleteDoc` 来更新内存缓存。

### `NewVersionManager()`
-   **功能**: 为指定用户创建一个版本管理器实例。
-   **逻辑**: 打开或创建用户目录下的 `versions.db` (BoltDB) 文件，并确保用于存储备份的 bucket 存在。

### `VersionManager.CreateBackup()`
-   **功能**: 为一次文件修改创建版本记录。
-   **逻辑**:
    1.  判断是进行全量备份还是差量备份。
    2.  如果是差量备份，使用 `diffmatchpatch` 库生成 patch 文本。
    3.  将版本信息（SHA1、patch、时间等）打包成 `VersionRecord` 结构体。
    4.  序列化为 JSON 并存入 BoltDB。

### `getSafeAttachmentPath(r, mdPath, attachPath)`
-   **功能**: **核心安全函数**。解析并验证附件路径。
-   **参数**:
    -   `r`: `*http.Request`，用于获取用户信息。
    -   `mdPath`: 关联的 Markdown 文件路径（相对于用户根目录）。
    -   `attachPath`: 附件相对于 `mdPath` 的路径。
-   **逻辑**:
    1.  获取用户的绝对根目录路径 (`userBasePath`)。
    2.  获取 `mdPath` 所在目录的绝对路径 (`mdDirAbsPath`)。
    3.  将 `mdDirAbsPath` 和 `attachPath` 拼接起来，得到一个初步的绝对路径。
    4.  使用 `filepath.Clean()` 解析路径中的 `..` 等，得到最终的绝对路径 `cleanedPath`。
    5.  **安全校验**: 检查 `cleanedPath` 是否仍然在 `userBasePath` 的子路径下。如果不是，说明发生了目录遍历攻击，立即返回错误。
    6.  返回安全的、在文件系统上真实存在的绝对路径。

### `handle*` 函数
所有以 `handle` 开头的函数都是 Chi 路由的 HTTP handler。它们负责：
1.  解析请求参数（URL query, JSON body, form data）。
2.  调用核心业务逻辑函数（如 `getSafeAttachmentPath`, `vm.CreateBackup` 等）。
3.  调用 `respondJSON` 或 `respondError` 向客户端返回标准格式的响应。

