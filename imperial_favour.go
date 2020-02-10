package confucius

import (
	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"golang.org/x/net/context"
)

func (g *Game) imperialFavourPhase(ctx context.Context) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	g.Phase = ImperialFavour
	g.ChiefMinister().ActionCubes += 1
	g.ChiefMinister().clearActions()
	g.SetCurrentPlayerers(g.ChiefMinister())
	g.ActionSpaces[ImperialFavourSpace].returnActionCubes()
	for _, player := range g.Players() {
		player.Passed = false
	}
	return
}
