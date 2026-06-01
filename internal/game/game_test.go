package game

import (
	"testing"
	"time"
)

func testDungeon() *Dungeon {
	openAt, _ := time.Parse("15:04:05", "14:05:00")
	return &Dungeon{
		Floors:   2,
		Monstres: 2,
		OpenAt:   openAt,
		Duration: 2 * time.Hour,
	}
}

func TestParseEvent(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		wantErr bool
		check   func(*testing.T, *Event)
	}{
		{
			name:    "valid registration",
			line:    "[14:00:00] 1 1",
			wantErr: false,
			check: func(t *testing.T, e *Event) {
				if e.PlayerId != 1 || e.EventId != 1 || e.ExtraStr != "" {
					t.Errorf("unexpected event: %+v", e)
				}
				if e.Time.Format("15:04:05") != "14:00:00" {
					t.Errorf("wrong time: %v", e.Time)
				}
			},
		},
		{
			name:    "valid damage",
			line:    "[14:27:00] 2 11 60",
			wantErr: false,
			check: func(t *testing.T, e *Event) {
				if e.EventId != 11 || e.ExtraUint != 60 {
					t.Errorf("damage event: %+v", e)
				}
			},
		},
		{
			name:    "cannot continue with reason",
			line:    "[15:00:00] 5 9 connection lost",
			wantErr: false,
			check: func(t *testing.T, e *Event) {
				if e.EventId != 9 || e.ExtraStr != "connection lost" {
					t.Errorf("cannot continue: %+v", e)
				}
			},
		},
		{
			name:    "empty line",
			line:    "",
			wantErr: true,
		},
		{
			name:    "not enough fields",
			line:    "[14:00:00] 1",
			wantErr: true,
		},
		{
			name:    "invalid time format",
			line:    "[14:00] 1 1",
			wantErr: true,
		},
		{
			name:    "negative player ID",
			line:    "[14:00:00] -1 1",
			wantErr: true,
		},
		{
			name:    "event ID out of range",
			line:    "[14:00:00] 1 12",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseEvent(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}
func TestRegister(t *testing.T) {
	d := testDungeon()
	player := &Player{}
	event := &Event{Time: d.OpenAt.Add(-1 * time.Hour), PlayerId: 1, EventId: 1}
	Register(player, d, event)
	if !player.Registered {
		t.Error("registration failed")
	}
	msg2 := Register(player, d, event)
	if msg2 == "" {
		t.Error("re-registration should be impossible move")
	}
	player2 := &Player{}
	event2 := &Event{Time: d.OpenAt.Add(d.Duration).Add(time.Minute), PlayerId: 2, EventId: 1}
	Register(player2, d, event2)
	if !player2.Disqualified {
		t.Error("registration after close should disqualify")
	}
}

func TestEnterDungeon(t *testing.T) {
	d := testDungeon()
	player := &Player{
		Floors: make([]Floor, d.Floors),
	}
	eventBefore := &Event{Time: d.OpenAt.Add(-time.Minute), PlayerId: 1, EventId: 2}
	EnterDungeon(player, d, eventBefore)
	if !player.Floors[0].Entry.IsZero() {
		t.Error("enter before open should be impossible move")
	}
	eventValid := &Event{Time: d.OpenAt.Add(time.Minute), PlayerId: 1, EventId: 2}
	EnterDungeon(player, d, eventValid)
	if player.Floors[0].Entry.IsZero() {
		t.Error("valid enter failed")
	}
	msg := EnterDungeon(player, d, eventValid)
	if msg == "" {
		t.Error("second enter should be impossible")
	}
}

func TestKillMonster(t *testing.T) {
	d := testDungeon()
	player := &Player{
		Floor:    0,
		Monsters: d.Monstres,
		Floors:   make([]Floor, d.Floors),
	}
	player.Floors[0].Monsters = d.Monstres
	player.Floors[0].Entry = time.Now()
	event := &Event{Time: time.Now(), PlayerId: 1, EventId: 3}
	msg := KillMonster(player, d, event)
	if msg == "" || player.Floors[0].Monsters != d.Monstres-1 || player.Monsters != d.Monstres-1 {
		t.Error("kill monster failed")
	}
	for i := uint(0); i < d.Monstres-1; i++ {
		KillMonster(player, d, event)
	}
	if player.Floors[0].Monsters != 0 || player.Floors[0].LastKillTime.IsZero() {
		t.Error("last monster not recorded")
	}
	player.Floor = d.Floors - 1
	msg = KillMonster(player, d, event)
	if msg == "" {
		t.Error("kill on boss floor should be impossible")
	}
}
func TestNextFloor(t *testing.T) {
	d := testDungeon()
	player := &Player{
		Floor:    0,
		Floors:   make([]Floor, d.Floors),
		Monsters: d.Monstres,
	}
	player.Floors[0].Monsters = 0
	player.Floors[0].Entry = time.Now()
	event := &Event{Time: time.Now(), PlayerId: 1, EventId: 4}
	if d.Floors > 2 {
		msg := NextFloor(player, d, event)
		if msg == "" || player.Floor != 1 {
			t.Error("next floor failed")
		}
	}
	if d.Floors == 2 {
		msg := NextFloor(player, d, event)
		if msg == "" || player.Floor != 1 {
			t.Error("next floor to boss should succeed")
		}
	}
}

func TestPreviousFloor(t *testing.T) {
	d := testDungeon()
	player := &Player{
		Floor:  1,
		Floors: make([]Floor, d.Floors),
	}
	player.Floors[0].Entry = time.Now()
	event := &Event{Time: time.Now(), PlayerId: 1, EventId: 5}
	msg := PreviousFloor(player, d, event)
	if msg == "" || player.Floor != 0 {
		t.Error("previous floor failed")
	}
	msg = PreviousFloor(player, d, event)
	if msg == "" {
		t.Error("previous from floor 0 should be impossible")
	}
	player.EnteredToBoss = true
	msg = PreviousFloor(player, d, event)
	if msg == "" {
		t.Error("previous from boss should be impossible")
	}
}

func TestEnterToBoss(t *testing.T) {
	d := testDungeon()
	player := &Player{
		Floor:  d.Floors - 1,
		Floors: make([]Floor, d.Floors),
	}
	event := &Event{Time: time.Now(), PlayerId: 1, EventId: 6}
	msg := EnterToBoss(player, d, event)
	if msg == "" || !player.EnteredToBoss || player.Floors[d.Floors-1].Entry.IsZero() {
		t.Error("enter boss floor failed")
	}
	msg = EnterToBoss(player, d, event)
	if msg == "" {
		t.Error("second enter boss should be impossible")
	}
}

func TestKillBoss(t *testing.T) {
	d := testDungeon()
	player := &Player{
		EnteredToBoss: true,
		Floor:         d.Floors - 1,
		Floors:        make([]Floor, d.Floors),
	}
	event := &Event{Time: time.Now(), PlayerId: 1, EventId: 7}
	msg := KillBoss(player, d, event)
	if msg == "" || !player.BossKilled || player.Floors[d.Floors-1].LastKillTime.IsZero() {
		t.Error("kill boss failed")
	}
	player2 := &Player{}
	msg2 := KillBoss(player2, d, event)
	if msg2 == "" {
		t.Error("kill boss without entering should be impossible")
	}
}

func TestLeaveDungeon(t *testing.T) {
	d := testDungeon()
	player := &Player{
		Floors: make([]Floor, d.Floors),
	}
	player.Floors[0].Entry = time.Now()
	event := &Event{Time: time.Now(), PlayerId: 1, EventId: 8}
	msg := LeaveDungeon(player, d, event)
	if msg == "" || !player.Finished || player.ExitTime.IsZero() {
		t.Error("leave dungeon failed")
	}
}

func TestCannotContinue(t *testing.T) {
	player := &Player{}
	event := &Event{Time: time.Now(), PlayerId: 1, EventId: 9}
	msg := CannotContinue(player, nil, event)
	if msg == "" || !player.Disqualified || !player.Finished || player.ExitTime.IsZero() {
		t.Error("cannot continue failed")
	}
}

func TestRestoreHealth(t *testing.T) {
	d := testDungeon()
	player := &Player{
		Hp:     50,
		Floors: make([]Floor, d.Floors),
	}
	player.Floors[0].Entry = time.Now()
	event := &Event{Time: time.Now(), PlayerId: 1, EventId: 10, ExtraUint: 30}
	msg := RestoreHealth(player, d, event)
	if msg == "" || player.Hp != 80 {
		t.Error("restore health failed")
	}
	event.ExtraUint = 30
	RestoreHealth(player, d, event)
	if player.Hp != 100 {
		t.Error("health cap failed")
	}
	msg = RestoreHealth(player, d, event)
	if msg == "" {
		t.Error("restore at max should be impossible")
	}
}

func TestReceiveDamage(t *testing.T) {
	d := testDungeon()
	player := &Player{
		Hp:     100,
		Floors: make([]Floor, d.Floors),
	}
	player.Floors[0].Entry = time.Now()
	event := &Event{Time: time.Now(), PlayerId: 1, EventId: 11, ExtraUint: 60}
	msg := ReceiveDamage(player, d, event)
	if msg == "" || player.Hp != 40 {
		t.Error("damage failed")
	}
	event.ExtraUint = 50
	msg = ReceiveDamage(player, d, event)
	if msg == "" || player.Hp != 0 || !player.Finished {
		t.Error("lethal damage not handled")
	}
}

func TestProcessEvent(t *testing.T) {
	d := testDungeon()
	players := make(map[uint]*Player)
	handlers := []EventHandler{
		nil,
		Register,
		EnterDungeon,
		KillMonster,
		NextFloor,
		PreviousFloor,
		EnterToBoss,
		KillBoss,
		LeaveDungeon,
		CannotContinue,
		RestoreHealth,
		ReceiveDamage,
	}
	ev, _ := ParseEvent("[14:00:00] 1 1")
	out := ProcessEvent(players, d, ev, handlers)
	if out == "" || players[1] == nil || !players[1].Registered {
		t.Error("processEvent registration failed")
	}
	evClose, _ := ParseEvent("[16:10:00] 1 2")
	out = ProcessEvent(players, d, evClose, handlers)
	if out != "" || !players[1].Finished {
		t.Error("event after close should finish player")
	}
}
