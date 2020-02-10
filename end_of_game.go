package confucius

import (
	"fmt"
	"html/template"

	"bitbucket.org/SlothNinja/slothninja-games/sn/contest"
	"bitbucket.org/SlothNinja/slothninja-games/sn/game"
	"bitbucket.org/SlothNinja/slothninja-games/sn/log"
	"bitbucket.org/SlothNinja/slothninja-games/sn/restful"
	"bitbucket.org/SlothNinja/slothninja-games/sn/send"
	"go.chromium.org/gae/service/mail"
	"golang.org/x/net/context"
)

func (g *Game) endOfGamePhase(ctx context.Context) (cs contest.Contests) {
	log.Debugf(ctx, "Entering")
	defer log.Debugf(ctx, "Exiting")

	if !g.endGame() {
		g.newRoundPhase(ctx)
		g.countGiftsPhase(ctx)
		g.chooseChiefMinisterPhase(ctx)
		return
	}

	if !g.Ministries.allResolved() {
		if completed := g.ministryResolutionPhase(ctx, true); !completed {
			return
		}
	}
	return g.endGameScoring(ctx)
}

func (g *Game) endGame() bool {
	return g.Ministries.allResolved() || len(g.Candidates) <= 0 || g.Wall >= 9
}

func (ms Ministries) allResolved() bool {
	for _, m := range ms {
		if !m.Resolved {
			return false
		}
	}
	return true
}

func (g *Game) endGameScoring(ctx context.Context) contest.Contests {
	g.Phase = EndGameScoring
	g.ScoreChiefMinister()
	g.ScoreAdmiral()
	g.ScoreGeneral()
	places := g.determinePlaces(ctx)
	g.SetWinners(places[0])
	return contest.GenContests(ctx, places)
}

func toIDS(places []Players) [][]int64 {
	sids := make([][]int64, len(places))
	for i, players := range places {
		for _, p := range players {
			sids[i] = append(sids[i], p.User().ID)
		}
	}
	return sids
}

func (g *Game) SendEndGameNotifications(ctx context.Context) error {
	g.Phase = GameOver
	g.Status = game.Completed

	ms := make([]*mail.Message, len(g.Players()))
	sender := "webmaster@slothninja.com"
	subject := fmt.Sprintf("SlothNinja Games: Confucius #%d Has Ended", g.ID)

	var body string
	for _, p := range g.Players() {
		body += fmt.Sprintf("%s scored %d points.\n", g.NameFor(p), p.Score)
	}

	var names []string
	for _, p := range g.Winners() {
		names = append(names, g.NameFor(p))
	}
	body += fmt.Sprintf("\nCongratulations to: %s.", restful.ToSentence(names))

	for i, p := range g.Players() {
		ms[i] = &mail.Message{
			To:      []string{p.User().Email},
			Sender:  sender,
			Subject: subject,
			Body:    body,
		}
	}
	return send.Message(ctx, ms...)
}

type playerCounts map[int]int

func (pcs playerCounts) For(player *Player) int {
	return pcs[player.ID()]
}

func (pcs playerCounts) SetFor(player *Player, value int) {
	pcs[player.ID()] = value
}

func (pcs playerCounts) IncrementFor(player *Player, by ...int) {
	increment := 1
	if len(by) == 1 {
		increment = by[0]
	}
	pcs.SetFor(player, pcs.For(player)+increment)
}

func (g *Game) ScoreChiefMinister() {
	g.Phase = AwardChiefMinister

	counts := make(playerCounts)
	for _, ministry := range g.Ministries {
		for _, official := range ministry.Officials {
			if official.Bribed() {
				counts.IncrementFor(official.Player())
			}
		}
		if minister := ministry.Minister(); minister != nil {
			counts[minister.ID()] += 1
		}
		if secretary := ministry.Secretary(); secretary != nil {
			counts[secretary.ID()] += 1
		}
	}

	players := Players{}
	max := 0
	for _, player := range g.Players() {
		switch {
		case counts.For(player) == max:
			players = append(players, player)
		case counts.For(player) > max:
			max = counts.For(player)
			players = Players{player}
			g.SetChiefMinister(player)
		}
	}

	if len(players) > 1 {
		g.SetChiefMinister(g.Ministries[Hubu].Minister())
	}

	if chief := g.ChiefMinister(); chief != nil {
		chief.Score += 1
		g.NewScoreChiefMinisterEntry(chief)
	}
}

type scoreChiefMinisterEntry struct {
	*Entry
}

func (g *Game) NewScoreChiefMinisterEntry(p *Player) *scoreChiefMinisterEntry {
	e := new(scoreChiefMinisterEntry)
	e.Entry = p.newEntry()
	p.Log = append(p.Log, e)
	g.Log = append(g.Log, e)
	return e
}

func (e *scoreChiefMinisterEntry) HTML() template.HTML {
	return restful.HTML("%s awarded title of Chief Minister and 1 point.", e.Player().Name())
}

func (g *Game) ScoreAdmiral() {
	g.Phase = AwardAdmiral
	counts := make(playerCounts)

	// 5 for each junk at a distant land
	for _, land := range g.DistantLands {
		for _, player := range land.Players() {
			counts.IncrementFor(player, 5)
		}
	}

	players := Players{}
	max := 0
	for _, player := range g.Players() {
		// 1 for each junk at sea
		counts[player.ID()] += player.OnVoyage

		// Find max or those tied with max
		switch {
		case counts[player.ID()] == max:
			players = append(players, player)
		case counts[player.ID()] > max:
			max = counts[player.ID()]
			players = Players{player}
			g.SetAdmiral(player)
		}
	}

	// Only 1 then Admiral is already found. Otherwise Gongbu Minister is Admiral.
	if len(players) > 1 {
		g.SetAdmiral(g.Ministries[Gongbu].Minister())
	}

	if admiral := g.Admiral(); admiral != nil {
		admiral.Score += 1
		g.NewScoreAdmiralEntry(admiral)
	}
}

type scoreAdmiralEntry struct {
	*Entry
}

func (g *Game) NewScoreAdmiralEntry(p *Player) *scoreAdmiralEntry {
	e := new(scoreAdmiralEntry)
	e.Entry = p.newEntry()
	p.Log = append(p.Log, e)
	g.Log = append(g.Log, e)
	return e
}

func (e *scoreAdmiralEntry) HTML() template.HTML {
	return restful.HTML("%s awarded title of Admiral and 1 point.", e.Player().Name())
}

func (g *Game) ScoreGeneral() {
	g.Phase = AwardGeneral

	counts := make(playerCounts)

	// 1 for each army at a foreign land
	for _, land := range g.ForeignLands {
		for _, box := range land.Boxes {
			if box.Player() != nil {
				counts.IncrementFor(box.Player())
			}
		}
	}

	players := Players{}
	max := 0
	for _, player := range g.Players() {
		// 1 for avenging emperor
		if player.Equal(g.Avenger()) {
			counts.IncrementFor(player)
		}

		// 1 for each army in military colonies
		counts[player.ID()] += player.RecruitedArmies

		// Find max or those tied with max
		switch {
		case counts[player.ID()] == max:
			players = append(players, player)
		case counts[player.ID()] > max:
			max = counts[player.ID()]
			players = Players{player}
			g.SetGeneral(player)
		}
	}

	// Only 1 then General is already found. Otherwise Bingbu Minister is Admiral.
	if len(players) > 1 {
		g.SetGeneral(g.Ministries[Bingbu].Minister())
	}

	if general := g.General(); general != nil {
		general.Score += 1
		general.newScoreGeneralEntry()
	}
}

type scoreGeneralEntry struct {
	*Entry
}

func (p *Player) newScoreGeneralEntry() *scoreGeneralEntry {
	g := p.Game()
	e := new(scoreGeneralEntry)
	e.Entry = p.newEntry()
	p.Log = append(p.Log, e)
	g.Log = append(g.Log, e)
	return e
}

func (e *scoreGeneralEntry) HTML() template.HTML {
	return restful.HTML("%s awarded title of General and 1 point.", e.Player().Name())
}

func (g *Game) SetWinners(rmap contest.ResultsMap) {
	g.Phase = AnnounceWinners
	g.Status = game.Completed

	g.SetCurrentPlayerers()
	for key := range rmap {
		p := g.PlayerByUserID(key.IntID())
		g.WinnerIDS = append(g.WinnerIDS, p.ID())
	}

	g.newAnnounceWinnersEntry()
}

//func (g *Game) SetWinners(winners Players) {
//	g.Phase = AnnounceWinners
//	g.Status = game.Completed
//
//	g.SetCurrentPlayerers()
//	g.WinnerIDS = game.UserIndices{}
//
//	for _, winner := range winners {
//		g.WinnerIDS = append(g.WinnerIDS, winner.ID())
//	}
//
//	g.newAnnounceWinnersEntry()
//}

type announceWinnersEntry struct {
	*Entry
}

func (g *Game) newAnnounceWinnersEntry() *announceWinnersEntry {
	e := new(announceWinnersEntry)
	e.Entry = g.newEntry()
	g.Log = append(g.Log, e)
	return e
}

func (e *announceWinnersEntry) HTML() template.HTML {
	names := make([]string, len(e.Winners()))
	for i, winner := range e.Winners() {
		names[i] = winner.Name()
	}
	return restful.HTML("Congratulations to: %s.", restful.ToSentence(names))
}

func (g *Game) Winners() (ps Players) {
	switch length := len(g.WinnerIDS); length {
	case 0:
	default:
		ps = make(Players, length)
		for i, pid := range g.WinnerIDS {
			player := g.PlayerByID(pid)
			ps[i] = player
		}
	}
	return
}
