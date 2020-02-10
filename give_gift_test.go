package confucius

//import (
//	"fmt"
//
//	. "gopkg.in/check.v1"
//)
//
//func (s *MySuite) TestGiveGiftEntry(c *C) {
//	cp := g1.Players()[0]
//	op := g1.Players()[1]
//	gift := cp.GiftsBought[0]
//
//	e := cp.newGiveGiftEntry(op, gift, false)
//	c.Check(e.HTML(), Equals, fmt.Sprintf("%s gave value %d gift (%s) to %s.",
//		cp.Name(), gift.Value, gift.Name(), op.Name()))
//
//	e = cp.newGiveGiftEntry(op, gift, true)
//	c.Check(e.HTML(), Equals, fmt.Sprintf("%s gave value %d gift (%s) to %s and canceled gift from %s.",
//		cp.Name(), gift.Value, gift.Name(), op.Name(), op.Name()))
//}
