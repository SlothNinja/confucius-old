package confucius

//import (
//	"fmt"
//
//	. "gopkg.in/check.v1"
//)
//
//func (s *MySuite) TestTutorStudent(c *C) {
//	p0 := g1.Players()[0]
//	p1 := g1.Players()[1]
//	candidate := g1.Candidate()
//	candidate.setPlayer(p0)
//	candidate.setOtherPlayer(p1)
//
//	c.Check(len(p0.ConCardHand), Equals, 3)
//
//	entry := p0.tutorStudent(ConCards{}, nil, false)
//	c.Check(entry.String(), Equals,
//		fmt.Sprintf("%s has no cards to tutor a student.", p0.Name()))
//	c.Check(p0.PerformedAction, Equals, true)
//	c.Check(len(p0.ConCardHand), Equals, 3)
//	c.Check(len(candidate.PlayerCards), Equals, 0)
//	c.Check(len(candidate.OtherPlayerCards), Equals, 0)
//
//	entry = p0.tutorStudent(ConCards{&ConCard{Coins: 1}}, p1, false)
//	c.Check(entry.String(), Equals,
//		fmt.Sprintf("%s spent 1 card to tutor student of %s.", p0.Name(), p1.Name()))
//	c.Check(p0.PerformedAction, Equals, true)
//	c.Check(len(p0.ConCardHand), Equals, 2)
//	c.Check(len(candidate.PlayerCards), Equals, 0)
//	c.Check(len(candidate.OtherPlayerCards), Equals, 1)
//
//	p1.GiftsReceived.Append(p0.GiftsBought[0])
//	c.Check(len(p1.GiftsReceived), Equals, 1)
//	entry = p1.tutorStudent(ConCards{&ConCard{Coins: 1}, &ConCard{Coins: 2}, &ConCard{Coins: 3}}, p0, false)
//	c.Check(entry.String(), Equals,
//		fmt.Sprintf("%s spent 3 cards to tutor student of %s and canceled gift received from %s.",
//			p1.Name(), p0.Name(), p0.Name()))
//	c.Check(p1.PerformedAction, Equals, true)
//	c.Check(len(p1.ConCardHand), Equals, 0)
//	c.Check(len(candidate.PlayerCards), Equals, 3)
//	c.Check(len(candidate.OtherPlayerCards), Equals, 1)
//}
//
//func (s *MySuite) TestTutorStudentsPhaseNextPlayer(c *C) {
//	g1.Phase = ImperialExamination
//	g1.Round = 2
//
//	p0, p1, p2 := g1.Players()[0], g1.Players()[1], g1.Players()[2]
//
//	// No Args
//	g1.SetCurrentPlayers(p0)
//	np, phase := g1.tutorStudentsPhaseNextPlayer()
//	c.Check(p1, Equals, np)
//	c.Check(phase, Equals, AwaitPlayerInput)
//
//	g1.SetCurrentPlayers(p1)
//	np, phase = g1.tutorStudentsPhaseNextPlayer()
//	c.Check(p2, Equals, np)
//	c.Check(phase, Equals, AwaitPlayerInput)
//
//	g1.SetCurrentPlayers(p2)
//	np, phase = g1.tutorStudentsPhaseNextPlayer()
//	c.Check(p0, Equals, np)
//	c.Check(phase, Equals, AwaitPlayerInput)
//
//	// Once around
//	np, phase = g1.tutorStudentsPhaseNextPlayer(p0)
//	c.Check(p1, Equals, np)
//	c.Check(phase, Equals, AwaitPlayerInput)
//
//	np, phase = g1.tutorStudentsPhaseNextPlayer(p1)
//	c.Check(p2, Equals, np)
//	c.Check(phase, Equals, AwaitPlayerInput)
//
//	np, phase = g1.tutorStudentsPhaseNextPlayer(p2)
//	c.Check(p0, Equals, np)
//	c.Check(phase, Equals, AwaitPlayerInput)
//
//	// Some Performed Action
//	p0.PerformedAction = true
//	np, phase = g1.tutorStudentsPhaseNextPlayer(p0)
//	c.Check(p1, Equals, np)
//	c.Check(phase, Equals, AwaitPlayerInput)
//
//	// All Performed Action
//	p0.PerformedAction, p1.PerformedAction, p2.PerformedAction = true, true, true
//	np, phase = g1.tutorStudentsPhaseNextPlayer(p0)
//	c.Check(np, IsNil)
//	c.Check(phase, Equals, ExaminationResolution)
//
//	// AutoTutor No Cards
//	p0.PerformedAction, p1.PerformedAction, p2.PerformedAction = false, false, false
//	p0.ConCardHand, p1.ConCardHand, p2.ConCardHand = ConCards{}, ConCards{}, ConCards{}
//	np, phase = g1.tutorStudentsPhaseNextPlayer(p0)
//	c.Check(p1, Equals, np)
//	c.Check(phase, Equals, AutoTutor)
//
//	np, phase = g1.tutorStudentsPhaseNextPlayer(p1)
//	c.Check(p2, Equals, np)
//	c.Check(phase, Equals, AutoTutor)
//
//	np, phase = g1.tutorStudentsPhaseNextPlayer(p2)
//	c.Check(p0, Equals, np)
//	c.Check(phase, Equals, AutoTutor)
//}
