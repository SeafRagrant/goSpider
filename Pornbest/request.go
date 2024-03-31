package Pornbest

import (
	"net/http"
	"regexp"
	"spider/Common"
)

func GetCode(url string, client *http.Client, header map[string]string) (string, error) {
	body, err := Common.Request(url, client, header)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile("eval.+")
	code := re.FindString(string(body))

	return code, nil
}

func M3u8List(url string, client *http.Client, header map[string]string) ([]string, error) {
	body, err := Common.Request(url, client, header)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile("v.+(.ts)")
	list := re.FindAllString(string(body), -1)

	return list, nil
}
