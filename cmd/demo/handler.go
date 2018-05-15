package main

import (
	"net/http"
	"dev.rubetek.com/go-archetype-project/pkg/logger"
	_ "image/jpeg"
	_ "image/gif"
	_ "image/png"
	"image"
	"image/png"
	"bytes"
	"encoding/base64"
	"github.com/nfnt/resize"
	"encoding/json"
	"strings"
)

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

type TResponse []string

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

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<html><body>"))
		for {
			nr, err := reader.NextPart()
			if nr == nil {
				s.l.Info("nr", err)
				break
			}
			s.l.Info("nr", err, nr.FileName())

			im, _, err := image.Decode(nr)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500"))
				return
			}
			thumbnail := resize.Resize(100, 100, im, resize.NearestNeighbor)

			var buff bytes.Buffer
			png.Encode(&buff, thumbnail)
			w.Write([]byte(`<img src="data:image/png;base64,`))
			encoder := base64.NewEncoder(base64.StdEncoding, w)
			encoder.Write(buff.Bytes())
			encoder.Close()
			w.Write([]byte(`"/>`))
			nr.Close()
		}
		w.Write([]byte("</html></body>"))
	}
	if ok && strings.HasPrefix(hh[0], "application/json") {
		req := make(TResponse, 0)
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500"))
			s.l.Info(err)
			return
		}

		resp := make(TResponse, 0)
		w.Header().Set("Content-Type", "application/json")
		for _, v := range req {
			ii := strings.Index(v, "base64,")
			if ii < 0 {
				continue
			}
			v = v[ii+len("base64,"):]
			input := base64.NewDecoder(base64.StdEncoding, strings.NewReader(v))
			im, _, err := image.Decode(input)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500"))
				return
			}
			thumbnail := resize.Resize(100, 100, im, resize.NearestNeighbor)

			var buff bytes.Buffer
			png.Encode(&buff, thumbnail)
			encoded := `data:image/png;base64,`
			encoded += base64.StdEncoding.EncodeToString(buff.Bytes())

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