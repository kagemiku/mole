package mole

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
)

type Mole struct {
	cmds     []*exec.Cmd
	Stdin    io.Reader
	Stdout   io.Writer
	executed bool
}

func NewMole() *Mole {
	return &Mole{
		cmds:     []*exec.Cmd{},
		executed: false,
	}
}

func (this *Mole) Add(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	this.cmds = append(this.cmds, cmd)
}

func (this *Mole) Run() error {
	return this.exec()
}

func (this *Mole) Output() ([]byte, error) {
	if this.Stdout != nil {
		return nil, errors.New("mole: Stdout already set")
	}

	var buf bytes.Buffer
	this.Stdout = &buf

	err := this.exec()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (this *Mole) exec1() error {
	cmd := this.cmds[0]
	cmd.Stdin = this.Stdin
	cmd.Stdout = this.Stdout

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

func (this *Mole) exec() error {
	if len(this.cmds) == 0 {
		return errors.New("mole: no commands are set")
	}

	if this.executed {
		return errors.New("mole: already started")
	}
	defer func(this *Mole) {
		this.executed = true
	}(this)

	if len(this.cmds) == 1 {
		return this.exec1()
	}

	var buf bytes.Buffer
	firstCmd := this.cmds[0]
	firstCmd.Stdin = this.Stdin
	firstCmd.Stdout = &buf

	if err := firstCmd.Start(); err != nil {
		return err
	}

	if err := firstCmd.Wait(); err != nil {
		return err
	}

	lastIdx := len(this.cmds) - 1
	for _, cmd := range this.cmds[1:lastIdx] {
		var bufcp bytes.Buffer
		io.Copy(&bufcp, &buf)

		cmd.Stdin = &bufcp
		cmd.Stdout = &buf

		if err := cmd.Start(); err != nil {
			return err
		}

		if err := cmd.Wait(); err != nil {
			return err
		}
	}

	lastCmd := this.cmds[lastIdx]
	lastCmd.Stdin = &buf
	lastCmd.Stdout = this.Stdout

	if err := lastCmd.Start(); err != nil {
		return err
	}

	return lastCmd.Wait()
}
