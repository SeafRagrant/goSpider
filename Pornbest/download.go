package Pornbest

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
)

func GetImage(url string, client *http.Client) error {
	body, err := Request(url, client)
	if err != nil {
		return err
	}
	re := regexp.MustCompile("(t.jpg).+")
	name := re.FindString(url)[6:] + ".jpg"
	filename := "./Video/" + name
	err = os.WriteFile(filename, body, 0777)
	if err != nil {
		return err
	}
	fmt.Println("图片下载成功")
	return nil
}

func GetMp4Slice(urls []string, filenames []string, client *http.Client) error {
	n := len(urls)
	for index, url := range urls {
		fmt.Printf("下载中...(%d / %d)\n", index+1, n)
		body, err := Request(url, client)
		if err != nil {
			return err
		}
		filename := "./Video/" + filenames[index]
		err = os.WriteFile(filename, body, 0777)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("ts片段下载成功")
	return nil
}

func RemoveMp4Slice(filenames []string) error {
	for _, name := range filenames {
		err := os.Remove("./Video/" + name)
		if err != nil {
			return err
		}
	}
	fmt.Println("ts片段删除成功")
	return nil
}

func MergeMp4(name string) error {
	cmd := exec.Command("cmd", "/C", fmt.Sprintf("copy /b *.ts %s.ts", name))
	cmd.Dir = "./Video"
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误

	err := cmd.Run()
	if err != nil {
		return errors.New(string(stderr.Bytes()))
	}
	//outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	//fmt.Printf("out:\n%s\n err:\n%s\n", outStr, errStr)
	fmt.Println("ts片段合并成功")
	return nil
}

func FfmpegToh264(name string) error {
	//cmd := exec.Command("ffmpeg", fmt.Sprintf("-i ./%s.ts -c:v libx264 -crf 18 ./%s.mp4", name, name))
	cmd := exec.Command("ffmpeg", "-i", "./"+name+".ts", "-c:v", "libx264", "-crf", "18", "./"+name+".mp4")
	cmd.Dir = "./Video"
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误

	fmt.Println("正在转码...")
	err := cmd.Run()
	if err != nil {
		return errors.New(string(stderr.Bytes()))
	}

	//outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	//fmt.Printf("out:\n%s\n err:\n%s\n", outStr, errStr)
	fmt.Println("ts转码成功")
	return nil
}
