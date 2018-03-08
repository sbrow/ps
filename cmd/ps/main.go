package main

import (
	"fmt"
	"github.com/sbrow/ps"
	"os"
)

func main() {
	args := []string{}
	cmd := ""
	switch {
	case len(os.Args) > 1:
		args = os.Args[2:]
		fallthrough
	case len(os.Args) > 0:
		cmd = os.Args[1]
	}

	fmt.Println(os.Args, cmd, args)
	if cmd == "action" {
		err := ps.DoAction(args[0], args[1])
		if err != nil {
			panic(err)
		}
	}
}
