package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type FileInfo struct {
	Name  string
	IsDir bool
	Mode  os.FileMode
}

const (
	filePrefix = "/music/"
	root       = "./music"
)

var PORT = os.Getenv("PORT")

func main() {
	http.HandleFunc("/", playerMainFrame)
	http.HandleFunc(filePrefix, File)
	http.HandleFunc("/music/search", DownloadRequest)
	http.ListenAndServe(":"+PORT, nil)
}

func playerMainFrame(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./player.html")
}

func File(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(root, r.URL.Path[len(filePrefix):])
	stat, err := os.Stat(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if stat.IsDir() {
		serveDir(w, r, path)
		return
	}
	http.ServeFile(w, r, path)
}

func serveDir(w http.ResponseWriter, r *http.Request, path string) {
	defer func() {
		if err, ok := recover().(error); ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	if err != nil {
		panic(err)
	}

	files, err := file.Readdir(-1)
	if err != nil {
		panic(err)
	}

	fileinfos := make([]FileInfo, len(files))

	for i := range files {
		fileinfos[i].Name = files[i].Name()
		fileinfos[i].IsDir = files[i].IsDir()
		fileinfos[i].Mode = files[i].Mode()
	}

	j := json.NewEncoder(w)

	if err := j.Encode(&fileinfos); err != nil {
		panic(err)
	}
}

func DownloadRequest(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("search")
	err := YoutubeDL(query)
        check(err)
	w.WriteHeader(http.StatusOK)
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func YoutubeDL(query string) error {
	cmd := `yt-dlp "ytsearch:`+query+`" --format bestaudio --output "music/%(title)s.%(ext)s"`
	proc := exec.Command("bash", "-c", cmd)
	return proc.Run()
}
