package main

import (
	"flag"
	"log"
	"os"
	"time"

	mc "github.com/sky0621/work-go-movieconverter"
)

func main() {
	var inputDir string
	var outputDir string
	var sleep time.Duration
	var scale string
	var logDir string

	flag.StringVar(&inputDir, "i", "in", "監視対象ディレクトリ")
	flag.StringVar(&outputDir, "o", "out", "変換結果出力先ディレクトリ")
	flag.DurationVar(&sleep, "t", 600, "監視間隔（秒）")
	flag.StringVar(&scale, "s", "640", "ffmpeg変換時スケール")
	flag.StringVar(&logDir, "l", ".", "ログ出力先ディレクトリ")
	flag.Parse()

	logfile, err := mc.SetupLog(logDir)
	if err != nil {
		os.Exit(1)
	}
	defer logfile.Close()

	log.Println("[START]movieconverter")
	mc.Run(&mc.ConvertInfo{
		InputDir:  inputDir,
		OutputDir: outputDir,
		Filename:  "",
		Scale:     scale},
		sleep)
	log.Println("[END]movieconverter")
}
