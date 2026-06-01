package report

import (
	"dungeon-challenge/internal/game"
	"fmt"
	"sort"
	"time"
)

func PrintFinalReport(players map[uint]*game.Player, dungeon *game.Dungeon) {
	format := func(d time.Duration) string {
		h := int(d / time.Hour)
		m := int((d % time.Hour) / time.Minute)
		s := int((d % time.Minute) / time.Second)
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	ids := make([]uint, 0, len(players))
	for id := range players {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	fmt.Println("Final report:")
	for _, id := range ids {
		p := players[id]
		var status string
		if p.Disqualified {
			status = "DISQUAL"
		} else if p.BossKilled {
			status = "SUCCESS"
		} else {
			status = "FAIL"
		}
		var total time.Duration
		if !p.Floors[0].Entry.IsZero() {
			if !p.ExitTime.IsZero() {
				total = p.ExitTime.Sub(p.Floors[0].Entry)
			} else {
				total = dungeon.OpenAt.Add(dungeon.Duration).Sub(p.Floors[0].Entry)
			}
		}
		var avg time.Duration
		if p.Monsters == 0 {
			var sum time.Duration
			for i := 0; i < int(dungeon.Floors-1); i++ {
				f := &p.Floors[i]
				sum += f.Accumulated + f.LastKillTime.Sub(f.Entry)
			}
			avg = sum / time.Duration(dungeon.Floors-1)
		}
		var boss time.Duration
		if p.BossKilled {
			f := &p.Floors[dungeon.Floors-1]
			boss = f.LastKillTime.Sub(f.Entry)
		}
		fmt.Printf("[%s] %d [%s, %s, %s] HP:%d\n",
			status, id, format(total), format(avg), format(boss), p.Hp)
	}
}
