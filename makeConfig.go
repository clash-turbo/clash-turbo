package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const baseSelect = "🚀 节点选择"

func makeConfig() {

	slog.Info("生成 config.yaml ")
	initConfig := getInitConfig()
	merge(initConfig)
	groupProxies := makeProfile(initConfig)
	mergeAppConfig(initConfig)
	makeCustomProxyGroup(initConfig, groupProxies)
	writeConfig(initConfig)

}

func makeProfile(config map[string]any) []string {

	profiles := _appConfig.Profiles
	proxyProviders := make(map[string]any)
	var proxyGroups []map[string]any

	var firstGroupProxies []string

	for _, profile := range profiles {
		proxyProviderMap := make(map[string]any)
		proxyProviderMap["name"] = profile.Name
		proxyProviderMap["type"] = "http"
		proxyProviderMap["url"] = _appConfig.HTTPProxy + profile.URL
		proxyProviderMap["interval"] = 86400
		proxyProviderMap["path"] = "./profiles/" + profile.Name + ".yaml"
		proxyProviderMap["override"] = map[string]string{
			"additional-prefix": profile.Name,
		}

		healthCheckMap := make(map[string]any)
		healthCheckMap["enable"] = true
		healthCheckMap["interval"] = 3600
		healthCheckMap["lazy"] = true
		healthCheckMap["url"] = "http://cp.cloudflare.com/generate_204"
		proxyProviderMap["health-check"] = healthCheckMap
		proxyProviders[profile.Name] = proxyProviderMap

		name := "✈️ [" + profile.Name + "]"
		filters := profile.Filters
		proxyGroups = append(proxyGroups, proxyGroup(name, profile.Name, "", "url-test"))
		firstGroupProxies = append(firstGroupProxies, name)

		for _, filter := range filters {
			subName := "✈️ [" + profile.Name + "]" + filter.Name
			proxyGroups = append(proxyGroups, proxyGroup(subName, profile.Name, filter.Filter, filter.FType))
			firstGroupProxies = append(firstGroupProxies, subName)
		}

	}

	firstGroupProxies = append(firstGroupProxies, "全部", "DIRECT")

	// 第二个 漏网之鱼
	matchGroup := make(map[string]any)
	matchGroup["name"] = "🐟 漏网之鱼"
	matchGroup["type"] = "select"
	matchGroup["proxies"] = append([]string{baseSelect}, firstGroupProxies...)
	proxyGroups = append([]map[string]any{matchGroup}, proxyGroups...)

	// 第一个 全部
	firstGroup := make(map[string]any)
	firstGroup["name"] = baseSelect
	firstGroup["type"] = "select"
	firstGroup["proxies"] = firstGroupProxies
	proxyGroups = append([]map[string]any{firstGroup}, proxyGroups...)

	allGroupMap := make(map[string]any)
	allGroupMap["name"] = "全部"
	allGroupMap["type"] = "select"
	allGroupMap["include-all"] = true

	proxyGroups = append(proxyGroups, allGroupMap)

	config["proxy-groups"] = proxyGroups
	config["proxy-providers"] = proxyProviders

	return firstGroupProxies
}

func makeCustomProxyGroup(config map[string]any, groupProxies []string) {

	// 切换成读取prepend_rules
	prependRules := _appConfig.PrependRules

	for _, i := range prependRules {

		initGroup := []string{"🐟 漏网之鱼", "🚀 节点选择", "DIRECT"}
		s := strings.Split(i, ",")
		if len(s) != 3 || slices.Contains(initGroup, s[2]) {
			continue
		}

		proxyGroupMap := make(map[string]any)

		proxyGroupMap["type"] = "select"
		proxyGroupMap["name"] = s[2]
		proxyGroupMap["proxies"] = append(groupProxies, baseSelect)
		//config["proxy-groups"] = append(config["proxy-groups"].([]map[string]any), proxyGroupMap)
		config["proxy-groups"] = append(config["proxy-groups"].([]map[string]any)[:2], append([]map[string]any{proxyGroupMap}, config["proxy-groups"].([]map[string]any)[2:]...)...)

	}

}

func proxyGroup(name, use, filter, fType string) map[string]any {
	proxyGroupMap := make(map[string]any)

	if len(filter) > 0 {
		proxyGroupMap["filter"] = filter
	}
	proxyGroupMap["type"] = fType
	proxyGroupMap["use"] = []string{use}
	proxyGroupMap["name"] = name
	return proxyGroupMap

}

func mergeAppConfig(config map[string]any) {
	config["external-controller"] = _appConfig.Core.ExternalController
	config["mixed-port"] = _appConfig.Core.MixedPort
	config["port"] = _appConfig.Core.Port
	config["secret"] = _appConfig.Core.Secret

	proxyType := _appConfig.ProxyType
	tunConfig := config["tun"].(map[string]any)
	if proxyType == 2 {
		tunConfig["enable"] = true
	}
	if _appConfig.GitHubProxy != "" {
		config["external-ui-url"] = _dashboardUrl
		geoxUrl := config["geox-url"].(map[string]any)
		geoxUrl["geoip"] = _geoipUrl
		geoxUrl["geosite"] = _geositeUrl
	}

	prependRules := _appConfig.PrependRules

	oldRules := config["rules"].([]any)
	oldRulesStr := make([]string, len(oldRules))
	for i, v := range oldRules {
		oldRulesStr[i] = v.(string)
	}
	config["rules"] = append(prependRules, oldRulesStr...)

}

func writeConfig(config map[string]any) {

	data, err := yaml.Marshal(config)
	if err != nil {
		fmt.Println("Error marshalling YAML:", err)
		return
	}

	yamlStr := replaceUnicodeEscapes(string(data))
	err = os.WriteFile(filepath.Join(getWorkPath(), "config.yaml"), []byte(yamlStr), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}
