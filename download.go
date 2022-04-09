package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

var (
	IDre = regexp.MustCompile("[youtube](.*): Downloading webpage")
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

func SearchRequest(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("search")
	ID := SearchYT(query)
	data, err := ioutil.ReadFile("./search/" + ID + ".info.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(data)
}

func YoutubeDL(query string) error {
	proc, err, a := Bash([]string{"yt-dlp", "ytsearch:" + query, "--format", "bestaudio", "--output", "music/%(title)s.%(ext)s"})
	fmt.Println(proc, err)
	return a
}

func SearchYT(query string) string {
	proc, _, _ := Bash([]string{"yt-dlp", "ytsearch:" + query, "--format", "bestaudio", "--output", "search/%(id)s.%(ext)s", "--write-info-json", "--write-thumbnail", "--skip-download"})
	var ID string
	for _, b := range IDre.FindAllStringSubmatch(proc, -1) {
		UID := b[1]
		ID = strings.SplitAfter(UID, " ")[1]
	}
	return ID
}

func Bash(cmd []string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	proc := exec.Command(cmd[0], cmd[1:]...)
	proc.Stdout = &stdout
	proc.Stderr = &stderr
	err := proc.Run()
	return stdout.String(), stderr.String(), err
}
