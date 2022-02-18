package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
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
    if out.String() == "No players where found" {
        return ""
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

func formatSpotify(player string) string {
    cmd := exec.Command("playerctl", "-p", player, "metadata", "--format", "{{artist}}  {{title}}")
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        log.Fatalf("failed to format metadata: %v", err)
    }
    return out.String()
}

func formatFirefox(player string) string {
    cmd := exec.Command("playerctl", "-p", player, "metadata", "xesam:title")
    var titleOut bytes.Buffer
    cmd.Stdout = &titleOut
    err := cmd.Run()
    if err != nil {
        log.Fatalf("failed to format metadata: %v", err)
    }
    cmd = exec.Command("playerctl", "-p", player, "metadata", "xesam:artist")
    var artistOut bytes.Buffer
    cmd.Stdout = &artistOut
    err = cmd.Run()
    if err != nil {
        log.Fatalf("failed to format metadata: %v", err)
    }

    artist := artistOut.String()
    title := titleOut.String()
    result  := ""
    if artist == "\n" {
        first := strings.Split(title, " - ")[0]
        second := strings.Split(title, " - ")[1]
        result = first + "  " + second
        return result
    }
    result = artist  + "  " + title
    return result
}

func main() {
    player := getCurrentPlayer()
    if player == "" {
        log.Fatalf("no current player playing")
    }
    if strings.Contains(player, "instance") {
        player = strings.Split(player, ".")[0]
    }
    // TODO: select correct format for formating player for display
    switch player {
    case "spotify":
        fmt.Printf(formatSpotify(player))
    case "firefox":
        fmt.Printf(formatFirefox(player))
    default:
        os.Exit(0)
    }
}
