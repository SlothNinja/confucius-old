package confucius

import (
	"net/http"
	"time"

	"bitbucket.org/SlothNinja/slothninja-games/sn"
	"bitbucket.org/SlothNinja/slothninja-games/sn/contest"
	"bitbucket.org/SlothNinja/slothninja-games/sn/game"
	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"bitbucket.org/SlothNinja/slothninja-games/sn/user/stats"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

func Finish(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := restful.ContextFrom(c)
		log.Debugf(ctx, "Entering")
		defer log.Debugf(ctx, "Exiting")
		defer c.Redirect(http.StatusSeeOther, showPath(c, prefix))

		g := gameFrom(ctx)
		oldCP := g.CurrentPlayer()

		var (
			s   *stats.Stats
			cs  contest.Contests
			err error
		)

		switch g.Phase {
		case Actions:
			s, err = g.actionsPhaseFinishTurn(ctx)
		case ImperialFavour:
			s, cs, err = g.imperialFavourFinishTurn(ctx)
		case ChooseChiefMinister:
			s, err = g.chooseChiefMinisterPhaseFinishTurn(ctx)
		case Discard:
			s, cs, err = g.discardPhaseFinishTurn(ctx)
		case ImperialExamination:
			s, err = g.tutorStudentsPhaseFinishTurn(ctx)
		case ExaminationResolution:
			s, cs, err = g.examinationResolutionFinishTurn(ctx)
		case MinistryResolution:
			s, cs, err = g.ministryResolutionFinishTurn(ctx)
		default:
			err = sn.NewVError("Improper Phase for finishing turn.")
		}

		if err != nil {
			log.Errorf(ctx, err.Error())
			return
		}

		// cs != nil then game over
		if cs != nil {
			g.Phase = GameOver
			g.Status = game.Completed
			if err = g.save(ctx, wrap(s.GetUpdate(ctx, time.Time(g.UpdatedAt)), cs)...); err == nil {
				err = g.SendEndGameNotifications(ctx)
			}
		} else {
			if err = g.save(ctx, s.GetUpdate(ctx, time.Time(g.UpdatedAt))); err == nil {
				if newCP := g.CurrentPlayer(); newCP != nil && oldCP.ID() != newCP.ID() {
					err = g.SendTurnNotificationsTo(ctx, newCP)
				}
			}
		}

		if err != nil {
			log.Errorf(ctx, err.Error())
		}

		return
	}
}

func (g *Game) validateFinishTurn(ctx context.Context) (s *stats.Stats, err error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	cp := g.CurrentPlayer()

	switch s = stats.Fetched(ctx); {
	case !g.CUserIsCPlayerOrAdmin(ctx):
		err = sn.NewVError("Only the current player may finish a turn.")
	case !cp.PerformedAction:
		err = sn.NewVError("%s has yet to perform an action.", g.NameFor(cp))
	}
	return
}

// ps is an optional parameter.
// If no player is provided, assume current player.
func (g *Game) nextPlayer(ps ...*Player) *Player {
	i := game.IndexFor(g.CurrentPlayer(), g.Playerers) + 1
	if len(ps) == 1 {
		i = game.IndexFor(ps[0], g.Playerers) + 1
	}
	return g.Players()[i%g.NumPlayers]
}

func (p *Player) canAutoPass() bool {
	return p.canPass() && !p.canTransferInfluence() && !p.canEmperorReward()
}

func (g *Game) actionsPhaseFinishTurn(ctx context.Context) (s *stats.Stats, err error) {
	if s, err = g.validateFinishTurn(ctx); err != nil {
		return
	}

	cp := g.CurrentPlayer()
	restful.AddNoticef(ctx, "%s finished turn.", g.NameFor(cp))

	// Reveal Cards
	cp.ConCardHand.Reveal()
	cp.EmperorHand.Reveal()

	// Reset Extra Action
	g.ExtraAction = false

	if p := g.actionPhaseNextPlayer(); p != nil {
		g.SetCurrentPlayerers(p)
	} else {
		g.imperialFavourPhase(ctx)
	}
	return
}

func (g *Game) actionPhaseNextPlayer(players ...*Player) *Player {
	ps := g.Players()
	p := g.nextPlayer(players...)
	for !ps.allPassed() {
		if p.Passed {
			p = g.nextPlayer(p)
		} else {
			p.beginningOfTurnReset()
			if p.canAutoPass() {
				p.autoPass()
				p = g.nextPlayer(p)
			} else {
				return p
			}
		}
	}
	return nil
}

func (g *Game) imperialFavourFinishTurn(ctx context.Context) (s *stats.Stats, cs contest.Contests, err error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	if s, err = g.validateFinishTurn(ctx); err != nil {
		return
	}

	cp := g.CurrentPlayer()
	restful.AddNoticef(ctx, "%s finished turn.", g.NameFor(cp))

	// Reveal Cards
	cp.ConCardHand.Reveal()
	cp.EmperorHand.Reveal()

	g.buildWallPhase(ctx)
	if completed := g.examinationPhase(ctx); completed {
		if completed := g.ministryResolutionPhase(ctx, false); completed {
			g.invasionPhase(ctx)
			cs = g.endOfRoundPhase(ctx)
		}
	}
	return
}

func (g *Game) chooseChiefMinisterPhaseFinishTurn(ctx context.Context) (s *stats.Stats, err error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	if s, err = g.validateFinishTurn(ctx); err != nil {
		return
	}

	for _, player := range g.Players() {
		player.PerformedAction = false
	}
	g.SetCurrentPlayerers(g.nextPlayer(g.ChiefMinister()))
	g.actionsPhase(ctx)
	return
}

func (g *Game) examinationResolutionFinishTurn(ctx context.Context) (s *stats.Stats, cs contest.Contests, err error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	if s, err = g.validateFinishTurn(ctx); err != nil {
		return
	}

	// Place New Candidate
	if len(g.Candidates) > 0 {
		g.Candidates = g.Candidates[1:]
	}

	var i int
	for index, candidate := range g.Candidates {
		i = index
		if candidate.Playable() {
			break
		}
	}
	g.Candidates = g.Candidates[i:]
	if completed := g.ministryResolutionPhase(ctx, false); completed {
		g.invasionPhase(ctx)
		cs = g.endOfRoundPhase(ctx)
	}
	return
}
