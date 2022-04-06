package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	yt "github.com/SherlockYigit/youtube-go"
	youtube "github.com/kkdai/youtube/v2"
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
	http.HandleFunc("/music/search", XD)
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

func XD(w http.ResponseWriter, r *http.Request) {
	q := r.FormValue("search")
	fmt.Println(q)
	DownloadSong(q)
	fmt.Println("Downloaded")
	w.WriteHeader(http.StatusOK)
}

func DownloadSong(q string) {
        YoutubeDL(q)
        return
	result := SearchVideos(q, 1)
	Result := result[0]
	youtube := youtube.Client{}
	video, err := youtube.GetVideo("https://www.youtube.com/watch?v=" + Result.ID)
	check(err)
	stream, _, err := youtube.GetStream(video, video.Formats.FindByQuality("tiny"))
	check(err)
	defer stream.Close()
	outFile, err := os.Create(root + "/song.mp3")
	check(err)
	defer outFile.Close()
	io.Copy(outFile, stream)

}

func SearchVideos(q string, limit int) []YoutubeVideo {
	vids := yt.Search(q, yt.SearchOptions{
		Limit: limit,
		Type:  "video",
	})
	var results = make([]YoutubeVideo, limit)
	for i, v := range vids {
		results[i] = YoutubeVideo{
			ID:        v.Video.Id,
			Title:     v.Title,
			URL:       v.Video.Url,
			Thumbnail: v.Video.Thumbnail.Url,
			Channel:   v.Channel.Name,
		}
	}
	return results
}

type YoutubeVideo struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
	Channel   string `json:"channel"`
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func YoutubeDL(query string) bool {
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "music")
	cmd := `yt-dlp "ytsearch:`+query+`" --format bestaudio --output "music/%(title)s.%(ext)s"`
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	proc := exec.Command("bash", "-c", cmd)
	proc.Stdout = &stdout
	proc.Stderr = &stderr
	err := proc.Run()
	check(err)
        if err != nil {
return true
} 
return false
}
