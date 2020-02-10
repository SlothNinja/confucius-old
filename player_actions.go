package confucius

import (
	"bitbucket.org/SlothNinja/slothninja-games/sn"
	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"golang.org/x/net/context"
)

func (g *Game) validatePlayerAction(ctx context.Context) (cbs int, err error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	var sid SpaceID
	if sid, err = g.getSpaceID(ctx); err != nil {
		return
	}

	if cbs, err = g.validatePlaceCubesFor(sid); err != nil {
		return
	}

	c := restful.GinFrom(ctx)
	switch a, cp := c.PostForm("action"), g.CurrentPlayer(); {
	case !g.CUserIsCPlayerOrAdmin(ctx):
		err = sn.NewVError("Only the current player may perform the player action %q.", a)
	case (a == "pass" || IsEmperorRewardAction(a)) && g.Phase != Actions:
		err = sn.NewVError("You cannot perform a %q action during the %s phase.", a, g.PhaseName())
	case !g.inActionsOrImperialFavourPhase():
		err = sn.NewVError("You cannot perform a %q action during the %s phase.", a, g.PhaseName())
	case cp.Passed:
		err = sn.NewVError("You cannot perform a player action after passing.")
	}
	return
}

func (g *Game) validatePlaceCubesFor(id SpaceID) (cbs int, err error) {
	cp := g.CurrentPlayer()
	if cbs = cp.RequiredCubesFor(id); !cp.hasEnoughCubesFor(id) {
		err = sn.NewVError("You must have at least %d Action Cubes to perform this action.", cbs)
	}
	return
}

func IsEmperorRewardAction(s string) bool {
	return s == "Take Cash" || s == "Take Gift" || s == "Extra Action" || s == "Bribery Reward" ||
		s == "Avenge Emperor" || s == "Take Army"
}
