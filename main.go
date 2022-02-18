package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func main() {
    cmd := exec.Command("playerctl", "-l")

    var out bytes.Buffer
    cmd.Stdout = &out

    err := cmd.Run()

    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("playerctl output: %q\n", out.String())
}
