package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"example.com/pafcorp/pokedex/internal/pokeapi"
)

type config struct {
	pokeapiClient       pokeapi.Client
	pokedex             map[string]*pokeapi.Pokemon
	nextLocationURL     *string
	previousLocationURL *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func getAvailableCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays Pokedex usage and commands",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Display the names of the next 20 location areas in the Pokemon world",
			callback:    commandMapf,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the names of the previous 20 location areas in the Pokemon world",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Return all Pokemons present in an location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon by name",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a Pokemon caught in the Pokedex",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List Pokemons caught in your Pokedex",
			callback:    commandPokedex,
		},
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func startRepl(cfg *config) {
	commands := getAvailableCommands()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			log.Fatalf("error reading input: %v", err)
		}
		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}
		cmd, ok := commands[words[0]]
		if ok {
			err := cmd.callback(cfg, words[1:]...)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("unknown command")
		}
	}
}

func commandExit(c *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, args ...string) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	commands := getAvailableCommands()
	for key, value := range commands {
		fmt.Printf("%s: %s\n", key, value.description)
	}
	fmt.Println()
	return nil
}

func commandMapf(c *config, args ...string) error {
	// TODO(Paf): we need to check if we have a Next url in the config and get that.
	// then update the Next and Previous urls in the config after that.
	// if we don't have a Next in the config, we call the basic url.
	locations, err := c.pokeapiClient.ListLocations(c.nextLocationURL)
	if err != nil {
		return err
	}
	c.nextLocationURL = locations.Next
	c.previousLocationURL = locations.Previous
	for _, result := range locations.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandMapb(c *config, args ...string) error {
	// Check if config.Previous is not empty and get the Previous url
	if c.previousLocationURL == nil {
		return errors.New("you are already on the first page of locations areas")
	}
	locations, err := c.pokeapiClient.ListLocations(c.previousLocationURL)
	if err != nil {
		return err
	}
	c.nextLocationURL = locations.Next
	c.previousLocationURL = locations.Previous
	for _, result := range locations.Results {
		fmt.Println(result.Name)
	}
	return nil
}

func commandExplore(c *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("you must provide a location name")
	}
	area := args[0]
	encounters, err := c.pokeapiClient.ListPokemons(area)
	if err != nil {
		return err
	}
	for _, encounter := range encounters {
		fmt.Println(encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(c *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("you must provide a Pokemon name")
	}
	name := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", name)
	pokemon, err := c.pokeapiClient.GetPokemon(name)
	if err != nil {
		return err
	}
	catched := c.pokeapiClient.Catch(&pokemon)
	if catched {
		fmt.Printf("%s was caught!\n", name)
		fmt.Println("You may now inspect it with the inspect command.")
		c.pokedex[pokemon.Name] = &pokemon
	} else {
		fmt.Printf("%s escaped!\n", name)
	}
	return nil
}

func commandInspect(c *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("you must provide a Pokemon name")
	}
	name := args[0]
	p, ok := c.pokedex[name]
	if !ok {
		return fmt.Errorf("%s was not caught yet", name)
	}
	fmt.Println("Name:", p.Name)
	fmt.Println("Height:", p.Height)
	fmt.Println("Weight:", p.Weight)
	fmt.Println("BaseXP:", p.BaseExperience)
	fmt.Println("Stats:")
	for _, stat := range p.Stats {
		fmt.Printf("\t- %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, typ := range p.Types {
		fmt.Printf("\t- %s\n", typ.Type.Name)
	}
	return nil
}

func commandPokedex(c *config, args ...string) error {
	fmt.Println("Your Pokedex:")
	for k := range c.pokedex {
		fmt.Printf("\t- %s\n", k)
	}
	return nil
}
