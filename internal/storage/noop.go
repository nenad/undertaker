package storage

type NoopGravedigger struct {
	funcs []string
}

// Bury buries the given funcs
func (g *NoopGravedigger) Bury(funcs []string) error {
	g.funcs = funcs
	return nil
}

// Dig returns all buried funcs
func (g *NoopGravedigger) Dig() ([]string, error) {
	return g.funcs, nil
}
