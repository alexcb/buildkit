package llb

func LocalHost() State {
	attrs := map[string]string{}
	cons := Constraints{}

	source := NewSource("localhost://", attrs, cons)
	return NewState(source.Output())
}
