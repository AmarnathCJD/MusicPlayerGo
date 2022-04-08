package main

import (
	"fmt"
	"net/http"
	"os/exec"
)

func DownloadRequest(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("search")
	err := YoutubeDL(query)
	if err == nil {
		fmt.Printf("Downloading %s", query)
		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Printf("Error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func YoutubeDL(query string) error {
	cmd := `yt-dlp "ytsearch:` + query + `" --format bestaudio --output "music/%(title)s.%(ext)s"`
	proc := exec.Command("bash", "-c", cmd)
	return proc.Run()
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
