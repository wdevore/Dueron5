package cell

// StraightConnection has no delay. On each time mark data immediately
// appears on the output.
type StraightConnection struct {
	baseConn
}

func NewStraightConnection() IConnection {
	sc := new(StraightConnection)
	sc.baseConn.initialize()
	return sc
}

// IConnection implementations.

// Input ORs the data value to the connection
func (sc *StraightConnection) Input(b int) {
	sc.value = sc.value | b
}

func (sc *StraightConnection) Output() int {
	return sc.value
}

func (sc *StraightConnection) Post() {
	sc.value = 0
}
