package mole

import (
	"log"
	"testing"
)

func TestSingleRun(t *testing.T) {
	mole := NewMole()

	if err := mole.Run(); err == nil {
		t.Fatal("It should be occur error when Run is called with no commands")
	}

	mole.Add("ls", "-la")
	if err := mole.Run(); err != nil {
		t.Fatal("Failed to run Run() with 1 command")
	}

	if err := mole.Run(); err == nil {
		t.Fatal("It should be occur error when Run is called multiple times")
	}
}

func TestSingleOutput(t *testing.T) {
	mole := NewMole()

	mole.Add("ls", "-la")
	out, err := mole.Output()
	if err != nil {
		t.Fatal("Failed to run Output() with 1 command")
	}
	log.Println(string(out))

	out, err = mole.Output()
	if err == nil {
		t.Fatal("It should be occur error when Output is called multiple times")
	}
}

func TestMultipleRun(t *testing.T) {
	mole := NewMole()
	mole.Add("ls", "-la")
	mole.Add("head", "-5")
	mole.Add("wc", "-l")
	mole.Add("tr", "-d", " ")

	if err := mole.Run(); err != nil {
		t.Fatal("Failed to run Run() with 4 command")
	}

	if err := mole.Run(); err == nil {
		t.Fatal("It should be occur error when Run is called multiple times")
	}
}

func TestMultipleOutput(t *testing.T) {
	mole := NewMole()
	mole.Add("ls", "-la")
	mole.Add("head", "-5")
	mole.Add("wc", "-l")
	mole.Add("tr", "-d", " ")

	out, err := mole.Output()
	if err != nil {
		t.Fatal("Failed to run Output() with 1 command")
	}
	log.Println(string(out))

	out, err = mole.Output()
	if err == nil {
		t.Fatal("It should be occur error when Output is called multiple times")
	}
}
