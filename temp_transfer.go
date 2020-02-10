package confucius

import (
	"encoding/gob"
	"html/template"

	"bitbucket.org/SlothNinja/slothninja-games/sn"
	"bitbucket.org/SlothNinja/slothninja-games/sn/game"
	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"golang.org/x/net/context"
)

func init() {
	gob.RegisterName("*game.transferTempInfluenceInEntry", new(transferTempInfluenceInEntry))
	gob.RegisterName("*game.autoTransferTempInfluenceInEntry", new(autoTransferTempInfluenceInEntry))
}

func (g *Game) tempTransfer(ctx context.Context) (string, game.ActionType, error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	p, err := g.validateTempTransfer(ctx)
	if err != nil {
		return "", game.None, err
	}

	cp := g.CurrentPlayer()
	// Transfer Temporary Influence
	gift := cp.transferTempInfluenceTo(p)
	entry := cp.newTransferTempInfluenceInEntry(p, gift)

	// Set flash message
	restful.AddNoticef(ctx, string(entry.HTML()))
	return "", game.Cache, nil
}

func (g *Game) autoTempTransferInfluence(ctx context.Context, from, to *Player) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	canceledGift := from.transferTempInfluenceTo(to)
	from.newAutoTransferTempInfluenceInEntry(to, canceledGift)
}

type transferTempInfluenceInEntry struct {
	*Entry
	MinistryName string
	GiftName     string
}

func (p *Player) newTransferTempInfluenceInEntry(player *Player, gift *GiftCard) *transferTempInfluenceInEntry {
	g := p.Game()
	e := new(transferTempInfluenceInEntry)
	e.Entry = p.newEntry()
	e.OtherPlayerID = player.ID()
	e.MinistryName = g.ministryInProgress().Name()
	if gift != nil {
		e.GiftName = gift.Name()
	}
	p.Log = append(p.Log, e)
	g.Log = append(g.Log, e)
	return e
}

func (e *transferTempInfluenceInEntry) HTML() template.HTML {
	if e.GiftName == "" {
		return restful.HTML("%s temporarily transfered influence in %s ministry to %s.",
			e.Player().Name(), e.MinistryName, e.OtherPlayer().Name())
	}
	return restful.HTML("%s temporarily transfered influence in %s ministry to %s, and removed gift %s from play.",
		e.Player().Name(), e.MinistryName, e.OtherPlayer().Name(), e.GiftName)
}

type autoTransferTempInfluenceInEntry struct {
	*Entry
	MinistryName string
	GiftName     string
}

func (p *Player) newAutoTransferTempInfluenceInEntry(player *Player, gift *GiftCard) *autoTransferTempInfluenceInEntry {
	g := p.Game()
	e := new(autoTransferTempInfluenceInEntry)
	e.Entry = p.newEntry()
	e.OtherPlayerID = player.ID()
	e.MinistryName = g.ministryInProgress().Name()
	if gift != nil {
		e.GiftName = gift.Name()
	}
	p.Log = append(p.Log, e)
	g.Log = append(g.Log, e)
	return e
}

func (e *autoTransferTempInfluenceInEntry) HTML() template.HTML {
	if e.GiftName == "" {
		return restful.HTML("System auto-transfered influence in %s ministry temporarily from %s to %s.",
			e.MinistryName, e.Player().Name(), e.OtherPlayer().Name())
	}
	return restful.HTML("System auto-transfered influence in %s ministry temporarily from %s to %s, and removed gift %s from play.",
		e.MinistryName, e.Player().Name(), e.OtherPlayer().Name(), e.GiftName)
}

func (g *Game) validateTempTransfer(ctx context.Context) (*Player, error) {
	p, err := g.getPlayer(ctx, "temp-transfer-player")
	if err != nil {
		return nil, err
	}

	cp := g.CurrentPlayer()
	m := g.ministryInProgress()

	switch {
	case m == nil:
		return nil, sn.NewVError("No ministry resolution in progress.")
	case !g.CUserIsCPlayerOrAdmin(ctx):
		return nil, sn.NewVError("Only the current player may perform g action.")
	case !(g.Phase == MinistryResolution || g.Phase == FinalMinistryResolution):
		return nil, sn.NewVError("You cannot transfer influence during the %s phase.", g.PhaseName())
	case !cp.TempPlayers().Include(p):
		return nil, sn.NewVError("You cannot temporarily transfer influence in %s ministry to %s.",
			m.Name(), g.NameFor(p))
	}
	return p, nil
}

func (g *Game) EnableTempTransfer(ctx context.Context) bool {
	return g.CUserIsCPlayerOrAdmin(ctx) && g.Phase == MinistryResolution || g.Phase == FinalMinistryResolution
}

func (p *Player) TempPlayers() Players {
	m := p.Game().ministryInProgress()
	if m == nil {
		return nil
	}
	var ps Players

	// Restricted based on gift obligations
	giftValue := NoGiftValue
	for _, gift := range p.GiftsReceived {
		player := gift.Player()
		if player.hasTempInfluence() {
			switch {
			case gift.Value == giftValue:
				ps = append(ps, player)
			case gift.Value > giftValue:
				ps = Players{player}
				giftValue = gift.Value
			}
		}
	}

	if len(ps) > 0 {
		return ps
	}

	// No gift obligation so any other players in ministry.
	for _, o := range m.Officials {
		player := o.TempPlayer()
		if player != nil && player.NotEqual(p) && !ps.Include(player) {
			ps = append(ps, player)
		}
	}
	return ps
}

func (p *Player) transferTempInfluenceTo(player *Player) *GiftCard {
	m := p.Game().ministryInProgress()

	// Cancel Gift
	canceledGift := p.cancelGiftFrom(player)
	p.PerformedAction = true

	for _, o := range m.Officials {
		if p.Equal(o.TempPlayer()) {
			o.setTempPlayer(player)
		}
	}
	return canceledGift
}
