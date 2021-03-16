package customer

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

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

type Response struct {
	Rows    []Customer `json:"rows"`
	Page    int        `json:"page"`
	Total   int        `json:"total"`
	Records int        `json:"records"`
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	to := ViewPath
	http.Redirect(w, r, to, http.StatusFound)
}

func (h *Handler) ViewCustomers(w http.ResponseWriter, r *http.Request) {
	h.tplSearch.Execute(w, nil)
}

type CustomerFormView struct {
	*Customer
	Errors []string
}

func copyFormData(r *http.Request, f *CustomerFormView) {
	c := f.Customer
	frm := r.Form

	c.FirstName = frm.Get("first_name")
	c.LastName = frm.Get("last_name")
	dt, err := time.Parse(dateFormat, frm.Get("birth_date"))
	if err != nil {
		f.Errors = append(f.Errors, "birth_date")
	}
	c.BirthDate = JSONTime(dt)
	c.Gender = frm.Get("gender")
	c.Email = frm.Get("email")
	c.Address = frm.Get("address")
}

func (h *Handler) New(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case "GET":
		err := h.tplNew.Execute(w, &CustomerFormView{&Customer{}, nil})
		fmt.Println("err>", err)

	case "POST":
		r.ParseForm()

		f := CustomerFormView{
			Customer: &Customer{},
			Errors:   make([]string, 0),
		}
		c := f.Customer

		copyFormData(r, &f)
		c.Validate(&f.Errors)

		if len(f.Errors) == 0 {
			h.r.CustomersManager().SaveCustomer(ctx, c)

			id := strconv.FormatInt(c.ID, 10)
			http.Redirect(w, r, "/edit?id="+id, http.StatusFound)
			return
		}

		h.tplNew.Execute(w, &f)
	}
}

func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case "GET":
		id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
		c, err := h.r.CustomersManager().GetCustomerById(ctx, id)

		err = h.tplEdit.Execute(w, &CustomerFormView{c, nil})
		fmt.Println("err>", err)

	case "POST":
		err := r.ParseForm()

		f := CustomerFormView{
			Customer: &Customer{},
			Errors:   make([]string, 0),
		}
		c := f.Customer
		c.ID, err = strconv.ParseInt(r.Form.Get("id"), 10, 64)
		if err != nil {
			f.Errors = append(f.Errors, "id")
		}

		copyFormData(r, &f)
		c.Validate(&f.Errors)

		if len(f.Errors) == 0 {
			h.r.CustomersManager().SaveCustomer(ctx, c)
			http.Redirect(w, r, "/edit?id="+r.Form.Get("id"), http.StatusFound)
			return
		}

		h.tplEdit.Execute(w, &f)
	}
}

func (h *Handler) GetCustomers(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	o := CustomersQueryOptions{
		Page:      1,
		Limit:     10,
		OrderBy:   r.URL.Query().Get("sidx"),
		Order:     r.URL.Query().Get("sord"),
		FirstName: r.URL.Query().Get("fname"),
		LastName:  r.URL.Query().Get("lname"),
	}
	p, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if p > 0 {
		o.Page = p
	}
	rows, _ := strconv.Atoi(r.URL.Query().Get("rows"))
	if rows > 0 {
		o.Limit = rows
	}

	c, err := h.r.CustomersManager().GetCustomers(ctx, &o)
	if err != nil {
		w.WriteHeader(500)
	}

	json.NewEncoder(w).Encode(Response{c, o.ResultPage, o.ResultPages, o.ResultRecs})
}

func (h *Handler) logOrAudit(err error, r *http.Request) {
	x.LogAudit(r, err, h.r.Logger())
}
