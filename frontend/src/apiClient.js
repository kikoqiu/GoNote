// apiClient.js

/**
 * A client for interacting with the Markdown Web Service API.
 */
class ApiClient {
  /**
   * @param {string} baseURL - The base URL of the API server (e.g., 'http://localhost:8080').
   * @param {string} username - The username for Basic Authentication.
   * @param {string} password - The password for Basic Authentication.
   */
  constructor(baseURL, username, password) {
    this.baseURL = baseURL
    this.authHeader = 'Basic ' + btoa(`${username}:${password}`)
  }

  /**
   * A helper function to make authenticated API requests.
   * @private
   */
  async _request(endpoint, options = {}, isJson = true) {
    const url = `${this.baseURL}${endpoint}`

    const headers = {
      Authorization: this.authHeader,
      ...options.headers,
    }

    if (isJson) {
      headers['Content-Type'] = 'application/json'
    }

    headers['Api-Version'] = '0.1'

    const config = {
      ...options,
      headers,
    }

    try {
      const response = await fetch(url, config)

      if (!response.ok) {
        let errorData
        try {
          errorData = await response.json()
        } catch (e) {
          errorData = { error: 'An unknown error occurred', status: response.status }
        }
        throw new Error(errorData.error || `HTTP error! Status: ${response.status}`)
      }

      if (response.status === 204 || response.headers.get('content-length') === '0') {
        return null
      }

      if (!isJson) {
        if (options.body instanceof FormData) {
          return await response.json()
        }
        return response
      }

      return await response.json()
    } catch (error) {
      console.error('API Request Failed:', error)
      throw error
    }
  }

  // --- Directory and Listing Operations ---

  /**
   * Lists the contents of a directory.
   * @param {string} [path=''] - The path of the directory to list, relative to the user's root.
   * @param {boolean} [recursive=false] - If true, fetches the entire directory tree recursively.
   * @returns {Promise<Array<object>>} A list of file/directory objects.
   *
   * @example
   * // When recursive = false (or omitted), the return format is a flat array:
   * [
   *   {
   *     "name": "Notes",
   *     "is_dir": true,
   *     "size": 4096,
   *     "mod_time": "2023-10-27T10:00:00Z"
   *   },
   *   {
   *     "name": "README.md",
   *     "is_dir": false,
   *     "size": 1234,
   *     "mod_time": "2023-10-26T15:30:00Z",
   *     "attach_count": 2
   *   }
   * ]
   *
   * @example
   * // When recursive = true, the return format is a nested tree structure:
   * [
   *   {
   *     "name": "Notes",
   *     "is_dir": true,
   *     "size": 4096,
   *     "mod_time": "2023-10-27T10:00:00Z",
   *     "children": [
   *       {
   *         "name": "ProjectA",
   *         "is_dir": true,
   *         // ... other properties
   *         "children": [
   *           {
   *             "name": "status.md",
   *             "is_dir": false,
   *             // ... other properties
   *           }
   *         ]
   *       }
   *     ]
   *   }
   * ]
   */
  listDirectory(path = '', recursive = false) {
    const params = new URLSearchParams({ path })
    if (recursive) {
      params.append('recursive', 'true')
    }
    return this._request(`/api/list?${params.toString()}`)
  }

  /**
   * Creates a new directory.
   * @param {string} path - The path of the directory to create.
   * @returns {Promise<object>} A confirmation object.
   * @example { "status": "success" }
   */
  createDirectory(path) {
    return this._request('/api/dir', {
      method: 'POST',
      body: JSON.stringify({ action: 'create', path }),
    })
  }

  /**
   * Deletes a directory.
   * @param {string} path - The path of the directory to delete.
   * @returns {Promise<object>} A confirmation object.
   * @example { "status": "success" }
   */
  deleteDirectory(path) {
    return this._request('/api/dir', {
      method: 'POST',
      body: JSON.stringify({ action: 'delete', path }),
    })
  }

  /**
   * Renames a directory.
   * @param {string} oldPath - The current path of the directory.
   * @param {string} newPath - The new path for the directory.
   * @returns {Promise<object>} A confirmation object.
   * @example { "status": "success" }
   */
  renameDirectory(oldPath, newPath) {
    return this._request('/api/dir', {
      method: 'POST',
      body: JSON.stringify({ action: 'rename', path: oldPath, new_path: newPath }),
    })
  }

  // --- File Operations ---

  /**
   * Reads the content of a markdown file.
   * @param {string} path - The path of the file.
   * @returns {Promise<object>} An object containing the file's details.
   * @example
   * {
   *   "Path": "notes/MyNote.md",
   *   "SHA1": "a1b2c3d4e5f6...",
   *   "Content": "# My Note\n\nThis is the content of the note."
   * }
   */
  readFile(path) {
    const params = new URLSearchParams({ path })
    return this._request(`/api/file?${params.toString()}`)
  }

  /**
   * Gets the URL to download a file directly from the user's file space.
   * @param {string} path - The path of the file.
   * @returns {string} The full, usable URL to the file.
   */
  getDownloadUrl(path) {
    return this.getAttachmentURL('', path)
  }

  /**
   * Creates or updates a markdown file.
   * @param {string} path - The path of the file.
   * @param {string} content - The new content of the file.
   * @param {string} [comment=''] - An optional comment for the version history.
   * @returns {Promise<object>} A confirmation object with the new SHA1 hash, or a "no change" status.
   * @example { "status": "success", "sha1": "f6e5d4c3b2a1..." }
   * @example { "status": "no change" }
   */
  writeFile(path, content, comment = '') {
    if (!path.toLowerCase().endsWith('.md')) {
      return Promise.reject(new Error('File path must end with .md'))
    }
    return this._request('/api/file', {
      method: 'POST',
      body: JSON.stringify({ path, content, comment }),
    })
  }

  /**
   * Renames a markdown file.
   * @param {string} oldPath - The current path of the file.
   * @param {string} newPath - The new path for the file.
   * @returns {Promise<object>} A confirmation object.
   * @example { "status": "success" }
   */
  renameFile(oldPath, newPath) {
    if (!newPath.toLowerCase().endsWith('.md')) {
      return Promise.reject(new Error('New file path must end with .md'))
    }
    return this._request('/api/file', {
      method: 'PATCH',
      body: JSON.stringify({ action: 'rename', path: oldPath, new_path: newPath }),
    })
  }

  /**
   * Deletes a markdown file (moves it to a recycle bin on the server).
   * @param {string} path - The path of the file to delete.
   * @returns {Promise<object>} A confirmation object.
   * @example { "status": "success" }
   */
  deleteFile(path) {
    return this._request('/api/file', {
      method: 'PATCH',
      body: JSON.stringify({ action: 'delete', path }),
    })
  }

  // --- Attachment Operations ---

  /**
   * Uploads an attachment for a markdown file.
   * @param {string} mdPath - The path of the associated markdown file.
   * @param {File} fileObject - The file object to upload.
   * @returns {Promise<object>} A confirmation object with the attachment's path.
   * @example
   * {
   *   "status": "success",
   *   "mdPath": "notes/MyNote.md",
   *   "attachPath": "MyNote.md.attach/image.png"
   * }
   */
  uploadAttachment(mdPath, fileObject) {
    const formData = new FormData()
    formData.append('path', mdPath)
    formData.append('attachment', fileObject)

    return this._request(
      '/api/attach/upload',
      {
        method: 'POST',
        body: formData,
      },
      false
    )
  }

  /**
   * Lists all attachments for a markdown file.
   * @param {string} mdPath - The path of the associated markdown file.
   * @returns {Promise<object>} An object containing the list of attachments.
   * @example
   * {
   *   "mdPath": "notes/MyNote.md",
   *   "attachments": [
   *     {
   *       "name": "image.png",
   *       "attachPath": "MyNote.md.attach/image.png",
   *       "size": 51200,
   *       "mod_time": "2023-10-27T12:00:00Z"
   *     }
   *   ]
   * }
   */
  listAttachments(mdPath) {
    const params = new URLSearchParams({ path: mdPath })
    return this._request(`/api/attach/list?${params.toString()}`)
  }

  /**
   * Gets the direct URL for an attachment, suitable for download or src attributes.
   * @param {string} mdPath - The path of the associated markdown file.
   * @param {string} attachPath - The path of the attachment relative to the markdown file.
   * @returns {string} The full, usable URL to the attachment.
   */
  getAttachmentURL(mdPath, attachPath) {
    return `${this.baseURL}/api/attach/get/${mdPath}${attachPath}`
  }

  /**
   * Gets the base URL for resolving relative links within a markdown file.
   * @param {string} mdPath - The path of the associated markdown file.
   * @returns {string} The base URL to the attachment directory.
   */
  getAttachmentBase(mdPath) {
    return `${this.baseURL}/api/attach/get/${mdPath}/../`
  }

  /**
   * Deletes an attachment.
   * @param {string} mdPath - The path of the associated markdown file.
   * @param {string} attachPath - The path of the attachment to delete, relative to the markdown file.
   * @returns {Promise<object>} A confirmation object.
   * @example { "status": "success" }
   */
  deleteAttachment(mdPath, attachPath) {
    return this._request('/api/attach/delete', {
      method: 'POST',
      body: JSON.stringify({ mdPath, attachPath }),
    })
  }

  // --- Versioning ---

  /**
   * Fetches the version history for a specific file.
   * @param {string} path - The path of the file.
   * @returns {Promise<Array<object>>} A list of version records.
   * @example
   * [
   *   {
   *     "id": 2,
   *     "old_sha1": "a1b2c3...",
   *     "new_sha1": "d4e5f6...",
   *     "patch": "@@ -1,4 +1,4 @@\n # My Note\n \n-This is the content of the note.\n+This is the MODIFIED content of the note.",
   *     "type": "patch",
   *     "comment": "Updated the content.",
   *     "timestamp": "2023-10-27T14:00:00Z"
   *   },
   *   {
   *     "id": 1,
   *     "old_sha1": "",
   *     "new_sha1": "a1b2c3...",
   *     "patch": "# My Note\n\nThis is the content of the note.",
   *     "type": "full",
   *     "comment": "Initial creation.",
   *     "timestamp": "2023-10-27T13:00:00Z"
   *   }
   * ]
   */
  getFileHistory(path) {
    const params = new URLSearchParams({ path })
    return this._request(`/api/history?${params.toString()}`)
  }

  /**
   * Fetches the content of a specific version of a file.
   * @param {string} path - The path of the file.
   * @param {number} versionId - The ID of the version to retrieve.
   * @returns {Promise<object>} An object containing the historical content.
   * @example { "content": "# My Note\n\nThis is the old content." }
   */
  getFileVersion(path, versionId) {
    const params = new URLSearchParams({ path, id: versionId })
    return this._request(`/api/version?${params.toString()}`)
  }

  // --- Search ---

  /**
   * Searches for files based on a query.
   * @param {string} query - The search query.
   * @param {boolean} [useRegex=false] - Whether to treat the query as a regular expression.
   * @returns {Promise<Array<object>>} A list of search results.
   * @example
   * [
   *   {
   *     "path": "notes/MyNote.md",
   *     "context": [
   *       "2: This is the content of the note.",
   *       "4: Another line with note."
   *     ]
   *   }
   * ]
   */
  searchFiles(query, useRegex = false) {
    const params = new URLSearchParams({ q: query })
    if (useRegex) {
      params.append('regex', 'true')
    }
    return this._request(`/api/search?${params.toString()}`)
  }
}

export default ApiClient