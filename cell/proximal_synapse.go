package cell

type ProximalSynapse struct {
	baseSynapse
}

// func NewProximalSynapse() ISynapse {
// 	ps := new(ProximalSynapse)
// 	return ps
// }

func (n *ProximalSynapse) Integrate(dt float64) float64 {
	return 0.0
}

func (n *ProximalSynapse) Process() {
}
