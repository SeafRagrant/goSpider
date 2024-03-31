package Common

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

func Request(url string, client *http.Client, header map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range header {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetClient() *http.Client {
	return &http.Client{}
}

func GetClientWithProxy() *http.Client {
	ProxyUrl, err := url.Parse("http://127.0.0.1:7890") //clash代理
	if err != nil {
		log.Fatal(err)
		return GetClient()
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(ProxyUrl),
	}
	return &http.Client{
		Transport: transport,
	}
}
