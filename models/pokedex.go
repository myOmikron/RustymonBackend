package models

import "github.com/myOmikron/echotools/utilitymodels"

type PokedexEntry struct {
	utilitymodels.CommonSoftDelete
	PokemonID        uint    `json:"pokemon_id"`
	Pokemon          Pokemon `json:"-"`
	CaughtCount      uint16  `json:"caught_count"`
	SeenCount        uint16  `json:"seen_count"`
	ShinySeenCount   uint16  `json:"shiny_seen_count"`   // ShinySeenCount is included in SeenCount
	ShinyCaughtCount uint16  `json:"shiny_caught_count"` //ShinyCaughtCount is included in CaughtCount
	PlayerID         uint    `json:"player_id"`
}
