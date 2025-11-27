package test

import (
	"order-service/env"
	"order-service/src/common/log"
	"os"
	"path"

	"github.com/joho/godotenv"
)

func init() {
	dir, _ := os.Getwd()
	dir = path.Join(path.Dir(dir), "../../../")

	err := os.Chdir(dir)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	godotenv.Load()
	env.Load()
}
