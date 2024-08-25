package Njav

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"spider/Common"
	"strings"
)

type VideosRes struct {
	Data struct {
		Download []struct {
			Host  string `json:"host"`
			Index int64  `json:"index"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"download"`
		Watch []struct {
			Index int64  `json:"index"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"watch"`
	} `json:"data"`
	Messages struct {
		All   []interface{} `json:"all"`
		Keyed []interface{} `json:"keyed"`
	} `json:"messages"`
	Status int64 `json:"status"`
}

func GetMovieIDAndPoster(content string) (string, string) {

	idre := regexp.MustCompile("{id(.+?)}")
	id := idre.FindString(content)

	pre := regexp.MustCompile("data-poster(.+?)>")
	p := pre.FindString(content)

	return id[6 : len(id)-2], p[13 : len(p)-2]
}

func GetVideoPageUrl(h *Common.HttpRequest, id string) (string, error) {
	url := fmt.Sprintf("https://njav.tv/zh/ajax/v/%s/videos", id)
	h.SetHeader("Priority", "u=1, i")
	body, err := h.SendRequest(url)
	if err != nil {
		return "", err
	}

	var m VideosRes
	err = json.Unmarshal(body, &m)
	if err != nil {
		return "", err
	}
	return m.Data.Watch[0].URL, nil
}

func GetM3u8Url(pre string, poster string, h *Common.HttpRequest) (string, error) {
	url := fmt.Sprintf("%s&poster=%s", pre, poster)
	h.SetHeader("Priority", "u=0, i")
	h.SetHeader("Referer", "https://njav.tv/")

	body, err := h.SendRequest(url)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile("https(.+?)m3u8")
	u := re.FindString(string(body))
	return strings.ReplaceAll(u, "\\/", "/"), nil
}

func GetM3u8List(h *Common.HttpRequest, url string) ([]string, error) {
	h.SetHeader("Referer", "https://javplayer.me/")

	body, err := h.SendRequest(url)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile("v(.+?).jpg")
	u := re.FindAllString(string(body), -1)

	return u, nil
}

func Start(url string) {
	splist := strings.Split(url, "/")
	Mp4Name := strings.Replace(splist[len(splist)-1], "-", "", -1)
	client := Common.GetClientWithProxy()
	header := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36",
		"Cookies":    "", //设置cookies
		"Priority":   "u=0, i",
		"Referer":    url,
		"Origin":     url,
	}

	body, h, err := Common.RequestAndGetHttpRequest(url, client, header)
	if err != nil {
		log.Fatal(err)
		return
	}

	id, poster := GetMovieIDAndPoster(string(body))
	fmt.Println(id)
	fmt.Println(poster)

	//fmt.Scanln()
	pageUrl, err := GetVideoPageUrl(h, id)
	if err != nil {
		log.Fatal(err)
		return
	}
	//fmt.Scanln()
	fmt.Println(pageUrl)

	murl, err := GetM3u8Url(pageUrl, poster, h)
	if err != nil {
		log.Fatal(err)
		return
	}
	//fmt.Scanln()
	fmt.Println(murl)

	ulist, err := GetM3u8List(h, murl)
	if err != nil {
		log.Fatal(err)
		return
	}

	//fmt.Scanln()
	fmt.Println(Mp4Name)

	err = DownloadVideo(murl, ulist, Mp4Name, "./Video", h)
	if err != nil {
		log.Fatal(err)
		return
	}

}
