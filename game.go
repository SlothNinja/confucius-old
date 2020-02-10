package confucius

import (
	"encoding/gob"
	"errors"
	"fmt"

	"bitbucket.org/SlothNinja/slothninja-games/sn"
	"bitbucket.org/SlothNinja/slothninja-games/sn/color"
	"bitbucket.org/SlothNinja/slothninja-games/sn/game"
	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"bitbucket.org/SlothNinja/slothninja-games/sn/type"
	"bitbucket.org/SlothNinja/slothninja-games/sn/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

func Register(t gType.Type, r *gin.Engine) {
	gob.Register(new(Game))
	game.Register(t, newGamer, PhaseNames, nil)
	AddRoutes(t.Prefix(), r)
}

var ErrMustBeGame = errors.New("Resource must have type *Game.")

type IDS []int64

func (g *Game) GetPlayerers() game.Playerers {
	return g.Playerers
}

type Game struct {
	*game.Header
	*State
}

type State struct {
	Playerers game.Playerers
	Log       game.GameLog
	Junks     int `form:"junks"`

	ChiefMinisterID int `form:"chief-minister-id"`
	AdmiralID       int `form:"admiral-id"`
	GeneralID       int `form:"general-id"`
	AvengerID       int `form:"avenger-id"`

	ActionSpaces ActionSpaces

	Candidates     CandidateTiles
	OfficialsDeck  OfficialsDeck
	ConDeck        ConCards
	ConDiscardPile ConCards
	EmperorDeck    EmperorCards
	EmperorDiscard EmperorCards

	DistantLands DistantLands
	ForeignLands ForeignLands

	Ministries Ministries

	Wall        int  `form:"wall"`
	ExtraAction bool `form:"extra-action"`

	BasicGame      bool `form:"basic-game"`
	AdmiralVariant bool `form:"admiral-variant"`
}

func (g *Game) ChiefMinister() *Player {
	if g.ChiefMinisterID != NoPlayerID {
		return g.PlayerByID(g.ChiefMinisterID)
	}
	return nil
}

func (g *Game) SetChiefMinister(p *Player) {
	switch {
	case p == nil:
		g.ChiefMinisterID = NoPlayerID
	default:
		g.ChiefMinisterID = p.ID()
	}
}

func (g *Game) Admiral() *Player {
	if g.AdmiralID != NoPlayerID {
		return g.PlayerByID(g.AdmiralID)
	}
	return nil
}

func (g *Game) SetAdmiral(p *Player) {
	switch {
	case p == nil:
		g.AdmiralID = NoPlayerID
	default:
		g.AdmiralID = p.ID()
	}
}

func (g *Game) General() *Player {
	if g.GeneralID != NoPlayerID {
		return g.PlayerByID(g.GeneralID)
	}
	return nil
}

func (g *Game) SetGeneral(p *Player) {
	switch {
	case p == nil:
		g.GeneralID = NoPlayerID
	default:
		g.GeneralID = p.ID()
	}
}

func (g *Game) Avenger() *Player {
	if g.AvengerID != NoPlayerID {
		return g.PlayerByID(g.AvengerID)
	}
	return nil
}

func (g *Game) SetAvenger(p *Player) {
	switch {
	case p == nil:
		g.AvengerID = NoPlayerID
	default:
		g.AvengerID = p.ID()
	}
}

func (g *Game) Candidate() *CandidateTile {
	if len(g.Candidates) > 0 {
		return g.Candidates[0]
	}
	return nil
}

func (g *Game) ForeignLand(index int) *ForeignLand {
	return g.ForeignLands[index]
}

func (ids *IDS) Remove(id int64) error {
	for j, i := range *ids {
		if i == id {
			*ids = append((*ids)[:j], (*ids)[j+1:]...)
		}
	}
	return errors.New(fmt.Sprintf("ID: %d not found.", id))
}

func (as *ActionSpace) CubeCount() int {
	var count int
	for _, cubes := range as.Cubes {
		count += cubes
	}
	return count
}

type Games []*Game

func (g *Game) Colors() color.Colors {
	return color.Colors{color.Yellow, color.Purple, color.Green, color.White, color.Black}
}

func (g *Game) Start(ctx context.Context) error {
	g.Status = game.Running
	g.Phase = Setup
	g.Junks = 25

	g.ChiefMinisterID = NoPlayerID
	g.AdmiralID = NoPlayerID
	g.GeneralID = NoPlayerID
	g.AvengerID = NoPlayerID

	for _, u := range g.Users {
		g.addNewPlayer(u)
	}

	g.OfficialsDeck = NewOfficialsDeck()
	g.ConDeck = NewConDeck(g.NumPlayers)
	g.EmperorDeck = NewEmperorDeck()
	g.ActionSpaces = ActionSpaces{
		BribeSecureSpace:    &ActionSpace{ID: BribeSecureSpace, Cubes: Cubes{}},
		NominateSpace:       &ActionSpace{ID: NominateSpace, Cubes: Cubes{}},
		ForceSpace:          &ActionSpace{ID: ForceSpace, Cubes: Cubes{}},
		JunksVoyageSpace:    &ActionSpace{ID: JunksVoyageSpace, Cubes: Cubes{}},
		RecruitArmySpace:    &ActionSpace{ID: RecruitArmySpace, Cubes: Cubes{}},
		BuyGiftSpace:        &ActionSpace{ID: BuyGiftSpace, Cubes: Cubes{}},
		GiveGiftSpace:       &ActionSpace{ID: GiveGiftSpace, Cubes: Cubes{}},
		PetitionSpace:       &ActionSpace{ID: PetitionSpace, Cubes: Cubes{}},
		CommercialSpace:     &ActionSpace{ID: CommercialSpace, Cubes: Cubes{}},
		TaxIncomeSpace:      &ActionSpace{ID: TaxIncomeSpace, Cubes: Cubes{}},
		NoActionSpace:       &ActionSpace{ID: NoActionSpace, Cubes: Cubes{}},
		ImperialFavourSpace: &ActionSpace{ID: ImperialFavourSpace, Cubes: Cubes{}},
	}

	g.CreateMinistries()
	g.CreateDistantLands()
	g.CreateForeignLands()
	g.CreateCandidates()
	g.start(ctx)
	return nil
}

func (g *Game) addNewPlayer(u *user.User) {
	p := CreatePlayer(g, u)
	g.Playerers = append(g.Playerers, p)
}

func (g *Game) ColorMap() color.Colors {
	return color.Colors{color.Yellow, color.Purple, color.Green, color.White, color.Black}
}

func (g *Game) start(ctx context.Context) {
	g.Phase = StartGame
	g.Round = 1
	g.countGiftsPhase(ctx)
	g.chooseChiefMinisterPhase(ctx)
}

func (g *Game) Players() Players {
	ps := g.GetPlayerers()
	if length := len(ps); length > 0 {
		players := make(Players, length)
		for i, p := range ps {
			players[i] = p.(*Player)
		}
		return players
	}
	return nil
}

func (g *Game) setPlayers(players Players) {
	if length := len(players); length > 0 {
		ps := make(game.Playerers, length)
		for i, p := range players {
			ps[i] = p
		}
		g.Playerers = ps
	}
}

func (g *Game) actionsPhase(ctx context.Context) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	g.Phase = Actions
}

func (g *Game) inActionsOrImperialFavourPhase() bool {
	return g.Phase == Actions || g.Phase == ImperialFavour
}

func (g *Game) resetTurn(ctx context.Context) (string, game.ActionType, error) {
	cp := g.CurrentPlayer()

	if !g.CUserIsCPlayerOrAdmin(ctx) {
		return "", game.None, sn.NewVError("Only the current player may perform this action.")
	}
	restful.AddNoticef(ctx, "%s reset turn.", g.NameFor(cp))
	return "", game.Reset, nil
}

func (g *Game) CurrentPlayer() *Player {
	if p := g.CurrentPlayerer(); p != nil {
		return p.(*Player)
	}
	return nil
}

func (g *Game) PlayerByID(id int) *Player {
	if p := g.PlayererByID(id); p != nil {
		return p.(*Player)
	}
	return nil
}

func (g *Game) PlayerBySID(sid string) *Player {
	if p := g.Header.PlayerBySID(sid); p != nil {
		return p.(*Player)
	}
	return nil
}

func (g *Game) PlayerByUserID(id int64) *Player {
	if p := g.PlayererByUserID(id); p != nil {
		return p.(*Player)
	} else {
		return nil
	}
}

func (g *Game) PlayerByIndex(index int) *Player {
	if p := g.PlayererByIndex(index); p != nil {
		return p.(*Player)
	}
	return nil
}

func (g *Game) DrawConCard() *ConCard {
	if len(g.ConDeck) == 0 {
		g.ConDeck = g.ConDiscardPile
		g.ConDiscardPile = ConCards{}
	}
	return g.ConDeck.Draw()
}

func (g *Game) EnableActions(ctx context.Context) bool {
	return g.CUserIsCPlayerOrAdmin(ctx) && (g.Phase == Actions || g.Phase == ImperialFavour)
}

type JunkVoyages map[string][]int

func (g *Game) OnVoyage() JunkVoyages {
	jv := make(map[string][]int, 5)
	jv["white"] = []int{1, 2, 3, 4}
	jv["yellow"] = []int{1, 2, 3, 4}
	jv["black"] = []int{1, 2, 3, 4}
	jv["green"] = []int{1, 2, 3, 4}
	jv["purple"] = []int{1, 2, 3, 4}

	for _, player := range g.Players() {
		jv[player.Color().String()] = []int{}
		for i := 1 + player.OnVoyage; i <= 4; i++ {
			jv[player.Color().String()] = append(jv[player.Color().String()], i)
		}
	}
	return jv
}

func (g *Game) options() (s string) {
	if g.BasicGame {
		s = "Basic"
	} else {
		s = "Advanced"
	}
	if g.AdmiralVariant {
		s += " with Admiral Variant"
	} else {
		s += " without Admiral Variant"
	}
	return
}
