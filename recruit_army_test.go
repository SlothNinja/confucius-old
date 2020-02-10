package confucius

//import . "gopkg.in/check.v1"
//
////func (p *Player) canAffordAnArmy() bool {
////        return p.ConCardHand.Licenses() > p.armyCost()
////}
//
//func (s *MySuite) TestCanAffordAnArmy(c *C) {
//	p0 := g1.Players()[0]
//
//	// Has No Bingbu Influence but 6 licences
//	c.Check(p0.InfluenceIn(g1.Ministries[Bingbu]), Equals, 0)
//	c.Check(p0.ConCardHand.Licenses(), Equals, 6)
//	c.Check(p0.canAffordAnArmy(), Equals, true)
//
//	// Has No Bingbu Influence but only 4 licences
//	p0.ConCardHand.Remove(&ConCard{Coins: 2})
//	c.Check(p0.InfluenceIn(g1.Ministries[Bingbu]), Equals, 0)
//	c.Check(p0.ConCardHand.Licenses(), Equals, 4)
//	c.Check(p0.canAffordAnArmy(), Equals, false)
//
//	// Has Bingbu Influence and 4 licences
//	var seniority Seniority = 3
//	g1.Ministries[Bingbu].Officials[seniority].setPlayer(p0)
//	c.Check(p0.InfluenceIn(g1.Ministries[Bingbu]), Equals, 1)
//	c.Check(p0.ConCardHand.Licenses(), Equals, 4)
//	c.Check(p0.canAffordAnArmy(), Equals, true)
//
//	// Has Bingbu Influence but only 1 licence
//	p0.ConCardHand.Remove(&ConCard{Coins: 1})
//	c.Check(p0.InfluenceIn(g1.Ministries[Bingbu]), Equals, 1)
//	c.Check(p0.ConCardHand.Licenses(), Equals, 1)
//	c.Check(p0.canAffordAnArmy(), Equals, false)
//}
