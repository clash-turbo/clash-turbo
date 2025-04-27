package main

import _ "embed"

// 全局变量
var (
	_appExe       = ""
	_appConfig    appConfig
	_exePath      = ""
	_geoipUrl     = "https://github.com/MetaCubeX/meta-rules-dat/releases/download/latest/geoip.dat"
	_geositeUrl   = "https://github.com/MetaCubeX/meta-rules-dat/releases/download/latest/geosite.dat"
	_dashboardUrl = "https://github.com/Zephyruso/zashboard/releases/latest/download/dist.zip"
)

const (
	_appName      = "ct"
	_workPath     = "work/"
	_taskRunName  = "ct_run"
	_taskStopName = "ct_stop"
	_mihomoRes    = "mihomo-windows-amd64.exe"
)

//go:embed file/init.yaml
var initYaml string

//go:embed file/run.vbs
var runVbs string

//go:embed file/task.txt
var taskXml string
