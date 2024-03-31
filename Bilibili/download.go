package Bilibili

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"spider/Common"
	"strings"
)

func DownloadImage(name, url string, req *Common.HttpRequest) error {
	res, err := req.SendRequest(url)
	if err != nil {
		return err
	}
	ulist := strings.Split(url, "/")
	temp := ulist[len(ulist)-1]
	if name == "" {
		name = temp
	} else {
		ext := strings.Split(temp, ".")
		name = name + "." + ext[len(ext)-1]
	}
	err = os.WriteFile(name, res, 0777)
	if err != nil {
		return err
	}
	return nil
}

func DownloadVideoOrAudio(name, url string, req *Common.HttpRequest) error {
	body, err := req.SendRequest(url)
	if err != nil {
		return err
	}
	name = name + ".m4s"
	err = os.WriteFile(name, body, 0777)
	if err != nil {
		return err
	}
	return nil
}

//ffprobe -show_format input 	查看媒体文件信息

func FfmpegM4sToAac(audio string) error {
	//原始音频
	AudioName := path.Base(audio)
	NameAndExt := strings.Split(AudioName, ".")
	//ffmpeg -i input -vn -c:a copy output
	cmd := exec.Command("ffmpeg", "-i", "./"+AudioName, "-vn", "-c:a", "copy", "./"+NameAndExt[0]+".aac")
	cmd.Dir = path.Dir(audio)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误

	err := cmd.Run()
	if err != nil {
		return errors.New(string(stderr.Bytes()))
	}

	//outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	//fmt.Printf("out:\n%s\n err:\n%s\n", outStr, errStr)
	fmt.Println("转码成功")
	return nil
}

func FfmpegAacToFlac(audio string) error {
	AudioName := path.Base(audio)
	NameAndExt := strings.Split(AudioName, ".")
	//ffmpeg -i input -b:a 192k copy output
	cmd := exec.Command("ffmpeg", "-i", "./"+AudioName, "./"+NameAndExt[0]+".flac")
	cmd.Dir = path.Dir(audio)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误

	err := cmd.Run()
	if err != nil {
		return errors.New(string(stderr.Bytes()))
	}

	//outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	//fmt.Printf("out:\n%s\n err:\n%s\n", outStr, errStr)
	fmt.Println("转码成功")
	return nil
}

func FfmpegAddPicToFlac(audio, picture string) error {
	output := strings.Split(audio, ".")[0] + "withPic.flac"

	//ffmpeg -i audio.flac -i image.png -map 0:a -map 1 -codec copy -metadata:s:v title="Album cover" -metadata:s:v comment="Cover (front)" -disposition:v attached_pic output.flac
	cmd := exec.Command("ffmpeg", "-i", audio, "-i", picture, "-map", "0:a", "-map", "1", "-codec", "copy", "-disposition:v", "attached_pic", output)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误

	err := cmd.Run()
	if err != nil {
		return errors.New(string(stderr.Bytes()))
	}

	//outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	//fmt.Printf("out:\n%s\n err:\n%s\n", outStr, errStr)
	fmt.Println("添加成功")
	return nil
}

func MergeVideoAndAudio(video, audio, output string) error {
	videoName := path.Base(video)
	audioName := path.Base(audio)
	//ffmpeg -i video.m4s -i audio.m4s -codec copy output.mp4
	cmd := exec.Command("ffmpeg", "-i", "./"+videoName, "-i", "./"+audioName, "-codec", "copy", "./"+output+".mp4")
	cmd.Dir = path.Dir(video)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误

	fmt.Println("正在合成")
	err := cmd.Run()
	if err != nil {
		return errors.New(string(stderr.Bytes()))
	}

	//outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	//fmt.Printf("out:\n%s\n err:\n%s\n", outStr, errStr)
	fmt.Println("合成成功")
	return nil
}
