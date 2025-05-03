package pokeapi

import (
	"encoding/json"
	"net/http"
)

// List all Pokemons in an area
func (c *Client) ListPokemons(areaName string) ([]PokemonEncounter, error) {
	url := baseURL + "/" + "location-area" + "/" + areaName
	if val, ok := c.cache.Get(url); ok {
		var encounters []PokemonEncounter
		if err := json.Unmarshal(val, &encounters); err != nil {
			return nil, err
		}
		return encounters, nil
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var rawMsg map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&rawMsg); err != nil {
		return nil, err
	}
	dat := rawMsg["pokemon_encounters"]
	c.cache.Add(url, dat)
	var encounters []PokemonEncounter
	if err := json.Unmarshal(dat, &encounters); err != nil {
		return nil, err
	}
	return encounters, nil
}
