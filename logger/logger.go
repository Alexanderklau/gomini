package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

//定义日志级别
const (
	LevelError = iota
	LevelWarning
	LevelInfo
	LevelDebug
)

//定义日志结构体
type logging struct {
	level int
	file  *os.File
}

//
var logFile logging

//等级设置
func setLevel(level int) {
	logFile.level = level
}

//时间转字符串
func formatTime(t *time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

//当前文件
func getwd() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", errors.New(fmt.Sprintln("Getwd Error:", err))
	}
	filepathstr := filepath.Join(path, "logs")
	err = os.MkdirAll(filepathstr, 0755)
	if err != nil {
		return "", errors.New(fmt.Sprintln("Create Dirs Error:", err))
	}
	return filepathstr, nil
}

var loggerf *log.Logger

//创建日志文件
//制定创建路径和日志名
func Init(path, logname string, level int) error {
	if path == "" {
		filePath, _ := getwd()
		filepathstr := filepath.Join(filePath, logname)
		logfile, err := os.Create(filepathstr)
		if err != nil {
			return errors.New(fmt.Sprintln("Create File Error:", err))
		}
		logFile.file = logfile
		setLevel(level)
	} else if logname == "" {
		filepathstr := path
		err := os.MkdirAll(filepathstr, 0644)
		if err != nil {
			return errors.New(fmt.Sprintln("Create Dirs Error:", err))
		}
		logname = filepath.Join(filepathstr, "logs.log")
		logfile, err := os.Create(logname)
		if err != nil {
			return errors.New(fmt.Sprintln("Create File Error:", err))
		}
		logFile.file = logfile
		setLevel(level)
	} else {
		filepathstr := path
		err := os.MkdirAll(filepathstr, 0644)
		if err != nil {
			return errors.New(fmt.Sprintln("Create Dirs Error:", err))
		}
		logname = filepath.Join(filepathstr, logname)
		logfile, err := os.Create(logname)
		if err != nil {
			return errors.New(fmt.Sprintln("Create File Error:", err))
		}
		logFile.file = logfile
		setLevel(level)
	}
	//log.Llongfile|log.Ltime|log.Ldata 返回文件信息，时间和代码信息不需要重新写函数获取
	//初始化一个*log.Logger
	loggerf = log.New(logFile.file, "", log.Llongfile|log.Ltime|log.Ldate)
	return nil
}

//标准输出
func init() {
	logFile.level = LevelDebug
	loggerf = log.New(os.Stdout, "", log.Llongfile|log.Ltime|log.Ldate)
}

//获取当前函数所在的路径，文件名和行号
// func CallerFunc() string {
// 	filename, file, line, ok := runtime.Caller(0)
// 	if !ok {
// 		return ""
// 	}
// 	return "FileName:" + runtime.FuncForPC(filename).Name() + " Path:" + file + " Line:" + strconv.Itoa(line)
// }

//按照格式打印日志
func Debugf(format string, v ...interface{}) {
	if logFile.level >= LevelDebug {
		//log里的函数都自带有mutex
		//这里获取一次锁
		loggerf.Printf("[DEBUG] "+format, v...)
	}
}

//Infof
func Infof(format string, v ...interface{}) {
	if logFile.level >= LevelInfo {
		loggerf.Printf("[INFO] "+format, v...)
	}
}

//Warningf
func Warningf(format string, v ...interface{}) {
	if logFile.level >= LevelWarning {
		loggerf.Printf("[WARNING] "+format, v...)
	}
}

//Errorf
func Errorf(format string, v ...interface{}) {
	if logFile.level >= LevelError {
		loggerf.Printf("[ERROR] "+format, v...)
	}
}

//标准输出
//Println存在输出会换行
//这里采用Print输出方式
func Debug(v ...interface{}) {
	if logFile.level >= LevelDebug {
		loggerf.Print("[DEBUG] " + fmt.Sprintln(v...))
	}
}

//Info
func Info(v ...interface{}) {
	if logFile.level >= LevelInfo {
		loggerf.Print("[INFO] " + fmt.Sprintln(v...))
	}
}

//Warning
func Warning(v ...interface{}) {
	if logFile.level >= LevelWarning {
		loggerf.Print("[DEBUG] " + fmt.Sprintln(v...))
	}
}

//Error
func Error(v ...interface{}) {
	if logFile.level >= LevelError {
		loggerf.Print("[DEBUG] " + fmt.Sprintln(v...))
	}
}
