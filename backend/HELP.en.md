# GoNote - Your Private Markdown Knowledge Palace

**Secure your thoughts in a digital kingdom that belongs only to you.**

GoNote is a desktop-grade Markdown management and editing service designed for modern creators, developers, and knowledge workers. It provides a personal knowledge base with features comparable to cloud-based note-taking apps, but with data that is **100% private**, all in the form of a **single, green executable**.

Say goodbye to tedious installations and privacy concerns about the cloud. With GoNote, a single click launches a Markdown creation and management center that is entirely under your control, accessible across all your devices.

## Table of Contents

1.  [Core Features: Unleash Your Creativity](#1-core-features-unleash-your-creativity)
2.  [First Run and Configuration](#2-first-run-and-configuration)
    *   [Download GoNote](#download-gonote)
    *   [Launch GoNote](#launch-gonote)
    *   [Your First Account](#your-first-account)
3.  [User Guide](#3-user-guide)
    *   [Login and Main Interface](#login-and-main-interface)
    *   [File and Directory Management](#file-and-directory-management)
    *   [Attachment Management](#attachment-management)
    *   [Mobile Experience](#mobile-experience)
4.  [Deep Dive: The Creative Heart of GoNote](#4-deep-dive-the-creative-heart-of-gonote)
    *   [Three Modes, Switch at Will](#three-modes-switch-at-will)
    *   [More Than Text, It's Visual Expression](#more-than-text-its-visual-expression)
    *   [Smart Tools for Increased Efficiency](#smart-tools-for-increased-efficiency)
5.  [Advanced Configuration (Optional)](#5-advanced-configuration-optional)
6.  [Frequently Asked Questions (FAQ)](#6-frequently-asked-questions-faq)

---

## 1. Core Features: Unleash Your Creativity

GoNote's design philosophy is to merge ultimate convenience with uncompromising functionality.

*   **True Digital Sovereignty**
    All your notes, attachments, and configurations reside quietly on your own hard drive. No data uploads, no privacy leaks—your intellectual wealth is entirely yours.

*   **Seamless File Management**
    GoNote directly maps to your local file system. Every creation, rename, or move you make within the app is perfectly synchronized with operations in your computer's file explorer. What you see is what you get, not just in editing, but in management too.

*   **Worry-Free Data Safe**
    A built-in automatic backup mechanism can package your entire knowledge base on a schedule and automatically clean up old backups. Even if the unexpected happens, your hard work remains safe.

*   **Desktop-Class Powerful Editor**
    With the powerful Vditor editor core built-in, GoNote elevates your creative experience. From basic text formatting to complex diagram drawing, everything is at your fingertips.

*   **One-Click Start, Cross-Platform Roaming**
    Whether you are a Windows, macOS, or Linux user, a single file is all you need to start the full service. Its responsive web interface seamlessly adapts to your desktop, tablet, and phone, allowing you to capture inspiration anytime, anywhere.

*   **Open and Customizable**
    GoNote respects your individual needs. You can directly modify `config.json` to adjust service behavior or even replace the entire `www` front-end directory to create a unique, personalized skin.

## 2. First Run and Configuration

### Download GoNote

Download the latest version for your operating system from the project's release page. For example:
*   Windows: `gonote-windows-amd64.exe`
*   macOS (Apple Silicon): `gonote-darwin-arm64`
*   Linux: `gonote-linux-amd64`

### Launch GoNote

1.  Extract the downloaded file into a dedicated folder where you want to store your notes, for example, `D:\GoNoteData`.
2.  **Windows**: Simply double-click the `.exe` file.
3.  **macOS/Linux**: Open a terminal, navigate to the file's directory, and run it.
    ```bash
    # Example (macOS/Linux)
    cd /path/to/your/gonote
    chmod +x ./gonote-linux-amd64  # Grant execution permission on the first run
    ./gonote-linux-amd64
    ```

When you see `Starting server on localhost:8080` in the command line window, GoNote has started successfully.

Note: If you need to access it from devices other than your local machine (like your phone), you need to modify the `bind` parameter in `config.json`. This is explained in detail below.

### Your First Account

On the first run, GoNote will automatically create the necessary folders and an initial user account for you.

*   **Be sure to check the output in the command line window!**
*   The program will create an account named `user` and generate a **6-digit random password**.
    ```
    ========================= IMPORTANT =========================
    Creating a default user with the following credentials:
      Username: user
      Password: 654321  <-- (Your password will be randomly generated)
    =============================================================
    ```
*   **Write down this password immediately**, as it is your credential for the first login.

## 3. User Guide

### Login and Main Interface

1.  Open your web browser (Chrome, Firefox, Edge are recommended).
2.  Enter in the address bar: [http://localhost:8080](http://localhost:8080)
3.  In the authentication dialog that pops up, enter the username `user` and the random password you noted down.
4.  After logging in, you will see the GoNote main interface, which is typically divided into two columns:
    *   **Left Column**: A tree list of your files and folders.
    *   **Right Column**: The powerful Markdown editor area.

### File and Directory Management

You can perform all management operations in the file list in the left column:
*   **New**: Click the corresponding "+" button to create a new file or folder.
*   **Rename/Delete**: Right-click on a file or folder to bring up the context menu.
*   **Move (Drag and Drop)**: Simply drag and drop a file or folder to the target location.
*   **Recycle Bin**: Deleted files are moved to the recycle bin, where you can choose to restore or permanently delete them.

### Attachment Management

1.  While editing a note, drag and drop files like images directly into the editor.
2.  GoNote will automatically upload the attachment and insert the corresponding Markdown reference link at the cursor's position, such as `![image.png](image.png)`.
3.  All attachments are stored in the `.attach` directory associated with the note file, making management clear.
4.  Click on an image within a note to view a larger preview.

### Mobile Experience

When accessing GoNote on a mobile phone or tablet, the interface automatically switches to an optimized three-column sliding layout, allowing you to easily switch between the folder tree, file list, and editor to capture inspiration on the go.

## 4. Deep Dive: The Creative Heart of GoNote

GoNote's editing experience goes far beyond just typing text. It is a powerful creative toolbox designed to handle a variety of complex document needs.

### Three Modes, Switch at Will

You can freely choose the most suitable editing mode based on your current creative flow. Switch easily with the shortcut `⌘-⇧-M` (macOS) or `Ctrl-⇧-M` (Windows).

*   **WYSIWYG (What You See Is What You Get)**: For users who prefer rich-text editors, this mode offers the most intuitive experience. The format you type is rendered instantly, without being distracted by Markdown syntax. (`⌘-⌥-7` / `Ctrl-Alt-7`)
*   **Instant Rendering**: Perfectly combines the control of source code editing with the convenience of a real-time preview. When you finish typing a line of Markdown, it is immediately rendered in its final style, striking a perfect balance between efficiency and control. (`⌘-⌥-8` / `Ctrl-Alt-8`)
*   **Split View**: The classic source/preview dual-column layout, designed for professional users and developers who need precise control over Markdown syntax. (`⌘-⌥-9` / `Ctrl-Alt-9`)

### More Than Text, It's Visual Expression

Forget the hassle of switching between multiple applications. GoNote natively supports a variety of embedded diagrams and advanced formats:

*   **Diagramming**: Easily create **flowcharts, Gantt charts, sequence diagrams**, and **Echarts statistical charts** to make your data and logic clear at a glance.
*   **Task Management**: Use `[ ]` and `[x]` to create interactive **task lists** to efficiently track your project progress.
*   **Academic and Technical Writing**: Built-in **Katex math formula** support ensures that complex mathematical equations and chemical formulas are rendered elegantly.
*   **Music Composition**: Even supports the drawing of **musical staves**, providing a unique recording space for music creators.

### Smart Tools for Increased Efficiency

*   **One-Click Formatting**: Messy Markdown source code? Click the "Format" button to make it neat and orderly instantly.
*   **Content Conversion**: Content copied from web pages or Word is automatically converted to clean, standard Markdown when pasted into GoNote.
*   **Export Master**: With a single click, export your work into beautifully formatted **PDFs**, high-definition **PNG/JPEG** images, or even presentation-ready **PPT** files.
*   **WeChat Official Account Optimization**: A special "Copy to WeChat" feature solves the problem of messy formatting in the WeChat editor.

GoNote aims to be the final piece of your creative workflow puzzle—a powerful partner that lets you immerse yourself in thought and expression without being distracted by the tool itself.

## 5. Advanced Configuration (Optional)

You can edit the `config.json` file in the program directory for advanced settings (the program must be restarted after modification):
*   **Change Port**: Change the `bind` value from `"localhost:8080"` to `"0.0.0.0:8080"` to allow other devices on your local network (like phones or tablets) to access it.
*   **Enable Automatic Backups**: GoNote offers a powerful automatic backup feature. Find the `backup` section in `config.json` to configure it:
    *   `"enabled": true`: Set to `true` to enable this feature.
    *   `"dir": "backup"`: Specify the directory to store backup archives.
    *   `"cron": "0 0 1 * *"`: Set the [CRON expression](https://crontab.guru/) for the backup schedule. The default value means it runs at midnight on the 1st of every month.
    *   `"retention_days": 180`: The number of days to retain backup files. Old backups exceeding this duration will be automatically deleted.
*   **Manage Users**: Directly edit the `users.txt` file to add or modify users and passwords. After adding a new user, remember to manually create a folder with the same name for them under the `markdown` directory.

## 6. Frequently Asked Questions (FAQ)

**Q: What if I forget my password?**
**A:** Open the `users.txt` file in the program directory. You can see the password in plain text and can also modify it directly.

**Q: How do I back up my data?**
**A:** GoNote offers two backup methods to meet different needs:
1.  **Automatic Backup (Recommended)**: This is the best way to ensure daily data safety. You can enable and configure the automatic backup feature in the `config.json` file (see the "Advanced Configuration" section for details). The program will periodically package all your notes and store them in the specified backup directory.
2.  **Manual Full Backup**: For **system migration** or **one-time full archiving**, the original backup method is still effective: simply copy the entire GoNote program folder. The core data includes the `markdown/`, `backup/`, `users.txt`, and `config.json` files and folders.

**Q: Can I customize the interface?**
**A:** Of course. The `www` directory contains all the front-end files. If you are a developer, you can freely modify it to create a personalized GoNote.

**Q: Can I sync my notes across multiple computers?**
**A:** GoNote itself does not provide a cloud sync feature, to ensure absolute data privacy. However, you can achieve secure and reliable multi-device synchronization by placing the entire GoNote data folder in the sync directory of a third-party sync tool (such as Syncthing, Dropbox, OneDrive, etc.).

---

Thank you for choosing GoNote. Happy creating!
