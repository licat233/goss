package utils_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/google/uuid"
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

func TestUuid(t *testing.T) {
	u := uuid.New().String()
	uidWithoutDash := strings.ReplaceAll(u, "-", "")
	fmt.Println(uidWithoutDash)
}
