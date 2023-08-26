package main

import (
	"log"
	"path"
	"sync"
	"webscrapper/db"
	"webscrapper/funcs"
)

func truncateTextFromEnd(text string, maxLength int) string {
	if len(text) > maxLength {
		startIndex := len(text) - maxLength
		return "..." + text[startIndex:]
	}
	return text
}

func workerLoop(Url *funcs.MyURL, group *sync.WaitGroup, semaphore *funcs.Semaphore) {

	semaphore.Acquire()
	defer semaphore.Release()
	defer group.Done()

	if db.IsUrlInDatabase(Url.Parser.String()) {
		return
	}

	localPath := Url.GetLocalPath()
	response, _ := funcs.SendReq(Url.Parser.String()) // Todo -Change url
	defer funcs.CloseReqBody(response)
	isHtml := funcs.IsHtmlFile(response)

	if isHtml {
		htmlFName := path.Join(localPath, "index.html")
		isExits := funcs.IsExits(htmlFName)
		if !isExits {
			funcs.CreateSubDirsFromFile(htmlFName)
			log.Println("[H] Html File Downloading \t.:", truncateTextFromEnd(htmlFName, 75))
			err := funcs.SaveReqBody(response, htmlFName)
			if err != nil {
				log.Panic(err)
				return
			}
		}
		urls, err := funcs.PhraseHtmlATags(&htmlFName)
		for _, i := range *urls {
			url := &funcs.MyURL{}
			url.SetUrl(i)
			group.Add(1)
			go workerLoop(url, group, semaphore)
		}
		if err != nil {
			log.Panic(err)
			return
		}

	} else {
		// inDb := db.IsinDb(&db.TFile{
		// 	Url:  response.Request.URL.String(),
		// 	Size: response.ContentLength,
		// })
		// if !inDb {
		if !funcs.IsExits(path.Dir(localPath)) {
			funcs.CreateSubDirsFromFile(localPath)
		}
		log.Println("[D] File Downloading \t", filesCount, truncateTextFromEnd(localPath, 75))
		err := funcs.SaveReqBody(response, localPath)
		filesCount += 1
		if err != nil {
			log.Panic("[E] File Downloading Err \t.:", err)
			return
		} else {
			err := db.AddFileToDb(&db.TFile{Url: response.Request.URL.String(), Size: response.ContentLength})
			if err != nil {
				panic(err)
			}
		}
		// }
	}

}

var filesCount = db.DocumentsCount()

func main() {
	funcs.URL = &funcs.MyURL{}
	funcs.URL.SetUrl("https://iheartwatson.net/gallery/albums/images")
	log.Println("[S] Web Crawling Bot Started [S]")
	semaphore := funcs.NewSemaphore(150)

	var wg sync.WaitGroup
	wg.Add(1)
	workerLoop(funcs.URL, &wg, semaphore)

	wg.Wait()
}
