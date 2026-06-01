package report

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"dungeon-challenge/internal/game"
)

func TestPrintFinalReport(t *testing.T) {
	baseTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	dungeon := &game.Dungeon{
		Floors:   2,
		Monstres: 2,
		OpenAt:   baseTime.Add(14*time.Hour + 5*time.Minute),
		Duration: 2 * time.Hour,
	}
	player1 := &game.Player{
		Disqualified: false,
		BossKilled:   true,
		Hp:           35,
		Floors:       make([]game.Floor, dungeon.Floors),
		Monsters:     0,
		ExitTime:     baseTime.Add(15*time.Hour + 4*time.Minute),
	}
	player1.Floors[0].Entry = baseTime.Add(14*time.Hour + 40*time.Minute)
	player1.Floors[0].LastKillTime = baseTime.Add(14*time.Hour + 45*time.Minute)
	player1.Floors[0].Accumulated = 0
	player1.Floors[0].Monsters = 0
	player1.Floors[1].Entry = baseTime.Add(14*time.Hour + 48*time.Minute)
	player1.Floors[1].LastKillTime = baseTime.Add(14*time.Hour + 59*time.Minute)
	player2 := &game.Player{
		Disqualified: false,
		BossKilled:   false,
		Hp:           0,
		Floors:       make([]game.Floor, dungeon.Floors),
		Monsters:     1,
		ExitTime:     baseTime.Add(14*time.Hour + 29*time.Minute),
	}
	player2.Floors[0].Entry = baseTime.Add(14*time.Hour + 10*time.Minute)
	player3 := &game.Player{
		Disqualified: true,
		BossKilled:   false,
		Hp:           100,
		Floors:       make([]game.Floor, dungeon.Floors),
	}
	players := map[uint]*game.Player{
		1: player1,
		2: player2,
		3: player3,
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	PrintFinalReport(players, dungeon)
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()
	expected := `Final report:
[SUCCESS] 1 [00:24:00, 00:05:00, 00:11:00] HP:35
[FAIL] 2 [00:19:00, 00:00:00, 00:00:00] HP:0
[DISQUAL] 3 [00:00:00, 00:00:00, 00:00:00] HP:100
`
	if output != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, output)
	}
}
