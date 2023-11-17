package cmd

import (
	upload "github.com/licat233/goss/modules/upload"
	"github.com/licat233/goss/utils"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:     "upload",
	GroupID: "modules",
	Short:   upload.Name,
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			utils.Success("done.")
		}()
		run(upload.Run, upload.CheckoutFileExts)
	},
}

func init() {
	// uploadCmd.AddCommand(uploadStartCmd)
	uploadCmd.SetHelpTemplate(setColorizeHelp(uploadCmd.HelpTemplate()))
	rootCmd.AddCommand(uploadCmd)
}
