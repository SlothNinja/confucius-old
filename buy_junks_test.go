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

var _ = Describe("g.buyGift", func() {
	var (
		ctx                               context.Context
		g                                 *Game
		cp                                *Player
		tmpl                              string
		action                            game.ActionType
		err                               error
		form                              url.Values
		serr                              string
		hl, dl, gl, pl, sc, cbs, gjs, pjs int
	)

	BeforeEach(func() {
		ctx = restful.ContextFrom(tgc)
		g, err = get(ctx, 1)
		cp = g.CurrentPlayer()
		hl = len(cp.ConCardHand)
		dl = len(g.ConDiscardPile)
		pl = len(cp.Log)
		gl = len(g.Log)
		gjs = g.Junks
		pjs = cp.Junks

		sc = g.ActionSpaces[JunksVoyageSpace].Cubes[cp.ID()]
		cbs = cp.ActionCubes

		form = make(url.Values)
		form.Set("action", "buy-junks")
		form.Set("buy-junks-coins1", "1")
		form.Set("buy-junks-coins2", "1")
		form.Set("buy-junks-coins3", "0")
		form.Set("junks", "2")

		restful.SetRequest(ctx, httptest.NewRequest("POST", "/confucius/game/", nil)).PostForm = form
	})

	JustBeforeEach(func() {
		tmpl, action, err = g.buyJunks(ctx)
	})

	It("should buy junks", func() {
		// Confirm game state update
		Ω(cp.PerformedAction).Should(BeTrue())
		Ω(g.ActionSpaces[JunksVoyageSpace].Cubes[cp.ID()]).Should(Equal(sc + 1))
		Ω(cp.ActionCubes).Should(Equal(cbs - 1))
		Ω(cp.ConCardHand).Should(HaveLen(hl - 2))
		Ω(g.ConDiscardPile).Should(HaveLen(dl + 2))
		Ω(g.Junks).Should(Equal(gjs - 2))
		Ω(cp.Junks).Should(Equal(pjs + 2))

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
		Ω(g.ActionSpaces[JunksVoyageSpace].Cubes[cp.ID()]).Should(Equal(sc))
		Ω(cp.ActionCubes).Should(Equal(cbs))
		Ω(cp.ConCardHand).Should(HaveLen(hl))
		Ω(g.ConDiscardPile).Should(HaveLen(dl))
		Ω(g.Junks).Should(Equal(gjs))
		Ω(cp.Junks).Should(Equal(pjs))

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
			form.Set("buy-junks-coins1", "z")
			restful.RequestFrom(ctx).PostForm = form
			serr = "Invalid value for Coin 1 cards received."
		})

		It("should error and not update game status", func() { ShouldError() })
	})

	Context("not in actions phase", func() {
		BeforeEach(func() {
			g.Phase = BuildWall
			serr = `You cannot perform a "buy-junks" action during the Build Wall phase.`
		})
		It("should error and not update game status", func() { ShouldError() })

	})

	Context("incorrect junks value", func() {
		BeforeEach(func() {
			form.Set("junks", "z")
			restful.RequestFrom(ctx).PostForm = form
			serr = `Form value for "junks" is invalid.`
		})

		It("should error and not update game status", func() { ShouldError() })
	})

	Context("not enough junks available for purchase", func() {
		BeforeEach(func() {
			gjs = 1
			g.Junks = 1
			serr = `You selected more junks than there are available in stock.`
		})

		It("should error and not update game status", func() { ShouldError() })
	})

	Context("did not pay enough for junks", func() {
		BeforeEach(func() {
			form.Set("buy-junks-coins2", "0")
			restful.RequestFrom(ctx).PostForm = form
			serr = `You selected cards having 1 total coins, but you need 3 coins to buy the selected junks.`
		})

		It("should error and not update game status", func() { ShouldError() })
	})

})

var _ = Describe("g.EnableBuyJunks", func() {
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
		b = g.EnableBuyJunks(ctx)
	})

	It("enable buying junks", func() {
		Ω(b).Should(BeTrue())
	})

	Context("with no money", func() {
		BeforeEach(func() {
			cp.ConCardHand = nil
		})

		It("disables buying junks", func() {
			Ω(b).Should(BeFalse())
		})
	})
})
