package customer

import (
	"net/http"
	"time"
)

type EditFormView struct {
	*Customer
	Errors []ErrorKV
}

func NewEditFormView() *EditFormView {
	return &EditFormView{
		Customer: &Customer{},
		Errors:   make([]ErrorKV, 0),
	}
}

func (f *EditFormView) copyFormData(r *http.Request) {
	c := f.Customer
	frm := r.Form

	c.FirstName = frm.Get("first_name")
	c.LastName = frm.Get("last_name")
	dt, err := time.Parse(dateFormat, frm.Get("birth_date"))
	if err != nil {
		f.Errors = append(f.Errors, ErrorKV{"birth_date", "Wrong format"})
	}
	c.BirthDate = JSONTime(dt)
	c.Gender = frm.Get("gender")
	c.Email = frm.Get("email")
	c.Address = frm.Get("address")
}
