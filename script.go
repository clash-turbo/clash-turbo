package main

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"log"
	"log/slog"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func runMihomo() {

	if !isProgramRun(_mihomoRes) {
		slog.Info("开始启动 核心 ")
		if hasTask(_taskRunName) {
			runTask(_taskRunName)
		} else {
			runCore()
		}
		slog.Info("启动 核心 成功 ")
	}
}
func stopMihomo() {

	if isProgramRun(_mihomoRes) {
		slog.Info("开始关闭 核心 ")
		if hasTask(_taskStopName) {
			runTask(_taskStopName)
		} else {
			stopCore()
		}
		slog.Info("关闭 核心 成功 ")
	}
}

func runCore() {
	runCommand("cscript", []string{filepath.Join(getWorkPath(), "run.vbs")})
}
func stopCore() {
	runCommand("cscript", []string{filepath.Join(getWorkPath(), "run.vbs"), "stop"})
}

func isProgramRun(programName string) bool {

	// 使用 wmic 命令获取指定进程的详细信息
	cmd := exec.Command("cmd.exe", "/c", "tasklist | findstr /I "+programName)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// 将输出转换为字符串并检查是否包含指定的进程名称
	outputStr := string(output)
	return strings.Contains(outputStr, programName)
}

// HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Run
func setAutoRun(i bool) {
	// 打开注册表键
	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	defer key.Close()

	if i {
		err = key.SetStringValue(_appName, "\""+filepath.Join(getExePath(), _appExe)+"\" -path "+getExePath())
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = key.DeleteValue(_appName)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func isAutoRun() bool {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
	if err != nil {
		return false
	}
	defer key.Close()
	_, _, err = key.GetStringValue(_appName)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
