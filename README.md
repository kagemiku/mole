# mole
[![Build Status](https://travis-ci.org/kagemiku/mole.svg?branch=master)](https://travis-ci.org/kagemiku/mole)

mole is a wrapper of **exec.Cmd** suite for executing piped commands.

# Usage
* **Mole.Output()**
```go
mole := mole.NewMole()
mole.Add("ls", "-la")
mole.Add("head", "-5")
mole.Add("wc", "-l")
mole.Add("tr", "-d", " ")

out, err := mole.Output()
if err != nil {
	log.Fatal(err)
}
fmt.Println(string(out))
```

* **Mole.Run()**
```go
mole := mole.NewMole()
mole.Add("wc", "-c")

var buf bytes.Buffer
buf.WriteString("hoge")
mole.Stdin = &buf
mole.Stdout = os.Stdout
if err := mole.Run(); err != nil {
	log.Fatal(err)
}
```

# License
MIT License

# Author
kagemiku
