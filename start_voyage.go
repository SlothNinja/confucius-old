package confucius

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"strconv"

	"bitbucket.org/SlothNinja/slothninja-games/sn"
	"bitbucket.org/SlothNinja/slothninja-games/sn/game"
	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"golang.org/x/net/context"
)

func init() {
	gob.RegisterName("*game.startVoyageEntry", new(startVoyageEntry))
}

func (g *Game) startVoyage(ctx context.Context) (string, game.ActionType, error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	// Get Junks and Cards
	junks, cards, cubes, err := g.validateStartVoyage(ctx)
	if err != nil {
		return "", game.None, err
	}

	cp := g.CurrentPlayer()
	cp.PerformedAction = true

	// Place Action Cubes
	cp.PlaceCubesIn(JunksVoyageSpace, cubes)

	// Sail Junks
	completedVoyages := (cp.OnVoyage + junks) / 5
	cp.OnVoyage = (cp.OnVoyage + junks) % 5
	cp.Junks -= junks
	g.Junks += completedVoyages * 5

	lands := DistantLands{}
	points := []int{}
	emperorCards := make([]bool, completedVoyages)
	// For Each Completed Voyage Find Index of Land with Max VP
	for j := 0; j < completedVoyages; j++ {
		var max int
		land := new(DistantLand)
		for _, l := range g.DistantLands {
			if l.Chit.Value() > max {
				max = l.Chit.Value()
				land = l
			}
		}

		// All Tiles Taken.  Find First Land without current player
		if max == 0 {
			for _, l := range g.DistantLands {
				if !l.Players().Include(cp) {
					land = l
					break
				}
			}
		}

		scored := 0
		if land.Chit != NoChit {
			scored = land.Chit.Value()
		}
		cp.Score += scored
		points = append(points, scored)

		if len(g.EmperorDeck) > 0 {
			card := g.EmperorDeck.Draw()
			cp.EmperorHand.Append(card)
			emperorCards[j] = true
		}
		land.Chit = NoChit
		land.SetPlayers(append(land.Players(), cp))
		lands = append(lands, land)
	}

	// Move played cards from hand to discard pile
	cp.ConCardHand.Remove(cards...)
	g.ConDiscardPile.Append(cards...)

	// Create Action Object for logging
	e := cp.newStartVoyageEntry(cards, junks, lands, points, emperorCards)

	// Set flash message
	restful.AddNoticef(ctx, string(e.HTML()))
	return "", game.Cache, nil
}

type startVoyageEntry struct {
	*Entry
	Played       ConCards
	Junks        int
	DistantLands DistantLands
	MultiPoints  []int
	EmperorCards []bool
}

func (p *Player) newStartVoyageEntry(c ConCards, j int, l DistantLands, mp []int, ec []bool) *startVoyageEntry {
	g := p.Game()
	e := new(startVoyageEntry)
	e.Entry = p.newEntry()
	e.Played = c
	e.Junks = j
	e.DistantLands = l
	e.MultiPoints = mp
	e.EmperorCards = ec
	p.Log = append(p.Log, e)
	g.Log = append(g.Log, e)
	return e
}

func (e *startVoyageEntry) HTML() template.HTML {
	length := len(e.Played)
	licenses := e.Played.Licenses()

	s := fmt.Sprintf("<div>%s spent %d Confucius %s having %d %s to send %d %s on a voyage.</div>", e.Player().Name(), length, pluralize("card", length), licenses, pluralize("license", licenses), e.Junks, pluralize("junk", e.Junks))
	for i, land := range e.DistantLands {
		if e.EmperorCards[i] {
			s += fmt.Sprintf("<div>%s completed voyage to %s, scored %d points, and received an Emperor Reward card.</div>", e.Player().Name(), land.Name(), e.MultiPoints[i])
		} else {
			s += fmt.Sprintf("<div>%s completed voyage to %s, scored %d points, and did not receive an Emperor Reward card.", e.Player().Name(), land.Name(), e.MultiPoints[i])
		}
	}
	return template.HTML(s)
}

func (g *Game) validateStartVoyage(ctx context.Context) (int, ConCards, int, error) {
	cubes, err := g.validatePlayerAction(ctx)
	if err != nil {
		return 0, nil, 0, err
	}

	cards, err := g.getConCards(ctx, "start-voyage")
	if err != nil {
		return 0, nil, 0, err
	}

	cp := g.CurrentPlayer()
	c := restful.GinFrom(ctx)
	junks, err := strconv.Atoi(c.PostForm("junks"))
	licenses := cards.Licenses()
	switch {
	case err != nil:
		return 0, nil, 0, sn.NewVError("Invalid value for junks received.")
	case licenses < junks:
		return 0, nil, 0, sn.NewVError("You selected cards having %d total licenses, but you need %d licenses to start a voyage with %d junks.", licenses, junks, junks)
	case cp.Junks < junks:
		return 0, nil, 0, sn.NewVError("You have selected %d junks for the voyage, buy only have %d junks available.", junks, cp.Junks)
	case !g.hasDistantLandFor(cp):
		return 0, nil, 0, sn.NewVError("There are no distant lands to which you can voyage.")
	}

	return junks, cards, cubes, nil
}

func (g *Game) EnableStartVoyage(ctx context.Context) bool {
	cp := g.CurrentPlayer()
	return g.CUserIsCPlayerOrAdmin(ctx) && cp.canStartVoyage()
}

func (p *Player) canStartVoyage() bool {
	g := p.Game()
	return g.inActionsOrImperialFavourPhase() && !p.PerformedAction && p.hasEnoughCubesFor(JunksVoyageSpace) &&
		p.hasJunks() && p.hasLicenses()
}
