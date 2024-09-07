package main

import (
	"context"
	"log"

	"github.com/adelowo/go-mono/mono"
)

func main() {

	client, err := mono.New(mono.WithAPISecret("test_sk_oiyh2237kx59yd6ugugn"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Account.Unlink(context.Background(), "66dc25efde48d60efe9ebc39")
	if err != nil {
		log.Fatal(err)
	}

	// json.NewEncoder(os.Stdout).Encode(details)
}
