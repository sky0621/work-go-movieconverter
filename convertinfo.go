package movieconverter

import (
	"log"
	"path/filepath"
)

// ConvertInfo ...
type ConvertInfo struct {
	InputDir     string
	WatchListDir string
	OutputDir    string
	DeployDirVdo string
	DeployDirImg string
	Filename     string
	Scale        string
}

func (c *ConvertInfo) inputPath() string {
	// s := c.InputDir + "/" + c.Filename
	s := filepath.Join(c.InputDir, c.Filename)
	log.Println("[ConvertInfo.inputPath()]", s)
	return s
}

func (c *ConvertInfo) outputPath() string {
	// s := c.OutputDir + "/" + c.Filename
	s := filepath.Join(c.OutputDir, c.Filename)
	// log.Println("[outputPath]", s)
	return s
}

func (c *ConvertInfo) deployPath(targetFile string) string {
	// 動画と画像の振り分け（※拡張子完全固定・・・）
	switch filepath.Ext(targetFile) {
	case ".mp4":
		return joinPath(c.DeployDirVdo, targetFile)
	case ".jpg":
		return joinPath(c.DeployDirImg, targetFile)
	}
	return ""
}

func (c *ConvertInfo) cmdConvertVideo() string {
	// s := "cp " + c.inputPath() + " " + c.outputPath()
	s := "ffmpeg -i " + c.inputPath() + " -vf scale=" + c.Scale + ":-1 rotate=0 " + c.outputPath()
	// log.Println("[cmdConvertVideo]", s)
	return s
}

func (c *ConvertInfo) cmdCreateThumbnail() string {
	// s := "ls"
	s := "ffmpeg -i " + c.outputPath() + " -ss 1 -t 1 -r 1 -f image2 " + c.outputPath() + ".jpg"
	// log.Println("[cmdCreateThumbnail]", s)
	return s
}

func (c *ConvertInfo) cmdRotateThumbnail() string {
	// s := "ls"
	s := "convert -rotate 90 -resize 100x " + c.outputPath() + ".jpg " + c.outputPath() + "r.jpg"
	// log.Println("[cmdRotateThumbnail]", s)
	return s
}

func joinPath(paths ...string) string {
	var s string
	for _, path := range paths {
		if s == "" {
			s = path
		} else {
			s = s + "/" + path
		}
	}
	// s = filepath.Join(paths...)
	// log.Println("[joinPath]", s)
	return s
}
