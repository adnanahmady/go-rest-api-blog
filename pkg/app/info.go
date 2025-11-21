package app

import (
	"log"
	"os"
	"strings"
)

func GetRootPath() string {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get root path: %v", err)
	}
	root := strings.Split(currentPath, "/internal")[0]
	return root
}
