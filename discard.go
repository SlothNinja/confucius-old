package confucius

import (
	"encoding/gob"
	"html/template"

	"bitbucket.org/SlothNinja/slothninja-games/sn"
	"bitbucket.org/SlothNinja/slothninja-games/sn/contest"
	"bitbucket.org/SlothNinja/slothninja-games/sn/game"
	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"bitbucket.org/SlothNinja/slothninja-games/sn/user/stats"
	"golang.org/x/net/context"
)

func init() {
	gob.RegisterName("*game.discardEntry", new(discardEntry))
}

func (g *Game) discardPhase(ctx context.Context) (completed bool) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	g.Phase = Discard
	g.beginningOfPhaseReset()

	ps := make(game.Playerers, 0)
	for _, p := range g.Players() {
		log.Debugf(ctx, "ConCardHand: %#v", p.ConCardHand)
		log.Debugf(ctx, "len(ConCardHand): %#v", len(p.ConCardHand))
		if len(p.ConCardHand) > 4 {
			ps = append(ps, p)
		}
	}

	if len(ps) == 0 {
		completed = true

	}
	log.Debugf(ctx, "ps: %#v", ps)
	g.SetCurrentPlayerers(ps...)
	return
}

func (g *Game) discard(ctx context.Context) (string, game.ActionType, error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	cards, err := g.validateDiscard(ctx)
	if err != nil {
		return "", game.None, err
	}

	cp := g.CurrentPlayer()
	cp.discard(cards...)

	// Set flash message
	restful.AddNoticef(ctx, string(cp.newDiscardEntry(cards...).HTML()))
	return "", game.Cache, nil
}

func (p *Player) discard(cards ...*ConCard) {
	p.PerformedAction = true

	// Move played cards from hand to discard pile
	p.ConCardHand.Remove(cards...)
	p.Game().ConDiscardPile.Append(cards...)
}

type discardEntry struct {
	*Entry
	Discarded ConCards
}

func (p *Player) newDiscardEntry(cards ...*ConCard) *discardEntry {
	g := p.Game()
	e := new(discardEntry)
	e.Entry = p.newEntry()
	e.Discarded = cards
	p.Log = append(p.Log, e)
	g.Log = append(g.Log, e)
	return e
}

func (e *discardEntry) HTML() template.HTML {
	return restful.HTML("%s discarded %d cards.", e.Player().Name(), len(e.Discarded))
}

func (g *Game) validateDiscard(ctx context.Context) (ConCards, error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	cards, err := g.getConCards(ctx, "discard")
	if err != nil {
		return nil, err
	}

	cp := g.CurrentPlayer()
	newHandCount := len(cp.ConCardHand) - len(cards)
	switch {
	case !g.CUserIsCPlayerOrAdmin(ctx):
		return nil, sn.NewVError("Only a current player may discard cards.")
	case g.Phase != Discard:
		return nil, sn.NewVError("You cannot discard cards during the %s phase.", g.PhaseName())
	case newHandCount != 4:
		return nil, sn.NewVError("You must discard down to 4 cards.  You have discarded to %d cards.",
			newHandCount)
	}
	return cards, nil
}

func (g *Game) EnableDiscard(ctx context.Context) bool {
	return g.CUserIsCPlayerOrAdmin(ctx) && g.Phase == Discard && g.CurrentPlayer() != nil &&
		!g.CurrentPlayer().PerformedAction
}

func (g *Game) discardPhaseFinishTurn(ctx context.Context) (s *stats.Stats, cs contest.Contests, err error) {
	if s, err = g.validateFinishTurn(ctx); err != nil {
		return
	}

	cp := g.CurrentPlayer()
	g.RemoveCurrentPlayers(cp)

	if len(g.CurrentPlayerers()) == 0 {
		g.returnActionCubesPhase(ctx)
		cs = g.endOfGamePhase(ctx)
	}
	return
}
