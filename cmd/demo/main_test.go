package main

import (
	"bufio"
	"bytes"
	"dev.rubetek.com/go-archetype-project/pkg/logger"
	"encoding/base64"
	"encoding/json"
	"image"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path"
	"runtime/debug"
	"testing"
)

var l logger.Logger

func init() {
	InitServer()
	l, _ = logger.NewLogger(logger.LogLevelInfo, "_test_", logger.BackendConsole)
}

const serviceUrl = "http://localhost:8080/upload"

var imgsUrls = []string{
	"https://www.google.ru/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png",
	"https://www.google.ru/images/branding/googlelogo/2x/googlelogo_color_120x44dp.png",
}

func findImagesInDir(handler func(fileName string, imageFile *os.File)) (countImgs int) {
	testDir := "../../test_data/"
	countImgs = 0
	files, err := ioutil.ReadDir(testDir)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		imageFile, err := os.Open(testDir + f.Name())
		if err != nil {
			continue
		}
		_, _, err = image.Decode(imageFile)
		if err != nil {
			continue
		} else {
			countImgs++
			imageFile.Seek(0, 0)
			handler(testDir+f.Name(), imageFile)
		}
		imageFile.Close()
	}
	return
}

func handleResp(req *http.Request, countImgs int, t *testing.T) (err error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	//l.Info("response Status:", resp.Status)

	obj := TResponse{}
	err = json.NewDecoder(resp.Body).Decode(&obj)
	if err != nil {
		panic(err)
	}

	for _, v := range obj.Thumbs {
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			l.Info(err)
			continue
		}
		im, imageType, err := image.Decode(bytes.NewReader(b))
		if err != nil {
			l.Info(err)
			continue
		}
		l.Info("got", imageType, im.Bounds())
		countImgs--
	}
	l.Info("test> N_Sent - N_Got =", countImgs)
	if countImgs != 0 {
		t.Fail()
	}
	return
}

func handlePanic(t *testing.T) {
	if r := recover(); r != nil {
		l.Infof("%srv: %srv", r, debug.Stack()) // line 20
		t.Fail()
	}
}

func Test_Multipart(t *testing.T) {
	defer handlePanic(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	handler := func(fileName string, imageFile *os.File) {
		part, err := writer.CreateFormFile("file", path.Base(fileName))
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(part, imageFile)
		if err != nil {
			panic(err)
		}
	}
	countImgs := findImagesInDir(handler)
	jsonData := TRequest{}
	jsonData.Urls = imgsUrls
	countImgs += len(imgsUrls)

	hdr := make(textproto.MIMEHeader)
	hdr.Set("Content-Type", "text/json")
	hdr.Set("Content-Disposition", "form-data; name=\"imgsUrls\"")
	fileWriter, _ := writer.CreatePart(hdr)
	json.NewEncoder(fileWriter).Encode(jsonData)
	err := writer.Close()
	if err != nil {
		panic(err)
	}

	req, _ := http.NewRequest("POST", serviceUrl, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	handleResp(req, countImgs, t)
}

func readFile(file *os.File) []byte {
	fInfo, _ := file.Stat()
	buf := make([]byte, fInfo.Size())
	fReader := bufio.NewReader(file)
	fReader.Read(buf)
	return buf
}

func Test_Json(t *testing.T) {
	defer handlePanic(t)

	jsonData := TRequest{}
	handler := func(fileName string, imageFile *os.File) {
		imgBytes := readFile(imageFile)
		str := base64.StdEncoding.EncodeToString(imgBytes)
		jsonData.Imgs = append(jsonData.Imgs, str)
	}
	countImgs := findImagesInDir(handler)

	jsonData.Urls = imgsUrls
	countImgs += len(imgsUrls)
	b, _ := json.Marshal(jsonData)

	req, _ := http.NewRequest("POST", serviceUrl, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	handleResp(req, countImgs, t)
}
