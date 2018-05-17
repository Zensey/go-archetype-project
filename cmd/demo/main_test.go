package main

import (
	"testing"
	"fmt"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"image"
	"mime/multipart"
	"path/filepath"
	"io"
	"net/http"
	"golang.org/x/net/html"
)

func Test_Main(t *testing.T) {
	go Init()

	testDir := "../../test_/"
	files, err := ioutil.ReadDir(testDir)
	if err != nil {
		log.Fatal(err)
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	countImgs := 0

	for _, f := range files {
		imageFile, err := os.Open(testDir + f.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}
		_, imageType, err := image.Decode(imageFile)
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			fmt.Println(f.Name(), imageType)
			countImgs ++
			part, err := writer.CreateFormFile("file", filepath.Base(f.Name()))
			if err != nil {
				fmt.Println(err)
				t.Fail()
				return
			}
			imageFile.Seek(0,0)
			_, err = io.Copy(part, imageFile)
			if err != nil {
				fmt.Print(err)
				t.Fail()
				return
			}
		}
		imageFile.Close()
	}
	err = writer.Close()
	if err != nil {
		fmt.Print(err)
		t.Fail()
		return
	}

	url := "http://localhost:8080/"
	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err)
		t.Fail()
		return
	}

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	z := html.NewTokenizer(resp.Body)
	l: for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			break l

		case tt == html.SelfClosingTagToken:
			tk := z.Token()
			if tk.Data == "img" {
				fmt.Println("got img !", countImgs)
				countImgs--
				continue
			}
		}
	}
	if countImgs != 0 {
		t.Fail()
	}
}