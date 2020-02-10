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
	gob.RegisterName("*game.passEntry", new(passEntry))
	gob.RegisterName("*game.autoPassEntry", new(autoPassEntry))
}

func (g *Game) pass(ctx context.Context) (string, game.ActionType, error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	if _, err := g.validatePlayerAction(ctx); err != nil {
		return "", game.None, err
	}

	cp := g.CurrentPlayer()
	if err := cp.validatePass(ctx); err != nil {
		return "", game.None, err
	}

	cp.pass()

	// Create Action Object for logging
	e := cp.newPassEntry()

	// Set flash message
	restful.AddNoticef(ctx, string(e.HTML()))
	return "", game.Cache, nil
}

func (p *Player) pass() {
	// Pass
	p.Passed = true
	p.PerformedAction = true
}

type passEntry struct {
	*Entry
}

func (p *Player) newPassEntry() *passEntry {
	e := new(passEntry)
	g := p.Game()
	e.Entry = p.newEntry()
	p.Log = append(p.Log, e)
	g.Log = append(g.Log, e)
	return e
}

func (e *passEntry) HTML() template.HTML {
	return restful.HTML("%s passed.", e.Player().Name())
}

func (p *Player) autoPass() {
	p.pass()
	p.newAutoPassEntry()
}

type autoPassEntry struct {
	*Entry
}

func (p *Player) newAutoPassEntry() *autoPassEntry {
	e := new(autoPassEntry)
	g := p.Game()
	e.Entry = p.newEntry()
	p.Log = append(p.Log, e)
	g.Log = append(g.Log, e)
	return e
}

func (e *autoPassEntry) HTML() template.HTML {
	return restful.HTML("System auto passed for %s.", e.Player().Name())
}

func (p *Player) validatePass(ctx context.Context) (err error) {
	if _, err = p.Game().validatePlayerAction(ctx); err == nil && p.hasActionCubes() {
		err = sn.NewVError("You must use all of your action cubes before passing.")
	}
	return
}

func (g *Game) EnablePass(ctx context.Context) bool {
	cp := g.CurrentPlayer()
	return g.CUserIsCPlayerOrAdmin(ctx) && cp.canPass()
}

func (p *Player) canPass() bool {
	g := p.Game()
	return g.Phase == Actions && !p.PerformedAction && !p.Passed && !p.Game().ExtraAction && !p.hasActionCubes()
}
