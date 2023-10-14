package main

import (
	"github.com/licat233/goss/cmd"
	"github.com/licat233/goss/utils"
)

func main() {
	defer func() {
		utils.Success("done.")
	}()
	cmd.Execute()
}
