package customer_test

import (
	"strings"
	"testing"
	"time"

	"github.com/Zensey/go-archetype-project/pkg/customer"
)

var samples2 = []struct {
	customer.Customer
	valid bool
}{
	{
		customer.Customer{
			ID:        1,
			FirstName: strings.Repeat("A", 100),
			LastName:  strings.Repeat("A", 100),
			BirthDate: customer.JSONTime(time.Now().AddDate(-50, 0, 0)),
			Gender:    "Male",
			Email:     "a@b.c",
			Address:   "",
		},
		true,
	},
	{
		customer.Customer{
			ID:        2,
			FirstName: strings.Repeat("A", 100),
			LastName:  strings.Repeat("A", 100),
			BirthDate: customer.JSONTime(time.Now().AddDate(-10, 0, 0)),
			Gender:    "Male",
			Email:     "a@b.c",
			Address:   "",
		},
		false,
	},
	{
		customer.Customer{
			ID:        2,
			FirstName: strings.Repeat("A", 100),
			LastName:  strings.Repeat("A", 0),
			BirthDate: customer.JSONTime(time.Now().AddDate(-19, 0, 0)),
			Gender:    "Male",
			Email:     "a@b.c",
			Address:   "",
		},
		false,
	},
}

func TestValidator(t *testing.T) {
	for _, v := range samples2 {
		errors := make([]customer.ErrorKV, 0)
		v.Validate(&errors)
		//fmt.Println(v.ID, errors)

		if len(errors) > 0 && v.valid {
			t.Errorf(`id: "%v" => unexpected error: %v`, v.ID, errors)
		}
		if len(errors) == 0 && v.valid == false {
			t.Errorf(`id: "%v" => expected error`, v.ID)
		}
	}
}
