package confucius

import (
	"encoding/gob"
	"html/template"

	"bitbucket.org/SlothNinja/slothninja-games/sn"
	"bitbucket.org/SlothNinja/slothninja-games/sn/game"
	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"golang.org/x/net/context"
)

func init() {
	gob.RegisterName("*game.forceExamEntry", new(forceExamEntry))
}

func (g *Game) forceExam(ctx context.Context) (string, game.ActionType, error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	cards, cubes, err := g.validateForceExam(ctx)
	if err != nil {
		return "", game.None, err
	}

	cp := g.CurrentPlayer()
	cp.PerformedAction = true

	// Move played cards from hand to discard pile
	cp.ConCardHand.Remove(cards...)
	g.ConDiscardPile.Append(cards...)

	// Place Action Cubes
	cp.PlaceCubesIn(ForceSpace, cubes)

	// Create Action Object for logging
	e := cp.newForceExamEntry(cards)

	// Set flash message
	restful.AddNoticef(ctx, string(e.HTML()))
	return "", game.Cache, nil
}

type forceExamEntry struct {
	*Entry
	Played ConCards
}

func (p *Player) newForceExamEntry(cards ConCards) *forceExamEntry {
	g := p.Game()
	e := new(forceExamEntry)
	e.Entry = p.newEntry()
	e.Played = cards
	p.Log = append(p.Log, e)
	g.Log = append(g.Log, e)
	return e
}

func (e *forceExamEntry) HTML() template.HTML {
	length := len(e.Played)
	return restful.HTML("%s spent %d %s having %d coins to force an examination.",
		e.Player().Name(), length, pluralize("card", length), e.Played.Coins())
}

func (g *Game) validateForceExam(ctx context.Context) (ConCards, int, error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	cubes, err := g.validatePlayerAction(ctx)
	if err != nil {
		return nil, 0, err
	}

	cards, err := g.getConCards(ctx, "force-exam")
	if err != nil {
		return nil, 0, err
	}

	coinValue := cards.Coins()
	cp := g.CurrentPlayer()

	switch {
	case g.Round == 1:
		return nil, 0, sn.NewVError("You cannot force an examination during round %d.", g.Round)
	case !cp.canAffordForceExam():
		return nil, 0, sn.NewVError("You selected cards having %d total coins, but you need 2 coins to force an examination.", coinValue)
	}
	return cards, cubes, nil
}

func (g *Game) EnableForceExam(ctx context.Context) bool {
	cp := g.CurrentPlayer()
	return g.inActionsOrImperialFavourPhase() && cp != nil && g.CUserIsCPlayerOrAdmin(ctx) &&
		g.Round > 1 && !cp.PerformedAction && cp.hasEnoughCubesFor(ForceSpace) && cp.canAffordForceExam()
}

func (p *Player) canAffordForceExam() bool {
	return p.ConCardHand.Coins() >= 2
}
