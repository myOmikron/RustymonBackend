package models

import "github.com/myOmikron/echotools/utilitymodels"

type PokedexEntry struct {
	utilitymodels.Common
	CaughtCount    uint16 `json:"caught_count"`
	SeenCount      uint16 `json:"seen_count"`
	ShinySeenCount uint16 `json:"shiny_seen_count"` // ShinySeenCount is included in SeenCount
	PlayerID       uint   `json:"player_id"`
}
