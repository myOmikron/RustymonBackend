package models

type WeatherType struct {
	ID uint `gorm:"primarykey" json:"id"`
}

type MoonType struct {
	ID uint `gorm:"primarykey" json:"id"`
}

type TimeType struct {
	ID uint `gorm:"primarykey" json:"id"`
}

type SpawnArea struct {
	ID uint `gorm:"primarykey" json:"id"`
}

type Modifier struct {
	ID       uint    `json:"id" gorm:"primarykey"`
	Modifier float64 `json:"modifier" gorm:"not null"`
}

type Condition struct {
	ID         uint     `gorm:"primarykey"`
	Index      uint     `gorm:"not null;default:1"` // Greater is better xD. May be omitted
	ModifierID uint     `json:"modifier_id" gorm:"not null;constraint:OnDelete:CASCADE"`
	Modifier   Modifier `json:"-"`

	WeatherTypes []WeatherType `json:"weather_types" gorm:"many2many;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MoonTypes    []MoonType    `json:"moon_types" gorm:"many2many;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TimeTypes    []TimeType    `json:"time_types" gorm:"many2many;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type HeldItemCondition struct {
	ID          uint    `gorm:"primarykey"`
	ItemID      uint    `json:"item_id" gorm:"not null;constraint:OnDelete:CASCADE"`
	Item        Item    `json:"-"`
	Probability float64 `json:"probability" gorm:"not null"`
}

type PokemonSpawnRelation struct {
	ID                uint                `gorm:"primarykey"`
	PokemonID         uint                `json:"pokemon_id" gorm:"not null;constraint:OnDelete:CASCADE;"`
	Pokemon           Pokemon             `json:"-"`
	SpawnAreas        []SpawnArea         `json:"-" gorm:"many2many;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	MinLevel          uint8               `json:"min_level" gorm:"not null;default:1"`
	MaxLevel          uint8               `json:"max_level" gorm:"not null;default:100"`
	Probability       float64             `json:"probability" gorm:"not null"`
	FemaleProbability float64             `json:"female_probability" gorm:"default:0.5"`
	HeldItemCondition []HeldItemCondition `json:"held_item_condition" gorm:"many2many;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Conditions        []Condition         `json:"conditions" gorm:"many2many;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Pokemon struct {
	ID uint `gorm:"primarykey" json:"id"`
}
