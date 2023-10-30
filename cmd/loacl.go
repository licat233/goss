package cmd

import (
	_loacl "github.com/licat233/goss/modules/local"
	"github.com/licat233/goss/utils"
	"github.com/spf13/cobra"
)

var localCmd = &cobra.Command{
	Use:     "local",
	GroupID: "modules",
	Short:   _loacl.Name,
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			utils.Success("done.")
		}()
		run(_loacl.Run, _loacl.CheckoutFileExts)
	},
}

func init() {
	// localCmd.AddCommand(localStartCmd)
	localCmd.SetHelpTemplate(setColorizeHelp(localCmd.HelpTemplate()))
	rootCmd.AddCommand(localCmd)
}
