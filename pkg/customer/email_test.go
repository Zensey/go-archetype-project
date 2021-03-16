package customer_test

import (
	"testing"

	"github.com/Zensey/go-archetype-project/pkg/customer"
)

var (
	samples = []struct {
		mail   string
		format bool
	}{
		{mail: "florian@carrere.cc", format: true},
		{mail: "support@g2mail.com", format: true},
		{mail: " florian@carrere.cc", format: false},
		{mail: "florian@carrere.cc ", format: false},
		{mail: "test@912-wrong-domain902.com", format: true},
		{mail: "0932910-qsdcqozuioqkdmqpeidj8793@gmail.com", format: true},
		{mail: "@gmail.com", format: false},
		{mail: "test@gmail@gmail.com", format: false},
		{mail: "test test@gmail.com", format: false},
		{mail: " test@gmail.com", format: false},
		{mail: "test@wrong domain.com", format: false},
		{mail: "é&ààà@gmail.com", format: false},
		{mail: "admin@busyboo.com", format: true},
		{mail: "a@gmail.fi", format: true},
	}
)

func TestValidateFormat(t *testing.T) {
	for _, s := range samples {
		valid := customer.IsEmailValid(s.mail)

		if !valid && s.format == true {
			t.Errorf(`"%s" => unexpected error`, s.mail)
		}
		if valid && s.format == false {
			t.Errorf(`"%s" => expected error`, s.mail)
		}
	}
}
