package foundation

import (
	"github.com/sirupsen/logrus"
	"runtime"
	"os"
	"time"
	"fmt"
	"strings"
	"github.com/mholt/archiver"
)

type Logger struct {
	Subject      string
	Directory    string
	MetaLogger   *logrus.Logger
	fields       map[string]interface{}
	StdoutLogger *logrus.Logger
}

var dirs = map[string]bool{}

// 创建日志记录器
// subject string: 日志主题
// dir string: 日志存放目录，以/结尾
func NewLogger(subject string, dir string) *Logger {
	stdoutLogger := logrus.New()
	stdoutLogger.Out = os.Stdout
	stdoutLogger.Formatter = &logrus.TextFormatter{ForceColors: true, FullTimestamp: true}

	metaLogger := logrus.New()
	metaLogger.Formatter = &logrus.JSONFormatter{}

	logger := &Logger{subject, dir, metaLogger, map[string]interface{}{}, stdoutLogger}
	logger.executeTimingTasks(func() {
		// 切换日志文件输出点
		logger.setMetaLoggerOut()
	})
	if !dirs[dir] {
		// 当前目录下共有定时任务
		logger.executeTimingTasks(func() {
			// 历史日志整理
			logger.settleLogs()
			// 删除过期日志压缩包
			logger.removeOutdatedLogs()
		})
		dirs[dir] = true
	}

	return logger
}

// 执行定时任务，每日0点触发
func (l *Logger) executeTimingTasks(block func()) {
	block()
	go func() {
		duration := 24 * time.Hour
		//duration := 5 * time.Second
		now := time.Now()
		next := now.Add(duration)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		timer := time.NewTimer(next.Sub(now))
		<-timer.C
		block()
		c := time.Tick(duration)
		for {
			<-c
			block()
		}
	}()
}

// 获取当前时间节点下日志文件名称
func (l *Logger) logname() (logname string) {
	logname = fmt.Sprintf("%v%v_%v.logs", l.Directory, l.Subject, time.Now().Format("20060102"))
	//logname = fmt.Sprintf("%v%v_%v.logs", l.Subject, time.Now().Format("20060102150405"))
	return
}

// 设置日志文件输出点
func (l *Logger) setMetaLoggerOut() {
	logname := l.logname()
	fmt.Println(logname)
	file, err := os.OpenFile(logname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		l.MetaLogger.Fatal("Failed to logs to file, using default stderr")
		return
	}
	l.MetaLogger.Out = file
}

// 日志整理
func (l *Logger) settleLogs() {
	// 检查history子目录是否存在，不存在则创建
	historyDir := l.Directory + "history/"
	_, err := os.Open(historyDir)
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(historyDir, 0777)
		} else {
			l.Errorf("settleLogs", "fail to detect history directory, due to [%v]", err)
		}
	}
	// 进入日志目录，遍历各文件
	dir, err := os.Open(l.Directory)
	if err != nil {
		l.Errorf("settleLogs", "fail to open directory: %v, due to [%v]", l.Directory, err)
		return
	}
	files, err := dir.Readdir(0)
	if err != nil {
		l.Errorf("settleLogs", "fail to read files in directory: %v, due to [%v]", l.Directory, err)
		return
	}
	// 筛选出待打包压缩的日志组，按日期区分
	logs := map[string][]string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		logname := file.Name()
		// 筛除隐藏文件夹
		if logname[0:1] == "." {
			continue
		}
		// 从日志文件名中提取时间戳
		timestamp := strings.Split(strings.SplitAfter(logname, "_")[1], ".logs")[0]
		t, err := time.Parse("20060102", timestamp)
		if err != nil {
			l.Errorf("settleLogs", "fail to parse timestamp in %v/%v, due to [%v]", l.Directory, logname, err)
			continue
		}
		// 对一周前的日志文件进行筛选，同日期的日志归为一组
		if time.Now().Sub(t) >= 24*time.Hour*7 {
			paths := logs[timestamp]
			paths = append(paths, l.Directory+logname)
			logs[timestamp] = paths
		}
	}
	// 将各组日志打包压缩放入history子目录
	for date, paths := range logs {
		targzPath := fmt.Sprintf("%v%v.tar.gz", historyDir, date)
		err = archiver.TarGz.Make(targzPath, paths)
		if err != nil {
			l.Errorf("settleLogs", "fail to tar and gzip [%v]:[%v], due to [%v]", date, paths, err)
			continue
		}
		// 删除原始日志
		for _, p := range paths {
			err = os.Remove(p)
			if err != nil {
				l.Errorf("settleLogs", "fail to remove %v, due to [%v]", p, err)
			}
			continue
		}
	}
}

// 删除过期日志
func (l *Logger) removeOutdatedLogs() {
	historyDir := l.Directory + "history/"
	// 进入历史日志目录，遍历各文件
	dir, err := os.Open(historyDir)
	if err != nil {
		l.Errorf("removeOutdatedLogs", "fail to open directory: %v, due to [%v]", historyDir, err)
		return
	}
	files, err := dir.Readdir(0)
	if err != nil {
		l.Errorf("removeOutdatedLogs", "fail to read files in directory: %v, due to [%v]", historyDir, err)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		// 筛除隐藏文件夹
		if filename[0:1] == "." {
			continue
		}
		// 从压缩包文件名中提取时间戳
		timestamp := strings.Split(filename, ".tar.gz")[0]
		t, err := time.Parse("20060102", timestamp)
		if err != nil {
			l.Errorf("removeOutdatedLogs", "fail to parse timestamp in %v%v, due to [%v]", historyDir, filename, err)
			continue
		}
		// 删除6个月前的压缩包文件
		deadLine := time.Now().AddDate(0, -6, 0)
		if deadLine.Sub(t) >= 0 {
			err = os.Remove(historyDir + filename)
			if err != nil {
				l.Errorf("removeOutdatedLogs", "fail to remove file in %v%v, due to [%v]", historyDir, filename, err)
				continue
			}
		}
	}
}

type LogFields map[string]interface{}

type crashLocation struct {
	File     string
	Function string
	Line     int
}

func (l *Logger) Debug(anchor string, args ...interface{}) {
	l.MetaLogger.WithField("anchor", anchor).Debug(args...)
	l.StdoutLogger.WithField("anchor", anchor).Debug(args...)
}

func (l *Logger) Info(anchor string, args ...interface{}) {
	l.MetaLogger.WithField("anchor", anchor).Info(args...)
	l.StdoutLogger.WithField("anchor", anchor).Info(args...)
}

func (l *Logger) Warn(anchor string, args ...interface{}) {
	l.MetaLogger.WithField("anchor", anchor).Warn(args...)
	l.StdoutLogger.WithField("anchor", anchor).Warn(args...)
}

func (l *Logger) Error(anchor string, args ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc)
	location := crashLocation{file, function.Name(), line}
	l.MetaLogger.WithField("anchor", anchor).WithField("location", location).Error(args...)
	l.StdoutLogger.WithField("anchor", anchor).WithField("location", location).Error(args...)
}

func (l *Logger) Fatal(anchor string, args ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc)
	location := crashLocation{file, function.Name(), line}
	l.MetaLogger.WithField("anchor", anchor).WithField("location", location).Fatal(args...)
	l.StdoutLogger.WithField("anchor", anchor).WithField("location", location).Fatal(args...)
}

func (l *Logger) Panic(anchor string, args ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc)
	location := crashLocation{file, function.Name(), line}
	l.MetaLogger.WithField("anchor", anchor).WithField("location", location).Panic(args...)
	l.StdoutLogger.WithField("anchor", anchor).WithField("location", location).Panic(args...)
}

func (l *Logger) Debugf(anchor string, format string, args ...interface{}) {
	l.MetaLogger.WithField("anchor", anchor).Debugf(format, args...)
	l.StdoutLogger.WithField("anchor", anchor).Debugf(format, args...)
}

func (l *Logger) Infof(anchor string, format string, args ...interface{}) {
	l.MetaLogger.WithField("anchor", anchor).Infof(format, args...)
	l.StdoutLogger.WithField("anchor", anchor).Infof(format, args...)
}
func (l *Logger) Warnf(anchor string, format string, args ...interface{}) {
	l.MetaLogger.WithField("anchor", anchor).Warnf(format, args...)
	l.StdoutLogger.WithField("anchor", anchor).Warnf(format, args...)
}

func (l *Logger) Errorf(anchor string, format string, args ...interface{}) {
	//pc, file, line, _ := runtime.Caller(1)
	//function := runtime.FuncForPC(pc)
	//location := crashLocation{file, function.Name(), line}
	l.MetaLogger.WithField("anchor", anchor).Errorf(format, args...)
	l.StdoutLogger.WithField("anchor", anchor).Errorf(format, args...)
}

func (l *Logger) Fatalf(anchor string, format string, args ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc)
	location := crashLocation{file, function.Name(), line}
	l.MetaLogger.WithField("anchor", anchor).WithField("location", location).Fatalf(format, args...)
	l.StdoutLogger.WithField("anchor", anchor).WithField("location", location).Fatalf(format, args...)
}

func (l *Logger) Panicf(anchor string, format string, args ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc)
	location := crashLocation{file, function.Name(), line}
	l.MetaLogger.WithField("anchor", anchor).WithField("location", location).Panicf(format, args...)
	l.StdoutLogger.WithField("anchor", anchor).WithField("location", location).Panicf(format, args...)
}

type Entry struct {
	metaEntry   *logrus.Entry
	stdoutEntry *logrus.Entry
}

func (l *Logger) WithFields(fields LogFields) *Entry {
	metaEntry := l.MetaLogger.WithFields(logrus.Fields(fields))
	stdoutEntry := l.StdoutLogger.WithFields(logrus.Fields(fields))
	return &Entry{metaEntry, stdoutEntry}
}

func (l *Logger) WithField(key string, value interface{}) *Entry {
	metaEntry := l.MetaLogger.WithField(key, value)
	stdoutEntry := l.StdoutLogger.WithField(key, value)
	return &Entry{metaEntry, stdoutEntry}
}

func (e *Entry) WithFields(fields LogFields) *Entry {
	return &Entry{
		e.metaEntry.WithFields(logrus.Fields(fields)),
		e.stdoutEntry.WithFields(logrus.Fields(fields)),
	}
}

func (e *Entry) WithField(key string, value interface{}) *Entry {
	return &Entry{
		e.metaEntry.WithField(key, value),
		e.stdoutEntry.WithField(key, value),
	}
}

func (e *Entry) Debug(anchor string, args ...interface{}) {
	e.metaEntry.WithField("anchor", anchor).Debug(args...)
	e.stdoutEntry.WithField("anchor", anchor).Debug(args...)
}

func (e *Entry) Info(anchor string, args ...interface{}) {
	e.metaEntry.WithField("anchor", anchor).Info(args...)
	e.stdoutEntry.WithField("anchor", anchor).Info(args...)
}

func (e *Entry) Warn(anchor string, args ...interface{}) {
	e.metaEntry.WithField("anchor", anchor).Warn(args...)
	e.stdoutEntry.WithField("anchor", anchor).Warn(args...)
}


func (e *Entry) Error(anchor string, args ...interface{}) {
	//pc, file, line, _ := runtime.Caller(1)
	//function := runtime.FuncForPC(pc)
	//location := crashLocation{file, function.Name(), line}
	//e.metaEntry.WithField("anchor", anchor).WithField("location", location).Error(args...)
	//e.stdoutEntry.WithField("anchor", anchor).WithField("location", location).Error(args...)
	e.metaEntry.WithField("anchor", anchor).Error(args...)
	e.stdoutEntry.WithField("anchor", anchor).Error(args...)
}

func (e *Entry) Fatal(anchor string, args ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc)
	location := crashLocation{file, function.Name(), line}
	e.metaEntry.WithField("anchor", anchor).WithField("location", location).Fatal(args...)
	e.stdoutEntry.WithField("anchor", anchor).WithField("location", location).Fatal(args...)
}

func (e *Entry) Panic(anchor string, args ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc)
	location := crashLocation{file, function.Name(), line}
	e.metaEntry.WithField("anchor", anchor).WithField("location", location).Panic(args...)
	e.stdoutEntry.WithField("anchor", anchor).WithField("location", location).Panic(args...)
}

func (e *Entry) Debugf(anchor string, format string, args ...interface{}) {
	e.metaEntry.WithField("anchor", anchor).Debugf(format, args...)
	e.stdoutEntry.WithField("anchor", anchor).Debugf(format, args...)
}

func (e *Entry) Infof(anchor string, format string, args ...interface{}) {
	e.metaEntry.WithField("anchor", anchor).Infof(format, args...)
	e.stdoutEntry.WithField("anchor", anchor).Infof(format, args...)
}
func (e *Entry) Warnf(anchor string, format string, args ...interface{}) {
	e.metaEntry.WithField("anchor", anchor).Warnf(format, args...)
	e.stdoutEntry.WithField("anchor", anchor).Warnf(format, args...)
}

func (e *Entry) Errorf(anchor string, format string, args ...interface{}) {
	//pc, file, line, _ := runtime.Caller(1)
	//function := runtime.FuncForPC(pc)
	//location := crashLocation{file, function.Name(), line}
	//e.metaEntry.WithField("anchor", anchor).WithField("location", location).Errorf(format, args...)
	//e.stdoutEntry.WithField("anchor", anchor).WithField("location", location).Errorf(format, args...)
	e.metaEntry.WithField("anchor", anchor).Errorf(format, args...)
	e.stdoutEntry.WithField("anchor", anchor).Errorf(format, args...)
}

func (e *Entry) Fatalf(anchor string, format string, args ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc)
	location := crashLocation{file, function.Name(), line}
	e.metaEntry.WithField("anchor", anchor).WithField("location", location).Fatalf(format, args...)
	e.stdoutEntry.WithField("anchor", anchor).WithField("location", location).Fatalf(format, args...)
}

func (e *Entry) Panicf(anchor string, format string, args ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc)
	location := crashLocation{file, function.Name(), line}
	e.metaEntry.WithField("anchor", anchor).WithField("location", location).Panicf(format, args...)
	e.stdoutEntry.WithField("anchor", anchor).WithField("location", location).Panicf(format, args...)
}
