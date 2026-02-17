package main

import ( "flag" "fmt" "os")

func main() {
	cmd:= flag.String("cmd", "", "command: publish")
	flag.parse()

	switch *cmd {
	case "publish":
		fmt.Println("publish: not implemented yet (v0.1 wip)")
	default:
		fmt.Println("usage: blobrs -cmd publish [flags...]")
		os.Exit(2)
	}
}