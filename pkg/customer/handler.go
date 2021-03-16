package customer

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Zensey/go-archetype-project/pkg/driver/config"
	"github.com/Zensey/go-archetype-project/pkg/x"
)

const (
	CustomersPath = "/api/customers"
	IndexPath     = "/"
	ViewPath      = "/view"
	EditPath      = "/edit"
	NewPath       = "/new"
)

type Handler struct {
	r InternalRegistry
	c *config.Provider

	tplEdit, tplNew, tplSearch *template.Template
}

func NewHandler(r InternalRegistry, c *config.Provider) *Handler {
	return &Handler{r: r, c: c}
}

func (h *Handler) SetRoutes(public *x.RouterPublic, corsMiddleware func(http.Handler) http.Handler) {
	h.tplSearch = template.Must(template.ParseFiles("tpl/search.gohtml"))
	h.tplNew = template.Must(template.ParseFiles("tpl/new.gohtml"))
	h.tplEdit = template.Must(template.ParseFiles("tpl/edit.gohtml"))

	public.Handler("GET", CustomersPath, http.HandlerFunc(h.GetCustomers))
	public.Handler("GET", IndexPath, http.HandlerFunc(h.Index))
	public.Handler("GET", ViewPath, http.HandlerFunc(h.ViewCustomers))

	public.Handler("GET", EditPath, http.HandlerFunc(h.Edit))
	public.Handler("POST", EditPath, http.HandlerFunc(h.Edit))

	public.Handler("GET", NewPath, http.HandlerFunc(h.New))
	public.Handler("POST", NewPath, http.HandlerFunc(h.New))
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	to := ViewPath
	http.Redirect(w, r, to, http.StatusFound)
}

func (h *Handler) ViewCustomers(w http.ResponseWriter, r *http.Request) {
	h.tplSearch.Execute(w, nil)
}

func (h *Handler) New(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.tplNew.Execute(w, NewEditFormView())

	case "POST":
		r.ParseForm()

		f := NewEditFormView()
		f.copyFormData(r)

		f.Customer.Validate(&f.Errors)

		if len(f.Errors) == 0 {
			h.r.CustomersManager().SaveCustomer(r.Context(), f.Customer)

			id := strconv.FormatInt(f.Customer.ID, 10)
			http.Redirect(w, r, "/edit?id="+id, http.StatusFound)
			return
		}

		h.tplNew.Execute(w, &f)
	}
}

func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
		c, err := h.r.CustomersManager().GetCustomerById(r.Context(), id)
		if err != nil {
			http.Redirect(w, r, ViewPath, http.StatusFound)
		}
		h.tplEdit.Execute(w, &EditFormView{c, nil})

	case "POST":
		r.ParseForm()
		f := NewEditFormView()

		id, err := strconv.ParseInt(r.Form.Get("id"), 10, 64)
		if err != nil {
			f.Errors = append(f.Errors, ErrorKV{"id", err.Error()})
		}

		f.copyFormData(r)
		f.Customer.ID = id
		f.Customer.Validate(&f.Errors)

		if len(f.Errors) == 0 {
			h.r.CustomersManager().SaveCustomer(r.Context(), f.Customer)
			http.Redirect(w, r, "/edit?id="+r.Form.Get("id"), http.StatusFound)
			return
		}

		h.tplEdit.Execute(w, &f)
	}
}

type GetCustomersResponse struct {
	Rows    []Customer `json:"rows"`
	Page    int        `json:"page"`
	Total   int        `json:"total"`
	Records int        `json:"records"`
}

func (h *Handler) GetCustomers(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	po := PaginationOptions{
		Page:  1,
		Limit: 10,
	}
	p, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if p > 0 {
		po.Page = p
	}
	rows, _ := strconv.Atoi(r.URL.Query().Get("rows"))
	if rows > 0 {
		po.Limit = rows
	}

	qo := CustomersQueryOptions{
		OrderByCol: r.URL.Query().Get("sidx"),
		Order:      r.URL.Query().Get("sord"),
		FirstName:  r.URL.Query().Get("fname"),
		LastName:   r.URL.Query().Get("lname"),
	}
	c, err := h.r.CustomersManager().GetCustomers(ctx, &po, &qo)
	if err != nil {
		w.WriteHeader(500)
	}

	json.NewEncoder(w).Encode(GetCustomersResponse{c, po.ResultPage, po.ResultPages, po.ResultRecs})
}

func (h *Handler) logOrAudit(err error, r *http.Request) {
	x.LogAudit(r, err, h.r.Logger())
}
