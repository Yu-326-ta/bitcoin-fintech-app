package utils

import (
	"io"
	"log"
	"os"
)

func LoggingSetting(logFile string) {
	// 0666は読み取りと書き込みはいずれもOKで、実行だけできない
	logfile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}

	// 出力先を標準出力と作成したログファイルの２つともに書き込みを行う
	multiLogFile := io.MultiWriter(os.Stdout, logfile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(multiLogFile)
}