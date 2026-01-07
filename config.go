package rrmcp

type config struct {
	Host string
	Port int
}

func (c *config) InitDefaults() error {
	return nil
}
