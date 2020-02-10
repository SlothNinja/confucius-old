package confucius

import (
	"encoding/gob"

	"bitbucket.org/SlothNinja/slothninja-games/sn/contest"
	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"golang.org/x/net/context"
)

func init() {
	gob.RegisterName("*game.scoreChiefMinisterEntry", new(scoreChiefMinisterEntry))
	gob.RegisterName("*game.scoreAdmiralEntry", new(scoreAdmiralEntry))
	gob.RegisterName("*game.scoreGeneralEntry", new(scoreGeneralEntry))
	gob.RegisterName("*game.announceWinnersEntry", new(announceWinnersEntry))
}

func (g *Game) endOfRoundPhase(ctx context.Context) (cs contest.Contests) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	g.Phase = EndOfRound
	g.placeNewOfficialsPhase(ctx)
	if completed := g.discardPhase(ctx); completed {
		g.returnActionCubesPhase(ctx)
		cs = g.endOfGamePhase(ctx)
	}
	return
}

func (g *Game) placeNewOfficialsPhase(ctx context.Context) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	for _, m := range g.Ministries {
		g.placeNewOfficialIn(m)
	}
}

func (g *Game) placeNewOfficialIn(m *Ministry) {
	for _, s := range []Seniority{1, 2, 6, 7} {
		if _, ok := m.Officials[s]; !ok {
			o := g.OfficialsDeck.Draw()
			o.Seniority = s
			m.Officials[s] = o
			return
		}
	}
}

func (g *Game) newRoundPhase(ctx context.Context) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Entering")

	g.Round += 1
	for _, p := range g.Players() {
		p.TakenCommercial = false
	}
}
