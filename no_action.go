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
	gob.RegisterName("*game.noActionEntry", new(noActionEntry))
}

func (g *Game) noAction(ctx context.Context) (string, game.ActionType, error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	cubes, err := g.validatePlayerAction(ctx)
	if err != nil {
		return "", game.None, err
	}

	cp := g.CurrentPlayer()
	cp.PerformedAction = true

	// Place Action Cube In NoActionSpace
	cp.PlaceCubesIn(NoActionSpace, cubes)

	// Create Action Object for logging
	entry := cp.newNoActionEntry()

	// Set flash message
	restful.AddNoticef(ctx, string(entry.HTML()))
	return "", game.Cache, nil
}

type noActionEntry struct {
	*Entry
}

func (p *Player) newNoActionEntry() *noActionEntry {
	g := p.Game()
	e := new(noActionEntry)
	e.Entry = p.newEntry()
	p.Log = append(p.Log, e)
	g.Log = append(g.Log, e)
	return e
}

func (g *noActionEntry) HTML() template.HTML {
	return restful.HTML("%s performed no action.", g.Player().Name())
}

func (g *Game) EnableNoAction(ctx context.Context) bool {
	cp := g.CurrentPlayer()
	return g.inActionsOrImperialFavourPhase() && g.CurrentPlayer() != nil &&
		!cp.PerformedAction && g.CUserIsCPlayerOrAdmin(ctx) &&
		cp.hasEnoughCubesFor(NoActionSpace) && cp.hasActionCubes()
}

func (p *Player) hasActionCubes() bool {
	return p.ActionCubes > 0
}
