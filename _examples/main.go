package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/adelowo/go-mono/mono"
)

func main() {

	client, err := mono.New(mono.WithAPISecret("test_sk_oiyh22"))
	if err != nil {
		log.Fatal(err)
	}

	details, err := client.Account.Details(context.Background(), "66dc25efde48d60efe9ebc39")
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(os.Stdout).Encode(details)
}
