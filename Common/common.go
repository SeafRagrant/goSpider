package Common

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
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
		io.ReadAll(resp.Body)
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

func (h *HttpRequest) SetHeader(key string, value string) {
	h.r.Header.Set(key, value)
}

func (h *HttpRequest) ReSetClient() {
	h.c = &http.Client{}
}

func NewHttpRequest(client *http.Client, header map[string]string) (*HttpRequest, error) {
	req, err := http.NewRequest("GET", "https://www.baidu.com", nil)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header.Add(k, v)
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

// dir可以是绝顶路径也可以是相对路径,最后一级目录只用带目录名，不用加/
func MergeTS(name string, dir string) error {
	cmd := exec.Command("cmd", "/C", fmt.Sprintf("copy /b *.ts %s.ts", name))
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误

	err := cmd.Run()
	if err != nil {
		return errors.New(string(stderr.Bytes()))
	}
	//outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	//fmt.Printf("out:\n%s\n err:\n%s\n", outStr, errStr)
	fmt.Println("ts片段合并成功")
	return nil
}

// dir可以是绝顶路径也可以是相对路径,最后一级目录只用带目录名，不用加/
func FfmpegToh264(name string, dir string) error {
	//cmd := exec.Command("ffmpeg", fmt.Sprintf("-i ./%s.ts -c:v libx264 -crf 18 ./%s.mp4", name, name))
	cmd := exec.Command("ffmpeg", "-i", "./"+name+".ts", "-c:v", "libx264", "-crf", "18", "./"+name+".mp4")
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误

	fmt.Println("正在转码...")
	err := cmd.Run()
	if err != nil {
		return errors.New(string(stderr.Bytes()))
	}

	//outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	//fmt.Printf("out:\n%s\n err:\n%s\n", outStr, errStr)
	fmt.Println("ts转码成功")
	return nil
}
