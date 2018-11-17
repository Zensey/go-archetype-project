package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/Zensey/go-archetype-project/pkg/logger"
	"github.com/nfnt/resize"
)

const thumbDim = 100

func (r TResponse) decodeAndAddThumb(rr io.Reader, l logger.Logger) {
	im, imgType, err := image.Decode(rr)
	if err != nil {
		l.Info("srv> image.Decode", err)
	} else {
		l.Info("srv> i got", imgType)
		*r.thumbnails = append(*r.thumbnails, resize.Resize(thumbDim, thumbDim, im, resize.NearestNeighbor))
	}
}

func (r TResponse) writeResp(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	buff := bytes.Buffer{}
	for _, v := range *r.thumbnails {
		buff.Reset()
		png.Encode(&buff, v)
		encoded := base64.StdEncoding.EncodeToString(buff.Bytes())
		r.Thumbs = append(r.Thumbs, encoded)
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(r)
}

func ReadTRequest(r io.Reader) (TRequest, error) {
	obj := TRequest{}
	err := json.NewDecoder(r).Decode(&obj)
	return obj, err
}

func (r *TRequest) handleImgs(s *Handler, handleImg decodeAndAddThumb) error {
	for _, v := range r.Imgs {
		input := base64.NewDecoder(base64.StdEncoding, strings.NewReader(v))
		handleImg(input, s.Logger)
	}
	return nil
}

func (r *TRequest) handleUrls(s *Handler, handleImg decodeAndAddThumb) error {
	for _, u := range r.Urls {
		resp, err := http.Get(u)
		if err != nil {
			s.Info("get error:", err)
			continue
		}
		handleImg(resp.Body, s.Logger)
	}
	return nil
}

type Handler struct {
	logger.Logger
	mux      *http.ServeMux
	requests []int64
	sync.Mutex
}

func NewHandler(log logger.Logger) *Handler {
	s := &Handler{mux: http.NewServeMux()}
	s.Logger = log
	s.mux.HandleFunc("/upload", s.index)
	return s
}

func (s *Handler) index(w http.ResponseWriter, r *http.Request) {
	s.Info(r.RequestURI)
	defer func() {
		if r := recover(); r != nil {
			s.Infof("%srv: %srv", r, debug.Stack())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500"))
		}
	}()

	hh, ok := r.Header["Content-Type"]
	if ok {
		resp := NewTResponse()
		switch {
		case strings.HasPrefix(hh[0], "multipart/form-data;"):
			reader, err := r.MultipartReader()
			if err != nil {
				panic(err)
			}
			for {
				nr, err := reader.NextPart()
				if err != nil {
					break
				}
				contentType := nr.Header["Content-Type"][0]
				switch contentType {
				case "application/octet-stream":
					resp.decodeAndAddThumb(nr, s.Logger)

				case "text/json":
					req, _ := ReadTRequest(nr)
					req.handleImgs(s, resp.decodeAndAddThumb)
					req.handleUrls(s, resp.decodeAndAddThumb)
				}
				nr.Close()
			}

		case strings.HasPrefix(hh[0], "application/json"):
			req, err := ReadTRequest(r.Body)
			if err != nil {
				panic(err)
			}
			req.handleImgs(s, resp.decodeAndAddThumb)
			req.handleUrls(s, resp.decodeAndAddThumb)
		}
		resp.writeResp(w)
	}
}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
