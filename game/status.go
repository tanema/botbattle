package game

type Status struct {
  Id        int     `json:"id"`
  X         int     `json:"x"`
  Y         int     `json:"y"`
  Rotation  int     `json:"rotation"`
  Name      string  `json:"name"`
  Health    int     `json:"health"`
}
