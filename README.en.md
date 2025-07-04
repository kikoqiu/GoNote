# GoNote - Your Private Markdown Knowledge Palace

GoNote is a desktop-grade Markdown management and editing service designed for modern creators, developers, and knowledge workers. It provides a feature-rich personal knowledge base, comparable to cloud-based note-taking apps, but with data that is **100% private**, all in the form of a **single, green executable**.

Say goodbye to tedious installations and privacy concerns about the cloud. With GoNote, a single click launches a Markdown creation and management center that is entirely under your control, accessible across all your devices.


## ✨ Core Features

*   **🚀 One-Click Start, Green & Portable**: Whether on Windows, macOS, or Linux, a single executable file starts the complete service with no external dependencies.
*   **🔐 Absolute Data Privacy**: All your notes, attachments, and configurations are stored on your own device. No data uploads, no privacy leaks.
*   **🗂️ Native File System Management**: File operations within the app (create, rename, move) are perfectly synchronized with your computer's file explorer.
*   ** powerful Editor**: Deeply integrated with Vditor, supporting **WYSIWYG**, **Instant Rendering**, and **Split-Screen Preview** modes, with native support for advanced features like **flowcharts, Gantt charts, and mathematical formulas**.
*   **📱 Cross-Platform Responsive Experience**: Get an optimized experience whether you are on a desktop, tablet, or mobile phone.
*   **🛡️ Automatic Backups**: A built-in, configurable automatic backup mechanism periodically packages your entire knowledge base and automatically cleans up old backups.
*   **🎨 Open and Customizable**: Supports advanced settings through configuration file modification and even allows replacing the entire front-end to create your own custom skin.

## 🏁 Getting Started

### 1. Download

Download the latest version for your operating system from the project's **Releases** page.

### 2. Launch

1.  Extract the downloaded file to a dedicated folder where you want to store your notes (e.g., `D:\MyNotes`).
2.  **Windows**: Simply double-click the `.exe` file.
3.  **macOS/Linux**: Open a terminal, grant the file execution permissions, and run it.
    ```bash
    chmod +x ./gonote-linux-amd64
    ./gonote-linux-amd64
    ```
4.  On the first run, the program will output the default username `user` and a **random 6-digit password** in the command line window. Be sure to save this password.

### 3. Access

Open your browser, navigate to [http://localhost:8080](http://localhost:8080), and log in with the credentials from the previous step to start using GoNote.

## 🛠️ Development

This project is divided into a Go backend and a Vue.js frontend.

### Backend (Go)

The backend is the core of the project, responsible for file handling, the HTTP service, and user authentication.

1.  Navigate to the `backend` directory.
2.  Build:
    ```bash
    build.bat
    ```


### Frontend (Vue.js)

The frontend provides the entire user interface.

1.  Ensure you have [Node.js](https://nodejs.org/) (>= 16.*) and [Yarn](https://yarnpkg.com/) installed.
2.  Navigate to the `frontend` directory.
3.  Install dependencies:
    ```bash
    yarn install
    ```
4.  Start the development server:
    ```bash
    yarn start
    ```
5.  Build for production (files will be output to the `dist` directory):
    ```bash
    yarn build
    ```

## 🔍 Detailed Documentation

For more detailed information about the frontend and backend, advanced configurations, full feature lists, and FAQs, please refer to the original documents in their respective directories:

*   **Backend Details**: [backend/HELP.md](./backend/HELP.md)
*   **Frontend Details**: [frontend/README.md](./frontend/README.md)

## 🙏 Acknowledgements

*   The frontend of this project is modified and extended based on [Arya - Online Markdown Editor](https://github.com/nicejade/markdown-online-editor).
*   It deeply integrates the excellent browser-side Markdown editor **[Vditor](https://github.com/b3log/vditor)** and its Go backend **[Lute](https://github.com/88250/lute)**.

## 📄 License

This project is licensed under the MIT License.

Copyright (c) 2024 qiu.kiko@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
