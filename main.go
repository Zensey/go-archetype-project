package main

import (
	"github.com/Zensey/go-archetype-project/cmd"
	_ "github.com/Zensey/go-archetype-project/generated"

	"github.com/markbates/pkger"
	"github.com/ory/x/configx"
	"github.com/spf13/cobra"
)

var version string

var rootCmd = &cobra.Command{
	Use:   "sentinel",
	Short: "Run Sentinel",
}

func main() {
	pkger.Include("/.schema/config.schema.json")

	rootCmd.AddCommand(cmd.ServeCmd)
	configx.RegisterFlags(rootCmd.PersistentFlags())
	rootCmd.Execute()
	return
}
