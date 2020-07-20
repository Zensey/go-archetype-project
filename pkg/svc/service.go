package svc

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Zensey/slog"

	"github.com/ulule/deepcopier"

	"github.com/erply/api-go-wrapper/pkg/api/auth"

	"github.com/Zensey/go-archetype-project/pkg/cfg"
	"github.com/Zensey/go-archetype-project/pkg/domain"
	"github.com/erply/api-go-wrapper/pkg/api"
	"github.com/go-pg/pg/v10"
)

type CustomerService struct {
	cred *cfg.ErplyCredentials
	db   *pg.DB
	l    slog.Logger

	sessionExpire time.Time
	apiClient     *api.Client
	wg            sync.WaitGroup
}

func NewCustomerService(db *pg.DB, cred *cfg.ErplyCredentials, l slog.Logger) *CustomerService {
	return &CustomerService{
		db:   db,
		cred: cred,
		l:    l,
	}
}

func (s *CustomerService) refreshClient() {
	if s.sessionExpire.Add(1*time.Minute).After(time.Now()) || s.apiClient == nil {
		sessionKey, err := auth.VerifyUser(s.cred.Username, s.cred.Password, s.cred.ClientCode, http.DefaultClient)
		if err != nil {
			panic(err)
		}

		i, err := auth.GetSessionKeyInfo(sessionKey, s.cred.ClientCode, http.DefaultClient)
		if err != nil {
			panic(err)
		}
		t, _ := strconv.Atoi(i.ExpireUnixTime)
		s.sessionExpire = time.Unix(int64(t), 0)

		s.apiClient, err = api.NewClient(sessionKey, s.cred.ClientCode, nil)
		if err != nil {
			panic(err)
		}
	}
}

func (s *CustomerService) WaitSyncCustomersFinish() {
	s.wg.Wait()
	s.l.Info("SyncCustomers > Stop")
}

func (s *CustomerService) GetDB() *pg.DB {
	return s.db
}

func (s *CustomerService) SyncCustomersPeriodic(ctx context.Context) error {
	s.l.Info("SyncCustomers > Start")

	s.wg.Add(1)
	defer s.wg.Done()

	err := s.SyncCustomers(ctx)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			err := s.SyncCustomers(ctx)
			if err != nil {
				return err
			}
		}
	}
}

func (s *CustomerService) SyncCustomers(ctx context.Context) error {
	s.refreshClient()

	dirtyRecs := make([]domain.Customer, 0)
	var f = func(tx *pg.Tx) error {
		// lock dirty records for update
		err := tx.Model(&dirtyRecs).Where("COALESCE(dirty, false)=?", true).For("UPDATE").Select()
		if err != nil {
			return err
		}

		for _, c := range dirtyRecs {
			c.Dirty = false
			_, err = tx.Model(&c).WherePK().Update()
			if err != nil {
				return err
			}
		}
		return nil
	}
	err := s.db.RunInTransaction(f)

	customerCli := s.apiClient.CustomerManager
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// upload dirty recs
	for _, c := range dirtyRecs {
		m := mkStringMapFromStruct(&c, getCustomerFieldsBlacklist(&c))
		_, err := customerCli.SaveCustomer(ctx, m)
		if err != nil {
			return err
		}
	}

	// grab all customers
	customers, err := customerCli.GetCustomers(ctx, map[string]string{})
	if err != nil {
		return err
	}

	var ff = func(tx *pg.Tx) error {
		for _, c := range customers {
			cu := domain.Customer{}

			err := deepcopier.Copy(c).To(&cu)
			if err != nil {
				return err
			}
			err = s.SaveCustomer(&cu, tx)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return s.db.RunInTransaction(ff)
}

func (s *CustomerService) SaveCustomerInTx(cu *domain.Customer) error {
	return s.db.RunInTransaction(func(tx *pg.Tx) error {
		return s.SaveCustomer(cu, tx)
	})
}

func (s *CustomerService) SaveCustomer(cu *domain.Customer, tx *pg.Tx) error {
	//s.l.Info("SaveCustomer >", cu.ID, cu.FirstName)
	s.l.Tracef("SaveCustomer > %+v\n", cu)

	customerCli := s.apiClient.CustomerManager
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if cu.ID == 0 {
		m := mkStringMapFromStruct(cu, getCustomerFieldsBlacklist(cu))
		//s.l.Info("SaveCustomer m>", m)

		r, err := customerCli.SaveCustomer(ctx, m)
		if err != nil {
			return err
		}
		cu.ID = r.CustomerID
		cu.Dirty = false

		//s.l.Info("SaveCustomer>", cu)
		//s.l.Infof("SaveCustomer> %+v\n", cu)

		_, err = tx.Model(cu).WherePK().Insert()
		if err != nil {
			return err
		}

	} else {
		if cu.Dirty {
			m := mkStringMapFromStruct(cu, getCustomerFieldsBlacklist(cu))
			//s.l.Info("SaveCustomer m>", m)

			_, err := customerCli.SaveCustomer(ctx, m)
			if err != nil {
				return err
			}
			cu.Dirty = false
		}

		exists, err := tx.Model(cu).WherePK().Exists()
		if err != nil {
			return err
		}
		if exists {
			_, err := tx.Model(cu).WherePK().Update()
			if err != nil {
				return err
			}
		} else {
			// case for the SyncCustomers
			_, err := tx.Model(cu).Insert()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *CustomerService) GetCustomers() ([]domain.Customer, error) {
	a := make([]domain.Customer, 0)
	err := s.db.Model(&a).Select()
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (s *CustomerService) MakeDirtyForTest() error {
	c := domain.Customer{}
	err := s.db.Model(&c).Where("id=21").Select()
	if err != nil {
		return err
	}
	c.Dirty = true
	c.Email = "a@a.com"

	_, err = s.db.Model(&c).WherePK().Update()
	if err != nil {
		s.l.Error("MakeDirtyForTest err>", err)
		return err
	}
	return nil
}
