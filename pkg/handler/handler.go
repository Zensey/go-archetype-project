package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi"

	"github.com/Zensey/go-archetype-project/pkg/svc"

	"github.com/Zensey/go-archetype-project/pkg/domain"
	"github.com/Zensey/slog"
)

const apiKeyHeader = "api_key"

type Handler struct {
	l      slog.Logger
	s      *svc.CustomerService
	apiKey string
}

func NewHandler(l slog.Logger, s *svc.CustomerService, apiKey string) *Handler {
	return &Handler{l, s, apiKey}
}

func setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, PATCH, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, api_key")
}

func (h *Handler) SaveCustomer(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Header.Get(apiKeyHeader) != h.apiKey {
		w.WriteHeader(401)
	}

	err := func() error {
		cu := domain.Customer{}

		if err := json.NewDecoder(r.Body).Decode(&cu); err != nil {
			return err
		}
		cu.Dirty = true

		err := h.s.SaveCustomerInTx(&cu)
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(&cu)
	}()

	if err != nil {
		h.l.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) GetCustomers(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Header.Get(apiKeyHeader) != h.apiKey {
		w.WriteHeader(401)
	}

	err := func() error {
		c, err := h.s.GetCustomers()
		if err != nil {
			return err
		}

		return json.NewEncoder(w).Encode(&c)
	}()

	if err != nil {
		h.l.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		setCORS(w)

		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
