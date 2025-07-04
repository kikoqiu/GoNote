# GoNote - A Concise Markdown Manager and Editor

GoNote is the frontend part of a single-file Markdown management and editing solution developed with a Go backend. It aims to provide a lightweight and efficient experience for managing Markdown notes and documents.

This project is modified and extended based on [Arya - Online Markdown Editor](https://github.com/nicejade/markdown-online-editor).

## Feature Support

### Markdown Management
- [x] **Directly manage Markdown file directories**: Supports creating, renaming, deleting Markdown files and folders on the server's file system.
- [x] **File list and preview**: Displays a list of files in the current folder and supports clicking to preview Markdown file content.
- [x] **Attachment upload and management**: Supports uploading attachments like images in the Markdown editor and automatically generating Markdown reference links.
- [x] **Auto-save**: Content is saved automatically in real-time to prevent data loss.
- [x] **Recycle bin function**: Supports moving files and folders to a recycle bin, with options to restore or permanently delete.
- [x] **File download**: Supports downloading non-Markdown file types.

### UI Layout and Interaction
- [x] **Responsive layout**: Provides optimized layouts for desktop and mobile devices.
    - **Desktop layout**: A two-column design with the folder tree and file list on the left, and the Markdown editor on the right.
    - **Mobile layout**: A sliding three-column design for easy switching between the folder tree, file list, and editor.
- [x] **File/folder operation dialog**: A unified dialog for creating and renaming files or folders.
- [x] **Image preview**: Built-in image preview function; click on an image to view it larger.
- [x] **Internationalization support**: Interface text supports multiple languages.

### Markdown Editing and Exporting
- [x] **Vditor integration**: Deeply integrated with the Vditor editor for powerful Markdown editing and preview capabilities.
- [x] **Multiple editing modes**: Supports WYSIWYG, instant rendering, and split-screen rendering modes.
- [x] **Rich text features**: Supports flowcharts, Gantt charts, sequence diagrams, task lists, Echarts, and musical staves.
- [x] **Content conversion**: Supports pasting HTML and automatically converting it to Markdown.
- [x] **Shortcut key support**: Provides common editing shortcuts.
- [x] **Content export**: Supports exporting Markdown content to formats like PDF, PNG, JPEG, and PPT.
- [x] **Syntax checking and formatting**: Supports checking and formatting Markdown syntax.
- [x] **Additional features**: Added practical functions like copying to WeChat Official Accounts.
- [x] **Local file import**: Supports importing local Markdown files for editing.

## Key Features
- **Direct management of Markdown file directories**: Seamlessly integrates with the backend file system to directly manipulate files and directories.
- **Single-file application**: Lightweight design, easy to deploy and use.
- **Cross-platform**: Based on web technology, it can run on any device that supports modern browsers.
- **Mobile-friendly**: Provides a specially optimized mobile UI/UX for convenient use on phones or tablets.

## Future Plans (Todo)
- [ ] **Version control integration**: Plan to support version control systems for Markdown file version management.
- [ ] **Search function**: Support for a search function.
- [ ] **Automatic backup**: Implement automatic backup for Markdown files to ensure data security.
- [ ] **AI-assisted management**: Explore integrating AI capabilities to provide smart content recommendations, summary generation, grammar correction, etc.

## How to Use

Ready to use out of the box.

The default is [WYSIWYG](https://b3log.org/vditor/) mode, which can be toggled with `⌘-⇧-M` (`Ctrl-⇧-M`), or through the following:

- WYSIWYG: `⌘-⌥-7` (`Ctrl-alt-7`);
- Instant Rendering: `⌘-⌥-8` (`Ctrl-alt-8`);
- Split-screen Rendering: `⌘-⌥-9` (`Ctrl-alt-9`);

### PPT Preview

If you use it for `PPT` preview (entry point is in `Settings`), please note that it does not yet support rendering various charts. You can use `---` to define horizontal slides and `--` to define vertical slides. For more settings, please refer to the [Reveal.js Markdown documentation](https://revealjs.com/markdown/).

## How to Develop

### Prerequisites

Before installation and use, you need to have some prerequisites, for example: You need to install or upgrade [Node.js](https://nodejs.org/en/) (>= `16.*`, < `18.*`). It is recommended to use [Pnpm](https://pnpm.io/) or [Yarn](https://www.jeffjade.com/2017/12/30/135-npm-vs-yarn-detial-memo/) as the preferred package management tool.

```bash
cd markdown-online-editor

# Install dependencies
yarn

# Start development
yarn start

# Deploy to Github Pages (requires modifying commands/deploy.sh)
yarn deploy
```

## Special Thanks

This project is modified and extended based on [Arya - Online Markdown Editor](https://github.com/nicejade/markdown-online-editor).
Special thanks to the original author [nicejade](https://github.com/nicejade) and his team for their outstanding work, which provided a solid foundation for this project.

This project deeply integrates the following excellent open-source projects, and we express our sincere gratitude:

- **[Vditor](https://github.com/b3log/vditor)**: An excellent browser-side Markdown editor that provides powerful editing and preview capabilities for this project.
- **[Lute](https://github.com/88250/lute)**: The Go language backend Markdown processor for Vditor, which provides efficient and stable Markdown parsing and rendering services for this project.

The creation of [Arya](https://markdown.lovejade.cn/?utm_source=github.com) was also made possible by the support of open-source libraries like [Vue, Reveal.js, etc.](https://github.com/nicejade/markdown-online-editor/blob/master/package.json#L25-L64), and we thank them as well.

## License

[MIT](http://opensource.org/licenses/MIT)
