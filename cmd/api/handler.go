package main

import (
	"encoding/json"
	_ "image/gif"
	_ "image/jpeg"
	"net/http"

	"github.com/Zensey/go-archetype-project/pkg/app"
	"github.com/Zensey/go-archetype-project/pkg/logger"
	"github.com/Zensey/go-archetype-project/pkg/types"
	"github.com/Zensey/go-archetype-project/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	logger.Logger
	a   *app.App
	mux *http.ServeMux
}

func NewHandler(app *app.App) *Handler {
	h := &Handler{mux: http.NewServeMux(), Logger: app.Logger, a: app}
	h.mux.HandleFunc("/api/reviews", h.apiHandler)
	return h
}

func (h *Handler) apiHandler(w http.ResponseWriter, r *http.Request) {
	resp := types.ReviewResponse{}
	err := func() error {
		req := types.Review{}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			return err
		}

		err := utils.TxMutate(h.a.Db, func(tx *sqlx.Tx) error {
			err := tx.Get(&resp.ReviewID,
				"insert into production.productreview (productid, reviewername, emailaddress, comments) "+
					"values ($1,$2,$3,$4) returning productreviewid",
				req.ProductID, req.Name, req.Email, req.Review)
			return err
		})
		if err != nil {
			return err
		}

		m := types.MsgReview{Review: req.Review, ReviewID: resp.ReviewID}
		jm, err := json.Marshal(&m)
		if err != nil {
			return err
		}
		_, err = h.a.Redis.LPush(app.QueueWorker1, string(jm)).Result()
		return err
	}()

	if err != nil {
		h.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	resp.Success = err == nil
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
