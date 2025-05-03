package main

import (
	"time"

	"example.com/pafcorp/pokedex/internal/pokeapi"
)

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, 10*time.Second)
	cfg := &config{
		pokeapiClient: pokeClient,
		pokedex:       make(map[string]*pokeapi.Pokemon),
	}
	startRepl(cfg)
}
