package confucius

import (
	"net/http/httptest"
	"net/url"

	"bitbucket.org/SlothNinja/slothninja-games/sn/game"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"golang.org/x/net/context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("g.commercial", func() {
	var (
		ctx                               context.Context
		g                                 *Game
		gc                                *GiftCard
		cp                                *Player
		tmpl                              string
		action                            game.ActionType
		err                               error
		form                              url.Values
		serr                              string
		hl, dl, gl, pl, sc, cbs, gjs, pjs int
		tc                                bool
	)

	BeforeEach(func() {
		ctx = restful.ContextFrom(tgc)
		g, err = get(ctx, 1)
		cp = g.CurrentPlayer()
		gc = cp.GetGift(3)
		hl = len(cp.ConCardHand)
		dl = len(g.ConDiscardPile)
		pl = len(cp.Log)
		gl = len(g.Log)
		gjs = g.Junks
		pjs = cp.Junks
		tc = cp.TakenCommercial

		sc = g.ActionSpaces[CommercialSpace].Cubes[cp.ID()]
		cbs = cp.ActionCubes

		form = make(url.Values)
		form.Set("action", "commercial")
		form.Set("commercial-coins1", "1")
		form.Set("commercial-coins2", "0")
		form.Set("commercial-coins3", "1")

		restful.SetRequest(ctx, httptest.NewRequest("POST", "/confucius/game/", nil)).PostForm = form
	})

	JustBeforeEach(func() {
		tmpl, action, err = g.commercial(ctx)
	})

	It("should generate commercial income", func() {
		// Confirm game state update
		Ω(cp.PerformedAction).Should(BeTrue())
		Ω(cp.TakenCommercial).Should(BeTrue())
		Ω(g.ActionSpaces[CommercialSpace].Cubes[cp.ID()]).Should(Equal(sc + 1))
		Ω(cp.ActionCubes).Should(Equal(cbs - 1))
		Ω(cp.ConCardHand).Should(HaveLen(hl + 3))
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
		Ω(cp.TakenCommercial).Should(Equal(tc))
		Ω(g.ActionSpaces[CommercialSpace].Cubes[cp.ID()]).Should(Equal(sc))
		Ω(cp.ActionCubes).Should(Equal(cbs))
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
			form.Set("commercial-coins1", "z")
			restful.RequestFrom(ctx).PostForm = form
			serr = "Invalid value for Coin 1 cards received."
		})

		It("should error and not update game status", func() { ShouldError() })
	})

	Context("not in actions phase", func() {
		BeforeEach(func() {
			g.Phase = BuildWall
			serr = `You cannot perform a "commercial" action during the Build Wall phase.`
		})
		It("should error and not update game status", func() { ShouldError() })

	})

	Context("already taken commercial income", func() {
		BeforeEach(func() {
			tc = true
			cp.TakenCommercial = tc
			serr = "You have already taken the commercial income action this round."
		})
		It("should error and not update game status", func() { ShouldError() })

	})

	Context("paid too many coins", func() {
		BeforeEach(func() {
			form.Set("commercial-coins2", "1")
			restful.RequestFrom(ctx).PostForm = form
			serr = `You may only pay up to 4 coins. You paid 6 coins.`
		})
		It("should error and not update game status", func() { ShouldError() })

	})
})

var _ = Describe("g.EnableCommercial", func() {
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

		restful.SetRequest(ctx, httptest.NewRequest("POST", "/confucius/game/", nil))
	})

	JustBeforeEach(func() {
		b = g.EnableCommercial(ctx)
	})

	It("enable commericial income", func() {
		Ω(b).Should(BeTrue())
	})

	Context("with no money", func() {
		BeforeEach(func() {
			cp.ConCardHand = nil
		})

		It("disables commericial income", func() {
			Ω(b).Should(BeFalse())
		})
	})
})
