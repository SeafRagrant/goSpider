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

func NewHttpRequest(client *http.Client, header map[string]string) (*HttpRequest, error) {
	req, err := http.NewRequest("GET", "https://www.baidu.com", nil)
	if err != nil {
		return nil, err
	}
	return &HttpRequest{req, client}, nil
}

func request(url string, client *http.Client, header map[string]string) (*http.Request, []byte, error) {
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
	if err != nil {
		return nil, nil, err
	}
	return req, body, nil
}

// Request 只需用到请求结果可使用该函数
func Request(url string, client *http.Client, header map[string]string) ([]byte, error) {
	_, body, err := request(url, client, header)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// RequestAndGetHttpRequest 第一次request时调用该函数， 一般用于获取html主体并且获取 *HttpRequest（用于后续其他请求）
func RequestAndGetHttpRequest(url string, client *http.Client, header map[string]string) ([]byte, *HttpRequest, error) {
	req, body, err := request(url, client, header)
	if err != nil {
		return nil, nil, err
	}

	return body, &HttpRequest{req, client}, nil
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
