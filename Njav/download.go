package Njav

import (
	"fmt"
	"os"
	"spider/Common"
	"strconv"
	"strings"
	"time"
)

// dir可以是绝顶路径也可以是相对路径,最后一级目录只用带目录名，不用加/
func getTsSlice(url string, tslist []string, dir string, req *Common.HttpRequest) error {
	n := len(tslist)
	url = url[:len(url)-6]

	zero := make([]byte, len(strconv.Itoa(n)))
	for i := 0; i < n; i++ {
		zero[i] = '0'
	}

	var body []byte
	var err error

	req.SetHeader("Priority", "u=1, i")
	req.SetHeader("Origin", "https://javplayer.me")
	req.ReSetClient()

	for n != 0 {
		tmp := []string{}
		for index, ts := range tslist {
			tag := ts[1 : len(ts)-4]
			fmt.Printf("下载中..., 当前片段: %s (%d / %d)\n", ts, index, n-1)
			body, err = req.SendRequest(fmt.Sprintf("%s%s", url, ts))
			if err != nil {
				fmt.Println(err.Error())
				tmp = append(tmp, ts)
			}
			filename := fmt.Sprintf("%s/%s.ts", dir, fmt.Sprintf("%s%s", string(zero[:len(zero)-len(tag)]), tag))
			err = os.WriteFile(filename, body, 0777)
			if err != nil {
				return err
			}
			time.Sleep(500 * time.Millisecond)

		}
		n = len(tmp)
		tslist = tmp
	}

	fmt.Println("ts片段下载成功")
	return nil
}

func removeTsSlice(dir string, lens int) error {
	n := len(strconv.Itoa(lens))
	zero := make([]byte, n)
	for i := 0; i < n; i++ {
		zero[i] = '0'
	}
	for i := 0; i < lens; i++ {
		si := strconv.Itoa(i)
		u := fmt.Sprintf("%s/%s.ts", dir, fmt.Sprintf("%s%s", string(zero[:n-len(si)]), si))
		err := os.Remove(u)
		if err != nil {
			return err
		}
	}
	fmt.Println("ts片段删除成功")
	return nil
}

func Rename(path string) {
	dirInfo, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err.Error())
	}
	n := len(strconv.Itoa(len(dirInfo)))
	zero := make([]byte, n)
	for i := 0; i < n; i++ {
		zero[i] = '0'
	}
	for _, f := range dirInfo {
		name := f.Name()
		sp := strings.Split(name, ".")
		if sp[len(sp)-1] == "ts" {
			m := len(sp[0])
			newName := fmt.Sprintf("%s%s", string(zero[:n-m]), name)
			os.Rename(fmt.Sprintf("%s/%s", path, name), fmt.Sprintf("%s/%s", path, newName))
		}
	}

}

func DownloadVideo(url string, tslist []string, name string, dir string, req *Common.HttpRequest) error {
	n := len(tslist)
	err := getTsSlice(url, tslist, dir, req) //下载ts片段
	if err != nil {
		return err
	}

	//fmt.Scanln()
	fmt.Println(len(tslist))

	err = Common.MergeTS(name, dir) //将ts片段合成一个ts
	if err != nil {
		return err
	}

	//fmt.Scanln()
	err = removeTsSlice(dir, n) //删除ts片段
	if err != nil {
		return err
	}

	//fmt.Scanln()
	err = Common.FfmpegToh264(name, dir) //将ts视频转成h264编码的mp4视频
	if err != nil {
		return err
	}

	//fmt.Scanln()
	err = os.Remove(fmt.Sprintf("%s/%s.ts", dir, name)) //删除ts
	if err != nil {
		return err
	}

	return nil
}
