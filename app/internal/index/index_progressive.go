package index

import (
	"os"
	"path/filepath"
	"strings"

	webview "github.com/jchv/go-webview2"

	"winray-app/internal/models"
)

func BuildInitialProgressive(
	w webview.WebView,
	appendResults func(webview.WebView, []models.FileResult),
	setLoading func(webview.WebView, bool),
) {
	setLoading(w, true)

	home, _ := os.UserHomeDir()
	roots := []string{
		filepath.Join(home, "Desktop"),
		filepath.Join(home, "Documents"),
		filepath.Join(home, "Downloads"),
	}

	var folders []models.IndexedFile
	var files []models.IndexedFile
	seen := make(map[string]struct{}, 15000)
	var batch []models.FileResult
	batchSize := 50

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
				folderFile := models.IndexedFile{
					Path:           path,
					LastAccessTime: info.ModTime().Unix(),
				}
				folders = append(folders, folderFile)
				batch = append(batch, ToResult(path, info.ModTime()))
				if len(batch) >= batchSize {
					appendResults(w, batch)
					batch = batch[:0]
				}
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

			batch = append(batch, ToResult(path, info.ModTime()))
			if len(batch) >= batchSize {
				appendResults(w, batch)
				batch = batch[:0]
			}

			return nil
		})
	}

	if len(batch) > 0 {
		appendResults(w, batch)
	}

	indexMu.Lock()
	folderIndex = folders
	fileIndex = files
	indexMu.Unlock()

	setLoading(w, false)
}
