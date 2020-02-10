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
	gob.RegisterName("*game.invadeLandEntry", new(invadeLandEntry))
}

func (g *Game) invadeLand(ctx context.Context) (string, game.ActionType, error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	// Get Indices and Cards
	box, cards, cubes, err := g.validateInvadeLand(ctx)
	if err != nil {
		return "", game.None, err
	}

	cp := g.CurrentPlayer()
	cp.PerformedAction = true

	// Place Action Cubes
	cp.PlaceCubesIn(RecruitArmySpace, cubes)

	// Commit Recruited Army
	cp.RecruitedArmies -= 1
	box.setPlayer(cp)

	// Move played cards from hand to discard pile
	cp.ConCardHand.Remove(cards...)
	g.ConDiscardPile.Append(cards...)

	// Create Action Object for logging
	entry := cp.newInvadeLandEntry(cards, box)

	// Set flash message
	restful.AddNoticef(ctx, string(entry.HTML()))
	return "", game.Cache, nil
}

type invadeLandEntry struct {
	*Entry
	Played          ConCards
	ForeignLandName string
	Points          int
}

func (p *Player) newInvadeLandEntry(cards ConCards, box *ForeignLandBox) *invadeLandEntry {
	g := p.Game()
	e := new(invadeLandEntry)
	e.Entry = p.newEntry()
	e.Played = cards
	e.ForeignLandName = box.land.Name()
	e.Points = box.Points
	p.Log = append(p.Log, e)
	g.Log = append(g.Log, e)
	return e
}

func (e *invadeLandEntry) HTML() template.HTML {
	return restful.HTML("%s spent %d Confucius cards having a value of %d coins to invade the %d VP box of %s.",
		e.Player().Name(), len(e.Played), e.Played.Coins(), e.Points, e.ForeignLandName)
}

func (g *Game) validateInvadeLand(ctx context.Context) (*ForeignLandBox, ConCards, int, error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	cubes, err := g.validatePlayerAction(ctx)
	if err != nil {
		return nil, nil, 0, err
	}

	cards, err := g.getConCards(ctx, "invade-land")
	if err != nil {
		return nil, nil, 0, err
	}

	box, err := g.getForeignLandBox(ctx, "invade-land")
	if err != nil {
		return nil, nil, 0, err
	}

	coinValue := cards.Coins()
	land := box.land
	cost := land.Cost()
	cp := g.CurrentPlayer()

	switch {
	case coinValue < cost:
		return nil, nil, 0, sn.NewVError("You selected cards having %d total coins, but you need %d coins to invade the selected land.", coinValue, cost)
	case !cp.hasRecruitedArmies():
		return nil, nil, 0, sn.NewVError("You have no recruited armies for an invasion.")
	}

	return box, cards, cubes, nil
}

func (g *Game) EnableInvadeLand(ctx context.Context) (result bool) {
	cp := g.CurrentPlayer()
	return g.inActionsOrImperialFavourPhase() && g.CurrentPlayer() != nil &&
		!cp.PerformedAction && g.CUserIsCPlayerOrAdmin(ctx) &&
		cp.hasEnoughCubesFor(RecruitArmySpace) && cp.hasRecruitedArmies() && cp.canAffordInvasion()
}

func (p *Player) canAffordInvasion() bool {
	g := p.Game()
	coins := p.ConCardHand.Coins()

	for _, land := range g.ForeignLands {
		if coins >= land.Cost() {
			return true
		}
	}
	return false
}

func (p *Player) hasRecruitedArmies() bool {
	return p.RecruitedArmies > 0
}
