// main.go
package main

import (
	"archive/zip"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"embed"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math/big"
	rnd "math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/robfig/cron/v3" // 新增的依赖
	"github.com/sergi/go-diff/diffmatchpatch"
	"go.etcd.io/bbolt"
)

//go:embed all:embed
var embeddedFS embed.FS

//go:embed Help.md
var helpMarkdown []byte

//go:embed Help.en.md
var helpMarkdownEn []byte

// --- config.go ---

// 新增: 用于备份配置的结构体
type BackupConfig struct {
	Enabled       bool   `json:"enabled"`
	Dir           string `json:"dir"`
	Cron          string `json:"cron"`
	RetentionDays int    `json:"retention_days"`
}

type Config struct {
	Bind        string       `json:"bind"`
	TLS         bool         `json:"tls"`
	CertFile    string       `json:"cert_file"`
	KeyFile     string       `json:"key_file"`
	VisitLog    string       `json:"visit_log"`
	MarkdownDir string       `json:"markdown_dir"`
	WWWDir      string       `json:"www_dir"`
	UsersFile   string       `json:"users_file"`
	Backup      BackupConfig `json:"backup"` // 新增
}

var defaultConfig = Config{
	Bind:        "localhost:8080",
	TLS:         false,
	CertFile:    "cert.pem",
	KeyFile:     "key.pem",
	VisitLog:    "visit.log",
	MarkdownDir: "markdown",
	WWWDir:      "www",
	UsersFile:   "users.txt",
	// 新增: 备份的默认配置
	Backup: BackupConfig{
		Enabled:       false,
		Dir:           "backup",
		Cron:          "0 0 1 * *", // 每月1日午夜
		RetentionDays: 180,
	},
}

var AppConfig Config

func LoadConfig() {
	AppConfig = defaultConfig
	configFile := "config.json"
	if _, err := os.Stat(configFile); err == nil {
		file, err := os.ReadFile(configFile)
		if err == nil {
			// 先解码到临时变量，以防某些字段缺失
			tempConfig := defaultConfig
			if err := json.Unmarshal(file, &tempConfig); err == nil {
				AppConfig = tempConfig
			}
		}
	}

	// 无论如何，都重新序列化并写入，以确保新字段（如backup）存在于文件中
	data, err := json.MarshalIndent(AppConfig, "", "  ")
	if err == nil {
		os.WriteFile(configFile, data, 0644)
	}

	flag.StringVar(&AppConfig.Bind, "bind", AppConfig.Bind, "Address and port to bind the server")
	flag.BoolVar(&AppConfig.TLS, "tls", AppConfig.TLS, "Enable TLS (HTTPS)")
	flag.StringVar(&AppConfig.CertFile, "cert", AppConfig.CertFile, "Path to TLS certificate file")
	flag.StringVar(&AppConfig.KeyFile, "key", AppConfig.KeyFile, "Path to TLS key file")
	flag.StringVar(&AppConfig.VisitLog, "visitlog", AppConfig.VisitLog, "Path for visit log file (empty to disable)")
	flag.StringVar(&AppConfig.MarkdownDir, "markdown", AppConfig.MarkdownDir, "Path to markdown storage directory")
	flag.StringVar(&AppConfig.WWWDir, "www", AppConfig.WWWDir, "Path to static web assets")
	flag.StringVar(&AppConfig.UsersFile, "users", AppConfig.UsersFile, "Path to users file for basic auth")
	flag.Parse()
}

// --- auth.go ---

var userCredentials = make(map[string]string)
var userMutex = &sync.RWMutex{}

func LoadUsers() {
	userMutex.Lock()
	defer userMutex.Unlock()

	userCredentials = make(map[string]string)

	if _, err := os.Stat(AppConfig.UsersFile); os.IsNotExist(err) {
		defaultUser := "user"
		passwordNum := 100000 + rnd.Intn(900000)
		defaultPassword := fmt.Sprintf("%d", passwordNum)

		log.Printf("========================= IMPORTANT =========================")
		log.Printf("Users file '%s' not found.", AppConfig.UsersFile)
		log.Printf("Creating a default user with the following credentials:")
		log.Printf("  Username: %s", defaultUser)
		log.Printf("  Password: %s", defaultPassword)
		log.Printf("=============================================================")

		fileContent := fmt.Sprintf("%s %s\n", defaultUser, defaultPassword)
		if err := os.WriteFile(AppConfig.UsersFile, []byte(fileContent), 0644); err != nil {
			log.Fatalf("FATAL: Failed to create default users file: %v", err)
		}

		userMarkdownPath := filepath.Join(AppConfig.MarkdownDir, defaultUser)
		log.Printf("Creating markdown directory for new user at: %s", userMarkdownPath)
		if err := os.MkdirAll(userMarkdownPath, 0755); err != nil {
			log.Printf("WARNING: Failed to create markdown directory for user '%s': %v", defaultUser, err)
		}

		docDirPath := filepath.Join(userMarkdownPath, "Doc")
		log.Printf("Creating 'Doc' directory for user at: %s", docDirPath)
		if err := os.MkdirAll(docDirPath, 0755); err != nil {
			log.Printf("WARNING: Failed to create 'Doc' directory for user '%s': %v", defaultUser, err)
		}

		userHelpPath := filepath.Join(docDirPath, "Help.md")
		log.Printf("Creating user guide at: %s", userHelpPath)
		if err := os.WriteFile(userHelpPath, helpMarkdown, 0644); err != nil {
			log.Printf("WARNING: Failed to create user guide in 'Doc' directory: %v", err)
		}

		userHelpEnPath := filepath.Join(docDirPath, "Help.en.md")
		log.Printf("Creating user guide at: %s", userHelpPath)
		if err := os.WriteFile(userHelpEnPath, helpMarkdownEn, 0644); err != nil {
			log.Printf("WARNING: Failed to create user guide in 'Doc' directory: %v", err)
		}
	}

	file, err := os.ReadFile(AppConfig.UsersFile)
	if err != nil {
		log.Fatalf("Failed to read users file: %v", err)
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			log.Printf("Warning: malformed line in users file: %s", line)
			continue
		}
		userCredentials[parts[0]] = parts[1]
	}
	log.Printf("Loaded %d user(s)", len(userCredentials))
}

type contextKey string

const userContextKey = contextKey("user")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(userCredentials) == 0 {
			ctx := context.WithValue(r.Context(), userContextKey, "anonymous")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		sendUnauthorized := func(message string) {
			if r.Header.Get("Api-Version") != "" {
				respondError(w, http.StatusUnauthorized, message)
			} else {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		}
		user, pass, ok := r.BasicAuth()
		if !ok {
			sendUnauthorized("Authentication credentials required")
			return
		}
		userMutex.RLock()
		expectedPass, userExists := userCredentials[user]
		userMutex.RUnlock()
		if !userExists || expectedPass != pass {
			sendUnauthorized("Invalid username or password")
			return
		}
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// --- backup.go ---

func StartBackupScheduler() {
	if !AppConfig.Backup.Enabled {
		log.Println("Automatic backup is disabled.")
		return
	}

	log.Printf("Starting backup scheduler. Cron: '%s', Retention: %d days.", AppConfig.Backup.Cron, AppConfig.Backup.RetentionDays)

	c := cron.New()

	// Add the main backup job
	_, err := c.AddFunc(AppConfig.Backup.Cron, performBackup)
	if err != nil {
		log.Fatalf("FATAL: Invalid backup cron expression: %v", err)
	}

	// Add a daily cleanup job (runs at 1 AM every day)
	_, err = c.AddFunc("0 1 * * *", performBackupCleanup)
	if err != nil {
		log.Fatalf("FATAL: Could not schedule backup cleanup job: %v", err)
	}

	go c.Start()
}

func performBackup() {
	log.Println("Starting scheduled backup...")

	backupDir := AppConfig.Backup.Dir
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		log.Printf("ERROR: Could not create backup directory '%s': %v", backupDir, err)
		return
	}

	timestamp := time.Now().Format("2006-01-02T15-04-05")
	zipFileName := fmt.Sprintf("markdown-%s.zip", timestamp)
	zipFilePath := filepath.Join(backupDir, zipFileName)

	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		log.Printf("ERROR: Could not create zip file '%s': %v", zipFilePath, err)
		return
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	sourceDir := AppConfig.MarkdownDir

	err = filepath.Walk(sourceDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil // Skip directories themselves
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// Use forward slashes for zip compatibility
		zipPath := filepath.ToSlash(relPath)
		writer, err := zipWriter.Create(zipPath)
		if err != nil {
			return err
		}

		fileToZip, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fileToZip.Close()

		_, err = io.Copy(writer, fileToZip)
		return err
	})

	if err != nil {
		log.Printf("ERROR: Failed during backup zipping process: %v", err)
		// Attempt to clean up partially created zip file
		os.Remove(zipFilePath)
		return
	}

	log.Printf("Successfully created backup: %s", zipFilePath)
}

func performBackupCleanup() {
	log.Println("Starting backup cleanup task...")

	backupDir := AppConfig.Backup.Dir
	retentionDays := AppConfig.Backup.RetentionDays
	if retentionDays <= 0 {
		log.Println("Backup retention is disabled (retention_days <= 0).")
		return
	}

	cutoffTime := time.Now().Add(-time.Duration(retentionDays) * 24 * time.Hour)
	log.Printf("Deleting backups older than %s", cutoffTime.Format("2006-01-02"))

	entries, err := os.ReadDir(backupDir)
	if err != nil {
		log.Printf("ERROR: Could not read backup directory '%s' for cleanup: %v", backupDir, err)
		return
	}

	deletedCount := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".zip") || !strings.HasPrefix(entry.Name(), "markdown-") {
			continue
		}

		// Extract timestamp from filename: markdown-YYYY-MM-DDTHH-MM-SS.zip
		timestampStr := strings.TrimSuffix(strings.TrimPrefix(entry.Name(), "markdown-"), ".zip")
		backupTime, err := time.Parse("2006-01-02T15-04-05", timestampStr)
		if err != nil {
			log.Printf("WARNING: Could not parse timestamp from backup file '%s', skipping.", entry.Name())
			continue
		}

		if backupTime.Before(cutoffTime) {
			filePath := filepath.Join(backupDir, entry.Name())
			log.Printf("Deleting old backup: %s", filePath)
			if err := os.Remove(filePath); err != nil {
				log.Printf("ERROR: Failed to delete old backup '%s': %v", filePath, err)
			} else {
				deletedCount++
			}
		}
	}
	log.Printf("Backup cleanup task finished. Deleted %d old backup(s).", deletedCount)
}

// --- store.go ---
// ... (rest of the file is unchanged, so I will omit it for brevity and just show the main function)
// ... all store.go, file_monitor.go, versioning.go, search.go, handlers.go, utils.go code remains the same ...

// --- store.go ---

type Document struct {
	Path    string
	SHA1    string
	Content string
}

type InMemoryStore struct {
	sync.RWMutex
	docs map[string]Document
}

var store = InMemoryStore{docs: make(map[string]Document)}

func isSpecialPath(path string) bool {
	return strings.Contains(path, ".extra") || strings.HasSuffix(path, ".attach")
}

func (s *InMemoryStore) Scan() {
	s.Lock()
	defer s.Unlock()
	log.Println("Scanning markdown directory for initial cache...")

	s.docs = make(map[string]Document)

	users, err := os.ReadDir(AppConfig.MarkdownDir)
	if err != nil {
		log.Printf("Could not read markdown dir: %v. It may be created later.", err)
		return
	}

	for _, userEntry := range users {
		if !userEntry.IsDir() {
			continue
		}
		user := userEntry.Name()
		userPath := filepath.Join(AppConfig.MarkdownDir, user)

		filepath.WalkDir(userPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if isSpecialPath(path) {
				return nil
			}
			if d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Error reading file for cache %s: %v", path, err)
				return nil
			}

			relPath, _ := filepath.Rel(AppConfig.MarkdownDir, path)

			doc := Document{
				Path:    relPath,
				Content: string(content),
				SHA1:    calculateSHA1(content),
			}
			s.docs[relPath] = doc

			return nil
		})
	}
	log.Printf("Initial cache populated with %d documents.", len(s.docs))
}

func (s *InMemoryStore) UpdateDoc(relPath string, content []byte) {
	s.Lock()
	defer s.Unlock()

	doc := Document{
		Path:    relPath,
		Content: string(content),
		SHA1:    calculateSHA1(content),
	}
	s.docs[relPath] = doc
	log.Printf("Cache updated for: %s", relPath)
}

func (s *InMemoryStore) DeleteDoc(relPath string) {
	s.Lock()
	defer s.Unlock()
	delete(s.docs, relPath)
	log.Printf("Cache deleted for: %s", relPath)
}

// --- file_monitor.go ---

func WatchMarkdownDir() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create file watcher: %v", err)
	}
	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				relPath, err := filepath.Rel(AppConfig.MarkdownDir, event.Name)
				if err != nil || isSpecialPath(relPath) {
					continue
				}

				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					info, err := os.Stat(event.Name)
					if err == nil && !info.IsDir() {
						log.Printf("Watcher detected change in %s, updating cache.", relPath)
						content, err := os.ReadFile(event.Name)
						if err == nil {
							store.UpdateDoc(relPath, content)
						}
					}
				}
				if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
					store.DeleteDoc(relPath)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %v", err)
			}
		}
	}()

	filepath.WalkDir(AppConfig.MarkdownDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && !isSpecialPath(path) {
			watcher.Add(path)
		}
		return nil
	})
	log.Println("File watcher started.")
}

// --- versioning.go ---

const backupBucket = "versions"
const maxPatchChain = 50

type VersionRecord struct {
	ID        uint64    `json:"id"`
	OldSHA1   string    `json:"old_sha1"`
	NewSHA1   string    `json:"new_sha1"`
	Patch     string    `json:"patch"`
	Type      string    `json:"type"`
	Comment   string    `json:"comment"`
	Timestamp time.Time `json:"timestamp"`
}

type VersionManager struct {
	db *bbolt.DB
}

func NewVersionManager(user string) (*VersionManager, error) {
	dbPath := filepath.Join(AppConfig.MarkdownDir, user, ".extra", "versions.db")
	os.MkdirAll(filepath.Dir(dbPath), 0755)

	db, err := bbolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(backupBucket))
		return err
	})
	if err != nil {
		db.Close()
		return nil, err
	}

	return &VersionManager{db: db}, nil
}

func (vm *VersionManager) Close() {
	vm.db.Close()
}

func (vm *VersionManager) CreateBackup(filePath, oldSHA1, newSHA1, oldContent, newContent, comment string) error {
	return vm.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(backupBucket))
		fileBucket, err := b.CreateBucketIfNotExists([]byte(filePath))
		if err != nil {
			return err
		}

		stats := fileBucket.Stats()
		backupType := "patch"
		if stats.KeyN%maxPatchChain == 0 {
			backupType = "full"
		}

		patchContent := newContent
		if backupType == "patch" {
			dmp := diffmatchpatch.New()
			patches := dmp.PatchMake(oldContent, newContent)
			patchContent = dmp.PatchToText(patches)
		}

		id, _ := fileBucket.NextSequence()
		record := VersionRecord{
			ID:        id,
			OldSHA1:   oldSHA1,
			NewSHA1:   newSHA1,
			Patch:     patchContent,
			Type:      backupType,
			Comment:   comment,
			Timestamp: time.Now(),
		}

		buf, err := json.Marshal(record)
		if err != nil {
			return err
		}
		return fileBucket.Put(itob(id), buf)
	})
}

func (vm *VersionManager) GetHistory(filePath string) ([]VersionRecord, error) {
	var history []VersionRecord
	err := vm.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(backupBucket))
		fileBucket := b.Bucket([]byte(filePath))
		if fileBucket == nil {
			return nil
		}

		c := fileBucket.Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			var record VersionRecord
			if err := json.Unmarshal(v, &record); err == nil {
				history = append(history, record)
			}
		}
		return nil
	})
	return history, err
}

func (vm *VersionManager) GetVersionContent(filePath string, targetVersionID uint64) (string, error) {
	var recordsToApply []VersionRecord
	var baseContent string

	err := vm.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(backupBucket))
		fileBucket := b.Bucket([]byte(filePath))
		if fileBucket == nil {
			return fmt.Errorf("no history for file: %s", filePath)
		}

		c := fileBucket.Cursor()
		for k, v := c.Seek(itob(targetVersionID)); k != nil; k, v = c.Prev() {
			var record VersionRecord
			if err := json.Unmarshal(v, &record); err != nil {
				continue
			}
			recordsToApply = append(recordsToApply, record)
			if record.Type == "full" {
				baseContent = record.Patch
				break
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}
	if baseContent == "" {
		return "", fmt.Errorf("could not find a full backup base for version %d", targetVersionID)
	}

	dmp := diffmatchpatch.New()
	currentContent := baseContent
	for i := len(recordsToApply) - 2; i >= 0; i-- {
		record := recordsToApply[i]
		if record.Type == "patch" {
			patches, err := dmp.PatchFromText(record.Patch)
			if err != nil {
				return "", fmt.Errorf("error parsing patch for version %d: %v", record.ID, err)
			}
			newContent, _ := dmp.PatchApply(patches, currentContent)
			currentContent = newContent
		}
	}

	return currentContent, nil
}

// --- search.go ---

type SearchResult struct {
	Path    string   `json:"path"`
	Context []string `json:"context"`
}

func SearchInMemory(query string, useRegex bool, user string) []SearchResult {
	store.RLock()
	defer store.RUnlock()

	var results []SearchResult
	var re *regexp.Regexp
	var err error

	keywords := []string{}
	if useRegex {
		re, err = regexp.Compile(query)
		if err != nil {
			return results
		}
	} else {
		keywords = strings.Fields(strings.ToLower(query))
		if len(keywords) == 0 {
			return results
		}
	}

	userPrefix := user + string(filepath.Separator)
	for path, doc := range store.docs {
		if !strings.HasPrefix(path, userPrefix) {
			continue
		}

		contentLower := strings.ToLower(doc.Content)
		matched := false
		if useRegex {
			if re.MatchString(doc.Content) {
				matched = true
			}
		} else {
			allKeywordsFound := true
			for _, keyword := range keywords {
				if !strings.Contains(contentLower, keyword) {
					allKeywordsFound = false
					break
				}
			}
			if allKeywordsFound {
				matched = true
			}
		}

		if matched {
			context := getMatchContext(doc.Content, useRegex, re, keywords)
			results = append(results, SearchResult{
				Path:    strings.ReplaceAll(strings.TrimPrefix(doc.Path, userPrefix), string(filepath.Separator), "/"),
				Context: context,
			})
		}
	}
	return results
}

func getMatchContext(content string, useRegex bool, re *regexp.Regexp, keywords []string) []string {
	lines := strings.Split(content, "\n")
	var contextLines []string

	for i, line := range lines {
		lineLower := strings.ToLower(line)
		isMatch := false
		if useRegex {
			if re.MatchString(line) {
				isMatch = true
			}
		} else {
			for _, keyword := range keywords {
				if strings.Contains(lineLower, keyword) {
					isMatch = true
					break
				}
			}
		}

		if isMatch {
			contextLines = append(contextLines, fmt.Sprintf("%d: %s", i+1, line))
		}
		if len(contextLines) >= 5 {
			break
		}
	}
	return contextLines
}

// --- handlers.go ---

type TreeItem struct {
	Name        string      `json:"name"`
	IsDir       bool        `json:"is_dir"`
	Size        int64       `json:"size"`
	ModTime     time.Time   `json:"mod_time"`
	AttachCount int         `json:"attach_count,omitempty"`
	Children    []*TreeItem `json:"children,omitempty"`
}

type AttachmentInfo struct {
	Name       string    `json:"name"`
	AttachPath string    `json:"attachPath"`
	Size       int64     `json:"size"`
	ModTime    time.Time `json:"mod_time"`
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

func getUserPath(r *http.Request, subPath string) (basePath, fullPath, relPath string, err error) {
	user := r.Context().Value(userContextKey).(string)
	basePath = filepath.Join(AppConfig.MarkdownDir, user)

	cleanedSubPath := filepath.Clean(subPath)
	if strings.HasPrefix(cleanedSubPath, "..") || strings.Contains(cleanedSubPath, string(filepath.Separator)+"..") {
		return "", "", "", fmt.Errorf("invalid path: contains '..'")
	}

	fullPath = filepath.Join(basePath, cleanedSubPath)
	relPath = filepath.Join(user, cleanedSubPath)

	if !strings.HasPrefix(fullPath, basePath) {
		return "", "", "", fmt.Errorf("invalid path: outside of user directory")
	}
	return
}

func getSafeAttachmentPath(r *http.Request, mdPath, attachPath string) (string, error) {
	userBasePath, _, _, err := getUserPath(r, "")
	if err != nil {
		return "", err
	}

	_, mdFileAbsPath, _, err := getUserPath(r, mdPath)
	if err != nil {
		return "", err
	}
	mdDirAbsPath := filepath.Dir(mdFileAbsPath)

	resolvedPath := filepath.Join(mdDirAbsPath, attachPath)
	cleanedPath := filepath.Clean(resolvedPath)

	if !strings.HasPrefix(cleanedPath, userBasePath) {
		return "", fmt.Errorf("invalid attachment path: access denied, path escapes user root")
	}

	return cleanedPath, nil
}

func handleDirOp(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Action  string `json:"action"`
		Path    string `json:"path"`
		NewPath string `json:"new_path,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	_, fullPath, _, err := getUserPath(r, req.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	switch req.Action {
	case "create":
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to create directory: "+err.Error())
			return
		}
	case "delete":
		if err := os.RemoveAll(fullPath); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to delete directory: "+err.Error())
			return
		}
	case "rename":
		if req.NewPath == "" {
			respondError(w, http.StatusBadRequest, "Missing new_path for rename action")
			return
		}
		_, newFullPath, _, err := getUserPath(r, req.NewPath)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := os.Rename(fullPath, newFullPath); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to rename directory: "+err.Error())
			return
		}
	default:
		respondError(w, http.StatusBadRequest, "Invalid action")
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func handleFileWrite(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Path    string `json:"path"`
		Content string `json:"content"`
		Comment string `json:"comment,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if !strings.HasSuffix(strings.ToLower(req.Path), ".md") {
		respondError(w, http.StatusBadRequest, "File must have a .md extension")
		return
	}

	_, fullPath, relPath, err := getUserPath(r, req.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	os.MkdirAll(filepath.Dir(fullPath), 0755)

	var oldContent string
	var oldSHA1 string
	isNewFile := true

	store.RLock()
	if doc, exists := store.docs[relPath]; exists {
		oldContent = doc.Content
		oldSHA1 = doc.SHA1
		isNewFile = false
	}
	store.RUnlock()

	newContentBytes := []byte(req.Content)
	newSHA1 := calculateSHA1(newContentBytes)

	if !isNewFile && oldSHA1 == newSHA1 {
		respondJSON(w, http.StatusOK, map[string]string{"status": "no change"})
		return
	}

	if err := os.WriteFile(fullPath, newContentBytes, 0644); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to write file: "+err.Error())
		return
	}

	store.UpdateDoc(relPath, newContentBytes)

	if !isNewFile {
		user := r.Context().Value(userContextKey).(string)
		vm, err := NewVersionManager(user)
		if err != nil {
			log.Printf("Error creating version manager for %s: %v", user, err)
		} else {
			defer vm.Close()
			err := vm.CreateBackup(req.Path, oldSHA1, newSHA1, oldContent, req.Content, req.Comment)
			if err != nil {
				log.Printf("Error creating backup for %s: %v", relPath, err)
			}
		}
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "success", "sha1": newSHA1})
}

func handleFileRead(w http.ResponseWriter, r *http.Request) {
	pathParam := r.URL.Query().Get("path")
	if pathParam == "" {
		respondError(w, http.StatusBadRequest, "Missing path parameter")
		return
	}

	_, fullPath, relPath, err := getUserPath(r, pathParam)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	store.RLock()
	doc, exists := store.docs[relPath]
	store.RUnlock()

	if !exists {
		content, err := os.ReadFile(fullPath)
		if err != nil {
			respondError(w, http.StatusNotFound, "File not found")
			return
		}
		doc.Content = string(content)
		doc.Path = relPath
		doc.SHA1 = calculateSHA1(content)
	}

	respondJSON(w, http.StatusOK, doc)
}

func handleFileOp(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Action  string `json:"action"`
		Path    string `json:"path"`
		NewPath string `json:"new_path,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	_, fullPath, relPath, err := getUserPath(r, req.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	attachPath := fullPath + ".attach"

	switch req.Action {
	case "rename":
		if req.NewPath == "" {
			respondError(w, http.StatusBadRequest, "Missing new_path for rename action")
			return
		}
		if !strings.HasSuffix(strings.ToLower(req.NewPath), ".md") {
			respondError(w, http.StatusBadRequest, "New file name must have a .md extension")
			return
		}
		_, newFullPath, newRelPath, err := getUserPath(r, req.NewPath)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := os.Rename(fullPath, newFullPath); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to rename file: "+err.Error())
			return
		}
		if _, err := os.Stat(attachPath); err == nil {
			newAttachPath := newFullPath + ".attach"
			os.Rename(attachPath, newAttachPath)
		}

		store.DeleteDoc(relPath)
		content, _ := os.ReadFile(newFullPath)
		store.UpdateDoc(newRelPath, content)

	case "delete":
		doc, exists := store.docs[relPath]
		if !exists {
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				respondError(w, http.StatusNotFound, "File not found")
				return
			}
			content, _ := os.ReadFile(fullPath)
			doc.SHA1 = calculateSHA1(content)
		}

		user := r.Context().Value(userContextKey).(string)
		recycleDir := filepath.Join(AppConfig.MarkdownDir, user, ".extra", ".recycle", doc.SHA1)
		if err := os.MkdirAll(recycleDir, 0755); err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to create recycle dir")
			return
		}

		if err := os.Rename(fullPath, filepath.Join(recycleDir, filepath.Base(req.Path))); err != nil {
			os.Remove(fullPath)
		}
		if _, err := os.Stat(attachPath); err == nil {
			os.RemoveAll(attachPath)
		}

		store.DeleteDoc(relPath)
	default:
		respondError(w, http.StatusBadRequest, "Invalid action")
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func handleList(w http.ResponseWriter, r *http.Request) {
	pathParam := r.URL.Query().Get("path")
	recursive := r.URL.Query().Get("recursive") == "true"

	_, fullPath, _, err := getUserPath(r, pathParam)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var items []*TreeItem
	var listErr error

	if recursive {
		items, listErr = buildTree(fullPath)
	} else {
		entries, err := os.ReadDir(fullPath)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed to read directory: "+err.Error())
			return
		}

		items = make([]*TreeItem, 0, len(entries))
		for _, entry := range entries {
			name := entry.Name()
			if name == ".extra" || strings.HasSuffix(name, ".attach") {
				continue
			}

			info, err := entry.Info()
			if err != nil {
				log.Printf("Could not get file info for %s: %v", name, err)
				continue
			}

			item := &TreeItem{
				Name:    name,
				IsDir:   info.IsDir(),
				Size:    info.Size(),
				ModTime: info.ModTime(),
			}

			if !info.IsDir() && strings.HasSuffix(strings.ToLower(name), ".md") {
				attachDir := filepath.Join(fullPath, name+".attach")
				if attachEntries, err := os.ReadDir(attachDir); err == nil {
					item.AttachCount = len(attachEntries)
				}
			}
			items = append(items, item)
		}
	}

	if listErr != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list directory contents: "+listErr.Error())
		return
	}

	respondJSON(w, http.StatusOK, items)
}

func buildTree(dirPath string) ([]*TreeItem, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	items := make([]*TreeItem, 0)

	for _, entry := range entries {
		name := entry.Name()
		if name == ".extra" || strings.HasSuffix(name, ".attach") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			log.Printf("Could not get file info for %s in %s: %v", name, dirPath, err)
			continue
		}

		item := &TreeItem{
			Name:    name,
			IsDir:   info.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		}

		if info.IsDir() {
			children, err := buildTree(filepath.Join(dirPath, name))
			if err != nil {
				log.Printf("Error building subtree for %s: %v", name, err)
			}
			item.Children = children
		} else {
			if strings.HasSuffix(strings.ToLower(name), ".md") {
				attachDir := filepath.Join(dirPath, name+".attach")
				if attachEntries, err := os.ReadDir(attachDir); err == nil {
					item.AttachCount = len(attachEntries)
				}
			}
		}
		items = append(items, item)
	}

	return items, nil
}

func handleAttachUpload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	mdPath := r.FormValue("path")
	if mdPath == "" {
		respondError(w, http.StatusBadRequest, "Missing 'path' for markdown file")
		return
	}

	_, fullMdPath, _, err := getUserPath(r, mdPath)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if _, err := os.Stat(fullMdPath); os.IsNotExist(err) {
		respondError(w, http.StatusNotFound, "Markdown file does not exist")
		return
	}

	file, handler, err := r.FormFile("attachment")
	if err != nil {
		respondError(w, http.StatusBadRequest, "Missing 'attachment' file")
		return
	}
	defer file.Close()

	attachDir := fullMdPath + ".attach"
	if err := os.MkdirAll(attachDir, 0755); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create attachment directory")
		return
	}

	dstPath := filepath.Join(attachDir, handler.Filename)
	if filepath.Base(dstPath) != handler.Filename {
		respondError(w, http.StatusBadRequest, "Invalid attachment filename")
		return
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to save attachment")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to write attachment to disk")
		return
	}

	relativeAttachPath := filepath.ToSlash(filepath.Base(fullMdPath) + ".attach/" + handler.Filename)
	respondJSON(w, http.StatusOK, map[string]string{
		"status":     "success",
		"mdPath":     mdPath,
		"attachPath": relativeAttachPath,
	})
}

func handleAttachList(w http.ResponseWriter, r *http.Request) {
	mdPath := r.URL.Query().Get("path")
	if mdPath == "" {
		respondError(w, http.StatusBadRequest, "Missing 'path' for markdown file")
		return
	}

	_, fullMdPath, _, err := getUserPath(r, mdPath)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	attachDir := fullMdPath + ".attach"
	entries, err := os.ReadDir(attachDir)
	if err != nil {
		if os.IsNotExist(err) {
			respondJSON(w, http.StatusOK, map[string]interface{}{"mdPath": mdPath, "attachments": []AttachmentInfo{}})
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to read attachment directory")
		return
	}

	var attachments []AttachmentInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil || info.IsDir() {
			continue
		}
		relativeAttachPath := filepath.ToSlash(filepath.Base(fullMdPath) + ".attach/" + info.Name())
		attachments = append(attachments, AttachmentInfo{
			Name:       info.Name(),
			AttachPath: relativeAttachPath,
			Size:       info.Size(),
			ModTime:    info.ModTime(),
		})
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"mdPath":      mdPath,
		"attachments": attachments,
	})
}

func handleAttachGet(w http.ResponseWriter, r *http.Request) {
	requestedPath := chi.URLParam(r, "*")
	if requestedPath == "" {
		respondError(w, http.StatusBadRequest, "File path is required.")
		return
	}

	_, safeAbsPath, _, err := getUserPath(r, requestedPath)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	info, err := os.Stat(safeAbsPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
		} else {
			respondError(w, http.StatusInternalServerError, "Error checking file: "+err.Error())
		}
		return
	}

	if info.IsDir() {
		respondError(w, http.StatusBadRequest, "Serving directories is not allowed.")
		return
	}

	http.ServeFile(w, r, safeAbsPath)
}

func handleAttachDelete(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MdPath     string `json:"mdPath"`
		AttachPath string `json:"attachPath"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.MdPath == "" || req.AttachPath == "" {
		respondError(w, http.StatusBadRequest, "Missing 'mdPath' or 'attachPath' in request body")
		return
	}

	if !strings.Contains(filepath.ToSlash(req.AttachPath), ".md.attach/") {
		respondError(w, http.StatusBadRequest, "Deletion is only allowed for local attachments inside a '.attach' directory.")
		return
	}

	safeAbsPath, err := getSafeAttachmentPath(r, req.MdPath, req.AttachPath)
	if err != nil {
		if strings.Contains(err.Error(), "access denied") {
			respondError(w, http.StatusForbidden, err.Error())
		} else {
			respondError(w, http.StatusNotFound, err.Error())
		}
		return
	}

	if err := os.Remove(safeAbsPath); err != nil {
		if os.IsNotExist(err) {
			respondError(w, http.StatusNotFound, "Attachment not found")
		} else {
			respondError(w, http.StatusInternalServerError, "Failed to delete attachment: "+err.Error())
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

func handleHistory(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userContextKey).(string)
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		respondError(w, http.StatusBadRequest, "Missing path parameter")
		return
	}

	vm, err := NewVersionManager(user)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Could not open version db")
		return
	}
	defer vm.Close()

	history, err := vm.GetHistory(filePath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get history: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, history)
}

func handleVersionGet(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userContextKey).(string)
	filePath := r.URL.Query().Get("path")
	versionIDStr := r.URL.Query().Get("id")
	var versionID uint64
	fmt.Sscanf(versionIDStr, "%d", &versionID)

	if filePath == "" || versionID == 0 {
		respondError(w, http.StatusBadRequest, "Missing path or id parameter")
		return
	}

	vm, err := NewVersionManager(user)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Could not open version db")
		return
	}
	defer vm.Close()

	content, err := vm.GetVersionContent(filePath, versionID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get version content: "+err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"content": content})
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userContextKey).(string)
	query := r.URL.Query().Get("q")
	useRegex := r.URL.Query().Get("regex") == "true"

	if query == "" {
		respondJSON(w, http.StatusOK, []SearchResult{})
		return
	}

	results := SearchInMemory(query, useRegex, user)
	respondJSON(w, http.StatusOK, results)
}

// --- utils.go ---

func calculateSHA1(data []byte) string {
	h := sha1.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		b[i] = byte(v)
		v >>= 8
	}
	return b
}

func generateSelfSignedCert() error {
	if _, err := os.Stat(AppConfig.CertFile); err == nil {
		if _, err := os.Stat(AppConfig.KeyFile); err == nil {
			log.Println("Using existing certificate and key.")
			return nil
		}
	}
	log.Println("Generating self-signed certificate...")

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(time.Hour * 24 * 365),
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:    []string{"localhost"},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	certOut, err := os.Create(AppConfig.CertFile)
	if err != nil {
		return fmt.Errorf("failed to open cert.pem for writing: %w", err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, err := os.OpenFile(AppConfig.KeyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open key.pem for writing: %w", err)
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()

	log.Println("Self-signed certificate generated.")
	return nil
}

// --- main.go (entry point) ---

func unpackEmbeddedFS(destDir string) error {
	root := "embed"
	return fs.WalkDir(embeddedFS, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}
		destPath := filepath.Join(destDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		log.Printf("Unpacking: %s", relPath)
		fileData, err := embeddedFS.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(destPath, fileData, 0644)
	})
}

func main() {
	LoadConfig()

	if AppConfig.VisitLog != "" {
		logFile, err := os.OpenFile(AppConfig.VisitLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)
	}

	os.MkdirAll(AppConfig.MarkdownDir, 0755)
	os.MkdirAll(AppConfig.WWWDir, 0755)

	{
		exePath, err := os.Executable()
		if err != nil {
			log.Printf("WARNING: Could not determine executable path to write Help.md: %v", err)
		} else {
			programDir := filepath.Dir(exePath)
			rootHelpPath := filepath.Join(programDir, "Help.md")

			if _, err := os.Stat(rootHelpPath); os.IsNotExist(err) {
				log.Printf("Creating general user guide at: %s", rootHelpPath)
				if err := os.WriteFile(rootHelpPath, helpMarkdown, 0644); err != nil {
					log.Printf("WARNING: Failed to create user guide in program directory: %v", err)
				}
			}

			rootHelpEnPath := filepath.Join(programDir, "Help.en.md")
			if _, err := os.Stat(rootHelpEnPath); os.IsNotExist(err) {
				log.Printf("Creating general user guide at: %s", rootHelpPath)
				if err := os.WriteFile(rootHelpEnPath, helpMarkdownEn, 0644); err != nil {
					log.Printf("WARNING: Failed to create user guide in program directory: %v", err)
				}
			}
		}
	}

	LoadUsers()
	store.Scan()
	WatchMarkdownDir()
	StartBackupScheduler() // 新增: 启动备份调度器

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Api-Version", "Accept-Encoding"},
		ExposedHeaders:   []string{"Link", "Content-Encoding"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(corsMiddleware.Handler)
		r.Use(AuthMiddleware)
		r.Use(middleware.Compress(5, "application/json"))

		r.Post("/dir", handleDirOp)
		r.Get("/list", handleList)

		r.Post("/file", handleFileWrite)
		r.Get("/file", handleFileRead)
		r.Patch("/file", handleFileOp)

		r.Route("/attach", func(r chi.Router) {
			r.Post("/upload", handleAttachUpload)
			r.Get("/list", handleAttachList)
			r.Get("/get/*", handleAttachGet)
			r.Post("/delete", handleAttachDelete)
		})

		r.Get("/history", handleHistory)
		r.Get("/version", handleVersionGet)
		r.Get("/search", handleSearch)
	})

	if _, err := os.Stat(filepath.Join(AppConfig.WWWDir, "index.html")); os.IsNotExist(err) {
		log.Println("index.html not found in 'www' directory. Unpacking embedded assets...")
		if err := unpackEmbeddedFS(AppConfig.WWWDir); err != nil {
			log.Fatalf("Failed to unpack embedded assets: %v", err)
		}
		log.Println("Successfully unpacked assets to 'www' directory.")
	} else {
		log.Println("Found existing 'www/index.html'. Skipping asset unpacking.")
	}

	fs := http.FileServer(http.Dir(AppConfig.WWWDir))
	r.Handle("/*", fs)

	log.Printf("Starting server on %s", AppConfig.Bind)

	if AppConfig.TLS {
		if AppConfig.CertFile == "cert.pem" {
			if err := generateSelfSignedCert(); err != nil {
				log.Fatalf("Failed to generate self-signed certificate: %v", err)
			}
		}
		server := &http.Server{
			Addr:    AppConfig.Bind,
			Handler: r,
		}
		err := server.ListenAndServeTLS(AppConfig.CertFile, AppConfig.KeyFile)
		if err != nil {
			log.Fatalf("Failed to start HTTPS server: %v", err)
		}
	} else {
		err := http.ListenAndServe(AppConfig.Bind, r)
		if err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}
}
