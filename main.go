package main

import "github.com/blaggotech/gostcalc/cmd"

func main() {
	cmd.Register(cmd.RegisterFirestore)
	cmd.Execute()
}
