package pokeapi

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
)

func (c *Client) GetPokemon(pokemonName string) (Pokemon, error) {
	url := baseURL + "/" + "pokemon" + "/" + pokemonName
	if val, ok := c.cache.Get(url); ok {
		var pokemon Pokemon
		if err := json.Unmarshal(val, &pokemon); err != nil {
			return Pokemon{}, err
		}
		return pokemon, nil
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Pokemon{}, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Pokemon{}, err
	}
	defer resp.Body.Close()
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, err
	}
	var pokemon Pokemon
	if err := json.Unmarshal(dat, &pokemon); err != nil {
		return Pokemon{}, err
	}
	return pokemon, nil
}

func (c *Client) Catch(pokemon *Pokemon) bool {
	r := rand.Intn(pokemon.BaseExperience)
	if r >= 40 {
		// Catched !
		return true
	}
	return false
}
