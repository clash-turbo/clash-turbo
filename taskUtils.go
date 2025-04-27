package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func addTask() {

	if hasAdmin() {
		fmt.Println("当前用户有管理员权限，可以继续执行任务")
		if hasTask(_taskRunName) {
			delTask(_taskRunName)
		}
		if hasTask(_taskStopName) {
			delTask(_taskStopName)
		}
		addTaskDo(_taskRunName+".xml", _taskRunName, "")
		addTaskDo(_taskStopName+".xml", _taskStopName, "stop")

	}
}

func addTaskDo(xml string, name string, arg string) {

	taskXmlStr := strings.ReplaceAll(taskXml, "{{directory}}", getWorkPath())
	taskXmlStr = strings.ReplaceAll(taskXmlStr, "{{arg}}", arg)
	taskXmlStr = strings.ReplaceAll(taskXmlStr, "{{name}}", name)

	makeTextFile(taskXmlStr, getWorkPath(), xml)
	xmlPath := filepath.Join(getWorkPath(), xml)
	runCommand("schtasks", []string{"/Create", "/TN", name, "/XML", xmlPath})

	delFile(getWorkPath(), xml)

}

func runTask(taskName string) {
	runCommand("schtasks", []string{"/Run", "/TN", taskName})
}

func delTask(taskName string) {
	runCommand("schtasks", []string{"/Delete", "/TN", taskName, "/F"})
}

func hasTask(taskName string) bool {
	checkCmd := exec.Command("schtasks", "/Query", "/TN", taskName)
	checkCmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	_, err := checkCmd.CombinedOutput()
	return err == nil

}

func hasAdmin() bool {
	checkCmd := exec.Command("net", "session")
	checkCmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	_, err := checkCmd.CombinedOutput()
	return err == nil
}
