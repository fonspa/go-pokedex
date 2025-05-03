package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func (c *Client) ListLocations(pageURL *string) (LocationAreas, error) {
	url := baseURL + "/" + "location-area"
	if pageURL != nil {
		url = *pageURL
	}
	if val, ok := c.cache.Get(url); ok {
		var locations LocationAreas
		if err := json.Unmarshal(val, &locations); err != nil {
			return LocationAreas{}, err
		}
		return locations, nil
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationAreas{}, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return LocationAreas{}, err
	}
	defer resp.Body.Close()
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationAreas{}, err
	}
	c.cache.Add(url, dat)
	var locations LocationAreas
	if err := json.Unmarshal(dat, &locations); err != nil {
		return LocationAreas{}, err
	}
	return locations, nil
}
