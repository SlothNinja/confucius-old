package confucius

import (
	"encoding/gob"
	"html/template"

	"bitbucket.org/SlothNinja/slothninja-games/sn"
	"bitbucket.org/SlothNinja/slothninja-games/sn/game"
	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"bitbucket.org/SlothNinja/slothninja-games/sn/user/stats"
	"golang.org/x/net/context"
)

func init() {
	gob.RegisterName("*game.tutorStudentEntry", new(tutorStudentEntry))
}

func (g *Game) tutorStudent(ctx context.Context) (string, game.ActionType, error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	cards, player, err := g.validateTutorStudent(ctx)
	if err != nil {
		return "", game.None, err
	}

	cp := g.CurrentPlayer()
	e := cp.tutorStudent(cards, player, false)

	// Set flash message
	restful.AddNoticef(ctx, string(e.HTML()))
	return "", game.Cache, nil
}

func (p *Player) tutorStudent(cards ConCards, player *Player, auto bool) *tutorStudentEntry {
	g := p.Game()
	p.PerformedAction = true
	canceled := false

	if len(cards) > 0 {
		// Remove played cards from hand
		p.ConCardHand.Remove(cards...)
		g.ConDiscardPile.Append(cards...)

		// Apply cards to Candidate
		switch {
		case g.Candidate().Player().Equal(player):
			g.Candidate().PlayerCards.Append(cards...)
		default:
			g.Candidate().OtherPlayerCards.Append(cards...)
		}

		// Cancel Gift Obligation
		if len(cards) >= 3 && p.NotEqual(player) {
			canceled = (p.cancelGiftFrom(player) != nil)
		}
	}

	return p.newTutorStudentEntry(cards, player, canceled, auto)
}

type tutorStudentEntry struct {
	*Entry
	Played     ConCards
	CancelGift bool
	Auto       bool
}

func (p *Player) newTutorStudentEntry(cards ConCards, player *Player, canceled, auto bool) *tutorStudentEntry {
	e := new(tutorStudentEntry)
	g := p.Game()
	e.Entry = p.newEntry()
	e.Played = cards
	if player == nil {
		e.OtherPlayerID = NoPlayerID
	} else {
		e.OtherPlayerID = player.ID()
	}
	e.CancelGift = canceled
	e.Auto = auto
	p.Log = append(p.Log, e)
	g.Log = append(g.Log, e)
	return e
}

func (e *tutorStudentEntry) HTML() template.HTML {
	op := e.OtherPlayer()
	switch length := len(e.Played); {
	case e.OtherPlayer() != nil && !e.CancelGift:
		if e.Auto {
			return restful.HTML("%s auto-spent %d %s to tutor student of %s.",
				e.Player().Name(), length, pluralize("card", length), op.Name())
		} else {
			return restful.HTML("%s spent %d %s to tutor student of %s.",
				e.Player().Name(), length, pluralize("card", length), op.Name())
		}
	case e.OtherPlayer() != nil && e.CancelGift:
		name := op.Name()
		return restful.HTML("%s spent %d %s to tutor student of %s and canceled gift received from %s.",
			e.Player().Name(), length, pluralize("card", length), name, name)
	}
	return restful.HTML("%s has no cards to tutor a student.", e.Player().Name())
}

func (g *Game) validateTutorStudent(ctx context.Context) (cds ConCards, p *Player, err error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	if cds, err = g.getConCards(ctx, "tutor-student"); err != nil {
		return
	}

	c := restful.GinFrom(ctx)
	p = g.PlayerBySID(c.PostForm("player"))
	cp := g.CurrentPlayer()

	switch {
	case !cp.TutorPlayers().Include(p):
		err = sn.NewVError("You provided an incorrect player.")
	case !g.CUserIsCPlayerOrAdmin(ctx):
		err = sn.NewVError("Only the current player may pay to tutor a student.")
	case g.Phase != ImperialExamination:
		err = sn.NewVError("You cannot pay to tutor a student during the %s phase.", g.PhaseName())
	case len(cds) < 1 && len(cp.ConCardHand) > 0:
		err = sn.NewVError("You must play at least one Confucius Card.")
	}
	return
}

func (p *Player) TutorPlayers() Players {
	player := p.Game().Candidate().Player()
	otherPlayer := p.Game().Candidate().OtherPlayer()
	both_players := Players{}

	if player != nil {
		both_players = append(both_players, player)
	}

	if otherPlayer != nil {
		both_players = append(both_players, otherPlayer)
	}

	if len(both_players) != 2 {
		return both_players
	}

	var value GiftCardValue
	var ps Players
	for _, gift := range p.GiftsReceived {
		if gift.Player() != nil && (gift.Player().Equal(player) || gift.Player().Equal(otherPlayer)) {
			switch len(ps) {
			case 0:
				ps = append(ps, gift.Player())
				value = gift.Value
			case 1:
				switch {
				case gift.Value > value:
					return Players{gift.Player()}
				case gift.Value == value:
					return both_players
				}
			}
		}
	}

	if len(ps) == 1 {
		return ps
	}
	return both_players
}

func (p *Player) autoTutor(ctx context.Context) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	var student *Player
	if ps := p.TutorPlayers(); len(ps) == 1 {
		student = ps[0]
	}
	p.tutorStudent(p.ConCardHand, student, true)
}

func (p *Player) canAutoTutor() bool {
	l := len(p.ConCardHand)
	return (l == 1 && len(p.TutorPlayers()) == 1) || l == 0
}

func (g *Game) tutorStudentsPhaseFinishTurn(ctx context.Context) (s *stats.Stats, err error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	if s, err = g.validateFinishTurn(ctx); err != nil {
		return
	}

	if p := g.tutorStudentsPhaseNextPlayer(ctx); p != nil {
		g.SetCurrentPlayerers(p)
	} else {
		g.resolveExamination(ctx)
	}
	return
}

func (g *Game) tutorStudentsPhaseNextPlayer(ctx context.Context, ps ...*Player) (p *Player) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	for p = g.nextPlayer(ps...); !g.Players().allPerformedAction() && p.canAutoTutor(); p = g.nextPlayer() {
		g.SetCurrentPlayerers(p)
		p.autoTutor(ctx)
	}

	if p.PerformedAction {
		p = nil
	}
	return
}
