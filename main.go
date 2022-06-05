package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type BiliDownload struct {
	audio string
	video string
}

func main() {
	println("Searching ...")
	dirs, _ := findAllDirs(".")
	files := SearchForDownloadedFiles(dirs)
	// print statics
	println("Total: " + strconv.Itoa(len(files)))
	// if is not empty
	if len(files) > 0 {
		// make output if not exists
		if _, err := os.Stat("output"); os.IsNotExist(err) {
			err := os.Mkdir("output", 0755)
			if err != nil {
				panic(err)
			}
		}

		Merge(files)
	}
	println("Done!")
}

// Merge 合并m4s文件
func Merge(files []BiliDownload) {
	//	download the list
	for _, download := range files {
		// Get the title
		title := GetTitle(download.video)
		println("Merging " + title)
		//	merge the video
		cmd := exec.Command("ffmpeg", "-i", download.video, "-i", download.audio, "-codec", "copy", "output/"+title+".mp4")
		stdout, err := cmd.Output()

		if err != nil {
			fmt.Println("Merge " + download.video + " Error：" + err.Error())
			fmt.Println(string(stdout))
			println("Command as:" + cmd.String())
		}
	}
}

var errorCount = 0

// GetTitle 获取视频标题
func GetTitle(video string) string {
	dir := filepath.Dir(video)
	dir = filepath.Dir(dir)
	//	 open json
	file, err := os.Open(dir + "/entry.json")
	if err != nil {
		fmt.Println("Open Error: " + dir + "/entry.json")
		errorCount++
		return "Unknown Title" + strconv.Itoa(errorCount)
	}
	defer file.Close()
	//	read json
	jsonBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Read Error: " + dir + "/entry.json")
		errorCount++
		return "Unknown Title" + strconv.Itoa(errorCount)
	}
	//	parse json
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	if err != nil {
		fmt.Println("Parse Error: " + dir + "/entry.json")
		errorCount++
		return "Unknown Title" + strconv.Itoa(errorCount)
	}
	//	get title
	title := jsonMap["title"].(string)
	return title
}

// SearchForDownloadedFiles 查找文件夹列表中已经下载的视频
func SearchForDownloadedFiles(dirs []string) []BiliDownload {
	// make a BiliDownload list
	downloads := make([]BiliDownload, 0)
	//	print the list
	for _, dir := range dirs {
		println(dir)
		//	check if the path exists
		if _, err := os.Stat(dir + "\\audio.m4s"); err == nil {
			downloads = append(downloads, BiliDownload{dir + "\\audio.m4s", dir + "\\video.m4s"})
		}
	}
	return downloads
}

// findAllDirs 返回当前目录下所有文件夹列表
func findAllDirs(path string) ([]string, error) {
	dirs := make([]string, 0)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return dirs, nil
}
