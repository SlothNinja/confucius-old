package confucius

//import (
//	"bitbucket.org/SlothNinja/slothninja-games/sn/game"
//	"bitbucket.org/SlothNinja/slothninja-games/sn/user"
//	. "gopkg.in/check.v1"
//)
//
//type canPassTest struct {
//	passed      bool
//	phase       game.Phase
//	extraAction bool
//	cubes       int
//	result      bool
//}
//
//var canPassTests = []*canPassTest{
//	&canPassTest{false, NoPhase, false, 0, false},
//	&canPassTest{false, NoPhase, false, 1, false},
//	&canPassTest{false, NoPhase, true, 0, false},
//	&canPassTest{false, NoPhase, true, 1, false},
//	&canPassTest{false, Actions, false, 0, true},
//	&canPassTest{false, Actions, false, 1, false},
//	&canPassTest{false, Actions, true, 0, false},
//	&canPassTest{false, Actions, true, 1, false},
//	&canPassTest{true, NoPhase, false, 1, false},
//	&canPassTest{true, NoPhase, false, 1, false},
//	&canPassTest{true, NoPhase, true, 1, false},
//	&canPassTest{true, NoPhase, true, 1, false},
//	&canPassTest{true, Actions, false, 1, false},
//	&canPassTest{true, Actions, false, 1, false},
//	&canPassTest{true, Actions, true, 1, false},
//	&canPassTest{true, Actions, true, 1, false},
//}
//
//func (s *MySuite) TestCanPass(c *C) {
//	for _, t := range canPassTests {
//		g1.Phase = t.phase
//		g1.ExtraAction = t.extraAction
//		for _, p := range g1.Players() {
//			p.Passed = t.passed
//			p.ActionCubes = t.cubes
//			c.Check(p.canPass(), Equals, t.result,
//				Commentf("Passed: %v Phase: %s ExtraAction: %v ActionCubes: %v",
//					t.passed, PhaseNames[t.phase], t.extraAction, t.cubes))
//		}
//	}
//}
//
//func (s *MySuite) TestEnablePass(c *C) {
//	p0 := g1.Players()[0]
//	p1 := g1.Players()[1]
//	for _, t := range canPassTests {
//		g1.Phase = t.phase
//		g1.ExtraAction = t.extraAction
//		g1.SetCurrentPlayerers(p0)
//		user.WithCurrent(ctx, p0.User())
//		cp := g1.CurrentPlayer()
//		cp.Passed = t.passed
//		cp.ActionCubes = t.cubes
//		cp.PerformedAction = false
//		c.Check(g1.EnablePass(), Equals, t.result)
//		user.WithCurrent(ctx, p1.User())
//		c.Check(g1.EnablePass(), Equals, false)
//	}
//}
//
//func (s *MySuite) TestValidatePass(c *C) {
//	p0 := g1.Players()[0]
//	cu := p0.User()
//	user.WithCurrent(ctx, cu)
//	g1.SetCurrentPlayerers(p0)
//
//	//	v := sn.View{}
//	//	values := url.Values{}
//	//	values.Set("action", "pass")
//	err := p0.validatePass(ctx)
//	c.Check(err, Not(IsNil))
//	c.Check(err.Error(), Equals, "You must use all of your action cubes before passing.\n")
//
//	p0.ActionCubes = 0
//	v = sn.View{}
//	c.Check(p0.validatePass(values, v), Equals, nil)
//}
//
//func (s *MySuite) TestPass(c *C) {
//	p := g1.Players()[0]
//	p.Passed = false
//	p.PerformedAction = false
//	p.pass()
//	c.Check(p.Passed, Equals, true)
//	c.Check(p.PerformedAction, Equals, true)
//}
