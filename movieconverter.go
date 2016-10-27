package movieconverter

import (
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const listfile = "watch.list"

// Run ...
func Run(cInfo *ConvertInfo, sleep time.Duration) {
	// まず、リストアップ
	err := listup(cInfo.InputDir)
	if err != nil {
		return
	}

	// そして、リストアップしたファイルをコンバート
	err = convert(cInfo)
	if err != nil {
		return
	}
	rerr := os.Remove(joinPath(cInfo.InputDir, listfile))
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
	file, oerr := os.OpenFile(joinPath(inputDir, listfile), os.O_RDWR, 0644)
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
	_, serr := os.Stat(joinPath(inputDir, listfile))
	if serr == nil {
		return errors.New("listfile exists")
	}
	file, err := os.Create(joinPath(inputDir, listfile))
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

var isConverting bool

const maxProcesses = 5

func convert(cInfo *ConvertInfo) error {
	log.Println("[convert]START")
	isConverting = true

	sema := make(chan int, maxProcesses)

	notify := make(chan int)

	// リストアップファイルの中身を読み取る
	file, oerr := os.Open(joinPath(cInfo.InputDir, listfile))
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
		cInfo.Filename = lists[0]
		go runConvertMovies(*cInfo, sema, notify, no)
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

// TODO とりあえず、サイズ変換はした。　あと、サムネイル作成とローテ―ト！
// ゴルーチン実行、かつ、呼び元ループ中でファイル名を上書きしているので、ConvertInfoは参照渡しじゃダメ
func runConvertMovies(cInfo ConvertInfo, semaphore chan int, notify chan<- int, no int) {
	fname := "[runConvertMovies][" + cInfo.Filename + "]"
	semaphore <- 0
	log.Println(fname, "START")
	cmd := exec.Command(os.Getenv("SHELL"), "-c", cInfo.cmdConvertVideo())
	err := cmd.Run()
	if err != nil {
		log.Println(fname, err)
		<-semaphore
		return
	}
	log.Printf("ffmpeg fin [%s]\n", cInfo.Filename)
	err = os.Remove(cInfo.inputPath())
	if err != nil {
		log.Println(fname, err)
	}
	log.Println(fname, "END")
	<-semaphore
	notify <- no
}
