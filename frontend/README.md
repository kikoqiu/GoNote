[English Version](./README.en.md)

# GoNote - 简洁的 Markdown 管理与编辑器

GoNote 是一个基于 Go 后端开发的单文件 Markdown 管理和编辑解决方案的前端部分。它旨在提供一个轻量、高效的 Markdown 笔记和文档管理体验。

本项目基于 [Arya - 在线 Markdown 编辑器](https://github.com/nicejade/markdown-online-editor) 进行修改和扩展。

## 功能支持

### Markdown 管理
- [x] **直接管理 Markdown 文件目录**：支持对服务端文件系统的 Markdown 文件和文件夹进行创建、重命名、删除等操作。
- [x] **文件列表与预览**：显示当前文件夹下的文件列表，并支持点击预览 Markdown 文件内容。
- [x] **附件上传与管理**：支持在 Markdown 编辑器中上传图片等附件，并自动生成 Markdown 引用链接。
- [x] **自动保存**：编辑内容实时自动保存，防止数据丢失。
- [x] **回收站功能**：支持将文件和文件夹移动到回收站，提供恢复或永久删除选项。
- [x] **文件下载**：支持下载非 Markdown 类型的文件。

### UI 布局与交互
- [x] **响应式布局**：针对桌面和移动设备提供优化布局。
    - **桌面布局**：采用两栏设计，左侧为文件夹树和文件列表，右侧为 Markdown 编辑器。
    - **移动布局**：采用滑动式三栏设计，方便在文件夹树、文件列表和编辑器之间切换。
- [x] **文件/文件夹操作对话框**：统一的对话框用于新建、重命名文件或文件夹。
- [x] **图片预览**：内置图片预览功能，点击图片可放大查看。
- [x] **国际化支持**：界面文字支持多语言切换。

### Markdown 编辑与导出
- [x] **Vditor 集成**：深度集成 Vditor 编辑器，提供强大的 Markdown 编辑和预览能力。
- [x] **多种编辑模式**：支持所见即所得、即时渲染、分屏渲染等多种编辑模式。
- [x] **富文本特性**：支持流程图、甘特图、时序图、任务列表、Echarts 图表、五线谱等。
- [x] **内容转换**：支持粘贴 HTML 自动转换为 Markdown。
- [x] **快捷键支持**：提供常用编辑快捷键。
- [x] **内容导出**：支持将 Markdown 内容导出为 PDF、PNG、JPEG、PPT 等格式。
- [x] **语法检查与格式化**：支持检查并格式化 Markdown 语法。
- [x] **周边功能**：新增复制到微信公众号等实用功能。
- [x] **本地文件导入**：支持导入本地 Markdown 文件进行编辑。

## 突出特点
- **直接管理 Markdown 文件目录**：与后端文件系统无缝集成，直接对文件和目录进行操作。
- **单文件应用**：轻量级设计，易于部署和使用。
- **跨平台**：基于 Web 技术，可在任何支持现代浏览器的设备上运行。
- **移动端友好**：提供专门优化的移动端 UI/UX，方便在手机或平板上使用。

## 未来计划 (Todo)
- [ ] **版本控制集成**：计划支持版本控制系统，实现 Markdown 文件的版本管理。
- [ ] **搜索功能**：支持搜索功能。
- [ ] **自动备份**：实现 Markdown 文件的自动备份功能，确保数据安全。
- [ ] **AI 辅助管理**：探索集成 AI 能力，提供智能内容推荐、摘要生成、语法纠正等功能。

## 如何使用

即可使用。

默认为[所见即所得](https://b3log.org/vditor/)模式，可通过 `⌘-⇧-M`（`Ctrl-⇧-M`）进行切换；或通过以下方式：

- 所见即所得：`⌘-⌥-7`（`Ctrl-alt-7`）；
- 即时渲染：`⌘-⌥-8`（`Ctrl-alt-8`）；
- 分屏渲染：`⌘-⌥-9`（`Ctrl-alt-9`）；

### PPT 预览

如果您用作 `PPT` 预览（入口在`设置`中），需要注意，这里暂还不能支持各种图表的渲染；您可以使用 `---` 来定义水平方向上幻灯片，用 `--` 来定义垂直幻灯片；更多设定可以参见 [Reveal.js Markdown 文档](https://revealjs.com/markdown/)。


## 如何开发

### 先决条件

说明用户在安装和使用前，需要准备的一些先决条件，譬如：您需要安装或升级  [Node.js](https://nodejs.org/en/)（>= `16.*`，< `18.*`），推荐使用  [Pnpm](https://pnpm.io/)  或  [Yarn](https://www.jeffjade.com/2017/12/30/135-npm-vs-yarn-detial-memo/)  作为首选包管理工具。

```bash
cd markdown-online-editor

# 安装依赖
yarn

# 开始开发
yarn start

# 部署 Github Pages(需修改 commands/deploy.sh)
yarn deploy
```

## 特别鸣谢

本项目基于 [Arya - 在线 Markdown 编辑器](https://github.com/nicejade/markdown-online-editor) 进行修改和扩展。
特别感谢原项目作者 [nicejade](https://github.com/nicejade) 及其团队的杰出工作，为本项目提供了坚实的基础。

本项目深度集成了以下优秀开源项目，在此表示衷心感谢：

- **[Vditor](https://github.com/b3log/vditor)**: 一款卓越的浏览器端 Markdown 编辑器，为本项目提供了强大的编辑和预览能力。
- **[Lute](https://github.com/88250/lute)**: Vditor 的 Go 语言后端 Markdown 处理器，为本项目提供了高效稳定的 Markdown 解析和渲染服务。

[Arya](https://markdown.lovejade.cn/?utm_source=github.com) 的产生，还得益于 [Vue、Reveal.js 等开源库](https://github.com/nicejade/markdown-online-editor/blob/master/package.json#L25-L64)的支持，在此一并表示感谢。

## License

[MIT](http://opensource.org/licenses/MIT)