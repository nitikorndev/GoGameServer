package GoGameServer

const (
  BLACK = 0
  WHITE = 1
)

type Stone struct {
  Number   int
  Side     int
  Captured bool
}

