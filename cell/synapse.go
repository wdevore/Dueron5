package cell

// Weight behavior

// ISynapse is a common interface.
// A synapse is associated with a connection and one or more behaviors.
// One such behavior is STDP for evaluating the Long-term effect.
// Each synapse is part of a group called a Compartment.
type ISynapse interface {
	Id() int

	Connect(IConnection)
	GetConnection() IConnection

	// Evaluates the total effective weight for the synapse.
	Integrate(t float64, cell ICell) float64

	// Called before Integrate()
	Process()

	// Called at the end of a simulation run (not a pass)
	PostProcess()

	Reset()

	// The current input on the synapse. The input is feed from an IConnection's
	// output.
	Input() int

	IsExcititory() bool

	Load(json interface{})
	ToJSON() interface{}

	SetWeight(float64)
	SetWMax(float64)
	SetWMin(float64)

	SetField(field string, value string)
	SetTaoP(float64)
	SetTaoN(float64)
	SetAma(float64)
	SetAmb(float64)
	SetTsw(float64)
}

type SynapseType bool

const (
	Excititory SynapseType = true
	Inhibitory             = false
)

type baseSynapse struct {
	id int

	synType SynapseType

	// This is the weight at any point in time.
	w float64

	// This is the weight at which "w" decays towards both for short
	// and long term rules
	wi float64

	wMax float64
	wMin float64

	// A synapse will read this input on each integration pass
	conn IConnection

	// The compartment this synaspe resides in.
	comp ICompartment
}

func (bs *baseSynapse) initialize() {
}

func (bs *baseSynapse) IsExcititory() bool {
	return bs.synType == Excititory
}

func (bs *baseSynapse) Connect(con IConnection) {
	bs.conn = con
}

func (bs *baseSynapse) GetConnection() IConnection {
	return bs.conn
}

// Input returns the input feeding into this synapse.
// The synapse has an IConnection for input. This method returns
// that connection's output.
// i.e. the data entering the synapse from the connection.
func (bs *baseSynapse) Input() int {
	return bs.conn.Output()
}

// ----------------------------------------------------------------------
// Properties
// ----------------------------------------------------------------------

func (bs *baseSynapse) Id() int {
	return bs.id
}

func (bs *baseSynapse) SetId(id int) {
	bs.id = id
}

func (bs *baseSynapse) SetWeight(weight float64) {
	bs.w = weight
}

func (bs *baseSynapse) SetWMax(v float64) {
	bs.wMax = v
}

func (bs *baseSynapse) SetWMin(v float64) {
	bs.wMin = v
}
