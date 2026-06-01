package game

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Dungeon struct {
	Floors   uint
	Monstres uint
	OpenAt   time.Time
	Duration time.Duration
}

type Event struct {
	Time      time.Time
	PlayerId  uint
	EventId   uint
	ExtraStr  string
	ExtraUint uint
}

type Floor struct {
	Entry        time.Time
	Accumulated  time.Duration
	LastKillTime time.Time
	Monsters     uint
}

type Player struct {
	Registered    bool
	Floor         uint
	Floors        []Floor
	EnteredToBoss bool
	BossKilled    bool
	Disqualified  bool
	Hp            uint
	Finished      bool
	ExitTime      time.Time
	Monsters      uint
}

type EventHandler func(player *Player, dungeon *Dungeon, event *Event) string

func Register(player *Player, dungeon *Dungeon, event *Event) string {
	if !player.Registered {
		if event.Time.After(dungeon.OpenAt.Add(dungeon.Duration)) {
			player.Disqualified = true
			player.Finished = true
			return ""
		}
		player.Registered = true
		return fmt.Sprintf("[%s] Player [%d] registered", event.Time.Format("15:04:05"), event.PlayerId)
	}
	return fmt.Sprintf("[%s] Player [%d] makes imposible move [%d]", event.Time.Format("15:04:05"), event.PlayerId, event.EventId)
}

func EnterDungeon(player *Player, dungeon *Dungeon, event *Event) string {
	if player.Floors[0].Entry.IsZero() {
		if event.Time.Before(dungeon.OpenAt) {
			return fmt.Sprintf("[%s] Player [%d] makes imposible move [%d]", event.Time.Format("15:04:05"), event.PlayerId, event.EventId)
		}
		player.Floors[0].Entry = event.Time
		return fmt.Sprintf("[%s] Player [%d] entered the dungeon", event.Time.Format("15:04:05"), event.PlayerId)
	}
	return fmt.Sprintf("[%s] Player [%d] makes imposible move [%d]", event.Time.Format("15:04:05"), event.PlayerId, event.EventId)
}

func KillMonster(player *Player, dungeon *Dungeon, event *Event) string {
	if player.Floors[0].Entry.IsZero() || player.Floors[player.Floor].Monsters == 0 || player.Floor == dungeon.Floors-1 {
		return fmt.Sprintf("[%s] Player [%d] makes imposible move [%d]", event.Time.Format("15:04:05"), event.PlayerId, event.EventId)
	}
	player.Floors[player.Floor].Monsters--
	player.Monsters--
	if player.Floors[player.Floor].Monsters == 0 {
		player.Floors[player.Floor].LastKillTime = event.Time
	}
	return fmt.Sprintf("[%s] Player [%d] killed the monster", event.Time.Format("15:04:05"), event.PlayerId)
}

func NextFloor(player *Player, dungeon *Dungeon, event *Event) string {
	if player.Floors[0].Entry.IsZero() || player.Floors[player.Floor].Monsters != 0 || player.Floor == dungeon.Floors-1 {
		return fmt.Sprintf("[%s] Player [%d] makes imposible move [%d]", event.Time.Format("15:04:05"), event.PlayerId, event.EventId)
	}
	if player.Floors[player.Floor+1].Monsters != 0 {
		player.Floors[player.Floor+1].Entry = event.Time
	}
	player.Floor++
	return fmt.Sprintf("[%s] Player [%d] went to the next floor", event.Time.Format("15:04:05"), event.PlayerId)
}

func PreviousFloor(player *Player, dungeon *Dungeon, event *Event) string {
	if player.Floors[0].Entry.IsZero() || player.Floor == 0 || player.EnteredToBoss {
		return fmt.Sprintf("[%s] Player [%d] makes imposible move [%d]", event.Time.Format("15:04:05"), event.PlayerId, event.EventId)
	}
	if player.Floors[player.Floor].Monsters != 0 {
		player.Floors[player.Floor].Accumulated += event.Time.Sub(player.Floors[player.Floor].Entry)
	}
	player.Floor--
	return fmt.Sprintf("[%s] Player [%d] went to the previous floor", event.Time.Format("15:04:05"), event.PlayerId)
}

func EnterToBoss(player *Player, dungeon *Dungeon, event *Event) string {
	if player.Floor != dungeon.Floors-1 || player.EnteredToBoss {
		return fmt.Sprintf("[%s] Player [%d] makes imposible move [%d]", event.Time.Format("15:04:05"), event.PlayerId, event.EventId)
	}
	player.EnteredToBoss = true
	player.Floors[player.Floor].Entry = event.Time
	return fmt.Sprintf("[%s] Player [%d] entered the boss's floor", event.Time.Format("15:04:05"), event.PlayerId)
}

func KillBoss(player *Player, dungeon *Dungeon, event *Event) string {
	if !player.EnteredToBoss {
		return fmt.Sprintf("[%s] Player [%d] makes imposible move [%d]", event.Time.Format("15:04:05"), event.PlayerId, event.EventId)
	}
	player.BossKilled = true
	player.Floors[player.Floor].LastKillTime = event.Time
	return fmt.Sprintf("[%s] Player [%d] killed the boss", event.Time.Format("15:04:05"), event.PlayerId)
}

func LeaveDungeon(player *Player, dungeon *Dungeon, event *Event) string {
	if player.Floors[0].Entry.IsZero() {
		return fmt.Sprintf("[%s] Player [%d] makes imposible move [%d]", event.Time.Format("15:04:05"), event.PlayerId, event.EventId)
	}
	player.Finished = true
	player.ExitTime = event.Time
	return fmt.Sprintf("[%s] Player [%d] left the dungeon", event.Time.Format("15:04:05"), event.PlayerId)
}

func CannotContinue(player *Player, dungeon *Dungeon, event *Event) string {
	player.Disqualified = true
	player.Finished = true
	player.ExitTime = event.Time
	return fmt.Sprintf("[%s] Player [%d] disqualified", event.Time.Format("15:04:05"), event.PlayerId)
}

func RestoreHealth(player *Player, dungeon *Dungeon, event *Event) string {
	if player.Floors[0].Entry.IsZero() {
		return fmt.Sprintf("[%s] Player [%d] makes imposible move [%d]", event.Time.Format("15:04:05"), event.PlayerId, event.EventId)
	}
	health := event.ExtraUint
	if player.Hp == 100 {
		return fmt.Sprintf("[%s] Player [%d] makes imposible move [%d]", event.Time.Format("15:04:05"), event.PlayerId, event.EventId)
	}
	newHealth := player.Hp + health
	if newHealth > 100 {
		player.Hp = 100
	} else {
		player.Hp = newHealth
	}
	return fmt.Sprintf("[%s] Player [%d] has restored [%d] of health", event.Time.Format("15:04:05"), event.PlayerId, health)
}

func ReceiveDamage(player *Player, dungeon *Dungeon, event *Event) string {
	if player.Floors[0].Entry.IsZero() || player.BossKilled {
		return fmt.Sprintf("[%s] Player [%d] makes imposible move [%d]", event.Time.Format("15:04:05"), event.PlayerId, event.EventId)
	}
	damage := event.ExtraUint
	output := fmt.Sprintf("[%s] Player [%d] recieved [%d] of damage", event.Time.Format("15:04:05"), event.PlayerId, damage)
	if damage >= player.Hp {
		player.Hp = 0
		player.Finished = true
		player.ExitTime = event.Time
		return output + "\n" + fmt.Sprintf("[%s] Player [%d] is dead", event.Time.Format("15:04:05"), event.PlayerId)
	}
	player.Hp -= damage
	return output
}

func ParseEvent(line string) (*Event, error) {
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid format: expected at least 3 fields, got %d", len(parts))
	}
	if !strings.HasPrefix(parts[0], "[") || !strings.HasSuffix(parts[0], "]") {
		return nil, fmt.Errorf("time must be enclosed in brackets, got %q", parts[0])
	}
	if strings.Count(parts[0], "[") != 1 || strings.Count(parts[0], "]") != 1 {
		return nil, fmt.Errorf("time must be enclosed in a single pair of brackets, got %q", parts[0])
	}
	eventTime, err := time.Parse("15:04:05", strings.Trim(parts[0], "[]"))
	if err != nil {
		return nil, fmt.Errorf("invalid event time: expected HH:MM:SS, got %s", parts[0])
	}
	playerId, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("error of parsing playerId: expected positive integer, got %s", parts[1])
	}
	eventId, err := strconv.ParseUint(parts[2], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("error of parsing eventId: expected positive integer, got %s", parts[2])
	}
	if eventId == 0 || eventId > 11 {
		return nil, fmt.Errorf("event ID out of range [1,11]: %d", eventId)
	}
	switch eventId {
	case 9:
		if len(parts) < 4 {
			return nil, fmt.Errorf("event 9 requires at least 4 fields, got %d", len(parts))
		}
		return &Event{eventTime, uint(playerId), uint(eventId), strings.Join(parts[3:], " "), 0}, nil
	case 10, 11:
		if len(parts) != 4 {
			return nil, fmt.Errorf("event 10 - 11 requires 4 fields, got %d", len(parts))
		}
		extraId, err := strconv.ParseUint(parts[3], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error of parsing extraId: expected positive integer, got %s", parts[3])
		}
		return &Event{eventTime, uint(playerId), uint(eventId), "", uint(extraId)}, nil
	default:
		if len(parts) != 3 {
			return nil, fmt.Errorf("event 1 - 8 requires exactly 3 fields, got %d", len(parts))
		}
	}
	return &Event{eventTime, uint(playerId), uint(eventId), "", 0}, nil
}

func ProcessEvent(players map[uint]*Player, dungeon *Dungeon, event *Event, handlers []EventHandler) string {
	player, ok := players[event.PlayerId]
	if !ok {
		floors := make([]Floor, dungeon.Floors)
		for i := 0; i < int(dungeon.Floors-1); i++ {
			floors[i].Monsters = dungeon.Monstres
		}
		player = &Player{
			Hp:       100,
			Floors:   floors,
			Monsters: dungeon.Monstres * (dungeon.Floors - 1),
		}
		players[event.PlayerId] = player
		if event.EventId != 1 {
			player.Disqualified = true
			player.Finished = true
			return fmt.Sprintf("[%s] Player [%d] disqualified", event.Time.Format("15:04:05"), event.PlayerId)
		} else {
			return handlers[event.EventId](player, dungeon, event)
		}
	}
	if player.Finished == true {
		return ""
	}
	if event.Time.After(dungeon.OpenAt.Add(dungeon.Duration)) {
		player.Finished = true
		return ""
	}
	return handlers[event.EventId](player, dungeon, event)
}
