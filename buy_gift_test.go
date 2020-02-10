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
		ctx                     context.Context
		g                       *Game
		gc                      *GiftCard
		cp                      *Player
		tmpl                    string
		action                  game.ActionType
		err                     error
		form                    url.Values
		serr                    string
		hl, dl, gl, pl, sc, cbs int
	)

	BeforeEach(func() {
		ctx = restful.ContextFrom(tgc)
		g, err = get(ctx, 1)
		cp = g.CurrentPlayer()
		gc = cp.GetGift(Vase)
		hl = len(cp.ConCardHand)
		dl = len(g.ConDiscardPile)
		pl = len(cp.Log)
		gl = len(g.Log)

		sc = g.ActionSpaces[BuyGiftSpace].Cubes[cp.ID()]
		cbs = cp.ActionCubes

		form = make(url.Values)
		form.Set("action", "buy-gift")
		form.Set("buy-gift-coins1", "1")
		form.Set("buy-gift-coins2", "1")
		form.Set("buy-gift-coins3", "0")
		form.Set("buy-gift", "3")
		restful.SetRequest(ctx, httptest.NewRequest("POST", "/confucius/game/", nil)).PostForm = form
	})

	JustBeforeEach(func() {
		tmpl, action, err = g.buyGift(ctx)
	})

	It("should buy gift", func() {
		// Confirm game state update
		Ω(cp.PerformedAction).Should(BeTrue())
		Ω(g.ActionSpaces[BuyGiftSpace].Cubes[cp.ID()]).Should(Equal(sc + 1))
		Ω(cp.ActionCubes).Should(Equal(cbs - 1))
		Ω(cp.ConCardHand).Should(HaveLen(hl - 2))
		Ω(g.ConDiscardPile).Should(HaveLen(dl + 2))
		Ω(cp.GiftsBought).Should(ContainElement(GCEqual(gc)))
		Ω(cp.GiftCardHand).ShouldNot(ContainElement(GCEqual(gc)))

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
		Ω(g.ActionSpaces[BuyGiftSpace].Cubes[cp.ID()]).Should(Equal(sc))
		Ω(cp.ActionCubes).Should(Equal(cbs))
		Ω(cp.ConCardHand).Should(HaveLen(hl))
		Ω(g.ConDiscardPile).Should(HaveLen(dl))
		Ω(cp.GiftsBought).ShouldNot(ContainElement(GCEqual(gc)))
		Ω(cp.GiftCardHand).Should(ContainElement(GCEqual(gc)))

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
			form.Set("buy-gift-coins1", "z")
			restful.RequestFrom(ctx).PostForm = form
			serr = "Invalid value for Coin 1 cards received."
		})

		It("should error and not update game status", func() { ShouldError() })
	})

	Context("not in actions phase", func() {
		BeforeEach(func() {
			g.Phase = BuildWall
			serr = `You cannot perform a "buy-gift" action during the Build Wall phase.`
		})
		It("should error and not update game status", func() { ShouldError() })

	})

	Context("incorrect gift value", func() {
		BeforeEach(func() {
			form.Set("buy-gift", "z")
			restful.RequestFrom(ctx).PostForm = form
			serr = `You must select an gift card.`
		})

		It("should error and not update game status", func() { ShouldError() })
	})

	Context("card of gift value not available for purchase", func() {
		BeforeEach(func() {
			form.Set("buy-gift", "1")
			restful.RequestFrom(ctx).PostForm = form
			serr = `You don't have a gift of value 1 to buy.`
		})

		It("should error and not update game status", func() { ShouldError() })
	})

	Context("card costs too much", func() {
		BeforeEach(func() {
			form.Set("buy-gift-coins2", "0")
			restful.RequestFrom(ctx).PostForm = form
			serr = `You selected cards having 1 total coins, but the Vase gift costs 3 coins.`
		})

		It("should error and not update game status", func() { ShouldError() })
	})

})

var _ = Describe("g.EnableBuyGift", func() {
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

		c := restful.GinFrom(ctx)
		c.Request = httptest.NewRequest("POST", "/confucius/game/", nil)
	})

	JustBeforeEach(func() {
		b = g.EnableBuyGift(ctx)
	})

	It("enable bribing of official", func() {
		Ω(b).Should(BeTrue())
	})

	Context("with no money", func() {
		BeforeEach(func() {
			cp.ConCardHand = nil
		})

		It("disables buying gift card", func() {
			Ω(b).Should(BeFalse())
		})
	})
})
