package movieconverter

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Run ...
func Run(targetDir string, outputDir string) {
	for {
		fileInfoArray, err := ioutil.ReadDir(targetDir)
		if err != nil {
			log.Printf("指定ディレクトリ(%s)配下の動画ファイル一覧読み込み時にエラーが発生しました。 [ERROR]%s\n", targetDir, err)
			return
		}

		for _, fileInfo := range fileInfoArray {
			runConvertMovies(targetDir, fileInfo, outputDir)
		}

		time.Sleep(600 * time.Second)
	}
}

func runConvertMovies(targetDir string, fileInfo os.FileInfo, outputDir string) {
	cmdStr := "ffmpeg -i " + targetDir + "/" + fileInfo.Name() + " -vf scale=640:-1 " + outputDir + "/" + fileInfo.Name()
	cmd := exec.Command(os.Getenv("SHELL"), "-c", cmdStr)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("ffmpeg fin [%s]\n", fileInfo.Name())
	os.Remove(filepath.Join(targetDir, fileInfo.Name()))
}
