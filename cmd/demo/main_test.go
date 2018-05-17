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
	"io"
	"net/http"
	"golang.org/x/net/html"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"path"
)

const url = "http://localhost:8080/"
func init() {
	fmt.Println("Test init()")
	go InitServer()
}

func ReadFilesInDir(handler func(fileName string, imageFile *os.File)) (countImgs int) {
	countImgs = 0

	testDir := "../../test_data/"
	files, err := ioutil.ReadDir(testDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		imageFile, err := os.Open(testDir + f.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}
		_, _, err = image.Decode(imageFile)
		if err != nil {
			//fmt.Println(err)
			continue
		} else {
			//fmt.Println(testDir + f.Name(), imageType)
			countImgs ++
			imageFile.Seek(0,0)
			handler(testDir + f.Name(), imageFile)
		}
		imageFile.Close()
	}
	return
}

func Test_Multipart(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	handler := func (fileName string, imageFile *os.File) {
		part, err := writer.CreateFormFile("file", path.Base(fileName))
		if err != nil {
			fmt.Println(err)
			t.Fail()
			return
		}
		_, err = io.Copy(part, imageFile)
		if err != nil {
			fmt.Print(err)
			t.Fail()
			return
		}
	}
	countImgs := ReadFilesInDir(handler)

	err := writer.Close()
	if err != nil {
		fmt.Print(err)
		t.Fail()
		return
	}

	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err)
		t.Fail()
		return
	}

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	z := html.NewTokenizer(resp.Body)
	l: for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			break l

		case tt == html.SelfClosingTagToken:
			tk := z.Token()
			if tk.Data == "img" {
				fmt.Println("got")
				countImgs--
				continue
			}
		}
	}
	if countImgs != 0 {
		t.Fail()
	}
}

func Test_Json(t *testing.T) {
	images := make([]string, 0)
	handler := func (fileName string, imageFile *os.File) {
		imgFile, err := os.Open(fileName)
		if err != nil {
			fmt.Print(">", err)
			t.Fail()
			return
		}
		fInfo, _ := imgFile.Stat()
		buf := make([]byte, fInfo.Size())

		fReader := bufio.NewReader(imgFile)
		fReader.Read(buf)

		dat := base64.StdEncoding.EncodeToString(buf)
		images = append(images, dat)
	}
	countImgs := ReadFilesInDir(handler)
	fmt.Println(countImgs)

	b, _ := json.Marshal(images)
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err)
		t.Fail()
		return
	}
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)

	obj := make(TResponse, 0)
	err = json.NewDecoder(resp.Body).Decode(&obj)
	if err != nil {
		fmt.Print(err)
		t.Fail()
		return
	}
	for _,v := range obj {
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			continue
		}
		_, imageType, err := image.Decode(bytes.NewReader(b))
		if err != nil {
			continue
		}
		fmt.Println("got", imageType)
		countImgs--
	}
	fmt.Println("N_Sent - N_Got =", countImgs)
	if countImgs != 0 {
		t.Fail()
	}
}
