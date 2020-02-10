package confucius

//import (
//	"go.chromium.org/gae/service/datastore"
//	. "gopkg.in/check.v1"
//)
//
//func (s *MySuite) TestCompare(c *C) {
//	g1.SetChiefMinister(nil)
//	c.Assert(g1.Admiral(), IsNil)
//	c.Assert(g1.ChiefMinister(), IsNil)
//	c.Assert(g1.General(), IsNil)
//
//	// Players
//	p0 := g1.Players()[0]
//	p1 := g1.Players()[1]
//	p2 := g1.Players()[2]
//
//	// Based on score, all same score
//	c.Check(p0.compare(p1), Equals, sn.EqualTo)
//	c.Check(p0.compare(p2), Equals, sn.EqualTo)
//	c.Check(p1.compare(p0), Equals, sn.EqualTo)
//	c.Check(p1.compare(p2), Equals, sn.EqualTo)
//	c.Check(p2.compare(p0), Equals, sn.EqualTo)
//	c.Check(p2.compare(p1), Equals, sn.EqualTo)
//
//	// same score, but p1 is General
//	g1.SetGeneral(p1)
//	c.Check(p0.compare(p1), Equals, sn.LessThan)
//	c.Check(p0.compare(p2), Equals, sn.EqualTo)
//	c.Check(p1.compare(p0), Equals, sn.GreaterThan)
//	c.Check(p1.compare(p2), Equals, sn.GreaterThan)
//	c.Check(p2.compare(p0), Equals, sn.EqualTo)
//	c.Check(p2.compare(p1), Equals, sn.LessThan)
//
//	// same score, but p2 is ChiefMinister
//	g1.SetChiefMinister(p2)
//	c.Check(p0.compare(p1), Equals, sn.LessThan)
//	c.Check(p0.compare(p2), Equals, sn.LessThan)
//	c.Check(p1.compare(p0), Equals, sn.GreaterThan)
//	c.Check(p1.compare(p2), Equals, sn.LessThan)
//	c.Check(p2.compare(p0), Equals, sn.GreaterThan)
//	c.Check(p2.compare(p1), Equals, sn.GreaterThan)
//
//	// same score, but p0 is Admiral
//	g1.SetAdmiral(p0)
//	c.Check(p0.compare(p1), Equals, sn.GreaterThan)
//	c.Check(p0.compare(p2), Equals, sn.GreaterThan)
//	c.Check(p1.compare(p0), Equals, sn.LessThan)
//	c.Check(p1.compare(p2), Equals, sn.LessThan)
//	c.Check(p2.compare(p0), Equals, sn.LessThan)
//	c.Check(p2.compare(p1), Equals, sn.GreaterThan)
//
//	// different score
//	p2.Score = 10
//	c.Check(p0.compare(p1), Equals, sn.GreaterThan)
//	c.Check(p0.compare(p2), Equals, sn.LessThan)
//	c.Check(p1.compare(p0), Equals, sn.LessThan)
//	c.Check(p1.compare(p2), Equals, sn.LessThan)
//	c.Check(p2.compare(p0), Equals, sn.GreaterThan)
//	c.Check(p2.compare(p1), Equals, sn.GreaterThan)
//}
//
//func (s *MySuite) TestDeterminePlaces(c *C) {
//	g1.SetChiefMinister(nil)
//	c.Assert(g1.Admiral(), IsNil)
//	c.Assert(g1.ChiefMinister(), IsNil)
//	c.Assert(g1.General(), IsNil)
//
//	p0 := g1.Players()[0]
//	p1 := g1.Players()[1]
//	p2 := g1.Players()[2]
//
//	p0.Score = 20
//	p1.Score = 18
//	p2.Score = 17
//
//	places := g1.determinePlaces()
//	c.Check(places, DeepEquals, []Players{{p0}, {p1}, {p2}})
//	c.Check(toIDS(places), DeepEquals, []datastore.IDS{{p0.User().ID()}, {p1.User().ID()}, {p2.User().ID()}})
//
//	p1.Score = 20
//	places = g1.determinePlaces()
//	c.Check(places, DeepEquals, []Players{{p0, p1}, {p2}})
//	c.Check(toIDS(places), DeepEquals, []datastore.IDS{{p0.User().ID(), p1.User().ID()}, {p2.User().ID()}})
//}
//
////func (p *Player) clearActions() {
////	p.PerformedAction = false
////	p.Log = make(sn.GameLog, 0)
////}
//
//func (s *MySuite) TestClearActions(c *C) {
//	for _, p := range g1.Players() {
//		p.PerformedAction = true
//		p.Log = make(sn.GameLog, 2)
//	}
//
//	for _, p := range g1.Players() {
//		p.clearActions()
//	}
//
//	for _, p := range g1.Players() {
//		c.Check(p.PerformedAction, Equals, false)
//		c.Check(len(p.Log), Equals, 0)
//	}
//}
//
//func (s *MySuite) TestBeginningOfPhaseReset(c *C) {
//	for _, p := range g1.Players() {
//		p.PerformedAction = true
//		p.Log = make(sn.GameLog, 2)
//	}
//
//	g1.beginningOfPhaseReset()
//
//	for _, p := range g1.Players() {
//		c.Check(p.PerformedAction, Equals, false)
//		c.Check(len(p.Log), Equals, 0)
//	}
//}
