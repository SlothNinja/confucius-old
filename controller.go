package confucius

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"bitbucket.org/SlothNinja/slothninja-games/sn"
	"bitbucket.org/SlothNinja/slothninja-games/sn/codec"
	"bitbucket.org/SlothNinja/slothninja-games/sn/color"
	"bitbucket.org/SlothNinja/slothninja-games/sn/contest"
	"bitbucket.org/SlothNinja/slothninja-games/sn/game"
	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"bitbucket.org/SlothNinja/slothninja-games/sn/mlog"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"bitbucket.org/SlothNinja/slothninja-games/sn/type"
	"bitbucket.org/SlothNinja/slothninja-games/sn/user"
	"bitbucket.org/SlothNinja/slothninja-games/sn/user/stats"
	"github.com/gin-gonic/gin"
	"go.chromium.org/gae/service/datastore"
	"go.chromium.org/gae/service/info"
	"go.chromium.org/gae/service/memcache"
	"golang.org/x/net/context"
)

//type Action func(*Game, url.Values) (string, game.ActionType, error)
//
//var actionMap = map[string]Action{
//	"bribe-official":          bribeOfficial,
//	"secure-official":         secureOfficial,
//	"buy-gift":                buyGift,
//	"give-gift":               giveGift,
//	"nominate-student":        nominateStudent,
//	"force-exam":              forceExam,
//	"transfer-influence":      transferInfluence,
//	"temp-transfer-influence": tempTransfer,
//	"move-junks":              moveJunks,
//	"replace-student":         replaceStudent,
//	"swap-officials":          swapOfficials,
//	"redeploy-army":           redeployArmy,
//	"replace-influence":       replaceInfluence,
//	"place-student":           placeStudent,
//	"buy-junks":               buyJunks,
//	"start-voyage":            startVoyage,
//	"commercial":              commercial,
//	"tax-income":              taxIncome,
//	"recruit-army":            recruitArmy,
//	"invade-land":             invadeLand,
//	"no-action":               noAction,
//	"pass":                    pass,
//	"take-cash":               takeCash,
//	"take-gift":               takeGift,
//	"take-extra-action":       takeExtraAction,
//	"take-bribery-reward":     takeBriberyReward,
//	"avenge-emperor":          avengeEmperor,
//	"take-army":               takeArmy,
//	"discard":                 discard,
//	"choose-chief-minister":   chooseChiefMinister,
//	"tutor-student":           tutorStudent,
//	"reset":                   resetTurn,
//	"finish":                  finishTurn,
//	"game-state":              adminState,
//	"player":                  adminPlayer,
//	"ministry":                adminMinstry,
//	"official":                adminMinstryOfficial,
//	"candidate":               adminCandidate,
//	"foreign-land":            adminForeignLand,
//	"foreign-land-box":        adminForeignLandBox,
//	"action-space":            adminActionSpace,
//	"invoke-invade-phase":     invokeInvadePhase,
//	"distant-land":            adminDistantLand,
//}

const (
	gameKey   = "Game"
	homePath  = "/"
	jsonKey   = "JSON"
	statusKey = "Status"
)

func gameFrom(ctx context.Context) (g *Game) {
	g, _ = ctx.Value(gameKey).(*Game)
	return
}

func withGame(c *gin.Context, g *Game) *gin.Context {
	c.Set(gameKey, g)
	return c
}

func jsonFrom(ctx context.Context) (g *Game) {
	g, _ = ctx.Value(jsonKey).(*Game)
	return
}

func withJSON(c *gin.Context, g *Game) *gin.Context {
	c.Set(jsonKey, g)
	return c
}

func (g *Game) Update(ctx context.Context) (tmpl string, t game.ActionType, err error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	c := restful.GinFrom(ctx)
	switch a := c.PostForm("action"); a {
	case "bribe-official":
		return g.bribeOfficial(ctx)
	case "secure-official":
		return g.secureOfficial(ctx)
	case "buy-gift":
		return g.buyGift(ctx)
	case "give-gift":
		return g.giveGift(ctx)
	case "nominate-student":
		return g.nominateStudent(ctx)
	case "force-exam":
		return g.forceExam(ctx)
	case "transfer-influence":
		return g.transferInfluence(ctx)
	case "temp-transfer-influence":
		return g.tempTransfer(ctx)
	case "move-junks":
		return g.moveJunks(ctx)
	case "replace-student":
		return g.replaceStudent(ctx)
	case "swap-officials":
		return g.swapOfficials(ctx)
	case "redeploy-army":
		return g.redeployArmy(ctx)
	case "replace-influence":
		return g.replaceInfluence(ctx)
	case "place-student":
		return g.placeStudent(ctx)
	case "buy-junks":
		return g.buyJunks(ctx)
	case "start-voyage":
		return g.startVoyage(ctx)
	case "commercial":
		return g.commercial(ctx)
	case "tax-income":
		return g.taxIncome(ctx)
	case "recruit-army":
		return g.recruitArmy(ctx)
	case "invade-land":
		return g.invadeLand(ctx)
	case "no-action":
		return g.noAction(ctx)
	case "pass":
		return g.pass(ctx)
	case "take-cash":
		return g.takeCash(ctx)
	case "take-gift":
		return g.takeGift(ctx)
	case "take-extra-action":
		return g.takeExtraAction(ctx)
	case "take-bribery-reward":
		return g.takeBriberyReward(ctx)
	case "avenge-emperor":
		return g.avengeEmperor(ctx)
	case "take-army":
		return g.takeArmy(ctx)
	case "discard":
		return g.discard(ctx)
	case "choose-chief-minister":
		return g.chooseChiefMinister(ctx)
	case "tutor-student":
		return g.tutorStudent(ctx)
	case "reset":
		return g.resetTurn(ctx)
	case "game-state":
		return g.adminHeader(ctx)
	case "player":
		return g.adminPlayer(ctx)
		//	case "ministry":
		//		return g.adminMinstry(c)
		//	case "official":
		//		return g.adminMinstryOfficial(c)
		//	case "candidate":
		//		return g.adminCandidate(c)
		//	case "foreign-land":
		//		return g.adminForeignLand(c)
		//	case "foreign-land-box":
		//		return g.adminForeignLandBox(c)
		//	case "action-space":
		//		return g.adminActionSpace(c)
		//	case "invoke-invade-phase":
		//		return g.invokeInvadePhase(c)
		//	case "distant-land":
		//		return g.adminDistantLand(c)
	default:
		return "confucius/flash_notice", game.None, sn.NewVError("%v is not a valid action.", a)
	}
}

func Show(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := restful.ContextFrom(c)
		log.Debugf(ctx, "Entering")
		defer log.Debugf(ctx, "Exiting")

		g := gameFrom(ctx)
		cu := user.CurrentFrom(ctx)
		c.HTML(http.StatusOK, prefix+"/show", gin.H{
			"Context":    ctx,
			"VersionID":  info.VersionID(ctx),
			"CUser":      cu,
			"Game":       g,
			"IsAdmin":    user.IsAdmin(ctx),
			"Admin":      game.AdminFrom(ctx),
			"MessageLog": mlog.From(ctx),
			"ColorMap":   color.MapFrom(ctx),
			"Notices":    restful.NoticesFrom(ctx),
			"Errors":     restful.ErrorsFrom(ctx),
		})
	}
}

func Update(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := restful.ContextFrom(c)
		log.Debugf(ctx, "Entering")
		defer log.Debugf(ctx, "Exiting")

		g := gameFrom(ctx)
		if g == nil {
			log.Errorf(ctx, "Controller#Update Game Not Found")
			c.Redirect(http.StatusSeeOther, homePath)
			return
		}
		template, actionType, err := g.Update(ctx)
		log.Debugf(ctx, "template: %v actionType: %v err: %v", template, actionType, err)
		switch {
		case err != nil && sn.IsVError(err):
			restful.AddErrorf(ctx, "%v", err)
			withJSON(c, g)
		case err != nil:
			log.Errorf(ctx, err.Error())
			c.Redirect(http.StatusSeeOther, homePath)
			return
		case actionType == game.Cache:
			if err := g.cache(ctx); err != nil {
				restful.AddErrorf(ctx, "%v", err)
			}
		case actionType == game.Save:
			if err := g.save(ctx); err != nil {
				log.Errorf(ctx, "%s", err)
				restful.AddErrorf(ctx, "Controller#Update Save Error: %s", err)
				c.Redirect(http.StatusSeeOther, showPath(ctx, prefix))
				return
			}
		case actionType == game.Undo:
			mkey := g.UndoKey(ctx)
			if err := memcache.Delete(ctx, mkey); err != nil && err != memcache.ErrCacheMiss {
				log.Errorf(ctx, "memcache.Delete error: %s", err)
				c.Redirect(http.StatusSeeOther, showPath(ctx, prefix))
				return
			}
		}

		switch jData := jsonFrom(ctx); {
		case jData != nil && template == "json":
			c.JSON(http.StatusOK, jData)
		case template == "":
			c.Redirect(http.StatusSeeOther, showPath(ctx, prefix))
		default:
			cu := user.CurrentFrom(ctx)

			d := gin.H{
				"Context":   ctx,
				"VersionID": info.VersionID(ctx),
				"CUser":     cu,
				"Game":      g,
				"IsAdmin":   user.IsAdmin(ctx),
				"Notices":   restful.NoticesFrom(ctx),
				"Errors":    restful.ErrorsFrom(ctx),
			}
			c.HTML(http.StatusOK, template, d)
		}
	}
}

func (g *Game) save(ctx context.Context, es ...interface{}) (err error) {
	err = datastore.RunInTransaction(ctx, func(tc context.Context) (terr error) {
		oldG := New(tc)
		if ok := datastore.PopulateKey(oldG.Header, datastore.KeyForObj(tc, g.Header)); !ok {
			terr = fmt.Errorf("Unable to populate game with key.")
			return
		}

		if terr = datastore.Get(tc, oldG.Header); terr != nil {
			return
		}

		if oldG.UpdatedAt != g.UpdatedAt {
			terr = fmt.Errorf("Game state changed unexpectantly.  Try again.")
			return
		}

		if terr = g.encode(ctx); terr != nil {
			return
		}

		if terr = datastore.Put(tc, append(es, g.Header)); terr != nil {
			return
		}

		if terr = memcache.Delete(tc, g.UndoKey(tc)); terr == memcache.ErrCacheMiss {
			terr = nil
		}
		return
	}, &datastore.TransactionOptions{XG: true})
	return
}

func (g *Game) encode(ctx context.Context) (err error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	var encoded []byte
	if encoded, err = codec.Encode(g.State); err != nil {
		return
	}
	g.SavedState = encoded
	g.updateHeader()

	return
}

func (g *Game) cache(ctx context.Context) error {
	item := memcache.NewItem(ctx, g.UndoKey(ctx)).SetExpiration(time.Minute * 30)
	v, err := codec.Encode(g)
	if err != nil {
		return err
	}
	item.SetValue(v)
	return memcache.Set(ctx, item)
}

func wrap(s *stats.Stats, cs contest.Contests) (es []interface{}) {
	es = make([]interface{}, len(cs)+1)
	es[0] = s
	for i, c := range cs {
		es[i+1] = c
	}
	return
}

func showPath(ctx context.Context, prefix string) string {
	c := restful.GinFrom(ctx)
	return fmt.Sprintf("/%s/game/show/%s", prefix, c.Param("hid"))
}

func recruitingPath(prefix string) string {
	return fmt.Sprintf("/%s/games/recruiting", prefix)
}

func newPath(prefix string) string {
	return fmt.Sprintf("/%s/game/new", prefix)
}

func newGamer(ctx context.Context) game.Gamer {
	return New(ctx)
}

func Undo(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := restful.ContextFrom(c)
		log.Debugf(ctx, "Entering")
		defer log.Debugf(ctx, "Exiting")
		c.Redirect(http.StatusSeeOther, showPath(ctx, prefix))

		g := gameFrom(ctx)
		if g == nil {
			log.Errorf(ctx, "Controller#Update Game Not Found")
			return
		}
		mkey := g.UndoKey(ctx)
		if err := memcache.Delete(ctx, mkey); err != nil && err != memcache.ErrCacheMiss {
			log.Errorf(ctx, "Controller#Undo Error: %s", err)
		}
	}
}

func EndRound(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := restful.ContextFrom(c)
		log.Debugf(ctx, "Entering")
		defer log.Debugf(ctx, "Exiting")
		c.Redirect(http.StatusSeeOther, showPath(ctx, prefix))

		g := gameFrom(ctx)
		if g == nil {
			log.Errorf(ctx, "game not found")
			return
		}
		g.endOfRoundPhase(ctx)
		if err := g.save(ctx); err != nil {
			log.Errorf(ctx, "cache error: %s", err.Error())
		}
		return
	}
}

func Index(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := restful.ContextFrom(c)
		log.Debugf(ctx, "Entering")
		defer log.Debugf(ctx, "Exiting")

		gs := game.GamersFrom(ctx)
		switch status := game.StatusFrom(ctx); status {
		case game.Recruiting:
			c.HTML(http.StatusOK, "shared/invitation_index", gin.H{
				"Context":   ctx,
				"VersionID": info.VersionID(ctx),
				"CUser":     user.CurrentFrom(ctx),
				"Games":     gs,
				"Type":      gType.Confucius.String(),
			})
		default:
			c.HTML(http.StatusOK, "shared/games_index", gin.H{
				"Context":   ctx,
				"VersionID": info.VersionID(ctx),
				"CUser":     user.CurrentFrom(ctx),
				"Games":     gs,
				"Type":      gType.Confucius.String(),
				"Status":    status,
			})
		}
	}
}
func NewAction(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := restful.ContextFrom(c)
		log.Debugf(ctx, "Entering")
		defer log.Debugf(ctx, "Exiting")

		g := New(ctx)
		withGame(c, g)
		if err := g.FromParams(ctx, gType.GOT); err != nil {
			log.Errorf(ctx, err.Error())
			c.Redirect(http.StatusSeeOther, recruitingPath(prefix))
			return
		}

		c.HTML(http.StatusOK, prefix+"/new", gin.H{
			"Context":   ctx,
			"VersionID": info.VersionID(ctx),
			"CUser":     user.CurrentFrom(ctx),
			"Game":      g,
		})
	}
}

func Create(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := restful.ContextFrom(c)

		log.Debugf(ctx, "Entering")
		defer log.Debugf(ctx, "Exiting")
		defer c.Redirect(http.StatusSeeOther, recruitingPath(prefix))

		g := New(ctx)
		withGame(c, g)

		var err error
		if err = g.FromParams(ctx, g.Type); err == nil {
			err = g.fromForm(ctx)
		}

		if err == nil {
			err = g.encode(ctx)
		}

		if err == nil {
			err = datastore.RunInTransaction(ctx, func(tc context.Context) (err error) {
				if err = datastore.Put(tc, g.Header); err != nil {
					return
				}

				m := mlog.New()
				m.ID = g.ID
				return datastore.Put(tc, m)

			}, &datastore.TransactionOptions{XG: true})
		}

		if err == nil {
			restful.AddNoticef(ctx, "<div>%s created.</div>", g.Title)
		} else {
			log.Errorf(ctx, err.Error())
		}
	}
}

func Accept(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := restful.ContextFrom(c)
		log.Debugf(ctx, "Entering")
		defer log.Debugf(ctx, "Exiting")
		defer c.Redirect(http.StatusSeeOther, recruitingPath(prefix))

		g := gameFrom(ctx)
		if g == nil {
			log.Errorf(ctx, "game not found")
			return
		}

		var (
			start bool
			err   error
		)

		u := user.CurrentFrom(ctx)
		if start, err = g.Accept(ctx, u); err == nil && start {
			err = g.Start(ctx)
		}

		if err == nil {
			err = g.save(ctx)
		}

		if err == nil && start {
			g.SendTurnNotificationsTo(ctx, g.CurrentPlayer())
		}

		if err != nil {
			log.Errorf(ctx, err.Error())
		}

	}
}

func Drop(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := restful.ContextFrom(c)
		log.Debugf(ctx, "Entering")
		defer log.Debugf(ctx, "Exiting")
		defer c.Redirect(http.StatusSeeOther, recruitingPath(prefix))

		g := gameFrom(ctx)
		if g == nil {
			log.Errorf(ctx, "game not found")
			return
		}

		var err error

		u := user.CurrentFrom(ctx)
		if err = g.Drop(u); err == nil {
			err = g.save(ctx)
		}

		if err != nil {
			log.Errorf(ctx, err.Error())
			restful.AddErrorf(ctx, err.Error())
		}

	}
}

func Fetch(c *gin.Context) {
	ctx := restful.ContextFrom(c)
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	var (
		id  int64
		err error
	)

	// create Gamer
	if id, err = strconv.ParseInt(c.Param("hid"), 10, 64); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	g := New(ctx)
	g.ID = id

	switch action := c.PostForm("action"); {
	case action == "reset":
		// pull from memcache/datastore
		// same as undo & !MultiUndo
		fallthrough
	case action == "undo":
		// pull from memcache/datastore
		if err = dsGet(ctx, g); err != nil {
			c.Redirect(http.StatusSeeOther, homePath)
			return
		}
	default:
		if user.CurrentFrom(ctx) != nil {
			// pull from memcache and return if successful; otherwise pull from datastore
			if err := mcGet(ctx, g); err == nil {
				return
			}
		}
		if err = dsGet(ctx, g); err != nil {
			c.Redirect(http.StatusSeeOther, homePath)
			return
		}
	}
}

// pull temporary game state from memcache.  Note may be different from value stored in datastore.
func mcGet(ctx context.Context, g *Game) error {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	mkey := g.GetHeader().UndoKey(ctx)
	item, err := memcache.GetKey(ctx, mkey)
	if err != nil {
		return err
	}

	if err := codec.Decode(g, item.Value()); err != nil {
		return err
	}

	if err := g.AfterCache(); err != nil {
		return err
	}

	color.WithMap(withGame(restful.GinFrom(ctx), g), g.ColorMapFor(user.CurrentFrom(ctx)))
	return nil
}

// pull game state from memcache/datastore.  returned memcache should be same as datastore.
func dsGet(ctx context.Context, g *Game) (err error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	switch err = datastore.Get(ctx, g.Header); {
	case err != nil:
		restful.AddErrorf(ctx, err.Error())
		return
	case g == nil:
		err = fmt.Errorf("Unable to get game for id: %v", g.ID)
		restful.AddErrorf(ctx, err.Error())
		return
	}

	s := newState()
	if err = codec.Decode(&s, g.SavedState); err != nil {
		restful.AddErrorf(ctx, err.Error())
		return
	} else {
		g.State = s
	}

	if err = g.init(ctx); err != nil {
		restful.AddErrorf(ctx, err.Error())
		return
	}

	cm := g.ColorMapFor(user.CurrentFrom(ctx))
	color.WithMap(withGame(restful.GinFrom(ctx), g), cm)
	return
}

func JSON(ctx context.Context) {
	c := restful.GinFrom(ctx)
	c.JSON(http.StatusOK, gameFrom(ctx))
}

func JSONIndexAction(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := restful.ContextFrom(c)
		log.Debugf(ctx, "Entering")
		defer log.Debugf(ctx, "Exiting")

		game.JSONIndexAction(c)
	}
}

func (g *Game) updateHeader() {
	g.OptString = g.options()
	switch g.Phase {
	case GameOver:
		g.Progress = g.PhaseName()
	default:
		g.Progress = fmt.Sprintf("<div>Round: %d</div><div>Phase: %s</div>", g.Round, g.PhaseName())
	}
	if u := g.Creator; u != nil {
		g.CreatorSID = user.GenID(u.GoogleID)
		g.CreatorName = u.Name
	}

	if l := len(g.Users); l > 0 {
		g.UserSIDS = make([]string, l)
		g.UserNames = make([]string, l)
		g.UserEmails = make([]string, l)
		for i, u := range g.Users {
			g.UserSIDS[i] = user.GenID(u.GoogleID)
			g.UserNames[i] = u.Name
			g.UserEmails[i] = u.Email
		}
	}
}
