package cli

import (
	"maphoto/internal/app/store"
	"maphoto/internal/util"

	"github.com/spf13/cobra"
)

var dsUrl string

func init() {
	rootCmd.Flags().StringVarP(&dsUrl, "database-url", "d", util.ExcutePath()+"/maphoto.db", "database url")
	rootCmd.AddCommand(resetUserCmd)
}

var resetUserCmd = &cobra.Command{
	Use:     "reset",
	Short:   "reset user password",
	Example: "maphoto reset admin admin1234",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		store.InitDB(dsUrl)
		store.ResetUser(args[0], args[1])
	},
}
