package main

import (
	"fmt"
	"golang.org/x/sys/windows"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

var (
	Version string = "dev"
	Env     string
)

func init() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	}))

	if Env == "prod" {
		logger = slog.New(slog.NewTextHandler(&lumberjack.Logger{
			Filename:   filepath.Join(getWorkPath(), "app.log"),
			MaxSize:    10, // megabytes
			MaxBackups: 3,
			MaxAge:     2,    //days
			Compress:   true, // disabled by default
		}, &slog.HandlerOptions{
			AddSource:   false,
			Level:       slog.LevelInfo,
			ReplaceAttr: nil,
		}))
	}

	slog.SetDefault(logger)

}

/*
*
1. 初始化文件
2. 覆盖
3. 添加订阅
*/
func main() {

	// 防止重复运行
	singleton, err := checkSingleton()
	if err != nil {
		slog.Info("无法获取锁", err)
		log.Fatalf("无法获取锁: %v", err)
	}
	defer windows.CloseHandle(singleton)

	// 获取exe 文件名
	exePath, _ := os.Executable()
	_appExe = filepath.Base(exePath)

	args := os.Args
	if len(args) > 1 {
		switch args[1] {
		case "-path":
			_exePath = args[2]
		case "-v":
			fmt.Println(Version)
			return
		}
	}

	initConfig()

	// 生成配置
	makeConfig()

	initProxy()
	runMihomo()
	runSystray()

}
