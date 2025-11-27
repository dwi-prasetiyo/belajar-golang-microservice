package env

import (
	"context"
	"fmt"
	"email-service/src/common/log"
	"os"

	"github.com/hashicorp/vault/api"
)

func LoadFromVault() {
	client, err := api.NewClient(&api.Config{
		Address: os.Getenv("VAULT_ADDR"),
	})

	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	client.SetToken(os.Getenv("VAULT_TOKEN"))

	secret, err := client.KVv2("belajar-golang-microservice").Get(context.Background(), "email-service")
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	for key, val := range secret.Data {
		strVal := fmt.Sprintf("%v", val)
		os.Setenv(key, strVal)
	}
}
