package confucius

import (
	"encoding/gob"
	"html/template"
	"strconv"

	"bitbucket.org/SlothNinja/slothninja-games/sn"
	"bitbucket.org/SlothNinja/slothninja-games/sn/game"
	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"golang.org/x/net/context"
)

func init() {
	gob.RegisterName("*game.chooseChiefMinisterEntry", new(chooseChiefMinisterEntry))
}

func (g *Game) chooseChiefMinisterPhase(ctx context.Context) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	g.Phase = ChooseChiefMinister
	if g.Round == 1 {
		g.RandomTurnOrder()
		g.SetChiefMinister(g.CurrentPlayer())
		g.ChiefMinister().PlaceCubesIn(ImperialFavourSpace, 1)
		g.SetCurrentPlayerers(g.nextPlayer())
		g.actionsPhase(ctx)
	} else {
		g.SetCurrentPlayerers(g.ChiefMinister())
	}
}

func (g *Game) chooseChiefMinister(ctx context.Context) (string, game.ActionType, error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	recipient, err := g.validateChooseChiefMinister(ctx)
	if err != nil {
		return "", game.None, err
	}

	// Appoint New ChiefMinister
	g.SetChiefMinister(recipient)
	g.ChiefMinister().PlaceCubesIn(ImperialFavourSpace, 1)

	// Clear Actions
	cp := g.CurrentPlayer()
	cp.clearActions()
	cp.PerformedAction = true

	// Create Action Object for logging
	e := cp.newChooseChiefMinisterEntry(recipient)

	// Set flash message
	restful.AddNoticef(ctx, string(e.HTML()))
	return "", game.Cache, nil
}

type chooseChiefMinisterEntry struct {
	*Entry
}

func (p *Player) newChooseChiefMinisterEntry(op *Player) *chooseChiefMinisterEntry {
	g := p.Game()
	e := new(chooseChiefMinisterEntry)
	e.Entry = p.newEntry()
	e.OtherPlayerID = op.ID()
	p.Log = append(p.Log, e)
	g.Log = append(g.Log, e)
	return e
}

func (e *chooseChiefMinisterEntry) HTML() template.HTML {
	return restful.HTML("%s chose %s to be chief minister.", e.Player().Name(), e.OtherPlayer().Name())
}

func (g *Game) validateChooseChiefMinister(ctx context.Context) (*Player, error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	c := restful.GinFrom(ctx)
	recipientID, err := strconv.Atoi(c.PostForm("player"))
	if err != nil {
		return nil, err
	}

	recipient := g.PlayerByID(recipientID)
	cp := g.CurrentPlayer()
	switch {
	case recipient == nil:
		return nil, sn.NewVError("Recipient not found.")
	case !g.CUserIsCPlayerOrAdmin(ctx):
		return nil, sn.NewVError("Only the current player may choose a chief minister.")
	case g.Phase != ChooseChiefMinister:
		return nil, sn.NewVError("You cannot choose a chief minister during the %s phase.", g.Phase)
	case cp.NotEqual(g.ChiefMinister()):
		return nil, sn.NewVError("Only the current chief minister may select the succeeding chief minister.")
	case cp.Equal(recipient):
		return nil, sn.NewVError("You cannot appoint yourself chief minister.")
	}
	return recipient, nil
}

func (g *Game) EnableChooseChiefMinister(ctx context.Context) bool {
	return g.CUserIsCPlayerOrAdmin(ctx) && g.Phase == ChooseChiefMinister
}
