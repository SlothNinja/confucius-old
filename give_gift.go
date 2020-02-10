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
	gob.RegisterName("*game.giveGiftEntry", new(giveGiftEntry))
}

func (g *Game) giveGift(ctx context.Context) (string, game.ActionType, error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	// Get Recipient and Gift
	recipient, gift, cubes, err := g.validateGiveGift(ctx)
	if err != nil {
		return "", game.None, err
	}

	cp := g.CurrentPlayer()
	cp.PerformedAction = true

	// Place Action Cube(s) In GiveGiftSpace
	cp.PlaceCubesIn(GiveGiftSpace, cubes)

	// Give Gift
	canceledGift := cp.GiveGiftTo(gift, recipient)
	cp.GiftCardHand.Remove(gift)

	// Create Action Object for logging
	entry := cp.newGiveGiftEntry(recipient, gift, canceledGift)

	// Set flash message
	restful.AddNoticef(ctx, string(entry.HTML()))
	return "", game.Cache, nil
}

type giveGiftEntry struct {
	*Entry
	Gift         *GiftCard
	CanceledGift bool
}

func (p *Player) newGiveGiftEntry(op *Player, gift *GiftCard, canceled bool) *giveGiftEntry {
	g := p.Game()
	e := new(giveGiftEntry)
	e.Entry = p.newEntry()
	e.Gift = gift
	e.CanceledGift = canceled
	e.SetOtherPlayer(op)
	p.Log = append(p.Log, e)
	g.Log = append(g.Log, e)
	return e
}

func (e *giveGiftEntry) HTML() template.HTML {
	if !e.CanceledGift {
		return restful.HTML("%s gave value %d gift (%s) to %s.",
			e.Player().Name(), e.Gift.Value, e.Gift.Name(), e.OtherPlayer().Name())
	}
	return restful.HTML("%s gave value %d gift (%s) to %s and canceled gift from %s.",
		e.Player().Name(), e.Gift.Value, e.Gift.Name(), e.OtherPlayer().Name(), e.OtherPlayer().Name())
}

func (g *Game) validateGiveGift(ctx context.Context) (*Player, *GiftCard, int, error) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	cubes, err := g.validatePlayerAction(ctx)
	if err != nil {
		return nil, nil, 0, err
	}

	recipient, err := g.getPlayer(ctx, "give-gift-player")
	if err != nil {
		return nil, nil, 0, err
	}

	giftValue, err := g.getGiftValue(ctx, "give-gift")
	if err != nil {
		return nil, nil, 0, err
	}

	cp := g.CurrentPlayer()
	oldGift := recipient.giftFrom(cp)
	givenGift := cp.GetBoughtGift(giftValue)
	receivedGift := cp.giftFrom(recipient)

	switch {
	case recipient == nil:
		return nil, nil, 0, sn.NewVError("Recipient not found.")
	case givenGift == nil:
		return nil, nil, 0, sn.NewVError("You don't have a gift of value %d to give.", giftValue)
	case oldGift != nil && oldGift.Value > givenGift.Value:
		return nil, nil, 0, sn.NewVError("You must give a gift that is greater than your present gift to the player.")
	case cp.Equal(recipient):
		return nil, nil, 0, sn.NewVError("You can't give yourself a gift.")
	case receivedGift != nil && receivedGift.Value > givenGift.Value:
		return nil, nil, 0, sn.NewVError("You must give a gift that is greater than or equal to the gift the player gave you.")
	}
	return recipient, givenGift, cubes, nil
}

func (g *Game) EnableGiveGift(ctx context.Context) bool {
	cp := g.CurrentPlayer()
	requiredCubes := cp.RequiredCubesFor(GiveGiftSpace)
	return g.CUserIsCPlayerOrAdmin(ctx) && cp.ActionCubes >= requiredCubes && len(cp.GiftsBought) >= 1
}
