package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	bindAddr = env("BIND_ADDR", ":80")
	htmlRoot = env("HTML_ROOT", "./statics/")

	swaggerRoot = filepath.Join(htmlRoot, "/swagger/")
	markedRoot  = filepath.Join(htmlRoot, "/marked/")
	textRoot    = filepath.Join(htmlRoot, "/text/")

	docRoot = filepath.Join(htmlRoot, "/docs/")
)

type entry struct {
	APP string
	URL template.URL
}

func main() {
	docFS := http.FileServer(http.Dir(swaggerRoot))
	http.Handle("/swagger/", http.StripPrefix("/swagger/", docFS))
	markedFS := http.FileServer(http.Dir(markedRoot))
	http.Handle("/marked/", http.StripPrefix("/marked/", markedFS))
	textFS := http.FileServer(http.Dir(textRoot))
	http.Handle("/text/", http.StripPrefix("/text/", textFS))

	apiFS := http.FileServer(http.Dir(docRoot))
	http.Handle("/apis/", http.StripPrefix("/apis/", apiFS))

	indexPath := filepath.Join(htmlRoot, "index.html")
	t, err := template.ParseFiles(indexPath)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var entries []*entry
		err := filepath.Walk(docRoot, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			relativePath := strings.TrimPrefix(path, docRoot)
			apiURI := filepath.ToSlash(filepath.Join("/apis/", relativePath))
			ext := filepath.Ext(apiURI)

			var app string
			switch ext {
			case ".md":
				app = "marked"
			case ".json", ".yml", ".yaml":
				app = "swagger"
			default:
				app = "text"
			}
			entry := &entry{
				APP: app,
				URL: template.URL(apiURI),
			}
			entries = append(entries, entry)
			return nil
		})
		if err != nil {
			_, _ = w.Write([]byte(fmt.Sprintf("list apis failed: %s", err)))
			return
		}

		if err := t.Execute(w, entries); err != nil {
			_, _ = w.Write([]byte(fmt.Sprintf("list apis failed: %s", err)))
			return
		}
	})

	log.Fatal(http.ListenAndServe(bindAddr, nil))
}

func env(name string, defaultValue string) string {
	if v := os.Getenv(name); v != "" {
		return v
	} else {
		return defaultValue
	}
}
