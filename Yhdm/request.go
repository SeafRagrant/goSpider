package Yhdm

import (
	"net/http"
	"regexp"
	"spider/Common"
)

func GetM3u8(url string, client *http.Client, header map[string]string) (string, error) {
	body, err := Common.Request(url, client, header)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile("_player_x_")
	code := re.FindString(string(body))
	return code, nil
}
