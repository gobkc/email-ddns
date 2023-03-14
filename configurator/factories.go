package configurator

type Configurator interface {
	UnMarshal(path string, dest any) error
}

func Factory[T Configurator]() *T {
	v := new(T)
	return v
}
