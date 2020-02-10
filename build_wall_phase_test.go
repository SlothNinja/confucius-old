package confucius

import (
	"bitbucket.org/SlothNinja/slothninja-games/sn/game"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"golang.org/x/net/context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("g.buildWallPhase", func() {
	var (
		ctx     context.Context
		g, oldG *Game
		err     error
	)

	BeforeEach(func() {
		ctx = restful.ContextFrom(tgc)
		g, err = get(ctx, 1)
		Ω(err).Should(BeNil())
		Ω(g.Status).Should(Equal(game.Running))

		oldG, err = get(ctx, 1)
		Ω(err).Should(BeNil())
		Ω(g.Status).Should(Equal(game.Running))
	})

	JustBeforeEach(func() {
		g.buildWallPhase(ctx)
	})

	It("should set phase correctly.", func() {
		Ω(g.Phase).Should(Equal(BuildWall))
	})

	It("should increase wall", func() {
		Ω(g.Wall).Should(Equal(oldG.Wall + 1))
	})

})
