package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Zensey/go-archetype-project/pkg/driver"
	"github.com/Zensey/go-archetype-project/pkg/driver/config"
	"github.com/Zensey/go-archetype-project/pkg/utils"
	"github.com/Zensey/go-archetype-project/pkg/x"
	"github.com/oklog/run"
	"github.com/spf13/cobra"
	"go.uber.org/automaxprocs/maxprocs"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "",
	Long:  "",
	Run:   runServe,
}

func runServe(cmd *cobra.Command, args []string) {
	d := driver.New(cmd.Flags())

	if d.Config().CGroupsV1AutoMaxProcsEnabled() {
		_, err := maxprocs.Set(maxprocs.Logger(d.Logger().Infof))
		if err != nil {
			d.Logger().WithError(err).Fatal("Couldn't set GOMAXPROCS")
		}
	}
	d.Logger().Info("Customer service v." + config.Version)

	// Shutdown gracefully
	{
		var g run.Group

		sh := utils.NewSigTermHandler()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		public := setup(d)
		runHttpServer(&g, ctx, d, public.Router, d.Config().PublicListenOn())

		g.Add(func() error { return sh.Wait() }, func(err error) { cancel(); sh.Stop() })
		g.Run()
	}
}

func setup(d driver.Registry) (public *x.RouterPublic) {
	public = x.NewRouterPublic()

	d.RegisterRoutes(public)
	return
}

func runHttpServer(gr *run.Group, ctx context.Context, d driver.Registry, handler http.Handler, address string) {
	srv := &http.Server{
		Addr:    address,
		Handler: handler,
	}

	gr.Add(func() error {
		d.Logger().Infof("Setting up http server on %s", address)
		return srv.ListenAndServe()

	}, func(err error) {
		fmt.Println("shutdown http server...")
		srv.Shutdown(ctx)
	})
}
