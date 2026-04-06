package cli

import (
	"fmt"
	"maphoto/internal/env"
	"maphoto/internal/util"
	"os"

	"github.com/spf13/cobra"
)

var (
	Inject               *env.Injection
	databaseURL          string
	defaultAdminName     string
	defaultAdminPassword string
	domain               string
	port                 int
)

var rootCmd = &cobra.Command{
	Use:   "Transfeed",
	Short: "A minimalistic feed translator.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Transfeed is a web app, please type `maphoto -help` for more information")
	},
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8090, "port to listen on")
	rootCmd.PersistentFlags().StringVarP(&databaseURL, "database-url", "d", util.ExcutePath()+"/maphoto.db", "database url")
	rootCmd.PersistentFlags().StringVarP(&defaultAdminName, "admin-name", "a", "admin", "default admin name")
	rootCmd.PersistentFlags().StringVarP(&defaultAdminPassword, "admin-password", "w", "admin1234", "default admin password")
	rootCmd.PersistentFlags().StringVarP(&domain, "domain", "x", fmt.Sprintf("http://127.0.0.1:%d", port), "domain")
}

func Execute(injection *env.Injection) {
	Inject = injection
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
