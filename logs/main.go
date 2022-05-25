package logs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/imaimang/mgo/console"
)

/*
前景
 30    黑色
 31    红色
 32    绿色
 33    黄色
 34    蓝色
 35    紫红色
 36    青蓝色
 37    白色
*/

//Log 日志组件
type Logger struct {
	file *os.File
	lock *sync.Mutex

	fileName string
	name     string

	savePath    string
	maxFileSize int64
	saveTimeout time.Duration
}

//NewLogger
func NewLogger() *Logger {
	logger := new(Logger)
	logger.lock = new(sync.Mutex)
	logger.SetSaveOption("./logs/", 0, 10)
	return logger
}

//savePath 文件保存路径 默认 ./logs/
//saveTimeout 文件保存时间 默认为0无限制
//maxFileSize 单个文件存储大小 单位 MB 默认 10MB
func (l *Logger) SetSaveOption(savePath string, saveTimeout time.Duration, maxFileSize float32) {
	l.maxFileSize = int64(1024 * 1024 * maxFileSize)
	l.savePath = savePath
}

func (l *Logger) refresh() error {
	l.fileName = l.name + "_" + time.Now().Format("20060102_030405") + ".txt"
	path := path.Join(l.savePath, l.fileName)
	err := os.MkdirAll(l.savePath, os.ModePerm)
	if err == nil {
		l.file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			fmt.Println("创建Log文件失败", err)
		}
	} else {
		fmt.Println("创建Log目录失败", err)
	}
	return err
}

func (l *Logger) println(messageType string, content ...interface{}) string {
	funcName := ""
	title := time.Now().Format("2006-01-02 15:04:05") + " -------------- : " + messageType + "  "
	if pc, _, _, ok := runtime.Caller(3); ok {
		if f := runtime.FuncForPC(pc); f != nil {
			file, line := f.FileLine(pc)
			title += path.Base(file) + "  " + strconv.Itoa(line)
			funcName = f.Name()
		}
	}
	console.ColorPrint(messageType, title)
	messageBody := ""
	for i := 0; i < len(content); i++ {
		messageBody = messageBody + fmt.Sprint(content[i])
		if i != len(content)-1 {
			messageBody += ":"
		}
	}
	fmt.Println(messageBody)
	fmt.Println()

	return title + "\r\n" + funcName + "\r\n" + messageBody
}

//WriteDebug WriteDebug
func (l *Logger) WriteDebug(content ...interface{}) {
	l.save("DEBUG", content...)
}

func (l *Logger) save(logType string, content ...interface{}) {
	l.lock.Lock()
	var err error
	message := l.println(logType, content...)
	if l.file == nil {
		err = l.refresh()
	}
	if err == nil {
		_, err = l.file.WriteString(message + "\r\n\r\n")
		if err != nil {
			fmt.Println("Logger写入文件失败", err)
		}
		status, err := l.file.Stat()
		if err == nil {
			if status.Size() >= int64(l.maxFileSize) {
				files, err := ioutil.ReadDir(l.savePath)
				if err == nil {
					for _, f := range files {
						if l.saveTimeout != 0 {
							if time.Now().UTC().Add(-l.saveTimeout).Unix() > f.ModTime().UTC().Unix() && f.Name() != l.fileName {
								os.Remove(path.Join(l.savePath, f.Name()))
							}
						}
						if f.Name() != l.fileName && f.Size() == 0 {
							os.Remove(path.Join(l.savePath, f.Name()))
						}
						if f.Name() == l.fileName && f.Size() > l.maxFileSize {
							l.file.Close()
							l.refresh()
						}
					}
				}
			}
		}
		//l.file.Sync()
		l.lock.Unlock()
	}
}

//WriteInfo WriteInfo
func (l *Logger) WriteInfo(content ...interface{}) {
	l.save("INFO", content...)
}

//WriteWarn WriteWarn
func (l *Logger) WriteWarn(content ...interface{}) {
	l.save("WARN", content...)
}

//WriteError WriteError
func (l *Logger) WriteError(content ...interface{}) {
	l.save("ERROR", content...)
}
