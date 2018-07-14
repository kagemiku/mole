package mole

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
)

type Mole struct {
	cmds     []*exec.Cmd
	Stdin    io.Reader
	Stdout   io.Writer
	finished bool
}

func NewMole() *Mole {
	return &Mole{
		cmds:     []*exec.Cmd{},
		finished: false,
	}
}

func (this *Mole) Add(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	this.cmds = append(this.cmds, cmd)
}

func (this Mole) Start() error {
	err := this.exec()
	if err != nil {
		return err
	}

	return nil
}

func (this Mole) Output() ([]byte, error) {
	return this.cmds[0].Output()
}

func (this Mole) Debug() {
	fmt.Printf("Length: %d\n", len(this.cmds))
}

func (this Mole) exec() error {
	if len(this.cmds) < 3 {
		return errors.New("len < 3")
	}

	var buf bytes.Buffer
	firstCmd := this.cmds[0]
	firstCmd.Stdin = this.Stdin
	firstCmd.Stdout = &buf

	if err := firstCmd.Start(); err != nil {
		log.Println("Failed to start first cmd")
		return err
	}

	if err := firstCmd.Wait(); err != nil {
		log.Println("Failed to wait first cmd")
		return err
	}

	lastIdx := len(this.cmds) - 1
	for _, cmd := range this.cmds[1:lastIdx] {
		var bufcp bytes.Buffer
		io.Copy(&bufcp, &buf)

		cmd.Stdin = &bufcp
		cmd.Stdout = &buf

		if err := cmd.Start(); err != nil {
			log.Println("Failed to start cmd")
			return err
		}

		if err := cmd.Wait(); err != nil {
			log.Println("Failed to wait cmd")
			return err
		}
	}

	lastCmd := this.cmds[lastIdx]
	lastCmd.Stdin = &buf
	lastCmd.Stdout = this.Stdout

	if err := lastCmd.Start(); err != nil {
		log.Println("Failed to start last cmd")
		return err
	}

	if err := lastCmd.Wait(); err != nil {
		log.Println("Failed to wait last cmd")
		return err
	}

	return nil
}
