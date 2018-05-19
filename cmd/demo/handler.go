package main

import (
	"net/http"
	"dev.rubetek.com/go-archetype-project/pkg/logger"
	_ "image/jpeg"
	_ "image/gif"
	_ "image/png"
	"image"
	"strings"
	"bytes"
	"image/png"
	"encoding/base64"
	"encoding/json"
	"io"
	"runtime/debug"
	"github.com/nfnt/resize"
)

type TResponse struct {
   Thumbs []string  `json:"thumbs,omitempty"`
   thumbnails *[]image.Image
}

func NewTResponse() TResponse {
	return TResponse{thumbnails: &[]image.Image{}}
}

func (r TResponse) addThumb(im image.Image) {
	*r.thumbnails = append(*r.thumbnails, resize.Resize(thumbDim, thumbDim, im, resize.NearestNeighbor))
}

func (r TResponse) writeResp(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	for _, v := range *r.thumbnails {
		//s.l.Info("s> t")
		var buff bytes.Buffer
		png.Encode(&buff, v)
		encoded := base64.StdEncoding.EncodeToString(buff.Bytes())

		r.Thumbs = append(r.Thumbs, encoded)
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(r)
}

type TRequest struct {
	Urls []string `json:"imgsUrls,omitempty"`
	Imgs []string `json:"imgs,omitempty"`
}

func ReadTRequest(r io.Reader) (TRequest, error){
	obj := TRequest{}
	err := json.NewDecoder(r).Decode(&obj)
	return obj, err
}

func (r *TRequest) handleImgs(s *HandlerCtx, res TResponse) error {
	for _, v := range r.Imgs {
		input := base64.NewDecoder(base64.StdEncoding, strings.NewReader(v))
		im, imageType, err := image.Decode(input)
		if err != nil {
			s.l.Info("s> image.Decode", err)
		}
		s.l.Info("s> i got", imageType)
		res.addThumb(im)
	}
	return nil
}

func (r *TRequest) handleUrls(s *HandlerCtx, res TResponse) error {
	for _,u := range r.Urls {
		//s.l.Info("s> url", u)
		resp, err := http.Get(u)
		if err != nil {
			s.l.Info("s> get", err)
			continue
		}
		im, _, err := image.Decode(resp.Body)
		if err != nil {
			s.l.Info("s> image.Decode", err)
			continue
		}
		res.addThumb(im)
	}
	return nil
}

const thumbDim = 10

type HandlerCtx struct {
	mux    *http.ServeMux
	l logger.Logger
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
			s.l.Infof("%s: %s", r, debug.Stack()) // line 20
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500"))
		}
	}()

	hh, ok:= r.Header["Content-Type"]
	if ok && strings.HasPrefix(hh[0], "multipart/form-data;") {
		resp := NewTResponse()
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
			s.l.Info("s> fileName", nr.FileName(), contentType, ok, nr.FormName())

			switch contentType {
			case "application/octet-stream":
				im, _, err := image.Decode(nr)
				if err != nil {
					s.l.Info("s>", err)
					continue
				}
				resp.addThumb(im)

			case "text/json":
				req, _ := ReadTRequest(nr)
				req.handleUrls(s, resp)
			}
			nr.Close()
		}
		resp.writeResp(w)
	}

	if ok && strings.HasPrefix(hh[0], "application/json") {
		resp := NewTResponse()
		req, err := ReadTRequest(r.Body)
		if err != nil {
			panic(err)
		}
		req.handleImgs(s, resp)
		req.handleUrls(s, resp)

		resp.writeResp(w)
	}
}

func (s *HandlerCtx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "example Go server")
	s.mux.ServeHTTP(w, r)
}