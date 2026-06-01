package config

import (
	"dungeon-challenge/internal/game"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func Init(path string) (*game.Dungeon, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	aux := struct {
		Floors   uint   `json:"Floors"`
		Monstres uint   `json:"Monsters"`
		OpenAt   string `json:"OpenAt"`
		Duration uint   `json:"Duration"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return nil, err
	}
	if aux.Duration == 0 {
		return nil, fmt.Errorf("invalid Floors: must be at least 1, got %d", aux.Floors)
	}
	if aux.Floors < 2 {
		return nil, fmt.Errorf("invalid Floors: must be at least 2, got %d", aux.Floors)
	}
	if aux.Monstres == 0 {
		return nil, fmt.Errorf("invalid Monsters: must be greater than 0, got %d", aux.Monstres)
	}
	openTime, err := time.Parse("15:04:05", aux.OpenAt)
	if err != nil {
		return nil, fmt.Errorf("invalid OpenAt format: %w", err)
	}
	d := &game.Dungeon{
		Floors:   aux.Floors,
		Monstres: aux.Monstres,
		OpenAt:   openTime,
		Duration: time.Duration(aux.Duration) * time.Hour,
	}
	return d, nil
}
