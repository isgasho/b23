// create by d1y<chenhonzhou@gmail.com>
// 公用方法

package utils

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/d1y/b23/config"
)

// CheckBID 判断是不是 `bilibili` id
func CheckBID(bid string) bool {
	var worA = "av"
	var worB = "bv"
	var isA = strings.Index(bid, worA) > -1
	var isB = strings.Index(bid, worB) > -1
	return (isA || isB)
}

// CheckIDIsAvid 判断是不是 `avid`
func CheckIDIsAvid(id string) bool {
	var worx = "av"
	var index = strings.Index(id, worx)
	return index > -1
}

// ContentLength2MB 将 `ContentLength` 转为 `xxmb` 格式
func ContentLength2MB(n int64) string {
	var mb = n / 1024 / 1024
	f, err := strconv.ParseFloat(string(mb), 2)
	if err != nil {
		return "未知"
	}
	var r = fmt.Sprintf("%vmb", f)
	return r
}

// CheckFileIsExists 判断一个文件是否存在(或者文件夹)
func CheckFileIsExists(cmd string) bool {
	if _, err := os.Stat(cmd); os.IsNotExist(err) {
		return false
	}
	return true
}

// IsCommandAvailable 判断一个命令是否存在
func IsCommandAvailable(name string) bool {
	cmd := exec.Command("command", "-v", name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// GetB23ID 获取 `b23` ID
func GetB23ID(str string) (string, error) {
	URL, e := url.Parse(str)
	if e != nil {
		return "", errors.New("解析失败")
	}
	var host = URL.Host
	var isArrow = false
	for _, item := range config.B23Hosts {
		if item == host {
			isArrow = true
		}
	}
	if !isArrow {
		return "", errors.New("域名错误")
	}
	var b = URL.Path
	var videoWord = "/video/"
	if strings.Contains(b, videoWord) {
		b = b[len(videoWord):]
	}
	return b, nil
}

// IsValidURL 判断是否是个 `url`
func IsValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

// Unzip 解压 `zip`
func Unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}
