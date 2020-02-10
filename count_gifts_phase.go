package confucius

import (
	"encoding/gob"
	"fmt"
	"html/template"

	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"golang.org/x/net/context"
)

func init() {
	gob.RegisterName("*game.countGiftsEntry", new(countGiftsEntry))
}

type GiftCount struct {
	PlayerID      int
	GiftsGiven    int
	GiftsReceived int
	ActionCubes   int
}

type GiftCounts []*GiftCount

func (g *Game) countGiftsPhase(ctx context.Context) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	g.Phase = CountGifts
	g.beginningOfPhaseReset()

	counts := make(GiftCounts, g.NumPlayers)

	// Assign Action Cubes
	for i, player := range g.Players() {
		counts[i] = new(GiftCount)
		counts[i].PlayerID = player.ID()
		counts[i].GiftsGiven = player.GiftsGiven()
		counts[i].GiftsReceived = len(player.GiftsReceived)

		switch counts[i].GiftsGiven + counts[i].GiftsReceived {
		case 0:
			player.ActionCubes, counts[i].ActionCubes = 3, 3
		case 1, 2:
			player.ActionCubes, counts[i].ActionCubes = 4, 4
		default:
			player.ActionCubes, counts[i].ActionCubes = 5, 5
		}
	}

	// Create ActionLog Entry
	g.newCountGiftsEntry(counts)
}

type countGiftsEntry struct {
	*Entry
	Counts GiftCounts
}

func (g *Game) newCountGiftsEntry(counts GiftCounts) *countGiftsEntry {
	e := &countGiftsEntry{
		Entry:  g.newEntry(),
		Counts: counts,
	}
	g.Log = append(g.Log, e)
	return e
}

func (e *countGiftsEntry) HTML() template.HTML {
	g := e.Game().(*Game)
	var s string
	for _, count := range e.Counts {
		s += fmt.Sprintf("<div>%s received %d action cubes for giving %d gifts and receiving %d gifts.</div>",
			g.NameByPID(count.PlayerID), count.ActionCubes, count.GiftsGiven, count.GiftsReceived)
	}
	return template.HTML(s)
}
