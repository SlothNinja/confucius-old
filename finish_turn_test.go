package confucius

//import . "gopkg.in/check.v1"
//
//func (s *MySuite) TestNextPlayer(c *C) {
//	c.Assert(len(g1.Players()), Equals, 3)
//	c.Assert(g1.Players()[0], Not(IsNil))
//	c.Assert(g1.Players()[1], Not(IsNil))
//	c.Assert(g1.Players()[2], Not(IsNil))
//	c.Assert(g1.NumPlayers, Equals, 3)
//	c.Assert(g1.Phase, Equals, Actions)
//	c.Assert(g1.Round, Equals, 1)
//
//	p0, p1, p2 := g1.Players()[0], g1.Players()[1], g1.Players()[2]
//	c.Check(p0.canPass(), Equals, false)
//	c.Check(p1.canPass(), Equals, false)
//	c.Check(p2.canPass(), Equals, false)
//	c.Check(p1, Equals, g1.CurrentPlayer())
//
//	// Once Around
//	g1.SetCurrentPlayers(g1.nextPlayer())
//	c.Check(p2, Equals, g1.CurrentPlayer())
//	g1.SetCurrentPlayers(g1.nextPlayer())
//	c.Check(p0, Equals, g1.CurrentPlayer())
//	g1.SetCurrentPlayers(g1.nextPlayer())
//	c.Check(p1, Equals, g1.CurrentPlayer())
//	g1.SetCurrentPlayers(g1.nextPlayer())
//	c.Check(p2, Equals, g1.CurrentPlayer())
//
//	// With Player Argument
//	g1.SetCurrentPlayers(g1.nextPlayer(p0))
//	c.Check(p1, Equals, g1.CurrentPlayer())
//	g1.SetCurrentPlayers(g1.nextPlayer(p2))
//	c.Check(p0, Equals, g1.CurrentPlayer())
//	g1.SetCurrentPlayers(g1.nextPlayer(p1))
//	c.Check(p2, Equals, g1.CurrentPlayer())
//}
//
//func (s *MySuite) TestActionPhaseNextPlayer(c *C) {
//	c.Assert(len(g1.Players()), Equals, 3)
//	c.Assert(g1.Players()[0], Not(IsNil))
//	c.Assert(g1.Players()[1], Not(IsNil))
//	c.Assert(g1.Players()[2], Not(IsNil))
//	c.Assert(g1.NumPlayers, Equals, 3)
//	c.Assert(g1.Phase, Equals, Actions)
//	c.Assert(g1.Round, Equals, 1)
//
//	p0, p1, p2 := g1.Players()[0], g1.Players()[1], g1.Players()[2]
//	c.Check(p0.canPass(), Equals, false)
//	c.Check(p1.canPass(), Equals, false)
//	c.Check(p2.canPass(), Equals, false)
//	c.Check(p1, Equals, g1.CurrentPlayer())
//
//	// Once around
//	g1.SetCurrentPlayers(g1.actionPhaseNextPlayer())
//	c.Check(p2, Equals, g1.CurrentPlayer())
//	g1.SetCurrentPlayers(g1.actionPhaseNextPlayer())
//	c.Check(p0, Equals, g1.CurrentPlayer())
//	g1.SetCurrentPlayers(g1.actionPhaseNextPlayer())
//	c.Check(p1, Equals, g1.CurrentPlayer())
//	g1.SetCurrentPlayers(g1.actionPhaseNextPlayer())
//	c.Check(p2, Equals, g1.CurrentPlayer())
//
//	// All but one passed
//	p0.Passed, p1.Passed, p2.Passed = false, true, true
//	g1.SetCurrentPlayers(g1.actionPhaseNextPlayer(p1))
//	c.Check(p0, Equals, g1.CurrentPlayer(), Commentf("p0.ID: %v %v cp.ID: %v %v",
//		p0.ID(), p0.Passed, g1.CurrentPlayer().ID(), g1.CurrentPlayer().Passed))
//
//	// AutoPass
//	p0.Passed, p1.Passed, p2.Passed = false, false, false
//	g1.SetCurrentPlayers(p2)
//	p0.ActionCubes = 0
//	c.Check(p1, Equals, g1.actionPhaseNextPlayer())
//	c.Check(p0.Passed, Equals, true)
//
//	// All passed
//	p0.Passed, p1.Passed, p2.Passed = true, true, true
//	c.Check(g1.actionPhaseNextPlayer(), IsNil)
//}
