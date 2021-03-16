package sql

import (
	"context"
	"fmt"
	"strings"

	"github.com/Zensey/go-archetype-project/pkg/customer"
	"github.com/ory/x/sqlcon"
)

func (p *Persister) GetCustomerById(ctx context.Context, id int64) (*customer.Customer, error) {
	var cl customer.Customer
	return &cl, sqlcon.HandleError(p.Connection(ctx).Where("id = ?", id).First(&cl).Error)
}

func (p *Persister) GetCustomers(ctx context.Context, qo *customer.CustomersQueryOptions) ([]customer.Customer, error) {

	args := make([]interface{}, 0)
	pred := make([]string, 0)
	if qo.FirstName != "" || qo.LastName != "" {
		if qo.FirstName != "" {
			pred = append(pred, "first_name ilike ?")
			args = append(args, "%"+qo.FirstName+"%")
		}
		if qo.LastName != "" {
			pred = append(pred, "last_name ilike ?")
			args = append(args, "%"+qo.LastName+"%")
		}
	}
	qq := "select count(1) from customers"
	whereStr := ""
	if len(pred) > 0 {
		whereStr = " where " + strings.Join(pred, " and ")
	}

	cnt := int64(0)
	p.Connection(ctx).Raw(qq+whereStr, args...).Scan(&cnt)

	////////////////////////////////////////////
	pages := 1
	if qo.Limit > 0 {
		pages = int(cnt) / qo.Limit
		if int(cnt)%qo.Limit > 0 {
			pages += 1
		}
	}
	if qo.Page > pages {
		qo.Page = pages
	}
	qo.ResultPages = pages
	qo.ResultPage = qo.Page
	qo.ResultRecs = int(cnt)

	cl := make([]customer.Customer, 0)
	if qo.OrderBy == "" {
		qo.OrderBy = "id"
	}
	orderBy := qo.OrderBy
	if qo.Order != "" {
		orderBy += " " + qo.Order
	}

	q := fmt.Sprintf("select * from (select * from customers c %s order by %s) t limit ? offset ? ", whereStr, orderBy)

	args = append(args, qo.Limit)
	args = append(args, (qo.Page-1)*qo.Limit)
	err := p.Connection(ctx).Raw(q, args...).Scan(&cl).Error

	return cl, sqlcon.HandleError(err)
}

func (p *Persister) CreateCustomer(ctx context.Context, c *customer.Customer) error {
	err := p.Connection(ctx).Omit("id").Create(c).Error
	return sqlcon.HandleError(err)
}

func (p *Persister) DeleteCustomer(ctx context.Context, id int64) (int, error) {
	c := p.Connection(ctx).Delete(&customer.Customer{ID: id})
	return int(c.RowsAffected), c.Error
}

func (p *Persister) SaveCustomer(ctx context.Context, c *customer.Customer) error {
	e := p.Connection(ctx).Save(c).Error
	return e
}
