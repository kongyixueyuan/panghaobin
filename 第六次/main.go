package main

import "os"

func main() {
	os.Setenv("NODE_ID", "1")
	cli := PHBCLI{}
	cli.PHBRun()
}
