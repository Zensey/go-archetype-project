package customer

import (
	"fmt"
	"time"
)

const dateFormat = "2006-01-02"

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf(`"%s"`, time.Time(t).Format(dateFormat))
	return []byte(stamp), nil
}

func (t JSONTime) String() string {
	return time.Time(t).Format(dateFormat)
}

type Customer struct {
	ID        int64    `json:"id" gorm:"primaryKey;"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	BirthDate JSONTime `json:"birth_date"`
	Gender    string   `json:"gender"`
	Email     string   `json:"email"`
	Address   string   `json:"address"`

	// Version int // should be used for an optimistic locking to prevent lost updates
}

func (Customer) TableName() string {
	return "customers"
}

func (c Customer) Validate(errors *[]string) {
	if len(c.FirstName) == 0 || len(c.FirstName) > 100 {
		*errors = append(*errors, "first_name")
	}
	if len(c.LastName) == 0 || len(c.LastName) > 100 {
		*errors = append(*errors, "last_name")
	}
	switch c.Gender {
	case "Male":
	case "Female":
	default:
		*errors = append(*errors, "gender")
	}
	if len(c.Address) > 200 {
		*errors = append(*errors, "address")
	}
	if !IsEmailValid(c.Email) {
		*errors = append(*errors, "email")
	}

	if time.Time(c.BirthDate).Before(time.Now().AddDate(-60, 0, 0)) {
		*errors = append(*errors, "birth_day")
	}
	if time.Time(c.BirthDate).After(time.Now().AddDate(-18, 0, 0)) {
		*errors = append(*errors, "birth_day")
	}
}
