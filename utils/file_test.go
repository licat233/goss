package utils_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/licat233/goss/utils"
)

func TestGetFiles(t *testing.T) {
	list, err := utils.GetDirFiles("../example", "html")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(list)
}

func TestFileExt(t *testing.T) {
	name := "abcjpg"
	res := utils.FileExt(name)
	fmt.Println(res)
}
