package main

import "github.com/royge/gostcalc/cmd"

func main() {
	cmd.Register(cmd.RegisterFirestore)
	cmd.Execute()
}
