package movieconverter

import (
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const listfile = "watch.list"

// Run ...
func Run(inputDir string, outputDir string, scale string, sleep time.Duration) {
	// まず、リストアップ
	err := listup(inputDir)
	if err != nil {
		return
	}

	// そして、リストアップしたファイルをコンバート
	err = convert(inputDir, outputDir, scale)
	if err != nil {
		return
	}
	// rerr := os.Remove(filepath.Join(inputDir, listfile))
	rerr := os.Remove(inputDir + "/" + listfile)
	if rerr != nil {
		log.Println("[after convert]", rerr)
		return
	}

	// 最後に、コンバートしたファイルをデプロイ

}

func listup(inputDir string) error {
	log.Println("[listup]START")

	// リストアップファイルを作成（※作成済み＝後続タスクで未処理のため、以降の処理は行わない）
	cerr := createFileList(inputDir)
	if cerr != nil {
		log.Println("[listup]createFileList:", cerr)
		return cerr
	}

	// 監視ディレクトリ配下のファイル情報一覧を取得
	fileInfos, rerr := ioutil.ReadDir(inputDir)
	if rerr != nil {
		log.Println("[listup]ioutil.ReadDir(inputDir):", rerr)
		return rerr
	}

	// リストアップファイルにファイル情報（ファイル名、作成日時）を出力
	file, oerr := os.OpenFile(filepath.Join(inputDir, listfile), os.O_RDWR, 0644)
	if oerr != nil {
		log.Println("[listup]os.Open(inputDir + listfile):", oerr)
		return oerr
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, fileInfo := range fileInfos {
		if listfile == fileInfo.Name() {
			continue
		}
		_, err := writer.WriteString(fileInfo.Name() + "\t" + fileInfo.ModTime().Format(time.ANSIC) + "\r\n")
		if err != nil {
			log.Println("[listup]writer.WriteString:", oerr)
		}
	}
	writer.Flush()

	log.Println("[listup]END")
	return nil
}

func createFileList(inputDir string) error {
	_, serr := os.Stat(filepath.Join(inputDir, listfile))
	if serr == nil {
		return errors.New("listfile exists")
	}
	file, err := os.Create(filepath.Join(inputDir, listfile))
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

var isConverting bool

type convertInfo struct {
	inputDir  string
	outputDir string
	filename  string
	scale     string
}

const maxProcesses = 5

func convert(inputDir string, outputDir string, scale string) error {
	log.Println("[convert]START")
	isConverting = true

	sema := make(chan int, maxProcesses)

	notify := make(chan int)

	// リストアップファイルの中身を読み取る
	file, oerr := os.Open(filepath.Join(inputDir, listfile))
	if oerr != nil {
		log.Println("[convert]os.Open(inputDir + listfile):", oerr)
		return oerr
	}
	defer file.Close()

	no := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		no++
		list := scanner.Text()
		lists := strings.Split(list, "\t")
		log.Println(lists[0])
		fInfo := &convertInfo{inputDir: inputDir, outputDir: outputDir, filename: lists[0], scale: scale}
		go runConvertMovies(fInfo, sema, notify, no)
	}
	serr := scanner.Err()
	if serr != nil {
		log.Println("[convert]scanner.Err():", serr)
		return serr
	}

	for i := 0; i < no; i++ {
		<-notify // 他のゴルーチンからの終了通知を待つ
	}

	isConverting = false
	log.Println("[convert]END")
	return nil
}

func runConvertMovies(cInfo *convertInfo, semaphore chan int, notify chan<- int, no int) {
	semaphore <- 0
	log.Println("[runConvertMovies][" + cInfo.filename + "]START")
	// cmdStr := "ffmpeg -i " + filepath.Join(inputDir, filename) + " -vf scale=" + scale + ":-1 " + filepath.Join(outputDir, filename)
	cmdStr := "cp " + cInfo.inputDir + "/" + cInfo.filename + " " + cInfo.outputDir + "/" + cInfo.filename
	cmd := exec.Command(os.Getenv("SHELL"), "-c", cmdStr)
	err := cmd.Run()
	if err != nil {
		log.Println("[runConvertMovies]["+cInfo.filename+"]", err)
		<-semaphore
		return
	}
	log.Printf("ffmpeg fin [%s]\n", cInfo.filename)
	// err = os.Remove(filepath.Join(cInfo.inputDir, cInfo.filename))
	err = os.Remove(cInfo.inputDir + "/" + cInfo.filename)
	if err != nil {
		log.Println("[runConvertMovies]["+cInfo.filename+"]", err)
	}
	log.Println("[runConvertMovies][" + cInfo.filename + "]END")
	<-semaphore
	notify <- no
}
