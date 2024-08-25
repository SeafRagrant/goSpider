package Pornbest

import (
	"fmt"
	"os"
	"regexp"
	"spider/Common"
	"strconv"
)

func GetImage(url string, req *Common.HttpRequest) error {
	body, err := req.SendRequest(url)
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

// dir可以是绝顶路径也可以是相对路径,最后一级目录只用带目录名，不用加/
func getTsSlice(urls []string, dir string, req *Common.HttpRequest) error {
	n := len(urls)
	for index, url := range urls {
		fmt.Printf("下载中...(%d / %d)\n", index+1, n)
		body, err := req.SendRequest(url)
		if err != nil {
			return err
		}
		filename := fmt.Sprintf("%s/%s.ts", dir, strconv.Itoa(index))
		err = os.WriteFile(filename, body, 0777)
		if err != nil {
			return err
		}
	}
	fmt.Println("ts片段下载成功")
	return nil
}

func removeTsSlice(dir string, n int) error {
	for i := 0; i < n; i++ {
		u := fmt.Sprintf("%s/%s.ts", dir, strconv.Itoa(i))
		err := os.Remove(u)
		if err != nil {
			return err
		}
	}
	fmt.Println("ts片段删除成功")
	return nil
}

func DownloadVideo(urls []string, name string, dir string, req *Common.HttpRequest) error {
	err := getTsSlice(urls, dir, req) //下载ts片段
	if err != nil {
		return err
	}

	err = Common.MergeTS(name, dir) //将ts片段合成一个ts
	if err != nil {
		return err
	}

	err = removeTsSlice(dir, len(urls)) //删除ts片段
	if err != nil {
		return err
	}

	err = Common.FfmpegToh264(name, dir) //将ts视频转成h264编码的mp4视频
	if err != nil {
		return err
	}

	err = os.Remove(fmt.Sprintf("%s/%s.ts", dir, name)) //删除ts
	if err != nil {
		return err
	}

	return nil
}
