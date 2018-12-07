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
	"html/template"
)

type Handler struct {
	logger.Logger
	a   *app.App
	mux *http.ServeMux

	t *template.Template
}

func NewHandler(app *app.App) *Handler {
	h := &Handler{mux: http.NewServeMux(), Logger: app.Logger, a: app}
	h.mux.HandleFunc("/api/reviews", utils.Log(utils.BasicAuth(h.reviewHandler, "user", "pass"), app.Logger))
	h.mux.HandleFunc("/api/report", utils.Log(utils.BasicAuth(h.reportHandler, "admin", "admin"), app.Logger))

	h.t = template.Must(template.New("webpage").Parse(tpl))
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) reviewHandler(w http.ResponseWriter, r *http.Request) {
	resp := types.ReviewResponse{}
	err := func() error {
		req := types.Review{}
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			return err
		}
		if req.Name == "" || req.ProductID == "" || req.Review == "" {
			return types.NewErrorLogic("invalid request")
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
		switch err.(type) {
		case *types.ErrorLogic:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		h.Error(err)
	}
	resp.Success = err == nil
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) reportHandler(w http.ResponseWriter, r *http.Request) {
	report := types.ReportTable{
		Title: "ReportTable",
		Items: []types.RowReview{},
		Count: []types.RowCount{},
	}

	err := func() error {
		rows, err := h.a.Db.Queryx("select productid, reviewername, emailaddress, approved, reviewdate from production.productreview")
		if err != nil {
			return err
		}
		for rows.Next() {
			rev := types.RowReview{}
			err = rows.StructScan(&rev)
			if err != nil {
				return err
			}
			report.Items = append(report.Items, rev)
		}
		q := `select t.approved, COALESCE(count, 0) count
			from (VALUES (true), (false), (null)) as t (approved)
			left join (select count(*), approved from production.productreview group by (approved)) b 
			on t.approved=b.approved`

		rows, err = h.a.Db.Queryx(q)
		if err != nil {
			return err
		}
		for rows.Next() {
			rev := types.RowCount{}
			err = rows.StructScan(&rev)
			if err != nil {
				return err
			}
			report.Count = append(report.Count, rev)
		}
		return h.t.Execute(w, report)
	}()

	if err != nil {
		h.Info(err)
	}
}
