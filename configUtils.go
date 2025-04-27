package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func initAppConfigYaml() {
	config := appConfig{
		ProxyType:   0,
		GitHubProxy: "",
		HTTPProxy:   "",
		Core: core{
			Port:               10001,
			MixedPort:          10002,
			ExternalController: "127.0.0.1:10003",
			Secret:             "secret",
		},
		GuiPort: 11000,
		Profiles: []profile{
			{
				Name: "示例",
				URL:  "https://test.com",
				Filters: []filter{
					{
						Name:   "香港",
						Filter: "(?i)港|hk|hongkong|hong kong",
						FType:  "select",
					},
					{
						Name:   "日本",
						Filter: "(?i)台|tw|taiwan",
						FType:  "url-test",
					},
					{
						Name:   "新加坡",
						Filter: "(?i)新|sg|singapore",
						FType:  "select",
					},
				},
			},
		},
		PrependRules: nil,
		CoverConfig: coverConfig{
			Dns:          nil,
			Rules:        nil,
			Sniffer:      nil,
			Tun:          nil,
			Ipv6:         false,
			AllowLan:     false,
			UnifiedDelay: true,
		},
	}

	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		fmt.Println("Error writing YAML file:", err)
		return
	}

	remark :=
		"# 参考: https://raw.githubusercontent.com/MetaCubeX/mihomo/refs/heads/Meta/docs/config.yaml \n" +
			"# GitHubProxy: https://ghproxy.cfd/" + "\n" +
			"# proxy_type: 1->系统代理  2->tun 0->关闭" + "\n" +
			"# 自带这三种代理组 🐟 漏网之鱼  🚀 节点选择  DIRECT" + "\n"

	makeTextFile(remark+string(yamlData), getExePath(), "app.yaml")

}

func parseAppConfig() appConfig {

	file, err := os.ReadFile(filepath.Join(getExePath(), "app.yaml"))

	if err != nil {
		fmt.Println("读取文件失败", err)
	}
	var appConfig appConfig
	err = yaml.Unmarshal(file, &appConfig)

	if err != nil {
		fmt.Println("Error unmarshalling YAML app.yaml:", err)
	}
	return appConfig
}

func getInitConfig() map[string]any {

	var initConfig map[string]any
	err := yaml.Unmarshal([]byte(initYaml), &initConfig)
	if err != nil {
		fmt.Println("Error unmarshalling YAML init:", err)
		return nil
	}
	return initConfig

}

func initConfig() {

	exists := exePathFileExists("app.yaml")
	if !exists {
		initAppConfigYaml()
		slog.Info("初始化 app.yaml 成功")
	}

	makeTextFile(strings.ReplaceAll(runVbs, "{{directory}}", getWorkPath()), getWorkPath(), "run.vbs")

	_appConfig = parseAppConfig()

	_geoipUrl = _appConfig.GitHubProxy + "https://github.com/MetaCubeX/meta-rules-dat/releases/download/latest/geoip.dat"
	_geositeUrl = _appConfig.GitHubProxy + "https://github.com/MetaCubeX/meta-rules-dat/releases/download/latest/geosite.dat"
	_dashboardUrl = _appConfig.GitHubProxy + "https://github.com/Zephyruso/zashboard/releases/latest/download/dist.zip"
}
