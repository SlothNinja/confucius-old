package confucius

import (
	"encoding/gob"
	"html/template"

	"bitbucket.org/SlothNinja/slothninja-games/sn/game"
	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"golang.org/x/net/context"
)

func init() {
	gob.RegisterName("*game.taxIncomeEntry", new(taxIncomeEntry))
}

func (g *Game) taxIncome(ctx context.Context) (string, game.ActionType, error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	cubes, err := g.validatePlayerAction(ctx)
	if err != nil {
		return "", game.None, err
	}

	// Create Action Object for logging
	cp := g.CurrentPlayer()
	cp.PerformedAction = true

	// Place Action Cube(s) In BuyGiftSpace
	cp.PlaceCubesIn(TaxIncomeSpace, cubes)

	// Perform Tax Action
	cp.ConCardHand.Append(g.DrawConCard(), g.DrawConCard())

	entry := cp.newTaxIncomeEntry()

	// Set flash message
	restful.AddNoticef(ctx, string(entry.HTML()))
	return "", game.Cache, nil
}

type taxIncomeEntry struct {
	*Entry
}

func (p *Player) newTaxIncomeEntry() *taxIncomeEntry {
	g := p.Game()
	e := new(taxIncomeEntry)
	e.Entry = p.newEntry()
	p.Log = append(p.Log, e)
	g.Log = append(g.Log, e)
	return e
}

func (e *taxIncomeEntry) HTML() template.HTML {
	return restful.HTML("%s received two Confucius cards of tax income.", e.Player().Name())
}

func (g *Game) EnableTaxIncome(ctx context.Context) bool {
	cp := g.CurrentPlayer()
	return g.CUserIsCPlayerOrAdmin(ctx) && cp.canCollectTaxIncome()
}

func (p *Player) canCollectTaxIncome() bool {
	g := p.Game()
	return g.inActionsOrImperialFavourPhase() && !p.PerformedAction && p.hasEnoughCubesFor(TaxIncomeSpace)
}
