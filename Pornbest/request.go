package Pornbest

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

func Request(url string, client *http.Client) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetCode(url string, client *http.Client) (string, error) {
	body, err := Request(url, client)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile("eval.+")
	code := re.FindString(string(body))

	return code, nil
}

func M3u8List(url string, client *http.Client) ([]string, error) {
	body, err := Request(url, client)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile("v.+(.ts)")
	list := re.FindAllString(string(body), -1)

	return list, nil
}

func GetClient() *http.Client {
	ProxyUrl, err := url.Parse("http://127.0.0.1:7890") //clash代理
	if err != nil {
		log.Fatal(err)
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(ProxyUrl),
	}
	return &http.Client{
		Transport: transport,
	}
}
