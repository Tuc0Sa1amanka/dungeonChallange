package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"dungeon-challenge/internal/config"
	"dungeon-challenge/internal/game"
	"dungeon-challenge/internal/report"
)

func main() {
	dungeon, err := config.Init("config.json")
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Open("events")
	if err != nil {
		log.Fatalf("Failed to open events.txt: %v", err)
	}
	defer file.Close()
	handlers := []game.EventHandler{
		nil,
		game.Register,
		game.EnterDungeon,
		game.KillMonster,
		game.NextFloor,
		game.PreviousFloor,
		game.EnterToBoss,
		game.KillBoss,
		game.LeaveDungeon,
		game.CannotContinue,
		game.RestoreHealth,
		game.ReceiveDamage,
	}
	players := make(map[uint]*game.Player)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		event, err := game.ParseEvent(line)
		if err != nil {
			log.Fatal(err)
		}
		out := game.ProcessEvent(players, dungeon, event, handlers)
		if out != "" {
			fmt.Println(out)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("error reading file: %v\n", err)
	}
	report.PrintFinalReport(players, dungeon)
}
