package main

import (
	"fmt"
	"github.com/Trisia/gosysproxy"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func changeProxyState() {
	switch _appConfig.ProxyType {
	case 1:
		sysProxyState(true)
		tunState(false)
	case 2:
		sysProxyState(false)
		tunState(true)
	default:
		sysProxyState(false)
		tunState(false)
	}
}

func sysProxyState(state bool) {

	if state {
		pass := "localhost;127.*;192.168.*;10.*;172.16.*;172.17.*;172.18.*;172.19.*;172.20.*;172.21.*;172.22.*;172.23.*;172.24.*;172.25.*;172.26.*;172.27.*;172.28.*;172.29.*;172.30.*;172.31.*;<local>"

		err := gosysproxy.SetGlobalProxy("127.0.0.1:"+strconv.Itoa(_appConfig.Core.MixedPort), pass)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		err := gosysproxy.Off()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func tunState(state bool) {

	s := `{"tun":{"enable":false}}`
	if state {
		s = `{"tun":{"enable":true}}`
	}

	client := &http.Client{}
	var data = strings.NewReader(s)
	req, _ := http.NewRequest("PATCH", "http://"+_appConfig.Core.ExternalController+"/configs", data)

	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+_appConfig.Core.Secret)
	_, err := client.Do(req)
	if err != nil {
		return
	}

}

func reloadConfig() {
	client := &http.Client{}
	var data = strings.NewReader(`{"path":"","payload":""}`)
	req, err := http.NewRequest("PUT", "http://"+_appConfig.Core.ExternalController+"/configs", data)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+_appConfig.Core.Secret)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func initProxy() {

	switch _appConfig.ProxyType {
	case 1:
		sysProxyState(true)
	case 2:
		sysProxyState(false)
	case 0:
		sysProxyState(false)

	}
	slog.Info("初始化代理完成")

}
