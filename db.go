package confucius

import (
	"bitbucket.org/SlothNinja/slothninja-games/sn/game"
	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"bitbucket.org/SlothNinja/slothninja-games/sn/type"
	"github.com/gin-gonic/gin/binding"
	"go.chromium.org/gae/service/datastore"
	"golang.org/x/net/context"
)

const kind = "Game"

func New(ctx context.Context) *Game {
	g := new(Game)
	g.Header = game.NewHeader(ctx, g)
	g.State = newState()
	g.Parent = pk(ctx)
	g.Type = gType.Confucius
	return g
}

func newState() *State {
	return new(State)
}

func pk(ctx context.Context) *datastore.Key {
	return datastore.NewKey(ctx, gType.Confucius.SString(), "root", 0, game.GamesRoot(ctx))
}

func newKey(ctx context.Context, id int64) *datastore.Key {
	return datastore.NewKey(ctx, kind, "", id, pk(ctx))
}

func (g *Game) NewKey(ctx context.Context, id int64) *datastore.Key {
	return newKey(ctx, id)
}

func (g *Game) init(ctx context.Context) error {
	if err := g.Header.AfterLoad(g); err != nil {
		return err
	}

	for _, player := range g.Players() {
		player.init(g)
	}

	for _, entry := range g.Log {
		entry.Init(g)
	}

	for _, ministry := range g.Ministries {
		ministry.init(g)
	}

	for _, candidate := range g.Candidates {
		candidate.game = g
	}

	for _, land := range g.ForeignLands {
		land.init(g)
	}

	for _, land := range g.DistantLands {
		land.init(g)
	}
	return nil
}

func (g *Game) AfterCache() error {
	return g.init(g.CTX())
}

func (g *Game) fromForm(c context.Context) (err error) {
	log.Debugf(c, "Entering")
	defer log.Debugf(c, "Exiting")

	s := new(State)

	if err = restful.BindWith(c, s, binding.FormPost); err == nil {
		g.BasicGame = s.BasicGame
		g.AdmiralVariant = s.AdmiralVariant
	}
	return
}
