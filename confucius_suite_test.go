package confucius

import (
	"fmt"
	"testing"

	"bitbucket.org/SlothNinja/slothninja-games/sn/codec"
	"bitbucket.org/SlothNinja/slothninja-games/sn/game"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"bitbucket.org/SlothNinja/slothninja-games/sn/user"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"go.chromium.org/gae/service/datastore"
	gUser "go.chromium.org/gae/service/user"
	"golang.org/x/net/context"
)

var (
	tgc *gin.Context
)

var gusers = [3]*gUser.User{
	&gUser.User{
		Email: "sloth@example.com",
		ID:    "1",
	},
	&gUser.User{
		Email: "ninja@example.com",
		ID:    "2",
	},
	&gUser.User{
		Email: "jerry@example.com",
		ID:    "3",
	},
}

func createTestContext() (c *gin.Context) {
	c, _, _ = gin.CreateTestContext()
	c = restful.WithContext(c, restful.CreateTestContext(c))
	return
}

func getUsers(ctx context.Context) (us user.Users, err error) {
	us = make(user.Users, len(gusers))
	var nu *user.NUser
	for i := range gusers {
		nu, err = user.ByGoogleID(ctx, gusers[i].ID)
		us[i] = user.New(ctx)
		us[i].ID, us[i].Data = nu.OldID, nu.Data
	}
	return
}

func createRunningGame(ctx context.Context) (g *Game, err error) {
	g = New(ctx)
	g.ID = 1
	g.NumPlayers = 3
	g.Title = "game title"
	g.CreatorID = 1
	start := false

	var us user.Users
	if us, err = getUsers(ctx); err != nil {
		return
	}

	for _, u := range us {
		if start, err = g.Accept(ctx, u); err != nil {
			return
		}
	}

	if !start {
		fmt.Errorf("expected start to be true.")
	}

	err = g.Start(ctx)

	// Set Current User to that of Current Player
	user.WithCurrent(tgc, g.CurrentPlayer().User())
	return
}

func store(ctx context.Context, g *Game) (err error) {
	var encoded []byte
	if encoded, err = codec.Encode(g.State); err != nil {
		return
	}

	g.SavedState = encoded

	return datastore.Put(ctx, g.Header)
}

func get(ctx context.Context, id int64) (g *Game, err error) {
	g = New(ctx)
	g.ID = id
	err = dsGet(ctx, g)
	return
}

func storeUsers(ctx context.Context) (err error) {
	us := make(user.Users, len(gusers))
	nus := make([]*user.NUser, len(gusers))
	for i := range gusers {
		us[i] = user.FromGUser(ctx, gusers[i])
		us[i].ID = int64(i + 1)
		nus[i] = user.ToNUser(ctx, us[i])
	}

	if err = datastore.Put(ctx, us); err != nil {
		return
	}

	if err = datastore.Put(ctx, nus); err != nil {
		return
	}

	// Default current user
	user.WithCurrent(tgc, us[0])

	return
}

func TestConfucius(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Confucius Suite")
}

var _ = BeforeSuite(func() {
	var (
		ctx context.Context
		g   *Game
		err error
	)

	Describe("create gin context", func() {
		tgc = createTestContext()
		ctx = restful.ContextFrom(tgc)
		Ω(tgc).ShouldNot(BeNil())
		Ω(ctx).ShouldNot(BeNil())
	})

	Describe("store users", func() {
		err = storeUsers(ctx)
		Ω(err).ShouldNot(HaveOccurred())
	})

	Describe("store running game", func() {
		g, err = createRunningGame(ctx)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(g.Status).Should(Equal(game.Running))

		err = store(ctx, g)
		Ω(err).ShouldNot(HaveOccurred())
	})
})

func GCEqual(expected *GiftCard) types.GomegaMatcher {
	return &gcEqualMatcher{
		expected: expected,
	}
}

type gcEqualMatcher struct {
	expected *GiftCard
}

func (matcher *gcEqualMatcher) Match(actual interface{}) (success bool, err error) {
	if gc, ok := actual.(*GiftCard); !ok {
		err = fmt.Errorf("GCEqual matcher expects a GiftCard")
	} else {
		success = matcher.expected.Equal(gc)
	}
	return
}

func (matcher *gcEqualMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nto equal \n\t%#v", actual, matcher.expected)
}

func (matcher *gcEqualMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nnot to equal \n\t%#v", actual, matcher.expected)
}
