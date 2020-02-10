package confucius

import (
	"fmt"
	"net/http/httptest"
	"net/url"

	"bitbucket.org/SlothNinja/slothninja-games/sn/game"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"bitbucket.org/SlothNinja/slothninja-games/sn/user"
	"golang.org/x/net/context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("g.discard", func() {
	var (
		ctx            context.Context
		g              *Game
		cp             *Player
		tmpl           string
		action         game.ActionType
		err            error
		form           url.Values
		serr           string
		hl, dl, gl, pl int
	)

	BeforeEach(func() {
		ctx = restful.ContextFrom(tgc)
		g, err = get(ctx, 1)
		g.Phase = Discard
		cp = g.CurrentPlayer()
		pl = len(cp.Log)
		gl = len(g.Log)

		ncds := make(ConCards, 3)
		for i := range ncds {
			ncds[i] = g.DrawConCard()
		}
		cp.ConCardHand.Append(ncds...)

		hl = len(cp.ConCardHand)
		dl = len(g.ConDiscardPile)

		form = make(url.Values)
		form.Set("action", "discard")
		form.Set("discard-coins1", "1")
		form.Set("discard-coins2", "1")
		form.Set("discard-coins3", "0")

		restful.SetRequest(ctx, httptest.NewRequest("POST", "/confucius/game/", nil)).PostForm = form
	})

	JustBeforeEach(func() {
		tmpl, action, err = g.discard(ctx)
	})

	It("should discard Confucius cards", func() {
		// Confirm game state update
		Ω(cp.PerformedAction).Should(BeTrue())
		Ω(cp.ConCardHand).Should(HaveLen(hl - 2))
		Ω(g.ConDiscardPile).Should(HaveLen(dl + 2))

		// Confirm Game and Player Log updates
		Ω(g.Log).Should(HaveLen(gl + 1))
		Ω(cp.Log).Should(HaveLen(pl + 1))

		// Confirm Return values
		Ω(tmpl).Should(Equal(""))
		Ω(action).Should(Equal(game.Cache))
		Ω(err).ShouldNot(HaveOccurred())
	})

	ShouldError := func() {
		// Confirm game state update
		Ω(cp.PerformedAction).Should(BeFalse())
		Ω(cp.ConCardHand).Should(HaveLen(hl))
		Ω(g.ConDiscardPile).Should(HaveLen(dl))

		// Confirm Game and Player Log updates
		Ω(g.Log).Should(HaveLen(gl))
		Ω(cp.Log).Should(HaveLen(pl))

		// Confirm Return values
		Ω(tmpl).Should(Equal(""))
		Ω(action).Should(Equal(game.None))
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring(serr))
	}

	Context("incorrect value for Confucius Cards", func() {
		BeforeEach(func() {
			form.Set("discard-coins1", "z")
			restful.RequestFrom(ctx).PostForm = form
			serr = "Invalid value for Coin 1 cards received."
		})

		It("should error and not update game status", func() { ShouldError() })
	})

	Context("not in actions phase", func() {
		BeforeEach(func() {
			g.Phase = BuildWall
			serr = `You cannot discard cards during the Build Wall phase.`
		})
		It("should error and not update game status", func() { ShouldError() })

	})

	Context("not current player", func() {
		BeforeEach(func() {
			g.SetCurrentPlayerers(g.nextPlayer())
			serr = `Only a current player may discard cards.`
		})

		It("should error and not update game status", func() { ShouldError() })
	})

	Context("did not discard enough cards", func() {
		BeforeEach(func() {
			form.Set("discard-coins2", "0")
			restful.RequestFrom(ctx).PostForm = form
			serr = `You must discard down to 4 cards.  You have discarded to 5 cards.`
		})

		It("should error and not update game status", func() { ShouldError() })
	})

})

var _ = Describe("g.EnableDiscard", func() {
	var (
		ctx context.Context
		g   *Game
		cp  *Player
		err error
		b   bool
	)

	BeforeEach(func() {
		ctx = restful.ContextFrom(tgc)
		g, err = get(ctx, 1)
		cp = g.CurrentPlayer()
		g.Phase = Discard

		restful.SetRequest(ctx, httptest.NewRequest("POST", "/confucius/game/", nil))
	})

	JustBeforeEach(func() {
		b = g.EnableDiscard(ctx)
	})

	It("enable discarding", func() {
		Ω(b).Should(BeTrue())
	})

	Context("performed action", func() {
		BeforeEach(func() {
			cp.PerformedAction = true
		})

		It("disables discarding", func() {
			Ω(b).Should(BeFalse())
		})
	})
})

var _ = Describe("g.discardPhase", func() {
	var (
		ctx       context.Context
		g         *Game
		cp        *Player
		completed bool
		err       error
	)

	BeforeEach(func() {
		ctx = restful.ContextFrom(tgc)
		g, err = get(ctx, 1)
		cp = g.CurrentPlayer()
		g.Phase = Discard

		restful.SetRequest(ctx, httptest.NewRequest("POST", "/confucius/game/", nil))
	})

	JustBeforeEach(func() {
		completed = g.discardPhase(ctx)
	})

	Context("when player has too many cards", func() {
		var (
			p *Player
		)

		BeforeEach(func() {
			p = g.nextPlayer()
			ncds := make(ConCards, 3)
			for i := range ncds {
				ncds[i] = g.DrawConCard()
			}
			p.ConCardHand.Append(ncds...)
		})

		It("should add player to current players", func() {
			Ω(g.CurrentPlayerers()).Should(HaveLen(1))
			Ω(g.CurrentPlayer().Equal(p)).Should(BeTrue())
		})

		It("should not complete phase", func() {
			Ω(completed).Should(BeFalse())
		})
	})

	Context("when no player has too many cards", func() {
		var (
			p *Player
		)

		BeforeEach(func() {
			p = g.CurrentPlayer()
			p.PlaceCubesIn(BuyGiftSpace, 1)
			p.PlaceCubesIn(JunksVoyageSpace, 1)
		})

		It("should complete phase", func() {
			Ω(completed).Should(BeTrue())
		})
	})
})

var _ = Describe("g.discardPhaseFinishTurn", func() {
	var (
		ctx   context.Context
		g     *Game
		p, cp *Player
		err   error
		serr  string
	)

	BeforeEach(func() {
		ctx = restful.ContextFrom(tgc)
		g, err = get(ctx, 1)
		cp = g.CurrentPlayer()
		g.Phase = Discard
		p = g.CurrentPlayer()
		p.PlaceCubesIn(BuyGiftSpace, 1)
		p.PlaceCubesIn(JunksVoyageSpace, 1)
		p.PerformedAction = true

		restful.SetRequest(ctx, httptest.NewRequest("POST", "/confucius/game/", nil))
	})

	JustBeforeEach(func() {
		_, _, err = g.discardPhaseFinishTurn(ctx)
	})

	It("should finish turn", func() {
		// Remove current player
		Ω(g.CurrentPlayer().Equal(g.ChiefMinister())).Should(BeTrue())

		// Reset action cubes
		Ω(g.ActionSpaces[BuyGiftSpace].Cubes[p.ID()]).Should(BeZero())
		Ω(g.ActionSpaces[JunksVoyageSpace].Cubes[p.ID()]).Should(BeZero())

		// Confirm Return values
		Ω(err).ShouldNot(HaveOccurred())
	})

	ShouldError := func() {
		// Should not remove current player
		Ω(g.CurrentPlayer().ID()).Should(Equal(cp.ID()))

		// Should not reset action cubes
		Ω(g.ActionSpaces[BuyGiftSpace].Cubes[p.ID()]).Should(Equal(1))
		Ω(g.ActionSpaces[JunksVoyageSpace].Cubes[p.ID()]).Should(Equal(1))

		// Confirm Return values
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring(serr))
	}

	Context("has not performed action", func() {
		BeforeEach(func() {
			cp.PerformedAction = false
			serr = fmt.Sprintf("%s has yet to perform an action.", g.NameFor(cp))
		})

		It("should error and not update game status", func() { ShouldError() })
	})

	Context("is not current player", func() {
		BeforeEach(func() {
			p := g.nextPlayer()
			user.WithCurrent(tgc, p.User())
			serr = "Only the current player may finish a turn."
		})

		It("should error and not update game status", func() { ShouldError() })

		AfterEach(func() {
			user.WithCurrent(tgc, cp.User())
		})
	})
})
