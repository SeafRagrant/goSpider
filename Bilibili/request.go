package Bilibili

import (
	"fmt"
	"os"
	"regexp"
	"spider/Common"
)

var UrlRegexp = regexp.MustCompile("(\"baseUrl\")(.+?)(\"base_url\")")

var Cookies = ""

func GetTitle(content string) string {
	if content == "" {
		return content
	}
	re := regexp.MustCompile("(<title data-vue-meta=\"true\">)(.+?)(</title>)")
	TitleTag := re.FindString(content)
	if TitleTag == "" {
		return TitleTag
	}
	i := len(TitleTag) - 1
	for i >= 0 {
		if TitleTag[i] == '_' {
			break
		}
		i--
	}
	return TitleTag[28:i]
}

func GetImageUrl(content string) string {
	if content == "" {
		return content
	}
	re := regexp.MustCompile("(itemprop=\"image\" content=\")(.+?)(\">)")

	url := re.FindString(content)
	return "https:" + url[26:len(url)-19]
}

func getDash(content string) string {
	if content == "" {
		return content
	}
	re := regexp.MustCompile("(\"dash\")(.*)(\"support_formats\")")
	Dash := re.FindString(content)
	return Dash
}

func getVideoUrl(dash string) string {
	videoRe := regexp.MustCompile("(\"video\")(.+?)(\"audio\")")
	videoCode := videoRe.FindString(dash)
	if videoCode == "" {
		return ""
	}
	ReUrl := UrlRegexp.FindString(videoCode)
	if ReUrl == "" {
		return ""
	}
	return ReUrl[11 : len(ReUrl)-12]
}

func getAudioUrl(dash string) string {
	audioRe := regexp.MustCompile("(\"audio\")(.+)(\"dolby\")")
	audioCode := audioRe.FindString(dash)
	if audioCode == "" {
		return ""
	}

	ReUrl := UrlRegexp.FindString(audioCode)
	if ReUrl == "" {
		return ""
	}
	return ReUrl[11 : len(ReUrl)-12]
}

func GetVideoUrlAndAudioUrl(content string) []string {
	dash := getDash(content)
	if dash == "" {
		return nil
	}
	return []string{getVideoUrl(dash), getAudioUrl(dash)}
}

func Start(url string) {
	header := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
		"Referer":    "https://www.bilibili.com/",
		"Cookie":     Cookies,
	}
	client := Common.GetClient()
	body, h, err := Common.Request(url, client, header) //body为页面

	if err != nil {
		fmt.Println(err)
	} else {
		urls := GetVideoUrlAndAudioUrl(string(body))
		err = DownloadVideoOrAudio("test", urls[1], h)
		if err != nil {
			fmt.Println(err)
		}
	}
	err = DownloadImage("bilibilipic", "https://xxxxxx.jpg", h)
	if err != nil {
		fmt.Println(err)
	} else {
		p, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		} else {
			err = FfmpegAddPicToFlac(p+"\\test.flac", p+"\\bilibilipic.jpg")
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}
