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
	var watchListDir string
	var outputDir string
	var deployDirVdo string
	var deployDirImg string
	var sleep time.Duration
	var scale string
	var logDir string

	flag.StringVar(&inputDir, "i", "in", "監視対象ディレクトリ")
	flag.StringVar(&watchListDir, "w", "list", "watch.list格納ディレクトリ")
	flag.StringVar(&outputDir, "o", "out", "変換結果出力先ディレクトリ")
	flag.StringVar(&deployDirVdo, "dm", "deploy/vdo", "変換結果（動画）デプロイ先ディレクトリ")
	flag.StringVar(&deployDirImg, "di", "deploy/img", "変換結果（サムネイル画像）デプロイ先ディレクトリ")
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

	log.Println("[main]args")
	log.Printf("[main] -> 監視対象ディレクトリ「%s」\n", inputDir)
	log.Printf("[main] -> watch.list格納ディレクトリ「%s」\n", watchListDir)
	log.Printf("[main] -> 変換結果出力先ディレクトリ「%s」\n", outputDir)
	log.Printf("[main] -> 変換結果（動画）デプロイ先ディレクトリ「%s」\n", deployDirVdo)
	log.Printf("[main] -> 変換結果（サムネイル画像）デプロイ先ディレクトリ「%s」\n", deployDirImg)
	log.Printf("[main] -> 監視間隔（秒）「%d」\n", sleep)
	log.Printf("[main] -> ffmpeg変換時スケール「%v」\n", scale)
	log.Printf("[main] -> ログ出力先ディレクトリ「%s」\n", logDir)

	log.Println("[main]go loop!!!")
	for {
		if mc.IsRunning {
			log.Println("[main]convert is running... ... ...")
		} else {
			go mc.Run(&mc.ConvertInfo{
				InputDir:     inputDir,
				WatchListDir: watchListDir,
				OutputDir:    outputDir,
				DeployDirVdo: deployDirVdo,
				DeployDirImg: deployDirImg,
				Filename:     "",
				Scale:        scale})
		}

		time.Sleep(sleep * time.Second)
	}
}
