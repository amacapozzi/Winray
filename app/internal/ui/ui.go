package ui

import (
	"encoding/json"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	webview "github.com/jchv/go-webview2"

	"winray-app/internal/index"
	"winray-app/internal/models"
	"winray-app/internal/windows"
)

var (
	uiRunning atomic.Bool
	uiMu      sync.Mutex
	uiWv      webview.WebView
)

const UI_PATH = "https://winray.vercel.app/"

type UI struct{}

func New() *UI {
	return &UI{}
}

func (ui *UI) Open() {

	uiMu.Lock()
	defer uiMu.Unlock()

	if uiRunning.Load() {
		return
	}
	go func() {
		// Esperamos un poco (ej. 500ms o 1 segundo) para asegurar
		// que la ventana gráfica ya se creó y es visible.
		time.Sleep(1 * time.Second)

		// Esta función buscará tu ventana por el ID del proceso,
		// le quitará los bordes y la pondrá "Always on Top".
		windows.SetupAppWindow()
	}()

	go func() {
		runtime.LockOSThread()

		w := webview.NewWithOptions(webview.WebViewOptions{
			Debug:     false,
			AutoFocus: true,
			WindowOptions: webview.WindowOptions{
				Center: true,
			},
		})
		uiWv = w
		uiRunning.Store(true)

		w.SetSize(450, 420, webview.HintFixed)

		_ = w.Bind("goSearch", func(q string) {
			results := index.Search(q, 60)
			pushResults(w, results)
		})
		_ = w.Bind("goOpen", func(path string, kind string) {
			windows.ShellOpen(path)
		})
		_ = w.Bind("goHide", func() {
			w.Terminate()
		})
		_ = w.Bind("goStartIndexing", func() {
			recentFiles := index.GetRecentFiles(60)

			if len(recentFiles) > 0 {
				pushResults(w, recentFiles)
				setLoading(w, false)
			} else {
				setLoading(w, true)
				go buildInitialIndexProgressive(w)
			}
		})

		w.Navigate(UI_PATH)

		w.Init(`
			if (window.winrayActivate) window.winrayActivate();
		`)

		go func() {
			time.Sleep(400 * time.Millisecond)
			recentFiles := index.GetRecentFiles(60)
			if len(recentFiles) > 0 {
				pushResults(w, recentFiles)
				setLoading(w, false)
			} else {
				setLoading(w, true)
				buildInitialIndexProgressive(w)
			}
		}()

		w.Run()

		w.Destroy()

		uiMu.Lock()
		uiWv = nil
		uiRunning.Store(false)
		uiMu.Unlock()
	}()
}

func (ui *UI) Close() {
	uiMu.Lock()
	w := uiWv
	uiMu.Unlock()

	if w == nil {
		return
	}

	w.Dispatch(func() {
		w.Terminate()
	})
}

func (ui *UI) Toggle() {
	if uiRunning.Load() {
		ui.Close()
		return
	}
	ui.Open()
}

func pushResults(w webview.WebView, results []models.FileResult) {
	b, _ := json.Marshal(results)
	w.Dispatch(func() {
		w.Eval("window.setResults && window.setResults(" + string(b) + ");")
	})
}

func setLoading(w webview.WebView, loading bool) {
	w.Dispatch(func() {
		if loading {
			w.Eval("window.setLoading && window.setLoading(true);")
		} else {
			w.Eval("window.setLoading && window.setLoading(false);")
		}
	})
}

func appendResults(w webview.WebView, results []models.FileResult) {
	b, _ := json.Marshal(results)
	w.Dispatch(func() {
		w.Eval("window.appendResults && window.appendResults(" + string(b) + ");")
	})
}

func buildInitialIndexProgressive(w webview.WebView) {
	go func() {
		index.BuildInitialProgressive(w, appendResults, setLoading)
	}()
}
