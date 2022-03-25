package models

import (
	"github.com/myOmikron/echotools/utilitymodels"
)

type PlayerPokemonMove struct {
	utilitymodels.Common
	PlayerPokemonID uint  `json:"player_pokemon_id"`
	MoveID          uint  `json:"move_id"`
	Move            Move  `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
	Pp              uint8 `json:"pp"`
	PpUp            uint8 `json:"pp_up"`
}

type PlayerPokemon struct {
	utilitymodels.Common
	PlayerID uint `json:"player_id"`

	// Static content
	PokemonID uint    `json:"pokemon_id" gorm:"not null"`
	Pokemon   Pokemon `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
	Sex       uint8   `json:"sex"` // Sex: 0: Genderless, 1: Female, 2: Male
	Shiny     bool    `json:"shiny"`
	Nature    uint8   `json:"nature"`
	BallUsed  uint8   `json:"ball_used" gorm:"default:0"`

	// All EVs are values in the range [0, 255], starting at 0
	EvHp             uint8 `json:"ev_hp" gorm:"default:0"`
	EvAttack         uint8 `json:"ev_attack" gorm:"default:0"`
	EvDefense        uint8 `json:"ev_defense" gorm:"default:0"`
	EvSpeed          uint8 `json:"ev_speed" gorm:"default:0"`
	EvSpecialAttack  uint8 `json:"ev_special_attack" gorm:"default:0"`
	EvSpecialDefense uint8 `json:"ev_special_defense" gorm:"default:0"`

	// All IVs are random values in the range [0, 31]
	IvHp             uint8 `json:"iv_hp"`
	IvAttack         uint8 `json:"iv_attack"`
	IvDefense        uint8 `json:"iv_defense"`
	IvSpeed          uint8 `json:"iv_speed"`
	IvSpecialAttack  uint8 `json:"iv_special_attack"`
	IvSpecialDefense uint8 `json:"iv_special_defense"`

	// Variable content
	Name      string              `json:"name"`
	CurrentHP uint16              `json:"current_hp"`
	EXP       uint32              `json:"exp"`
	Form      uint8               `json:"form" gorm:"default:0"`
	Ability   uint8               `json:"ability"`
	Happiness uint8               `json:"happiness" gorm:"default:70"` // TODO some species don't use 70 as default
	Status    uint8               `json:"status" gorm:"default:0"`
	EggsSteps uint16              `json:"eggs_steps"`
	Moves     []PlayerPokemonMove `json:"moves" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type PlayerItem struct {
	utilitymodels.Common
	PlayerID uint   `json:"player_id" gorm:"not null"`
	Amount   uint16 `json:"amount" gorm:"not null"`
	ItemID   uint   `json:"item_id" gorm:"not null"`
	Item     Item   `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
}

type Player struct {
	utilitymodels.Common
	UserID      uint               `json:"user_id" gorm:"not null"`
	User        utilitymodels.User `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
	TrainerName string             `json:"trainer_name"`
	Language    uint8              `json:"language" gorm:"default:0"` // Language: 0: English; 1: German
	Female      bool               `json:"female"`
	Money       uint32             `json:"money" gorm:"default:0"`
	Friends     []*Player          `json:"friends" gorm:"many2many:player_friends;"`
	Items       []PlayerItem       `json:"items" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Pokedex     []PokedexEntry     `json:"pokedex" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PokeBox     []PlayerPokemon    `json:"poke_box" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Party       []PlayerPokemon    `json:"party"  gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
