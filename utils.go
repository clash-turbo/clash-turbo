package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strconv"
	"syscall"
)

func merge(b map[string]any) {
	a := _appConfig.CoverConfig
	val := reflect.ValueOf(a)
	// 获取a的类型
	typ := reflect.TypeOf(a)

	// 遍历a的字段
	for i := 0; i < val.NumField(); i++ {
		// 获取字段的yaml标签
		fieldType := typ.Field(i)
		yamlTag := fieldType.Tag.Get("yaml") // 获取yaml标签
		if yamlTag == "" {
			// 如果没有yaml标签，可以跳过或直接使用字段名
			yamlTag = fieldType.Name
		}

		// 获取字段值
		fieldValue := val.Field(i).Interface()

		if !isInterfaceEmpty(fieldValue) {

			b[yamlTag] = fieldValue
		}

	}
}

func isInterfaceEmpty(value any) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return v == ""
	case []any:
		return len(v) == 0
	case []string:
		return len(v) == 0
	case []int:
		return len(v) == 0
	case map[string]any:
		return len(v) == 0
	case map[string]string:
		return len(v) == 0
	// 其他类型的判断
	default:
		return false
	}
}

func replaceUnicodeEscapes(input string) string {
	// 定义正则表达式匹配转义的 Unicode 码点
	re := regexp.MustCompile(`\\U([0-9A-Fa-f]{8})`)

	// 替换转义的 Unicode 码点为原始字符
	replaced := re.ReplaceAllStringFunc(input, func(match string) string {
		// 提取 Unicode 码点
		codePoint := match[3:] // 去掉前缀 \U
		// 将 Unicode 码点转换为 rune
		r, _ := strconv.ParseInt(codePoint, 16, 32)
		// 将 rune 转换为字符
		char := string(r)
		return char
	})

	return replaced
}

func runCommand(name string, script []string) {
	delCmd := exec.Command(name, script...)
	delCmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	err := delCmd.Run()
	if err != nil {
		fmt.Printf("出错: %s\n", err, name, script)
		return
	}
}

func checkSingleton() (windows.Handle, error) {
	path, err := os.Executable()
	if err != nil {
		return 0, err
	}
	hashName := md5.Sum([]byte(path))
	name, err := syscall.UTF16PtrFromString("Local\\" + hex.EncodeToString(hashName[:]))
	if err != nil {
		return 0, err
	}
	return windows.CreateMutex(nil, false, name)
}
