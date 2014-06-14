package game

type triangle struct {
  p1, p2, p3 *point
}

func (self *triangle) pointIsInside(p *point) bool{
  b1 := sign(p, self.p1, self.p2) < 0.0
  b2 := sign(p, self.p2, self.p3) < 0.0
  b3 := sign(p, self.p3, self.p1) < 0.0

  return ((b1 == b2) && (b2 == b3));
}

func sign(p1, p2, p3 *point) float64 {
  return (p1.x - p3.x) * (p2.y - p3.y) - (p2.x - p3.x) * (p1.y - p3.y)
}
