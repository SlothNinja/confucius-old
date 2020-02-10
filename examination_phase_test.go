package confucius

//import (
//	"fmt"
//
//	. "gopkg.in/check.v1"
//)
//
//func (s *MySuite) TestCanHoldExamination(c *C) {
//	// Setup for true
//	p0 := g1.Players()[0]
//	p1 := g1.Players()[1]
//	g1.Phase = ImperialExamination
//	g1.Round = 2
//	g1.Candidate().setPlayer(p0)
//	g1.Candidate().setOtherPlayer(p1)
//
//	// Correct phase, round, and two different students
//	c.Check(g1.canHoldExamination(), Equals, true)
//
//	// Incorrect phase
//	g1.Phase = Actions
//	c.Check(g1.canHoldExamination(), Equals, false)
//
//	// Incorrect round
//	g1.Phase = ImperialExamination
//	g1.Round = 1
//	c.Check(g1.canHoldExamination(), Equals, false)
//
//	// Two students but of same player
//	g1.Round = 2
//	g1.Candidate().setOtherPlayer(p0)
//	c.Check(g1.canHoldExamination(), Equals, true)
//
//	// But not place to put candidate
//	for _, m := range g1.Ministries {
//		m.Resolved = true
//	}
//	c.Check(g1.canHoldExamination(), Equals, false)
//
//	// One student but not forced examination
//	for _, m := range g1.Ministries {
//		m.Resolved = false
//	}
//	g1.Candidate().setOtherPlayer(nil)
//	c.Check(g1.canHoldExamination(), Equals, false)
//
//	// One student and forced examination
//	g1.ActionSpaces[ForceSpace].Cubes[p0.ID()] = 1
//	c.Check(g1.canHoldExamination(), Equals, true)
//}
//
//func (s *MySuite) TestCanResolveExamination(c *C) {
//	// Setup for true
//	p0 := g1.Players()[0]
//	p1 := g1.Players()[1]
//	g1.Phase = ImperialExamination
//	g1.Round = 2
//	candidate := g1.Candidate()
//	candidate.setPlayer(p0)
//	candidate.setOtherPlayer(p0)
//
//	// Correct phase, round, and two students but of same player
//	c.Check(g1.canResolveExamination(), Equals, true)
//
//	// Incorrect phase
//	g1.Phase = Actions
//	c.Check(g1.canResolveExamination(), Equals, false)
//
//	// Incorrect round
//	g1.Phase = ImperialExamination
//	g1.Round = 1
//	c.Check(g1.canResolveExamination(), Equals, false)
//
//	// Two students of different players
//	g1.Round = 2
//	candidate.setOtherPlayer(p1)
//	c.Check(g1.canResolveExamination(), Equals, false)
//
//	// One student but not forced examination
//	candidate.setOtherPlayer(nil)
//	c.Check(g1.canResolveExamination(), Equals, false)
//
//	// One student and forced examination
//	g1.ActionSpaces[ForceSpace].Cubes[p0.ID()] = 1
//	c.Check(g1.canResolveExamination(), Equals, true)
//}
//
//func (s *MySuite) TestAnyPlaceForCandidate(c *C) {
//	// empty spots in ministries
//	c.Check(len(g1.Ministries[Bingbu].emptyCandidateSpots()) > 0, Equals, true)
//	c.Check(len(g1.Ministries[Hubu].emptyCandidateSpots()) > 0, Equals, true)
//	c.Check(len(g1.Ministries[Gongbu].emptyCandidateSpots()) > 0, Equals, true)
//	c.Check(g1.hasPlaceForCandidate(), Equals, true)
//
//	// No empty spots
//	g1.placeNewOfficialsPhase()
//	g1.placeNewOfficialsPhase()
//	g1.placeNewOfficialsPhase()
//	g1.placeNewOfficialsPhase()
//	for _, ministry := range g1.Ministries {
//		ministry.init(g1)
//	}
//	c.Check(len(g1.Ministries[Bingbu].emptyCandidateSpots()) > 0, Equals, false)
//	c.Check(len(g1.Ministries[Hubu].emptyCandidateSpots()) > 0, Equals, false)
//	c.Check(len(g1.Ministries[Gongbu].emptyCandidateSpots()) > 0, Equals, false)
//	c.Check(g1.hasPlaceForCandidate(), Equals, true)
//
//	// No unbribed officials
//	p1 := g1.Players()[1]
//	for _, m := range g1.Ministries {
//		for _, o := range m.Officials {
//			o.setPlayer(p1)
//		}
//	}
//	c.Check(g1.hasPlaceForCandidate(), Equals, true)
//
//	// No unsecured officials
//	for _, m := range g1.Ministries {
//		for _, o := range m.Officials {
//			o.Secured = true
//		}
//	}
//	c.Check(g1.hasPlaceForCandidate(), Equals, false)
//}
//
//func (s *MySuite) TestExaminationForced(c *C) {
//	c.Check(g1.examinationForced(), Equals, false)
//	g1.ActionSpaces[ForceSpace].Cubes[0] = 1
//	c.Check(g1.examinationForced(), Equals, true)
//}
//
//func (s *MySuite) TestResolveExamination(c *C) {
//	// One student
//	p0 := g1.Players()[0]
//	g1.Candidate().setPlayer(p0)
//	phase := g1.resolveExamination()
//	c.Check(phase, Equals, AwaitPlayerInput)
//	c.Check(p0.Log[len(p0.Log)-1].String(), Equals,
//		fmt.Sprintf("<div>%s won the Imperial Examination uncontested.</div>", p0.Name()))
//	c.Check(g1.CurrentPlayer(), Equals, p0)
//
//	// Two students same player
//	g1.Candidate().setOtherPlayer(p0)
//	phase = g1.resolveExamination()
//	c.Check(phase, Equals, AwaitPlayerInput)
//	c.Check(p0.Log[len(p0.Log)-1].String(), Equals,
//		fmt.Sprintf("<div>%s won the Imperial Examination uncontested.</div>", p0.Name()))
//	c.Check(g1.CurrentPlayer(), Equals, p0)
//
//	// Two students different players, no cards.
//	p1 := g1.Players()[1]
//	g1.Candidate().setPlayer(p1)
//	phase = g1.resolveExamination()
//	c.Check(phase, Equals, AwaitPlayerInput)
//	c.Check(p1.Log[len(p1.Log)-1].String(), Equals, fmt.Sprintf("<div>%s won the Imperial Examination.</div><div>The student of %s received 0 coins on 0 cards.</div><div>The student of %s received 0 coins on 0 cards.</div>", p1.Name(), p1.Name(), p0.Name()))
//	c.Check(g1.CurrentPlayer(), Equals, p1)
//
//	// Two students different players
//	g1.Candidate().OtherPlayerCards = ConCards{&ConCard{Coins: 2}}
//	phase = g1.resolveExamination()
//	c.Check(phase, Equals, AwaitPlayerInput)
//
//	e := (p0.Log[len(p0.Log)-1]).(*studentPromotionEntry)
//	msg := fmt.Sprintf("<div>%s won the Imperial Examination.</div>", e.Player().Name())
//	winningLength := len(e.WinningCards)
//	msg += fmt.Sprintf("<div>The student of %s received %d coins on %d %s.</div>",
//		e.Player().Name(), e.WinningCards.Coins(), winningLength, pluralize("card", winningLength))
//
//	losingLength := len(e.LosingCards)
//	msg += fmt.Sprintf("<div>The student of %s received %d coins on %d %s.</div>",
//		e.OtherPlayer().Name(), e.LosingCards.Coins(), losingLength, pluralize("card", losingLength))
//	c.Check(e.String(), Equals, msg)
//	c.Check(g1.CurrentPlayer(), Equals, p0)
//}
//
//func (s *MySuite) TestOneStudent(c *C) {
//	candidate := g1.Candidate()
//	c.Assert(candidate, Not(IsNil))
//	c.Assert(candidate.Player(), IsNil)
//	c.Assert(candidate.OtherPlayer(), IsNil)
//
//	// No students
//	c.Check(candidate.hasOnePlayer(), Equals, false)
//
//	// One student
//	candidate.setPlayer(g1.CurrentPlayer())
//	c.Check(candidate.hasOnePlayer(), Equals, true)
//}
//
//func (s *MySuite) TestTwoStudents(c *C) {
//	candidate := g1.Candidate()
//	cp := g1.CurrentPlayer()
//	c.Assert(candidate, Not(IsNil))
//	c.Assert(candidate.Player(), IsNil)
//	c.Assert(candidate.OtherPlayer(), IsNil)
//	c.Assert(cp, Not(IsNil))
//
//	// No students
//	c.Check(candidate.hasTwoPlayers(), Equals, false)
//
//	// One student
//	candidate.setPlayer(cp)
//	c.Check(candidate.hasTwoPlayers(), Equals, false)
//
//	// Two students
//	candidate.setOtherPlayer(cp)
//	c.Check(candidate.hasTwoPlayers(), Equals, true)
//}
//
//func (s *MySuite) TestBothStudentsSamePlayer(c *C) {
//	candidate := g1.Candidate()
//	p1 := g1.Players()[1]
//	candidate.setPlayer(p1)
//	c.Check(candidate.hasTwoSamePlayers(), Equals, false)
//
//	p0 := g1.Players()[0]
//	candidate.setOtherPlayer(p0)
//	c.Check(candidate.hasTwoSamePlayers(), Equals, false)
//
//	candidate.setPlayer(p0)
//	c.Check(candidate.hasTwoSamePlayers(), Equals, true)
//}
//
//func (s *MySuite) TestEnableTutorStudent(c *C) {
//	c.Check(g1.EnableTutorStudent(), Equals, false)
//
//	g1.Phase = ImperialExamination
//	c.Check(g1.EnableTutorStudent(), Equals, true)
//
//	g1.CurrentPlayer().PerformedAction = true
//	c.Check(g1.EnableTutorStudent(), Equals, false)
//}
//
//func (s *MySuite) TestExaminationPhase(c *C) {
//	// canResolveExamination true
//	g1.Phase = ImperialExamination
//	g1.Round = 2
//	p0 := g1.Players()[0]
//	g1.Candidate().setPlayer(p0)
//	g1.Candidate().setOtherPlayer(p0)
//
//	c.Check(g1.canResolveExamination(), Equals, true)
//	c.Check(g1.examinationPhase(), Equals, ExaminationResolution)
//
//	// canResolveExamination false and canHoldExamination true
//	g1.Phase = ImperialExamination
//	p1 := g1.Players()[1]
//	g1.Candidate().setPlayer(p1)
//
//	c.Check(g1.canResolveExamination(), Equals, false)
//	c.Check(g1.canHoldExamination(), Equals, true)
//	c.Check(g1.examinationPhase(), Equals, AwaitPlayerInput)
//
//	// canResolveExamination true and canHoldExamination true
//	g1.Phase = ImperialExamination
//	g1.Round = 1
//
//	c.Check(g1.canResolveExamination(), Equals, false)
//	c.Check(g1.canHoldExamination(), Equals, false)
//	c.Check(g1.examinationPhase(), Equals, MinistryResolution)
//}
