package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
    "strings"
)

// returns the current plater based on which players are open
// and wich player is actually playing something
func getCurrentPlayer() string {
    cmd := exec.Command("playerctl", "-l")
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        log.Fatalf("Error when checking wich player is running: %v", err)
    }
    players := strings.Split(out.String(), "\n")
    if len(players) > 0 {
        players = players[:len(players) - 1]
    }
    for i := range players {
        player := players[i]
        cmd := exec.Command("playerctl", "-p", player, "status")
        var statusOut bytes.Buffer
        cmd.Stdout = &statusOut
        err := cmd.Run()
        if err != nil {
            log.Fatalf("error checking wich open player is playing: %v", err)
        }
        status := statusOut.String()
        //fmt.Printf("\n status from %v: %v", player, statusOut.String())
        // fmt.Printf("the current players: %v\n", players[i])
        if strings.Contains(status, "Playing") {
            return player
        }
    }
    return ""
}

func main() {
    fmt.Printf("current playing player: %v", getCurrentPlayer())
}
