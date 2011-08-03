package GoGameServer

type Game struct {
  Id    int64
  White *User
  Black *User
  Stones []Stone

  GameOver bool

  Komi float
}

