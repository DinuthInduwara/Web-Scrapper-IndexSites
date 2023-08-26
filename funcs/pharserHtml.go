package funcs

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"os"
)

func PhraseHtmlATags(htmlFile *string) (*[]string, error) {
	file, err := os.Open(*htmlFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	document, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, err
	}

	excludedValues := []string{"Name", "Last Modified", "Size", "Parent Directory"}
	var urls []string
	document.Find("a").Each(func(index int, element *goquery.Selection) {
		linkText := element.Text()
		if !stringExistsInArray(linkText, &excludedValues) {
			href, exists := element.Attr("href")
			if exists {
				joinedPath := url.URL{
					Scheme: URL.Parser.Scheme,
					Host:   URL.Parser.Hostname(),
					Path:   href,
				}
				urls = append(urls, joinedPath.String())
			}
		}
	})
	return &urls, nil
}

func stringExistsInArray(target string, array *[]string) bool {
	for _, element := range *array {
		if element == target {
			return true
		}
	}
	return false
}
