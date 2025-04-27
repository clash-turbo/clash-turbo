package main

import (
	"fmt"
	"fyne.io/systray"
	"log/slog"
	"os/exec"
	"runtime"
)

func runSystray() {
	slog.Info("启动托盘")
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(ICON)
	systray.SetTitle("mihomo")

	subMenuTop := systray.AddMenuItem("更多", "")
	autoRun := subMenuTop.AddSubMenuItemCheckbox("开机启动", "开机启动", isAutoRun())
	task := subMenuTop.AddSubMenuItem("提升权限", "提升权限")
	reBuild := subMenuTop.AddSubMenuItem("重启核心", "重启核心")

	systray.AddSeparator()
	// 系统代理  tun 关闭
	sysProxy := systray.AddMenuItemCheckbox("系统代理", "系统代理", _appConfig.ProxyType == 1)
	tunProxy := systray.AddMenuItemCheckbox("tun", "系统代理", _appConfig.ProxyType == 2)
	closeProxy := systray.AddMenuItemCheckbox("关闭代理", "系统代理", _appConfig.ProxyType == 0)

	systray.AddSeparator()
	openCoreWeb := systray.AddMenuItem("打开网页监控", "打开网页监控")
	systray.AddSeparator()
	openGuiWeb := systray.AddMenuItem("打开配置", "打开配置")
	openGuiWeb.Hide()
	//systray.AddSeparator()
	mQuit := systray.AddMenuItem("关闭", "关闭")

	if !hasTask(_taskRunName) {
		tunProxy.Hide()
	} else {
		task.Hide()
	}

	for {
		select {
		case <-mQuit.ClickedCh:
			stopMihomo()
			systray.Quit()

		case <-sysProxy.ClickedCh:
			_appConfig.ProxyType = 1
			changeProxyState()
			changeCheck(sysProxy, tunProxy, closeProxy)
		case <-tunProxy.ClickedCh:
			_appConfig.ProxyType = 2
			changeProxyState()
			changeCheck(sysProxy, tunProxy, closeProxy)

		case <-closeProxy.ClickedCh:
			_appConfig.ProxyType = 0
			changeProxyState()
			changeCheck(sysProxy, tunProxy, closeProxy)

		case <-openCoreWeb.ClickedCh:
			err := openBrowser("http://" + _appConfig.Core.ExternalController + "/ui/#/proxies")
			if err != nil {
				fmt.Println("Error opening browser:", err)
			}
		case <-openGuiWeb.ClickedCh:
			err := openBrowser("http://" + _appConfig.Core.ExternalController + "/ui/#/proxies")
			if err != nil {
				fmt.Println("Error opening browser:", err)
			}

		case <-autoRun.ClickedCh:
			if !isAutoRun() {
				setAutoRun(true)
				autoRun.Check()
			} else {
				setAutoRun(false)
				autoRun.Uncheck()
			}
		case <-reBuild.ClickedCh:
			initConfig()
			makeConfig()
			reloadConfig()
			changeProxyState()
			changeCheck(sysProxy, tunProxy, closeProxy)
			fmt.Println("重新生成")
		case <-task.ClickedCh:
			addTask()

			if !hasTask(_taskRunName) {
				tunProxy.Hide()
			} else {
				tunProxy.Show()
			}
		}

	}

}

func changeCheck(sysProxy, tunProxy, closeProxy *systray.MenuItem) {
	switch _appConfig.ProxyType {
	case 1:
		sysProxy.Check()
		if hasTask(_taskRunName) {
			tunProxy.Uncheck()
		}
		closeProxy.Uncheck()
	case 2:
		sysProxy.Uncheck()
		if hasTask(_taskRunName) {
			tunProxy.Check()
		}
		closeProxy.Uncheck()
	default:
		sysProxy.Uncheck()
		if hasTask(_taskRunName) {
			tunProxy.Uncheck()
		}
		closeProxy.Check()
	}
}

func onExit() {
	// clean up here
	sysProxyState(false)

}
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	default:
		return fmt.Errorf("unsupported platform")
	}

	return exec.Command(cmd, args...).Start()
}
