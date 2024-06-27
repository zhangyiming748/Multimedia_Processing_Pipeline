package main

import (
	"Multimedia_Processing_Pipeline/constant"
	mylog "Multimedia_Processing_Pipeline/log"
	"Multimedia_Processing_Pipeline/merge"
	"Multimedia_Processing_Pipeline/replace"
	"Multimedia_Processing_Pipeline/sql"
	translateShell "Multimedia_Processing_Pipeline/translate"
	"Multimedia_Processing_Pipeline/util"
	"Multimedia_Processing_Pipeline/whisper"
	"Multimedia_Processing_Pipeline/ytdlp"
	"log"
	"os"
	"strings"
)

func initConfig(p *constant.Param) {
	if !util.IsExistCmd("ffmpeg") {
		log.Fatalln("ffmpeg未安装")
	}
	if !util.IsExistCmd("whisper") {
		log.Fatalln("whisper未安装")
	}
	if !util.IsExistCmd("trans") {
		log.Fatalln("trans未安装")
	}
	if !util.IsExistCmd("yt-dlp") {
		log.Fatalln("yt-dlp")
	}
	mylog.SetLog(p)
	sql.SetDatabase(p)
	util.ExitAfterRun()
	replace.SetSensitive(p)
}
func main() {
	p := &constant.Param{
		Root:     "/mnt/c/Users/zen",
		Language: "English",
		Pattern:  "mp4",
		Model:    "base",
		Location: "/mnt/c/Users/zen",
		Proxy:    "127.0.0.1:8889",
		Merge:    false,
	}
	initConfig(p)
	if root := os.Getenv("root"); root != "" {
		p.SetRoot(root)
	}
	if language := os.Getenv("language"); language != "" {
		p.SetLanguage(language)
	}
	if pattern := os.Getenv("pattern"); pattern != "" {
		p.SetPattern(pattern)
	}
	if model := os.Getenv("model"); model != "" {
		p.SetModel(model)
	}
	if location := os.Getenv("location"); location != "" {
		p.SetLocation(location)
	}
	if proxy := os.Getenv("proxy"); proxy != "" {
		p.SetProxy(proxy)
	}
	if combination := os.Getenv("merge"); combination == "1" {
		p.Merge = true
	}
	var c constant.Count
	uris := strings.Join([]string{p.GetRoot(), "urls.list"}, string(os.PathSeparator))
	lines := util.ReadByLine(uris)
	for _, line := range lines {
		video, err := ytdlp.DownloadVideo(line, p)
		if err != nil {
			continue
		}
		whisper.GetSubtitle(video, p)
		translateShell.Trans(video, p, &c)
		merge.MkvWithAss(video, p)
	}
	//replace.Replace(p)
}
