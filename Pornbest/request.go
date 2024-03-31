package Pornbest

import (
	"regexp"
	"spider/Common"
)

func GetCode(content string) string {
	re := regexp.MustCompile("eval.+")
	code := re.FindString(content)

	return code
}

func M3u8List(url string, req *Common.HttpRequest) ([]string, error) {
	body, err := req.SendRequest(url)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile("v.+(.ts)")
	list := re.FindAllString(string(body), -1)

	return list, nil
}
