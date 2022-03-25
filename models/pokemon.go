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
	ID       uint `gorm:"primarykey"`
	Modifier float64
}

type Condition struct {
	ID           uint `gorm:"primarykey"`
	Index        uint `gorm:"default:1"` // Greater is better xD
	ModifierID   uint
	Modifier     Modifier
	WeatherTypes []WeatherType
	MoonTypes    []MoonType
	TimeTypes    []TimeType
}

type PokemonSpawnRelation struct {
	PokemonID   uint
	Pokemon     Pokemon
	SpawnAreaID uint
	SpawnArea   SpawnArea
	Probability float64
	Conditions  []Condition
}

type Pokemon struct {
	ID uint `gorm:"primarykey" json:"id"`
}
