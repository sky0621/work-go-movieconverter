package main

import (
	"flag"
	"log"
	"os"

	mc "github.com/sky0621/work-go-movieconverter"
)

func main() {
	var targetDir string
	var outputDir string
	var logDir string

	flag.StringVar(&targetDir, "t", "in", "監視対象ディレクトリ")
	flag.StringVar(&outputDir, "o", "out", "変換結果出力先ディレクトリ")
	flag.StringVar(&logDir, "l", ".", "ログ出力先ディレクトリ")
	flag.Parse()

	logfile, err := mc.SetupLog(logDir)
	if err != nil {
		os.Exit(1)
	}
	defer logfile.Close()

	log.Println("[START]movieconverter")
	mc.Run(targetDir, outputDir)
	log.Println("[END]movieconverter")
}
