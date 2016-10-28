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

const maxProcesses = 5

// IsRunning ... 処理中かどうか
var IsRunning bool

// Run ...
func Run(cInfo *ConvertInfo) {
	const fname = "<>[Run]"

	IsRunning = true
	defer func() {
		IsRunning = false
	}()

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
	// rerr := os.Remove(joinPath(cInfo.InputDir, listfile)) // コンバート成功をもって、リストアップしたファイルも削除
	// if rerr != nil {
	// 	log.Println("[after convert]", rerr)
	// 	return
	// }

	// 最後に、コンバートしたファイルをデプロイ
	err = deploy(cInfo)
	if err != nil {
		return
	}
}

func listup(inputDir string) error {
	const fname = "[■ １ ■][listup]"
	// log.Println(fname, "START")

	// リストアップファイルを作成（※作成済み＝後続タスクで未処理のため、以降の処理は行わない）
	isExist, cerr := createFileList(inputDir)
	if cerr != nil {
		log.Println(fname, "createFileList:", cerr)
		return cerr
	}
	if isExist {
		return errors.New("listfile exists")
	}

	// 監視ディレクトリ配下のファイル情報一覧を取得
	fileInfos, rerr := ioutil.ReadDir(inputDir)
	if rerr != nil {
		log.Println(fname, "ioutil.ReadDir(inputDir):", rerr)
		return rerr
	}
	// インプット無いならそこで処理終了
	if len(fileInfos) == 0 {
		log.Println(fname, "input file not exists")
		return errors.New("input file not exists")
	}

	// リストアップファイルにファイル情報（ファイル名、作成日時）を出力
	file, oerr := os.OpenFile(joinPath(inputDir, listfile), os.O_RDWR, 0644)
	if oerr != nil {
		log.Println(fname, "os.Open(inputDir + listfile):", oerr)
		return oerr
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, fileInfo := range fileInfos {
		filename := fileInfo.Name()
		if listfile == filename {
			continue
		}
		filename = strings.Replace(filename, "_", "", -1)
		rerr := os.Rename(filepath.Join(inputDir, fileInfo.Name()), filepath.Join(inputDir, filename))
		if rerr != nil {
			log.Println(fname, rerr)
			return rerr
		}

		_, err := writer.WriteString(filename + "\t" + fileInfo.ModTime().Format(time.ANSIC) + "\r\n")
		if err != nil {
			log.Println(fname, "writer.WriteString:", oerr)
		}
	}
	writer.Flush()

	// log.Println(fname, "END")
	return nil
}

func createFileList(inputDir string) (bool, error) {
	const fname = "[■ １b ■][createFileList]"
	// log.Println(fname, "START")
	_, serr := os.Stat(joinPath(inputDir, listfile))
	if serr == nil {
		// log.Println(fname, "listfile exists")
		return true, nil
	}
	file, err := os.Create(joinPath(inputDir, listfile))
	if err != nil {
		return false, err
	}
	defer file.Close()
	// log.Println(fname, "END")
	return false, nil
}

func convert(cInfo *ConvertInfo) error {
	const fname = "[■ ２ ■][convert]"
	// log.Println(fname, "START")

	sema := make(chan int, maxProcesses)

	notify := make(chan int)

	// リストアップファイルの中身を読み取る
	file, oerr := os.Open(joinPath(cInfo.InputDir, listfile))
	if oerr != nil {
		log.Println(fname, "os.Open(inputDir + listfile):", oerr)
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
		log.Println(fname, "scanner.Err():", serr)
		return serr
	}

	for i := 0; i < no; i++ {
		<-notify // 他のゴルーチンからの終了通知を待つ
	}

	// log.Println(fname, "END")
	return nil
}

// ゴルーチン実行、かつ、呼び元ループ中でファイル名を上書きしているので、ConvertInfoは参照渡しじゃダメ
func runConvertMovies(cInfo ConvertInfo, semaphore chan int, notify chan<- int, no int) {
	fname := "[■ ２b ■][runConvertMovies][" + cInfo.Filename + "]"
	semaphore <- 0
	// log.Println(fname, "START")
	/*
	 * ffmpeg[動画サイズ削減]
	 */
	// log.Printf(fname, "ffmpeg[動画サイズ削減] START")
	cmd := exec.Command(os.Getenv("SHELL"), "-c", cInfo.cmdConvertVideo())
	err := cmd.Run()
	if err != nil {
		log.Println(fname, err)
		<-semaphore
		notify <- no
		return
	}
	// log.Printf(fname, "ffmpeg[動画サイズ削減] END")

	/*
	 * ffmpeg[サムネイル画像抽出]
	 */
	// log.Printf(fname, "ffmpeg[サムネイル画像抽出] START")
	cmd = exec.Command(os.Getenv("SHELL"), "-c", cInfo.cmdCreateThumbnail())
	err = cmd.Run()
	if err != nil {
		log.Println(fname, err)
		<-semaphore
		notify <- no
		return
	}
	// log.Printf(fname, "ffmpeg[サムネイル画像抽出] END")

	/*
	 * convert by ImageMagick [サムネイル画像をローテート]
	 */
	// log.Printf(fname, "ffmpeg[サムネイル画像をローテート] START")
	cmd = exec.Command(os.Getenv("SHELL"), "-c", cInfo.cmdRotateThumbnail())
	err = cmd.Run()
	if err != nil {
		log.Println(fname, err)
		<-semaphore
		notify <- no
		return
	}
	// log.Printf(fname, "ffmpeg[サムネイル画像をローテート] END")

	/*
	 * 処理済みの変換前動画（サイズ大）を削除
	 */
	// log.Printf(fname, "処理済みの変換前動画（サイズ大）を削除 START")
	err = os.Remove(cInfo.inputPath())
	if err != nil {
		log.Println(fname, err)
	}
	// log.Printf(fname, "処理済みの変換前動画（サイズ大）を削除 END")
	<-semaphore
	notify <- no
}

func deploy(cInfo *ConvertInfo) error {
	const fname = "[■ ３ ■][deploy]"
	// log.Println(fname, "START")
	// log.Println(fname, "walk target directory: ", cInfo.OutputDir)

	filepath.Walk(
		cInfo.OutputDir,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			// log.Println(fname, path)
			deployPath := cInfo.deployPath(info.Name())
			if deployPath == "" {
				log.Println(fname, "not target Ext: "+info.Name())
				return nil
			}
			lerr := os.Link(path, deployPath)
			if lerr != nil {
				log.Println(fname, lerr)
			}
			rerr := os.Remove(path)
			if rerr != nil {
				log.Println(fname, rerr)
			}
			return nil
		})
	// log.Println(fname, "END")
	return nil
}
