package confucius

import (
	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"golang.org/x/net/context"
)

func (g *Game) returnActionCubesPhase(ctx context.Context) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	for _, s := range g.ActionSpaces {
		s.returnActionCubes()
	}
}

func (s *ActionSpace) returnActionCubes() {
	s.Cubes = make(Cubes)
}
