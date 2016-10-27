package main

import (
	"flag"
	"log"
	"os"
	"time"

	mc "github.com/sky0621/work-go-movieconverter"
)

func main() {
	log.Println("[main]START")
	var inputDir string
	var outputDir string
	var deployDir string
	var sleep time.Duration
	var scale string
	var logDir string

	flag.StringVar(&inputDir, "i", "in", "監視対象ディレクトリ")
	flag.StringVar(&outputDir, "o", "out", "変換結果出力先ディレクトリ")
	flag.StringVar(&deployDir, "d", "deploy", "変換結果デプロイ先ディレクトリ")
	flag.DurationVar(&sleep, "t", 600, "監視間隔（秒）")
	flag.StringVar(&scale, "s", "640", "ffmpeg変換時スケール")
	flag.StringVar(&logDir, "l", ".", "ログ出力先ディレクトリ")
	flag.Parse()
	log.Println("[main]flag parse fin.")

	logfile, err := mc.SetupLog(logDir)
	if err != nil {
		os.Exit(1)
	}
	defer logfile.Close()
	log.Println("[main]setup log fin.")

	log.Println("[main]go loop!!!")
	for {
		if mc.IsRunning {
			log.Println("[main] convert is running... ... ...")
		} else {
			go mc.Run(&mc.ConvertInfo{
				InputDir:  inputDir,
				OutputDir: outputDir,
				DeployDir: deployDir,
				Filename:  "",
				Scale:     scale})
		}

		time.Sleep(sleep * time.Second)
	}
}
