package main

import (
	"bytes"
	"dev.rubetek.com/go-archetype-project/pkg/logger"
	"encoding/base64"
	"encoding/json"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
)

const thumbDim = 100

type TRequest struct {
	Urls []string `json:"imgsUrls,omitempty"`
	Imgs []string `json:"imgs,omitempty"`
}

type TResponse struct {
	thumbnails *[]image.Image
	Thumbs     []string `json:"thumbs,omitempty"`
}

func NewTResponse() TResponse {
	return TResponse{thumbnails: &[]image.Image{}}
}

type decodeAndAddThumb func(rr io.Reader, l logger.Logger)

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
	for _, v := range *r.thumbnails {
		//srv.log.Info("srv> t")
		var buff bytes.Buffer
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

func (r *TRequest) handleImgs(s *HandlerCtx, handleImg decodeAndAddThumb) error {
	for _, v := range r.Imgs {
		input := base64.NewDecoder(base64.StdEncoding, strings.NewReader(v))
		handleImg(input, s.l)
	}
	return nil
}

func (r *TRequest) handleUrls(s *HandlerCtx, handleImg decodeAndAddThumb) error {
	for _, u := range r.Urls {
		//srv.log.Info("srv> url", u)
		resp, err := http.Get(u)
		if err != nil {
			s.l.Info("srv> get", err)
			continue
		}
		handleImg(resp.Body, s.l)
	}
	return nil
}

type HandlerCtx struct {
	mux *http.ServeMux
	l   logger.Logger
}

func NewHandler(log logger.Logger) *HandlerCtx {
	s := &HandlerCtx{mux: http.NewServeMux()}
	s.l = log
	s.mux.HandleFunc("/upload", s.index)
	return s
}

func (s *HandlerCtx) index(w http.ResponseWriter, r *http.Request) {
	s.l.Info(r.RequestURI)
	defer func() {
		if r := recover(); r != nil {
			s.l.Infof("%srv: %srv", r, debug.Stack())
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
				//srv.log.Info("srv> fileName", nr.FileName(), contentType, ok, nr.FormName())
				switch contentType {
				case "application/octet-stream":
					resp.decodeAndAddThumb(nr, s.l)

				case "text/json":
					req, _ := ReadTRequest(nr)
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

func (s *HandlerCtx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
