package Common

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type HttpRequest struct {
	r *http.Request
	c *http.Client
}

func (h *HttpRequest) SendRequest(u string) ([]byte, error) {
	tmp, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	if strings.LastIndex(tmp.Host, ":") > strings.LastIndex(tmp.Host, "]") {
		tmp.Host = strings.TrimSuffix(tmp.Host, ":")
	}
	h.r.URL = tmp
	h.r.Host = tmp.Host

	resp, err := h.c.Do(h.r)
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
func (h *HttpRequest) ResetHeader(header map[string]string) {
	h.r.Header = make(http.Header, len(header))
	for k, v := range header {
		h.r.Header.Add(k, v)
	}
}

// Request 第一次request时调用该函数， 一般用于获取html主体
func Request(url string, client *http.Client, header map[string]string) ([]byte, *HttpRequest, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	for k, v := range header {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, nil, err
	}

	body, err := io.ReadAll(resp.Body)
	h := &HttpRequest{req, client}
	if err != nil {
		return nil, h, err
	}
	return body, h, nil
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
