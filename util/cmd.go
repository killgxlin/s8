package util

type Cmd struct {
	msg  string
	ctrl string
}

func NewCmd() *Cmd {
	return &Cmd{}
}

func (c *Cmd) Parse(cmd string) error {
	return nil
}

func Marshal(c *Cmd) ([]byte, error) {
	return nil, nil
}

func Unuarshal([]byte) (*Cmd, error) {
	return nil, nil
}
