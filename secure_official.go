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
	gob.RegisterName("*game.secureOfficialEntry", new(secureOfficialEntry))
}

func (g *Game) secureOfficial(ctx context.Context) (string, game.ActionType, error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	cards, ministry, official, cubes, err := g.validateSecureOfficial(ctx)
	if err != nil {
		return "", game.None, err
	}

	// Place Action Cubes
	cp := g.CurrentPlayer()
	cp.PerformedAction = true
	cp.PlaceCubesIn(BribeSecureSpace, cubes)

	// Secure Marker On Official
	official.Secured = true

	// Move played cards from hand to discard pile
	cp.ConCardHand.Remove(cards...)
	g.ConDiscardPile.Append(cards...)

	// Create Action Object for logging
	e := cp.newSecureOfficialEntry(ministry, official, cards)

	// Set flash message
	restful.AddNoticef(ctx, string(e.HTML()))
	return "", game.Cache, nil
}

type secureOfficialEntry struct {
	*Entry
	MinistryName string
	Seniority    Seniority
	Played       ConCards
}

func (p *Player) newSecureOfficialEntry(m *Ministry, o *OfficialTile, c ConCards) *secureOfficialEntry {
	g := p.Game()
	e := new(secureOfficialEntry)
	e.Entry = p.newEntry()
	e.MinistryName = m.Name()
	e.Seniority = o.Seniority
	e.Played = c
	p.Log = append(p.Log, e)
	g.Log = append(g.Log, e)
	return e
}

func (e *secureOfficialEntry) HTML() template.HTML {
	length := len(e.Played)
	return restful.HTML("%s spent %d %s having %d coins to secure %s official having level %d seniority.",
		e.Player().Name(), length, pluralize("card", length), e.Played.Coins(), e.MinistryName, e.Seniority)
}

func (g *Game) validateSecureOfficial(ctx context.Context) (ConCards, *Ministry, *OfficialTile, int, error) {
	cubes, err := g.validatePlayerAction(ctx)
	if err != nil {
		return nil, nil, nil, 0, err
	}

	cards, err := g.getConCards(ctx, "secure-official")
	if err != nil {
		return nil, nil, nil, 0, err
	}

	ministry, official, err := g.getMinistryAndOfficial(ctx, "secure-official")
	if err != nil {
		return nil, nil, nil, 0, err
	}

	cp := g.CurrentPlayer()
	coinValue := cards.Coins()
	cost := cp.CostFor(official)

	switch {
	case official.Player() == nil:
		return nil, nil, nil, 0, sn.NewVError("You must select an official with a marker.")
	case official.Player().NotEqual(cp):
		return nil, nil, nil, 0, sn.NewVError("You must have a marker on the official before securing it.")
	case official.Secured:
		return nil, nil, nil, 0, sn.NewVError("You must select an official without a secured marker.")
	case coinValue < cost:
		return nil, nil, nil, 0, sn.NewVError("You selected cards having %d total coins, but you need %d coins to secure the selected official.", coinValue, cost)
	}

	return cards, ministry, official, cubes, nil
}

func (g *Game) EnableSecureOfficial(ctx context.Context) bool {
	cp := g.CurrentPlayer()
	return g.CUserIsCPlayerOrAdmin(ctx) && cp.canSecureAnOfficial()
}

func (p *Player) canSecureAnOfficial() bool {
	g := p.Game()
	return g.inActionsOrImperialFavourPhase() && !p.PerformedAction && p.hasEnoughCubesFor(BribeSecureSpace) &&
		g.hasSecurableOfficialFor(p)
}

func (g *Game) hasSecurableOfficialFor(p *Player) bool {
	for _, m := range g.Ministries {
		if !m.Resolved {
			for _, o := range m.Officials {
				if p.canSecure(o) {
					return true
				}
			}
		}
	}
	return false
}

func (p *Player) canSecure(o *OfficialTile) bool {
	return !o.Secured && p.hasBribed(o) && p.canAffordToSecure(o)
}

func (p *Player) hasBribed(o *OfficialTile) bool {
	return p.Equal(o.Player())
}

func (p *Player) canAffordToSecure(o *OfficialTile) bool {
	return p.canAffordToBribe(o)
}
