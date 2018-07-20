package mole

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
)

// Mole is a structure for running piped command
type Mole struct {
	cmds     []*exec.Cmd
	Stdin    io.Reader
	Stdout   io.Writer
	executed bool
}

// NewMole generates a new Mole instance
func NewMole() *Mole {
	return &Mole{
		cmds:     []*exec.Cmd{},
		executed: false,
	}
}

// Add adds command to command list
func (m *Mole) Add(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	m.cmds = append(m.cmds, cmd)
}

// Run runs command list
//
// If you want to get the stdout content, please use Output instead, or pass some io.Writer instance to the Mole.Stdout
func (m *Mole) Run() error {
	return m.exec()
}

// Output runs command list, and returns stdout content as []byte
func (m *Mole) Output() ([]byte, error) {
	if m.Stdout != nil {
		return nil, errors.New("mole: Stdout already set")
	}

	var buf bytes.Buffer
	m.Stdout = &buf

	err := m.exec()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// exec1 executes command. This method is only used when command list length is 1.
func (m *Mole) exec1() error {
	cmd := m.cmds[0]
	cmd.Stdin = m.Stdin
	cmd.Stdout = m.Stdout

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

// exec executes commands.
func (m *Mole) exec() error {
	if len(m.cmds) == 0 {
		return errors.New("mole: no commands are set")
	}

	if m.executed {
		return errors.New("mole: already started")
	}
	defer func(m *Mole) {
		m.executed = true
	}(m)

	if len(m.cmds) == 1 {
		return m.exec1()
	}

	var buf bytes.Buffer
	firstCmd := m.cmds[0]
	firstCmd.Stdin = m.Stdin
	firstCmd.Stdout = &buf

	if err := firstCmd.Start(); err != nil {
		return err
	}

	if err := firstCmd.Wait(); err != nil {
		return err
	}

	lastIdx := len(m.cmds) - 1
	for _, cmd := range m.cmds[1:lastIdx] {
		var bufcp bytes.Buffer
		_, err := io.Copy(&bufcp, &buf)
		if err != nil {
			return err
		}

		cmd.Stdin = &bufcp
		cmd.Stdout = &buf

		if err := cmd.Start(); err != nil {
			return err
		}

		if err := cmd.Wait(); err != nil {
			return err
		}
	}

	lastCmd := m.cmds[lastIdx]
	lastCmd.Stdin = &buf
	lastCmd.Stdout = m.Stdout

	if err := lastCmd.Start(); err != nil {
		return err
	}

	return lastCmd.Wait()
}
