package Pornbest

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func GetImgAndM3u8Url(sUrl []string, name string, code string) (string, string, error) {

	_, err := os.Stat("./Video")
	if os.IsNotExist(err) {
		err := os.Mkdir("./Video", 0777)
		if err != nil {
			return "", "", err
		}
	}
	jsFile := "./Video/" + name + ".js"
	_, err = os.Stat(jsFile)
	if os.IsExist(err) {
		err = os.Remove(jsFile)
		if err != nil {
			return "", "", err
		}
	}

	f, err := os.Create(jsFile)
	if err != nil {
		return "", "", err
	}
	_, err = f.WriteString("console.log" + code[4:])
	if err != nil {
		f.Close()
		return "", "", err
	}
	f.Close()

	cmd := exec.Command("node", fmt.Sprintf("%s.js", name))
	cmd.Dir = "./Video"
	//fmt.Println(cmd.Path, cmd.Args, cmd.Dir, cmd.Env)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err = cmd.Run()

	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if err != nil {
		return "", "", errors.New(errStr)
	}
	//fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)

	pic := regexp.MustCompile(`poster:'(.*?)'`).FindString(outStr)
	picUrl := sUrl[0] + "//" + sUrl[2] + "/" + pic[8:len(pic)-1]
	m3u8 := regexp.MustCompile(`src:'(.*?)'`).FindString(outStr)
	m3u8Url := m3u8[5 : len(m3u8)-1]

	err = os.Remove(jsFile)
	if err != nil {
		return "", "", err
	}
	return picUrl, m3u8Url, nil
}

func UrlList(url string, ts []string) []string {
	ulist := strings.Split(url, "/")
	n := len(ulist)
	path := ""
	for i := 0; i < n-1; i++ {
		path = path + ulist[i] + "/"
	}
	m := len(ts)
	ans := make([]string, m)
	for j := 0; j < m; j++ {
		ans[j] = path + ts[j]
	}
	return ans
}

func Start(url string) {
	re := regexp.MustCompile("https://www.pornbest.org/.+")
	str := re.FindString(url)
	if str == "" {
		fmt.Println("请输入正确的网址!")
		return
	}
	uList := strings.Split(url, "/")
	Mp4Name := uList[len(uList)-1]

	client := GetClient()
	code, err := GetCode(url, client)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(code)
	imgUrl, m3u8Url, err := GetImgAndM3u8Url(uList, Mp4Name, code)

	if err != nil {
		log.Fatal(err)
		return
	}
	tsList, err := M3u8List(m3u8Url, client)
	if err != nil {
		log.Fatal(err)
		return
	}
	tsUrls := UrlList(m3u8Url, tsList)

	err = GetImage(imgUrl, client) //下载封面
	if err != nil {
		log.Fatal(err)
		return
	}
	err = GetMp4Slice(tsUrls, tsList, client) //下载ts片段
	if err != nil {
		log.Fatal(err)
		return
	}
	err = MergeMp4(Mp4Name) //将ts片段合成一个ts
	if err != nil {
		log.Fatal(err)
		return
	}
	err = RemoveMp4Slice(tsList) //删除ts片段
	if err != nil {
		log.Fatal(err)
		return
	}
	err = FfmpegToh264(Mp4Name) //将ts视频转成h264编码的mp4视频
	if err != nil {
		log.Fatal(err)
		return
	}
	err = os.Remove("./Video/" + Mp4Name + ".ts")
	if err != nil {
		log.Fatal(err)
		return
	}
}
