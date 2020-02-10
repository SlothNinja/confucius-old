package confucius

import (
	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"golang.org/x/net/context"
)

func (g *Game) buildWallPhase(ctx context.Context) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	g.Phase = BuildWall
	g.beginningOfPhaseReset()
	g.Wall += 1
}
