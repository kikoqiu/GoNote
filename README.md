[English Version](./README.en.md)

# GoNote - 您的私有化 Markdown 知识宫殿

GoNote 是一款专为现代创作者、开发者和知识工作者打造的桌面级 Markdown 管理与编辑服务。它以一个**绿色单文件**的形式，为您提供了一个功能媲美云端笔记应用、但数据却 **100% 私有**的个人知识库。

告别繁琐的安装和对云端隐私的担忧。有了 GoNote，您只需一次点击，即可开启一个完全由您掌控的、跨越所有设备的 Markdown 创作与管理中心。


## ✨ 核心特性

*   **🚀 一键启动，绿色便携**：无论 Windows, macOS 还是 Linux，单个可执行文件即可启动完整服务，无需任何外部依赖。
*   **🔐 数据绝对私有**：所有笔记、附件和配置都存储在您自己的设备上，没有数据上传，没有隐私泄露。
*   **🗂️ 原生文件系统管理**：应用内的文件操作（创建、重命名、移动）与您电脑的资源管理器完全同步。
*   **强大的编辑器**：深度集成 Vditor，支持**所见即所得**、**即时渲染**和**分屏预览**三种模式，并原生支持**流程图、甘特图、数学公式**等高级功能。
*   **📱 跨平台响应式体验**：无论桌面、平板还是手机，都能获得优化的访问体验。
*   **🛡️ 自动备份**：内置可配置的自动备份机制，按计划将您的整个知识库打包压缩，并自动清理旧备份。
*   **🎨 开放与可定制**：支持通过修改配置文件进行高级设置，甚至可以替换整个前端，打造您的专属皮肤。

## 🏁 开始使用

### 1. 下载

从项目的 **Releases** 页面下载适合您操作系统的最新版本。

### 2. 启动

1.  将下载的文件解压到一个您希望存放笔记的专属文件夹（例如 `D:\MyNotes`）。
2.  **Windows**: 直接双击 `.exe` 文件。
3.  **macOS/Linux**: 打开终端，赋予文件执行权限并运行它。
    ```bash
    chmod +x ./gonote-linux-amd64
    ./gonote-linux-amd64
    ```
4.  首次运行时，程序会在命令行窗口输出默认用户名 `user` 和一个**随机6位数密码**。请务必记下此密码。

### 3. 访问

打开浏览器，访问 [http://localhost:8080](http://localhost:8080)，使用上一步的凭据登录即可开始使用。

## 🛠️ 开发

本项目分为 Go 后端和 Vue.js 前端。

### Backend (Go)

后端是项目的核心，负责文件处理、HTTP 服务和用户认证。

1.  进入 `backend` 目录。
2.  构建：
    ```bash
    go build
    ```
3.  运行：
    ```bash
    ./backend.exe # Windows
    ./backend     # macOS/Linux
    ```

### Frontend (Vue.js)

前端提供了所有的用户交互界面。

1.  确保您已安装 [Node.js](https://nodejs.org/) (>= 16.*) 和 [Yarn](https://yarnpkg.com/)。
2.  进入 `frontend` 目录。
3.  安装依赖：
    ```bash
    yarn install
    ```
4.  启动开发服务器：
    ```bash
    yarn start
    ```
5.  构建生产版本（文件将输出到 `dist` 目录）：
    ```bash
    yarn build
    ```

## 🔍 详细文档

关于前后端更详细的说明、高级配置、完整功能列表和常见问题解答，请参阅各自目录中的原始文档：

*   **后端详细说明**: [backend/HELP.md](./backend/HELP.md)
*   **前端详细说明**: [frontend/README.md](./frontend/README.md)

## 🙏 特别鸣谢

*   本项目前端基于 [Arya - 在线 Markdown 编辑器](https://github.com/nicejade/markdown-online-editor) 修改和扩展。
*   深度集成了卓越的浏览器端 Markdown 编辑器 **[Vditor](https://github.com/b3log/vditor)** 及其 Go 后端 **[Lute](https://github.com/88250/lute)**。

## 📄 许可证

本项目采用 MIT 许可证。

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
