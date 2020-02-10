package confucius

//import . "gopkg.in/check.v1"

//func (cs *GiftCards) Remove(cards ...*GiftCard) {
//	*cs = cs.removeMulti(cards...)
//}
//
//func (cs GiftCards) removeMulti(cards ...*GiftCard) GiftCards {
//        gcs := cs
//	for _, c := range cs {
//		gcs = gcs.remove(c)
//	}
//	return gcs
//}
//
//func (cs GiftCards) remove(card *GiftCard) GiftCards {
//        cards := cs
//	for i, c := range cs {
//		if c.Equal(card) {
//			return cards.removeAt(i)
//		}
//	}
//	return cs
//}
//
//func (cs GiftCards) removeAt(i int) GiftCards {
//	return append(cs[:i], cs[i+1:]...)
//}

//func (s *MySuite) TestGiftCardsRemove(c *C) {
//	p0 := g1.Players()[0]
//	gift := p0.GiftCardHand[0]
//
//	// Has one GiftCard, but tries to remove nothing
//	c.Check(len(p0.GiftCardHand), Equals, 5)
//	c.Check(len(p0.GiftsBought), Equals, 1)
//	p0.GiftsBought.Remove()
//	c.Check(len(p0.GiftsBought), Equals, 1)
//
//	// Has one GiftCard, but tries to remove nil
//	c.Check(len(p0.GiftCardHand), Equals, 5)
//	c.Check(len(p0.GiftsBought), Equals, 1)
//	p0.GiftsBought.Remove(nil)
//	c.Check(len(p0.GiftsBought), Equals, 1)
//
//	// Has one GiftCard, but tries to remove a different card
//	c.Check(len(p0.GiftCardHand), Equals, 5)
//	c.Check(len(p0.GiftsBought), Equals, 1)
//	c.Check(p0.GiftsBought.include(gift), Equals, false)
//	p0.GiftsBought.Remove(gift)
//	c.Check(len(p0.GiftsBought), Equals, 1)
//
//	// Has one GiftCard and removes it
//	gift = p0.GiftsBought[0]
//	c.Check(len(p0.GiftCardHand), Equals, 5)
//	c.Check(len(p0.GiftsBought), Equals, 1)
//	c.Check(p0.GiftsBought.include(gift), Equals, true)
//	p0.GiftsBought.Remove(gift)
//	c.Check(p0.GiftsBought.include(gift), Equals, false)
//	c.Check(len(p0.GiftsBought), Equals, 0)
//
//	// Has multi cards and remove one in the middle
//	gift = p0.GiftCardHand[3]
//	c.Check(len(p0.GiftCardHand), Equals, 5)
//	c.Check(p0.GiftCardHand.include(gift), Equals, true)
//	p0.GiftCardHand.Remove(gift)
//	c.Check(p0.GiftCardHand.include(gift), Equals, false)
//	c.Check(len(p0.GiftCardHand), Equals, 4)
//
//	// Has multi cards and remove one in the end
//	gift = p0.GiftCardHand[3]
//	c.Check(len(p0.GiftCardHand), Equals, 4)
//	c.Check(p0.GiftCardHand.include(gift), Equals, true)
//	p0.GiftCardHand.Remove(gift)
//	c.Check(p0.GiftCardHand.include(gift), Equals, false)
//	c.Check(len(p0.GiftCardHand), Equals, 3)
//
//	// Has multi cards and remove one in the beginning
//	gift = p0.GiftCardHand[0]
//	c.Check(len(p0.GiftCardHand), Equals, 3)
//	c.Check(p0.GiftCardHand.include(gift), Equals, true)
//	p0.GiftCardHand.Remove(gift)
//	c.Check(p0.GiftCardHand.include(gift), Equals, false)
//	c.Check(len(p0.GiftCardHand), Equals, 2)
//}
