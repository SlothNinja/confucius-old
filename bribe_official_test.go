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

var _ = Describe("g.bribeOfficial", func() {
	var (
		ctx              context.Context
		g, og            *Game
		cp, op, ocp, oop *Player
		tmpl             string
		action           game.ActionType
		err              error
		o                *OfficialTile
		m, om            *Ministry
		form             url.Values
		serr             string
	)

	BeforeEach(func() {
		ctx = restful.ContextFrom(tgc)
		g, err = get(ctx, 1)
		Ω(err).Should(BeNil())
		Ω(g.Status).Should(Equal(game.Running))
		Ω(g.Phase).Should(Equal(Actions))

		og, err = get(ctx, 1)
		Ω(err).Should(BeNil())
		Ω(g.Status).Should(Equal(game.Running))

		cp = g.CurrentPlayer()
		op = g.nextPlayer()
		ocp = og.CurrentPlayer()
		oop = og.nextPlayer()
		m = g.Ministries[Bingbu]
		om = og.Ministries[Bingbu]
		o = m.Officials[3]
		o.Cost = 3

		form = make(url.Values)
		form.Set("action", "bribe-official")
		form.Set("bribe-official-coins1", "1")
		form.Set("bribe-official-coins2", "1")
		form.Set("bribe-official-coins3", "0")
		form.Set("bribe-official", "Bingbu-3")
		restful.SetRequest(ctx, httptest.NewRequest("POST", "/confucius/game/", nil)).PostForm = form
	})

	JustBeforeEach(func() {
		tmpl, action, err = g.bribeOfficial(ctx)
	})

	It("should bribe official", func() {
		Ω(err).ShouldNot(HaveOccurred())
		// Confirm game state update
		Ω(cp.PerformedAction).Should(BeTrue())
		Ω(g.ActionSpaces[BribeSecureSpace].Cubes[cp.ID()]).Should(Equal(
			og.ActionSpaces[BribeSecureSpace].Cubes[cp.ID()] + 1))
		Ω(o.PlayerID).Should(Equal(cp.ID()))
		Ω(cp.ConCardHand).Should(HaveLen(len(ocp.ConCardHand) - 2))
		Ω(g.ConDiscardPile).Should(HaveLen(len(og.ConDiscardPile) + 2))

		// Confirm Game and Player Log updates
		Ω(g.Log).Should(HaveLen(len(og.Log) + 1))
		Ω(cp.Log).Should(HaveLen(len(ocp.Log) + 1))

		// Confirm Return values
		Ω(tmpl).Should(Equal(""))
		Ω(action).Should(Equal(game.Cache))
		Ω(err).ShouldNot(HaveOccurred())
	})

	ShouldError := func() {
		Ω(err).Should(HaveOccurred())
		// Confirm game state unchanged
		Ω(cp.PerformedAction).Should(BeFalse())
		Ω(g.ActionSpaces[BribeSecureSpace].Cubes[cp.ID()]).Should(Equal(
			og.ActionSpaces[BribeSecureSpace].Cubes[cp.ID()]))
		Ω(o.PlayerID).ShouldNot(Equal(cp.ID()))
		Ω(cp.ConCardHand).Should(HaveLen(len(ocp.ConCardHand)))
		Ω(g.ConDiscardPile).Should(HaveLen(len(og.ConDiscardPile)))

		// Confirm Game and Player Log did not update
		Ω(g.Log).Should(HaveLen(len(og.Log)))
		Ω(cp.Log).Should(HaveLen(len(ocp.Log)))

		// Confirm Return values
		Ω(tmpl).Should(Equal(""))
		Ω(action).Should(Equal(game.None))
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring(serr))
	}

	Context("incorrect value for coins1", func() {
		Context("value for coins1 not a number", func() {
			BeforeEach(func() {
				form.Set("bribe-official-coins1", "z")
				restful.RequestFrom(ctx).PostForm = form
				serr = "Invalid value for Coin 1 cards received."
			})

			It("should error and not update game status", func() { ShouldError() })
		})

		Context("value for coins1 not a positive integer", func() {
			BeforeEach(func() {
				form.Set("bribe-official-coins1", "-1")
				restful.RequestFrom(ctx).PostForm = form
				serr = "Invalid value for Coin 1 cards received."
			})

			It("should error and not update game status", func() { ShouldError() })
		})

		Context("value for coins1 is greater than what player has", func() {
			BeforeEach(func() {
				form.Set("bribe-official-coins1", "2")
				restful.RequestFrom(ctx).PostForm = form
				serr = "You selected 2 cards with one coin, but only have 1 of such cards."
			})

			It("should error and not update game status", func() { ShouldError() })
		})
	})

	Context("incorrect value for coins2", func() {
		Context("value for coins2 not a number", func() {
			BeforeEach(func() {
				form.Set("bribe-official-coins2", "z")
				restful.RequestFrom(ctx).PostForm = form
				serr = "Invalid value for Coin 2 cards received."
			})

			It("should error and not update game status", func() { ShouldError() })
		})

		Context("value for coins2 not a positive integer", func() {
			BeforeEach(func() {
				form.Set("bribe-official-coins2", "-1")
				restful.RequestFrom(ctx).PostForm = form
				serr = "Invalid value for Coin 2 cards received."
			})

			It("should error and not update game status", func() { ShouldError() })
		})

		Context("value for coins2 is greater than what player has", func() {
			BeforeEach(func() {
				form.Set("bribe-official-coins2", "2")
				restful.RequestFrom(ctx).PostForm = form
				serr = "You selected 2 cards with two coins, but only have 1 of such cards."
			})

			It("should error and not update game status", func() { ShouldError() })
		})
	})

	Context("incorrect value for coins3", func() {
		Context("value for coins3 not a number", func() {
			BeforeEach(func() {
				form.Set("bribe-official-coins3", "z")
				restful.RequestFrom(ctx).PostForm = form
				serr = "Invalid value for Coin 3 cards received."
			})

			It("should error and not update game status", func() { ShouldError() })
		})

		Context("value for coins3 not a positive integer", func() {
			BeforeEach(func() {
				form.Set("bribe-official-coins3", "-1")
				restful.RequestFrom(ctx).PostForm = form
				serr = "Invalid value for Coin 3 cards received."
			})

			It("should error and not update game status", func() { ShouldError() })
		})

		Context("value for coins2 is greater than what player has", func() {
			BeforeEach(func() {
				form.Set("bribe-official-coins3", "2")
				restful.RequestFrom(ctx).PostForm = form
				serr = "You selected 2 cards with three coins, but only have 1 of such cards."
			})

			It("should error and not update game status", func() { ShouldError() })
		})
	})

	Context("invalid ministry/seniority param", func() {
		Context("param is 'None'", func() {
			BeforeEach(func() {
				form.Set("bribe-official", "None")
				restful.RequestFrom(ctx).PostForm = form
				serr = "You must select an official."
			})

			It("should error and not update game status", func() { ShouldError() })
		})

		Context("param is missing", func() {
			BeforeEach(func() {
				delete(form, "bribe-official")
				restful.RequestFrom(ctx).PostForm = form
				serr = "Invalid format for ministry/seniority param."
			})

			It("should error and not update game status", func() { ShouldError() })
		})

		Context("param is missing '-'", func() {
			BeforeEach(func() {
				form.Set("bribe-official", "Bingbu3")
				restful.RequestFrom(ctx).PostForm = form
				serr = "Invalid format for ministry/seniority param."
			})

			It("should error and not update game status", func() { ShouldError() })
		})

		Context("param provides incorrect seniority", func() {
			BeforeEach(func() {
				form.Set("bribe-official", "Bingbu-z")
				restful.RequestFrom(ctx).PostForm = form
				serr = "Invalid Official Seniority Provided."
			})

			It("should error and not update game status", func() { ShouldError() })
		})

		Context("param provides incorrect ministry", func() {
			BeforeEach(func() {
				form.Set("bribe-official", "z-3")
				restful.RequestFrom(ctx).PostForm = form
				serr = "Invalid Ministry Provided."
			})

			It("should error and not update game status", func() { ShouldError() })
		})

		Context("param provides invalid official", func() {
			BeforeEach(func() {
				form.Set("bribe-official", "Hubu-0")
				restful.RequestFrom(ctx).PostForm = form
				serr = "Invalid official selected."
			})

			It("should error and not update game status", func() { ShouldError() })
		})
	})

	Context("current user not current player", func() {
		BeforeEach(func() {
			cp = g.nextPlayer()
			g.SetCurrentPlayerers(cp)
			op = g.nextPlayer(op)
			serr = `Only the current player may perform the player action "bribe-official".`
		})

		It("should error and not update game status", func() { ShouldError() })
	})

	Context("current player has already passed", func() {
		BeforeEach(func() {
			cp.Passed = true
			serr = `You cannot perform a player action after passing.`
		})

		It("should error and not update game status", func() { ShouldError() })
	})

	Context("not in actions phase", func() {
		BeforeEach(func() {
			g.Phase = BuildWall
			serr = `You cannot perform a "bribe-official" action during the Build Wall phase.`
		})
		It("should error and not update game status", func() { ShouldError() })

		Context("try to pass", func() {
			BeforeEach(func() {
				form.Set("action", "pass")
				restful.RequestFrom(ctx).PostForm = form
				serr = `You cannot perform a "pass" action during the Build Wall phase.`
			})
			It("should error and not update game status", func() { ShouldError() })
		})

	})

	Context("current player does not have sufficient cubes", func() {
		BeforeEach(func() {
			cp.ActionCubes = 0
			serr = `You must have at least 1 Action Cubes to perform this action.`
		})

		It("should error and not update game status", func() { ShouldError() })
	})

	Context("invalid action", func() {
		BeforeEach(func() {
			form.Set("action", "brib-official")
			restful.RequestFrom(ctx).PostForm = form
			serr = `"brib-official" is an invalid action.`
		})

		It("should error and not update game status", func() { ShouldError() })
	})

	Context("has gift from player", func() {
		var (
			gift *GiftCard
		)

		BeforeEach(func() {
			gift = op.GiftsBought[0]
			cp.GiftsReceived = append(cp.GiftsReceived, gift)
		})

		Context("player has same influence in ministry", func() {
			BeforeEach(func() {
				m.Officials[4].PlayerID = op.ID()
				m.Officials[5].PlayerID = cp.ID()
				serr = fmt.Sprintf("You have a gift obligation to %s that prevents you from bribing another official in the Bingbu ministry", g.NameFor(op))
			})

			It("should error and not update game status", func() { ShouldError() })
		})

		Context("player has less influence in ministry", func() {
			BeforeEach(func() {
				g.placeNewOfficialIn(m)
				m.Officials[4].PlayerID = op.ID()
				m.Officials[5].PlayerID = cp.ID()
				m.Officials[1].PlayerID = cp.ID()
				serr = fmt.Sprintf("You have a gift obligation to %s that prevents you from bribing another official in the Bingbu ministry", g.NameFor(op))
			})

			It("should error and not update game status", func() { ShouldError() })
		})
	})

	Context("official already has a marker", func() {
		BeforeEach(func() {
			o.PlayerID = op.ID()
			serr = `You can't bribe an official that already has a marker.`
		})

		It("should error and not update game status", func() { ShouldError() })
	})

	Context("official cost too much", func() {
		BeforeEach(func() {
			o.Cost = 4
			serr = `You selected cards having 3 total coins, but you need 4 coins to bribe the selected official.`
		})

		It("should error and not update game status", func() { ShouldError() })
	})
})

var _ = Describe("g.EnableBribeOfficial", func() {
	var (
		ctx              context.Context
		g, og            *Game
		cp, op, ocp, oop *Player
		err              error
		o                *OfficialTile
		m, om            *Ministry
		form             url.Values
		b                bool
	)

	BeforeEach(func() {
		ctx = restful.ContextFrom(tgc)
		g, err = get(ctx, 1)
		Ω(err).Should(BeNil())
		Ω(g.Status).Should(Equal(game.Running))
		Ω(g.Phase).Should(Equal(Actions))

		og, err = get(ctx, 1)
		Ω(err).Should(BeNil())
		Ω(g.Status).Should(Equal(game.Running))

		cp = g.CurrentPlayer()
		op = g.nextPlayer()
		ocp = og.CurrentPlayer()
		oop = og.nextPlayer()
		m = g.Ministries[Bingbu]
		om = og.Ministries[Bingbu]
		o = m.Officials[3]
		o.Cost = 3

		form = make(url.Values)
		form.Set("action", "bribe-official")
		form.Set("bribe-official-coins1", "1")
		form.Set("bribe-official-coins2", "1")
		form.Set("bribe-official-coins3", "0")
		form.Set("bribe-official", "Bingbu-3")

		restful.SetRequest(ctx, httptest.NewRequest("POST", "/confucius/game/", nil)).PostForm = form
	})

	JustBeforeEach(func() {
		b = g.EnableBribeOfficial(ctx)
	})

	It("enable bribing of official", func() {
		Ω(b).Should(BeTrue())
	})

	Context("all ministries resolved", func() {
		BeforeEach(func() {
			for _, m := range g.Ministries {
				m.Resolved = true
			}
		})

		It("disables bribing of official", func() {
			Ω(b).Should(BeFalse())
		})
	})
})
