package svc

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/Zensey/go-archetype-project/pkg/domain"
	"github.com/Zensey/go-archetype-project/pkg/utils"
)

var (
	blackListCompany = []string{"firstName", "lastName", "fullName"}
	blackListPerson  = []string{"companyName", "companyName2", "companyTypeID"}
)

func getCustomerFieldsBlacklist(c *domain.Customer) []string {
	if c.CompanyName != "" {
		return blackListCompany
	}
	return blackListPerson
}

func mkStringMapFromStruct(f interface{}, blackList []string) map[string]string {
	m := make(map[string]string)

	val := reflect.ValueOf(f).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		fieldName := tag.Get("json")
		if fieldName == "" {
			fieldName = typeField.Name
		} else {
			fieldName = strings.Split(fieldName, ",")[0]
		}

		if !utils.HasString(blackList, fieldName) {
			switch typeField.Type.String() {
			case "int":
				m[fieldName] = strconv.FormatInt(valueField.Int(), 10)
			case "string":
				v := valueField.String()
				if v != "" {
					m[fieldName] = valueField.String()
				}
			}
		}
	}
	return m
}
