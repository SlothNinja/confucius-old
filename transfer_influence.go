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
	gob.RegisterName("*game.transferInfluenceEntry", new(transferInfluenceEntry))
}

func (g *Game) transferInfluence(ctx context.Context) (string, game.ActionType, error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	ministry, official, player, err := g.validateTransferInfluence(ctx)
	if err != nil {
		return "", game.None, err
	}

	cp := g.CurrentPlayer()
	cp.PerformedAction = true

	// Transfer Influence
	official.setPlayer(player)

	// Cancel Gift
	gift := cp.cancelGiftFrom(player)
	e := cp.newTransferInfluenceEntry(player, ministry, official, gift)

	// Set flash message
	restful.AddNoticef(ctx, string(e.HTML()))
	return "", game.Cache, nil
}

type transferInfluenceEntry struct {
	*Entry
	MinistryName string
	Seniority    Seniority
	Gift         *GiftCard
}

func (p *Player) newTransferInfluenceEntry(op *Player, m *Ministry, o *OfficialTile, gift *GiftCard) *transferInfluenceEntry {
	g := p.Game()
	e := new(transferInfluenceEntry)
	e.Entry = p.newEntry()
	e.OtherPlayerID = op.ID()
	e.MinistryName = m.Name()
	e.Seniority = o.Seniority

	p.Log = append(p.Log, e)
	g.Log = append(g.Log, e)
	return e
}

func (e *transferInfluenceEntry) HTML() template.HTML {
	if e.Gift != nil && e.Gift.Value > 0 {
		return restful.HTML("%s transferred influence on %s official with level %d seniority to %s, and removed %s gift of %s from game.", e.Player().Name(), e.MinistryName, e.Seniority, e.OtherPlayer().Name(), e.Gift.Name(), e.OtherPlayer().Name())
	}
	return restful.HTML("%s transferred influence on %s official with level %d seniority to %s.", e.Player().Name(), e.MinistryName, e.Seniority, e.OtherPlayer().Name())
}

func (g *Game) validateTransferInfluence(ctx context.Context) (*Ministry, *OfficialTile, *Player, error) {
	if _, err := g.validatePlayerAction(ctx); err != nil {
		return nil, nil, nil, err
	}

	ministry, official, err := g.getMinistryAndOfficial(ctx, "transfer-influence-official")
	if err != nil {
		return nil, nil, nil, err
	}

	player, err := g.getPlayer(ctx, "transfer-influence-player")
	if err != nil {
		return nil, nil, nil, err
	}

	cp := g.CurrentPlayer()

	switch {
	case official.Player() == nil:
		return nil, nil, nil, sn.NewVError("You don't have influence over the official having seniority level %d in the %s ministry.", official.Seniority, ministry.Name())
	case official.Player().NotEqual(cp):
		return nil, nil, nil, sn.NewVError("You don't have influence over the official having seniority level %d in the %s ministry.", official.Seniority, ministry.Name())
	case ministry.Resolved:
		return nil, nil, nil, sn.NewVError("You can't transfer influence in a resolved ministry.")
	}
	return ministry, official, player, nil
}

func (g *Game) EnableTransferInfluence(ctx context.Context) bool {
	cp := g.CurrentPlayer()
	return g.CUserIsCPlayerOrAdmin(ctx) && cp.canTransferInfluence()
}

func (p *Player) canTransferInfluence() bool {
	g := p.Game()
	return g.inActionsOrImperialFavourPhase() && !p.PerformedAction && p.hasInfluenceToTransfer()
}

func (p *Player) hasInfluenceToTransfer() bool {
	for _, m := range p.Game().Ministries {
		if !m.Resolved {
			for _, o := range m.Officials {
				if p.Equal(o.Player()) {
					return true
				}
			}
		}
	}
	return false
}
