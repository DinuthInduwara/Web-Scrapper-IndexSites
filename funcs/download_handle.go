package funcs

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

var URL *MyURL

type MyURL struct {
	Parser *url.URL
}

func (u *MyURL) SetUrl(Url string) {
	parser, _ := url.Parse(Url)
	u.Parser = parser
}

func (u *MyURL) GetLocalPath() string {
	return path.Join(u.Parser.Hostname(), u.Parser.Path)
}

func IsHtmlFile(response *http.Response) bool {
	if strings.Contains(response.Header.Get("Content-Type"), "text/html") {
		return true
	}
	return false
}

func IsExits(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	return false
}

func SendReq(url string) (*http.Response, error) {
	return http.Get(url)
}

func SaveReqBody(response *http.Response, destination string) error {
	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return err
}

func CloseReqBody(response *http.Response) {
	response.Body.Close()
}

func CreateSubDirsFromFile(filePath string) error {
	dir := path.Dir(filePath)
	log.Println("[C] Creating File Dir \t.:", dir)
	return os.MkdirAll(dir, 0755)
}
