package main

import (
	"net/http"
	"dev.rubetek.com/go-archetype-project/pkg/logger"
	_ "image/jpeg"
	_ "image/gif"
	_ "image/png"
	"image"
	"github.com/nfnt/resize"
	"strings"
	"bytes"
	"image/png"
	"encoding/base64"
	"encoding/json"
)

type TResponse []string
type TRequest  []string

type HandlerCtx struct {
	mux    *http.ServeMux
	l logger.Logger
}

func NewHandler(log logger.Logger) *HandlerCtx {
	s := &HandlerCtx{mux: http.NewServeMux()}
	s.l = log

	s.mux.HandleFunc("/", s.index)
	return s
}

func (s *HandlerCtx) index(w http.ResponseWriter, r *http.Request) {
	s.l.Info("/")
	hh, ok:= r.Header["Content-Type"]
	if ok && strings.HasPrefix(hh[0], "multipart/form-data;") {
		reader, err := r.MultipartReader()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500"))
			return
		}

		var thumbnails []image.Image
		for {
			nr, err := reader.NextPart()
			if err != nil {
				s.l.Info("s>", err)
				break
			}
			s.l.Info("s> fileName", nr.FileName())

			im, _, err := image.Decode(nr)
			if err != nil {
				s.l.Info("s>", err)
				break
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500"))
				return
			}
			thumbnails = append(thumbnails, resize.Resize(100, 100, im, resize.NearestNeighbor))
			nr.Close()
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<html><body>"))
		for _,t := range thumbnails {
			buff := bytes.Buffer{}
			png.Encode(&buff, t)

			w.Write([]byte(`<img src="data:image/png;base64,`))
			encoder := base64.NewEncoder(base64.StdEncoding, w)
			encoder.Write(buff.Bytes())
			encoder.Close()
			w.Write([]byte(`"/>`))
		}
		w.Write([]byte("</html></body>"))
	}

	if ok && strings.HasPrefix(hh[0], "application/json") {
		req := make(TRequest, 0)
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			s.l.Info("s>", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500"))
			s.l.Info(err)
			return
		}

		resp := make(TResponse, 0)
		w.Header().Set("Content-Type", "application/json")
		for _, v := range req {
			input := base64.NewDecoder(base64.StdEncoding, strings.NewReader(v))
			im, imageType, err := image.Decode(input)
			if err != nil {
				s.l.Info("s> image.Decode", err)

				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500"))
				return
			}
			s.l.Info("s> got", imageType)
			thumbnail := resize.Resize(2, 2, im, resize.NearestNeighbor)

			var buff bytes.Buffer
			png.Encode(&buff, thumbnail)
			encoded := base64.StdEncoding.EncodeToString(buff.Bytes())

			resp = append(resp, encoded)
		}
		encoder := json.NewEncoder(w)
		encoder.Encode(resp)
	}
}

func (s *HandlerCtx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "example Go server")
	s.mux.ServeHTTP(w, r)
}