package cli

import (
	"maphoto/internal/app/store"
	"maphoto/internal/app/web/server"
	"maphoto/internal/env"
	"maphoto/internal/util"

	"github.com/spf13/cobra"
)

var (
	urlPrefix      string
	disableSwagger bool
	debug          bool
)

func init() {
	serveCmd.Flags().StringVarP(&urlPrefix, "url-prefix", "u", "", "url prefix")
	serveCmd.Flags().BoolVarP(&disableSwagger, "disable-swagger", "s", false, "disable swagger")
	serveCmd.Flags().BoolVarP(&debug, "debug", "g", false, "debug mode")

	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start a webserver instance",
	Run: func(cmd *cobra.Command, args []string) {

		cfg := env.Options{
			Port:             port,
			DataBaseURL:      databaseURL,
			DefaultAdminName: defaultAdminName,
			DefaultAdminPass: defaultAdminPassword,
			UrlPrefix:        urlPrefix,
			Version:          Inject.Version,
			Debug:            debug,
			Domain:           domain,
		}
		util.InitLogger()
		store.InitDB(cfg.DataBaseURL)
		store.InitAdmin(cfg.DefaultAdminName, cfg.DefaultAdminPass)
		server.RunServer(&cfg)
	},
}
