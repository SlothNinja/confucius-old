package confucius

import (
	"fmt"
	"net/http/httptest"
	"net/url"

	"bitbucket.org/SlothNinja/slothninja-games/sn/game"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"golang.org/x/net/context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("g.takeCash", func() {
	var (
		ctx   context.Context
		g     *Game
		eCard *EmperorCard
		cp    *Player
		tmpl  string
		act   game.ActionType
		err   error
		form  url.Values
		serr  string
		cHand, eDisc, eHand,
		cDisc, gLog, pLog int
	)

	BeforeEach(func() {
		ctx = restful.ContextFrom(tgc)
		g, err = get(ctx, 1)
		cp = g.CurrentPlayer()

		cHand = len(cp.ConCardHand)
		cDisc = len(g.ConDiscardPile)
		eDisc = len(g.EmperorDiscard)

		pLog = len(cp.Log)
		gLog = len(g.Log)

		eCard = NewEmperorCard(Cash)
		cp.EmperorHand.Append(eCard)
		eHand = len(cp.EmperorHand)

		form = make(url.Values)
		form.Set("action", "take-cash")
		form.Set("reward-card", fmt.Sprintf("%v", eCard.Type))
		restful.SetRequest(ctx, httptest.NewRequest("POST", "/confucius/game/", nil)).PostForm = form
	})

	JustBeforeEach(func() {
		tmpl, act, err = g.takeCash(ctx)
	})

	It("should take cash", func() {
		// Confirm game state update
		Ω(cp.PerformedAction).Should(BeTrue())
		Ω(cp.ConCardHand).Should(HaveLen(cHand + 4))
		Ω(cp.EmperorHand).Should(HaveLen(eHand - 1))
		Ω(g.EmperorDiscard).Should(HaveLen(eDisc + 1))

		// Confirm Game and Player Log updates
		Ω(g.Log).Should(HaveLen(gLog + 1))
		Ω(cp.Log).Should(HaveLen(pLog + 1))

		// Confirm Return values
		Ω(tmpl).Should(Equal(""))
		Ω(act).Should(Equal(game.Cache))
		Ω(err).ShouldNot(HaveOccurred())
	})

	ShouldError := func() {
		// Confirm game state update
		Ω(cp.PerformedAction).Should(BeFalse())
		Ω(cp.ConCardHand).Should(HaveLen(cHand))
		Ω(cp.EmperorHand).Should(HaveLen(eHand))
		Ω(g.EmperorDiscard).Should(HaveLen(eDisc))

		// Confirm Game and Player Log updates
		Ω(g.Log).Should(HaveLen(gLog))
		Ω(cp.Log).Should(HaveLen(pLog))

		// Confirm Return values
		Ω(tmpl).Should(Equal(""))
		Ω(act).Should(Equal(game.None))
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring(serr))
	}

	Context("not in actions phase", func() {
		BeforeEach(func() {
			g.Phase = BuildWall
			serr = `You cannot perform a "take-cash" action during the Build Wall phase.`
		})
		It("should error and not update game status", func() { ShouldError() })

	})

	Context("incorrect value for reward-card", func() {
		BeforeEach(func() {
			form.Set("reward-card", "z")
			restful.RequestFrom(ctx).PostForm = form
			serr = `You must select an Emperor's Reward card.`
		})
		It("should error and not update game status", func() { ShouldError() })

	})

	Context("played incorrect card", func() {
		BeforeEach(func() {
			cd := NewEmperorCard(FreeGift)
			cp.EmperorHand.Append(cd)
			eHand = len(cp.EmperorHand)

			form.Set("reward-card", fmt.Sprintf("%v", FreeGift))
			restful.RequestFrom(ctx).PostForm = form
			serr = `You did not play the correct emperor's reward card for the selected action.`
		})
		It("should error and not update game status", func() { ShouldError() })

	})

	Context("do not have selected card", func() {
		BeforeEach(func() {
			cp.EmperorHand = nil
			eHand = len(cp.EmperorHand)
			serr = `You don't have the selected Emperor's Reward card.`
		})
		It("should error and not update game status", func() { ShouldError() })

	})

})

var _ = Describe("g.takeGift", func() {
	var (
		ctx   context.Context
		g     *Game
		gCard *GiftCard
		eCard *EmperorCard
		cp    *Player
		tmpl  string
		act   game.ActionType
		err   error
		form  url.Values
		serr  string
		gcHand, gcBought, eDisc, eHand,
		cDisc, gLog, pLog int
	)

	BeforeEach(func() {
		ctx = restful.ContextFrom(tgc)
		g, err = get(ctx, 1)
		cp = g.CurrentPlayer()

		gcHand = len(cp.GiftCardHand)
		gcBought = len(cp.GiftsBought)
		cDisc = len(g.ConDiscardPile)
		eDisc = len(g.EmperorDiscard)

		pLog = len(cp.Log)
		gLog = len(g.Log)

		eCard = NewEmperorCard(FreeGift)
		cp.EmperorHand.Append(eCard)
		eHand = len(cp.EmperorHand)

		gCard = cp.GetGift(Junk)

		form = make(url.Values)
		form.Set("action", "take-gift")
		form.Set("reward-card", fmt.Sprintf("%v", eCard.Type))
		form.Set("take-gift", fmt.Sprintf("%d", gCard.Value))
		restful.SetRequest(ctx, httptest.NewRequest("POST", "/confucius/game/", nil)).PostForm = form
	})

	JustBeforeEach(func() {
		tmpl, act, err = g.takeGift(ctx)
	})

	It("should take gift", func() {
		// Confirm game state update
		Ω(cp.PerformedAction).Should(BeTrue())
		Ω(cp.GiftCardHand).Should(HaveLen(gcHand - 1))
		Ω(cp.GiftsBought).Should(HaveLen(gcBought + 1))
		Ω(cp.EmperorHand).Should(HaveLen(eHand - 1))
		Ω(g.EmperorDiscard).Should(HaveLen(eDisc + 1))

		// Confirm Game and Player Log updates
		Ω(g.Log).Should(HaveLen(gLog + 1))
		Ω(cp.Log).Should(HaveLen(pLog + 1))

		// Confirm Return values
		Ω(tmpl).Should(Equal(""))
		Ω(act).Should(Equal(game.Cache))
		Ω(err).ShouldNot(HaveOccurred())
	})

	ShouldError := func() {
		// Confirm game state update
		Ω(cp.PerformedAction).Should(BeFalse())
		Ω(cp.GiftCardHand).Should(HaveLen(gcHand))
		Ω(cp.GiftsBought).Should(HaveLen(gcBought))
		Ω(cp.EmperorHand).Should(HaveLen(eHand))
		Ω(g.EmperorDiscard).Should(HaveLen(eDisc))

		// Confirm Game and Player Log updates
		Ω(g.Log).Should(HaveLen(gLog))
		Ω(cp.Log).Should(HaveLen(pLog))

		// Confirm Return values
		Ω(tmpl).Should(Equal(""))
		Ω(act).Should(Equal(game.None))
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring(serr))
	}

	Context("not in actions phase", func() {
		BeforeEach(func() {
			g.Phase = BuildWall
			serr = `You cannot perform a "take-gift" action during the Build Wall phase.`
		})
		It("should error and not update game status", func() { ShouldError() })

	})

	Context("incorrect value for reward-card", func() {
		BeforeEach(func() {
			form.Set("reward-card", "z")
			restful.RequestFrom(ctx).PostForm = form
			serr = `You must select an Emperor's Reward card.`
		})
		It("should error and not update game status", func() { ShouldError() })

	})

	Context("played incorrect card", func() {
		BeforeEach(func() {
			cd := NewEmperorCard(Cash)
			cp.EmperorHand.Append(cd)
			eHand = len(cp.EmperorHand)

			form.Set("reward-card", fmt.Sprintf("%v", Cash))
			restful.RequestFrom(ctx).PostForm = form
			serr = `You did not play the correct emperor's reward card for the selected action.`
		})
		It("should error and not update game status", func() { ShouldError() })

	})

	Context("do not have selected card", func() {
		BeforeEach(func() {
			cp.EmperorHand = nil
			eHand = len(cp.EmperorHand)
			serr = `You don't have the selected Emperor's Reward card.`
		})
		It("should error and not update game status", func() { ShouldError() })
	})

	Context("invalid gift selected", func() {
		BeforeEach(func() {
			form.Set("take-gift", "z")
			restful.RequestFrom(ctx).PostForm = form
			serr = `You must select an gift card.`
		})
		It("should error and not update game status", func() { ShouldError() })
	})

	Context("gift not available", func() {
		BeforeEach(func() {
			form.Set("take-gift", fmt.Sprintf("%d", Hanging))
			restful.RequestFrom(ctx).PostForm = form
			serr = `Selected gift card is not available.`
		})
		It("should error and not update game status", func() { ShouldError() })
	})
})

var _ = Describe("g.takeArmy", func() {
	var (
		ctx   context.Context
		g     *Game
		eCard *EmperorCard
		cp    *Player
		tmpl  string
		act   game.ActionType
		err   error
		form  url.Values
		serr  string
		armies, rArmies, eDisc,
		eHand, gLog, pLog int
	)

	BeforeEach(func() {
		ctx = restful.ContextFrom(tgc)
		g, err = get(ctx, 1)
		cp = g.CurrentPlayer()

		eDisc = len(g.EmperorDiscard)

		pLog = len(cp.Log)
		gLog = len(g.Log)

		eCard = NewEmperorCard(RecruitFreeArmy)
		cp.EmperorHand.Append(eCard)
		eHand = len(cp.EmperorHand)

		armies = cp.Armies
		rArmies = cp.RecruitedArmies

		form = make(url.Values)
		form.Set("action", "take-army")
		form.Set("reward-card", fmt.Sprintf("%v", eCard.Type))
		restful.SetRequest(ctx, httptest.NewRequest("POST", "/confucius/game/", nil)).PostForm = form
	})

	JustBeforeEach(func() {
		tmpl, act, err = g.takeArmy(ctx)
	})

	It("should take army", func() {
		// Confirm game state update
		Ω(cp.PerformedAction).Should(BeTrue())
		Ω(cp.EmperorHand).Should(HaveLen(eHand - 1))
		Ω(g.EmperorDiscard).Should(HaveLen(eDisc + 1))
		Ω(cp.Armies).Should(Equal(armies - 1))
		Ω(cp.RecruitedArmies).Should(Equal(rArmies + 1))

		// Confirm Game and Player Log updates
		Ω(g.Log).Should(HaveLen(gLog + 1))
		Ω(cp.Log).Should(HaveLen(pLog + 1))

		// Confirm Return values
		Ω(tmpl).Should(Equal(""))
		Ω(act).Should(Equal(game.Cache))
		Ω(err).ShouldNot(HaveOccurred())
	})

	ShouldError := func() {
		// Confirm game state update
		Ω(cp.PerformedAction).Should(BeFalse())
		Ω(cp.EmperorHand).Should(HaveLen(eHand))
		Ω(g.EmperorDiscard).Should(HaveLen(eDisc))
		Ω(cp.Armies).Should(Equal(armies))
		Ω(cp.RecruitedArmies).Should(Equal(rArmies))

		// Confirm Game and Player Log updates
		Ω(g.Log).Should(HaveLen(gLog))
		Ω(cp.Log).Should(HaveLen(pLog))

		// Confirm Return values
		Ω(tmpl).Should(Equal(""))
		Ω(act).Should(Equal(game.None))
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring(serr))
	}

	Context("not in actions phase", func() {
		BeforeEach(func() {
			g.Phase = BuildWall
			serr = `You cannot perform a "take-army" action during the Build Wall phase.`
		})
		It("should error and not update game status", func() { ShouldError() })

	})

	Context("incorrect value for reward-card", func() {
		BeforeEach(func() {
			form.Set("reward-card", "z")
			restful.RequestFrom(ctx).PostForm = form
			serr = `You must select an Emperor's Reward card.`
		})
		It("should error and not update game status", func() { ShouldError() })

	})

	Context("played incorrect card", func() {
		BeforeEach(func() {
			cd := NewEmperorCard(Cash)
			cp.EmperorHand.Append(cd)
			eHand = len(cp.EmperorHand)

			form.Set("reward-card", fmt.Sprintf("%v", Cash))
			restful.RequestFrom(ctx).PostForm = form
			serr = `You did not play the correct emperor's reward card for the selected action.`
		})
		It("should error and not update game status", func() { ShouldError() })

	})

	Context("do not have selected card", func() {
		BeforeEach(func() {
			cp.EmperorHand = nil
			eHand = len(cp.EmperorHand)
			serr = `You don't have the selected Emperor's Reward card.`
		})
		It("should error and not update game status", func() { ShouldError() })
	})

	Context("do not have armies to recruit", func() {
		BeforeEach(func() {
			cp.Armies = 0
			armies = cp.Armies
			serr = `You have no armies to recruit.`
		})
		It("should error and not update game status", func() { ShouldError() })
	})

})

var _ = Describe("g.takeExtraAction", func() {
	var (
		ctx         context.Context
		g           *Game
		eCard       *EmperorCard
		cp          *Player
		tmpl        string
		act         game.ActionType
		err         error
		form        url.Values
		serr        string
		exAct, pAct bool
		pEmpHand, emDisc,
		gLog, pLog int
	)

	BeforeEach(func() {
		ctx = restful.ContextFrom(tgc)
		g, err = get(ctx, 1)
		cp = g.CurrentPlayer()

		emDisc = len(g.EmperorDiscard)
		exAct = g.ExtraAction

		pLog = len(cp.Log)
		gLog = len(g.Log)

		eCard = NewEmperorCard(ExtraAction)
		cp.EmperorHand.Append(eCard)
		pEmpHand = len(cp.EmperorHand)

		pAct = cp.PerformedAction

		form = make(url.Values)
		form.Set("action", "take-extra-action")
		form.Set("reward-card", fmt.Sprintf("%v", eCard.Type))
		restful.SetRequest(ctx, httptest.NewRequest("POST", "/confucius/game/", nil)).PostForm = form
	})

	JustBeforeEach(func() {
		tmpl, act, err = g.takeExtraAction(ctx)
	})

	It("should enable taking extra action", func() {
		// Confirm game state update
		Ω(cp.PerformedAction).Should(Equal(pAct))
		Ω(g.ExtraAction).Should(BeTrue())
		Ω(cp.EmperorHand).Should(HaveLen(pEmpHand - 1))
		Ω(g.EmperorDiscard).Should(HaveLen(emDisc + 1))

		// Confirm Game and Player Log updates
		Ω(g.Log).Should(HaveLen(gLog + 1))
		Ω(cp.Log).Should(HaveLen(pLog + 1))

		// Confirm Return values
		Ω(tmpl).Should(Equal(""))
		Ω(act).Should(Equal(game.Cache))
		Ω(err).ShouldNot(HaveOccurred())
	})

	ShouldError := func() {
		// Confirm game state update
		Ω(cp.PerformedAction).Should(Equal(pAct))
		Ω(g.ExtraAction).Should(BeFalse())
		Ω(cp.EmperorHand).Should(HaveLen(pEmpHand))
		Ω(g.EmperorDiscard).Should(HaveLen(emDisc))

		// Confirm Game and Player Log updates
		Ω(g.Log).Should(HaveLen(gLog))
		Ω(cp.Log).Should(HaveLen(pLog))

		// Confirm Return values
		Ω(tmpl).Should(Equal(""))
		Ω(act).Should(Equal(game.None))
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring(serr))
	}

	Context("not in actions phase", func() {
		BeforeEach(func() {
			g.Phase = BuildWall
			serr = `You cannot perform a "take-extra-action" action during the Build Wall phase.`
		})
		It("should error and not update game status", func() { ShouldError() })

	})

	Context("incorrect value for reward-card", func() {
		BeforeEach(func() {
			form.Set("reward-card", "z")
			restful.RequestFrom(ctx).PostForm = form
			serr = `You must select an Emperor's Reward card.`
		})
		It("should error and not update game status", func() { ShouldError() })

	})

	Context("played incorrect card", func() {
		BeforeEach(func() {
			cd := NewEmperorCard(Cash)
			cp.EmperorHand.Append(cd)
			pEmpHand = len(cp.EmperorHand)

			form.Set("reward-card", fmt.Sprintf("%v", Cash))
			restful.RequestFrom(ctx).PostForm = form
			serr = `You did not play the correct emperor's reward card for the selected action.`
		})
		It("should error and not update game status", func() { ShouldError() })

	})

	Context("do not have selected card", func() {
		BeforeEach(func() {
			cp.EmperorHand = nil
			pEmpHand = len(cp.EmperorHand)
			serr = `You don't have the selected Emperor's Reward card.`
		})
		It("should error and not update game status", func() { ShouldError() })
	})

})

var _ = Describe("g.avengeEmperor", func() {
	var (
		ctx         context.Context
		g           *Game
		eCard       *EmperorCard
		cp          *Player
		tmpl        string
		act         game.ActionType
		err         error
		form        url.Values
		serr        string
		exAct, pAct bool
		pEmpHand, emDisc, rArmies,
		score, gLog, pLog int
	)

	BeforeEach(func() {
		ctx = restful.ContextFrom(tgc)
		g, err = get(ctx, 1)
		cp = g.CurrentPlayer()

		emDisc = len(g.EmperorDiscard)
		exAct = g.ExtraAction

		pLog = len(cp.Log)
		gLog = len(g.Log)

		eCard = NewEmperorCard(EmperorInsulted)
		cp.EmperorHand.Append(eCard)
		pEmpHand = len(cp.EmperorHand)

		pAct = cp.PerformedAction

		cp.RecruitedArmies, rArmies = 1, 1
		score = cp.Score

		form = make(url.Values)
		form.Set("action", "avenge-emperor")
		form.Set("reward-card", fmt.Sprintf("%v", eCard.Type))
		restful.SetRequest(ctx, httptest.NewRequest("POST", "/confucius/game/", nil)).PostForm = form
	})

	JustBeforeEach(func() {
		tmpl, act, err = g.avengeEmperor(ctx)
	})

	It("should enable taking extra action", func() {
		// Confirm game state update
		Ω(cp.PerformedAction).Should(BeTrue())
		Ω(cp.RecruitedArmies).Should(Equal(rArmies - 1))
		Ω(cp.Score).Should(Equal(score + 2))
		Ω(cp.EmperorHand).Should(HaveLen(pEmpHand - 1))
		Ω(g.EmperorDiscard).Should(HaveLen(emDisc + 1))

		// Confirm Game and Player Log updates
		Ω(g.Log).Should(HaveLen(gLog + 1))
		Ω(cp.Log).Should(HaveLen(pLog + 1))

		// Confirm Return values
		Ω(tmpl).Should(Equal(""))
		Ω(act).Should(Equal(game.Cache))
		Ω(err).ShouldNot(HaveOccurred())
	})

	ShouldError := func() {
		Ω(err.Error()).Should(ContainSubstring(serr))
		// Confirm game state update
		Ω(cp.PerformedAction).Should(Equal(pAct))
		Ω(cp.RecruitedArmies).Should(Equal(rArmies))
		Ω(cp.Score).Should(Equal(score))
		Ω(cp.EmperorHand).Should(HaveLen(pEmpHand))
		Ω(g.EmperorDiscard).Should(HaveLen(emDisc))

		// Confirm Game and Player Log updates
		Ω(g.Log).Should(HaveLen(gLog))
		Ω(cp.Log).Should(HaveLen(pLog))

		// Confirm Return values
		Ω(tmpl).Should(Equal(""))
		Ω(act).Should(Equal(game.None))
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring(serr))
	}

	Context("not in actions phase", func() {
		BeforeEach(func() {
			g.Phase = BuildWall
			serr = `You cannot perform a "avenge-emperor" action during the Build Wall phase.`
		})
		It("should error and not update game status", func() { ShouldError() })

	})

	Context("incorrect value for reward-card", func() {
		BeforeEach(func() {
			form.Set("reward-card", "z")
			restful.RequestFrom(ctx).PostForm = form
			serr = `You must select an Emperor's Reward card.`
		})
		It("should error and not update game status", func() { ShouldError() })

	})

	Context("played incorrect card", func() {
		BeforeEach(func() {
			cd := NewEmperorCard(Cash)
			cp.EmperorHand.Append(cd)
			pEmpHand = len(cp.EmperorHand)

			form.Set("reward-card", fmt.Sprintf("%v", Cash))
			restful.RequestFrom(ctx).PostForm = form
			serr = `You did not play the correct emperor's reward card for the selected action.`
		})
		It("should error and not update game status", func() { ShouldError() })

	})

	Context("do not have selected card", func() {
		BeforeEach(func() {
			cp.EmperorHand = nil
			pEmpHand = len(cp.EmperorHand)
			serr = `You don't have the selected Emperor's Reward card.`
		})
		It("should error and not update game status", func() { ShouldError() })
	})

	Context("do not have a recruited army", func() {
		BeforeEach(func() {
			cp.RecruitedArmies, rArmies = 0, 0
			serr = `You have no recruited armies with which to avenge the Emperor.`
		})
		It("should error and not update game status", func() { ShouldError() })
	})
})
