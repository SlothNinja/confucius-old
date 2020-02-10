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
	gob.RegisterName("*game.commercialEntry", new(commercialEntry))
}

func (g *Game) commercial(ctx context.Context) (tmpl string, a game.ActionType, err error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	var (
		cds, ncds ConCards
		cbs, cv   int
	)

	// Validate and get cards and cubes
	if cds, cbs, err = g.validateCommercial(ctx); err != nil {
		a = game.None
		return
	}

	cp := g.CurrentPlayer()
	cp.PerformedAction = true
	cp.TakenCommercial = true

	// Place Action Cubes
	cp.PlaceCubesIn(CommercialSpace, cbs)

	// Move played cards from hand to discard pile
	cp.ConCardHand.Remove(cds...)
	g.ConDiscardPile.Append(cds...)

	// Take Cards and Create Action Object for logging
	cv = cds.Coins()
	ncds = make(ConCards, cv+1)
	for i := range ncds {
		ncds[i] = g.DrawConCard()
	}
	cp.ConCardHand.Append(ncds...)

	// Create Action Object for logging
	entry := cp.newCommercialEntry(cds, ncds)

	// Set flash message
	restful.AddNoticef(ctx, string(entry.HTML()))
	a = game.Cache
	return
}

type commercialEntry struct {
	*Entry
	Played   ConCards
	Received ConCards
}

func (p *Player) newCommercialEntry(cds, ncds ConCards) *commercialEntry {
	g := p.Game()
	e := new(commercialEntry)
	e.Entry = p.newEntry()
	e.Played = cds
	e.Received = ncds
	p.Log = append(p.Log, e)
	g.Log = append(g.Log, e)
	return e
}

func (e *commercialEntry) HTML() template.HTML {
	length := len(e.Played)
	coins := e.Played.Coins()
	return restful.HTML("%s spent %d Confucius %s having %d %s to receive %d cards of commercial income.",
		e.Player().Name(), length, pluralize("card", length), coins, pluralize("coin", coins), len(e.Received))
}

func (g *Game) validateCommercial(ctx context.Context) (cds ConCards, cbs int, err error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	if cbs, err = g.validatePlayerAction(ctx); err != nil {
		return
	}

	if cds, err = g.getConCards(ctx, "commercial"); err != nil {
		return
	}

	switch cp, cv := g.CurrentPlayer(), cds.Coins(); {
	case cp.TakenCommercial:
		err = sn.NewVError("You have already taken the commercial income action this round.")
	case cv > 4:
		err = sn.NewVError("You may only pay up to 4 coins. You paid %d coins.", cv)
	}

	return
}

func (g *Game) EnableCommercial(ctx context.Context) bool {
	cp := g.CurrentPlayer()
	return g.inActionsOrImperialFavourPhase() && g.CurrentPlayer() != nil &&
		!cp.PerformedAction && g.CUserIsCPlayerOrAdmin(ctx) &&
		cp.hasEnoughCubesFor(CommercialSpace) && !cp.TakenCommercial && cp.hasConCards()
}
