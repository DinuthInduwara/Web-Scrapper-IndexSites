package main

import (
	"log"
	"path"
	"webscrapper/db"
	"webscrapper/funcs"
)

func workerLoop() {

}

func main() {
	funcs.URL = &funcs.MyURL{} //Todo  -  Change this
	funcs.URL.SetUrl("https://iheartwatson.net/gallery/albums/images/Appearances/2020/2023/0312-EltonJohnFoundation/iheartwatson-20230312-eltonjohnfoundation-001.jpeg")
	log.Println("[S] Web Crawling Bot Started \t.:")

	localPath := funcs.URL.GetLocalPath()

	response, _ := funcs.SendReq(funcs.URL.Parser.String()) // Todo -Change url
	defer funcs.CloseReqBody(response)

	isHtml := funcs.IsHtmlFile(response)

	if isHtml {
		htmlFName := path.Join(localPath, "index.html")
		isExits := funcs.IsExits(htmlFName)
		if !isExits {
			funcs.CreateSubDirsFromFile(htmlFName)
			log.Println("[H] Html File Downloading \t.:", htmlFName)
			err := funcs.SaveReqBody(response, htmlFName)
			if err != nil {
				log.Panic(err)
				return
			}
		}
		urls, err := funcs.PhraseHtmlATags(&htmlFName)
		for _, url := range *urls {
			// Todo - Loop over the function
			url = url
		}
		if err != nil {
			log.Panic(err)
			return
		}

	} else {
		inDb := db.IsinDb(&db.TFile{
			Url:  funcs.URL.Parser.String(), // Todo - change url
			Size: response.ContentLength,
		})
		if !inDb {
			funcs.CreateSubDirsFromFile(localPath)
			log.Println("[D] File Downloading \t.:", localPath)
			err := funcs.SaveReqBody(response, localPath)
			if err != nil {
				log.Panic("[E] File Downloading Err \t.:", err)
				return
			} else {
				err := db.AddFileToDb(&db.TFile{Url: response.Request.URL.String(), Size: response.ContentLength})
				log.Println("File adder and err is ", err)
				if err != nil {
					panic(err)
				}
			}
		}
	}

}
