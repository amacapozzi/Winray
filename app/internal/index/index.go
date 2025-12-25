package index

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"winray-app/internal/models"
)

var (
	indexMu     sync.RWMutex
	folderIndex []models.IndexedFile
	fileIndex   []models.IndexedFile
)

func BuildInitial() {
	home, _ := os.UserHomeDir()
	roots := []string{
		filepath.Join(home, "Desktop"),
		filepath.Join(home, "Documents"),
		filepath.Join(home, "Downloads"),
	}

	var folders []models.IndexedFile
	var files []models.IndexedFile
	seen := make(map[string]struct{}, 15000)

	for _, root := range roots {
		_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}

			if _, ok := seen[path]; ok {
				return nil
			}
			seen[path] = struct{}{}

			info, err := d.Info()
			if err != nil {
				return nil
			}

			if d.IsDir() {
				name := strings.ToLower(d.Name())
				if name == "node_modules" || name == ".git" || name == "dist" || name == "build" {
					return filepath.SkipDir
				}
				folders = append(folders, models.IndexedFile{
					Path:           path,
					LastAccessTime: info.ModTime().Unix(),
				})
				return nil
			}

			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".tmp" || ext == ".log" {
				return nil
			}

			files = append(files, models.IndexedFile{
				Path:           path,
				LastAccessTime: info.ModTime().Unix(),
			})

			return nil
		})
	}

	indexMu.Lock()
	folderIndex = folders
	fileIndex = files
	indexMu.Unlock()
}

func GetRecentFiles(limit int) []models.FileResult {
	indexMu.RLock()
	defer indexMu.RUnlock()

	var allFiles []models.IndexedFile

	allFiles = append(allFiles, fileIndex...)
	allFiles = append(allFiles, folderIndex...)

	if len(allFiles) == 0 {
		return nil
	}

	sort.Slice(allFiles, func(i, j int) bool {
		return allFiles[i].LastAccessTime > allFiles[j].LastAccessTime
	})

	n := min(limit, len(allFiles))
	results := make([]models.FileResult, 0, n)
	for i := 0; i < n; i++ {
		results = append(results, ToResult(allFiles[i].Path, time.Unix(allFiles[i].LastAccessTime, 0)))
	}

	return results
}

func Search(query string, limit int) []models.FileResult {
	q := strings.TrimSpace(strings.ToLower(query))

	indexMu.RLock()
	defer indexMu.RUnlock()

	var hits []models.FileResult

	if q == "" {
		return GetRecentFiles(limit)
	}

	for _, f := range fileIndex {
		name := strings.ToLower(filepath.Base(f.Path))
		if strings.Contains(name, q) || strings.Contains(strings.ToLower(f.Path), q) {
			lastAccess := time.Unix(f.LastAccessTime, 0)
			hits = append(hits, ToResult(f.Path, lastAccess))
			if len(hits) >= limit {
				break
			}
		}
	}

	for _, f := range folderIndex {
		name := strings.ToLower(filepath.Base(f.Path))
		if strings.Contains(name, q) || strings.Contains(strings.ToLower(f.Path), q) {
			lastAccess := time.Unix(f.LastAccessTime, 0)
			hits = append(hits, ToResult(f.Path, lastAccess))
			if len(hits) >= limit {
				break
			}
		}
	}

	return hits
}

func ToResult(path string, lastAccess time.Time) models.FileResult {
	kind := KindFromPath(path)
	var lastAccessTime int64
	if !lastAccess.IsZero() {
		lastAccessTime = lastAccess.Unix()
	}

	return models.FileResult{
		ID:             path,
		Name:           filepath.Base(path),
		Path:           path,
		Kind:           kind,
		MetaLeft:       filepath.Dir(path),
		LastAccessTime: lastAccessTime,
	}
}

func KindFromPath(path string) string {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return "Folder"
	}

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".exe":
		return "App"
	case ".lnk":
		return "Shortcut"
	default:
		return "File"
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
