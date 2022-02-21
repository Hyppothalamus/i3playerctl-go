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
        if strings.Contains(status, "Playing") {
            return player
        }
    }
    return ""
}

// format the current playing spotify song 
// to display with nerdfont glyph
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

// format current playing content in firefox
// to display with nerdfont glyph
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

    artist := strings.Split(artistOut.String(), "\n")[0]
    title := titleOut.String()
    result  := ""
    if artist == "\n" {
        if !strings.Contains(title, " - ") {
            return "  " + title
        }
        first := strings.Split(title, " - ")[0]
        second := strings.Split(title, " - ")[1]
        result = first + "  " + second
        return result
    }
    result = artist  + "  " + title
    return result
}

// format current playing content in vlc
// for display with nerdfont glyph
func formatVlc(player string) string {
    cmd := exec.Command("playerctl", "-p", player, "metadata", "xesam:url")
    var filePathOut bytes.Buffer
    cmd.Stdout = &filePathOut
    err := cmd.Run()
    if err != nil {
        log.Fatalf("failed to format metadata: %v", err)
    }
    // TODO: get part from url between last . and /
    filePath := filePathOut.String()
    lastIndex := strings.LastIndex(filePath, ".")
    firstIndex := strings.LastIndex(filePath, "/")
    file := filePath[(firstIndex + 1):lastIndex]
    return " 嗢" + file
}

func main() {
    // TODO: add flags and args so content is controllable
    // with this program widePeepoHappy
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
    case "vlc":
        fmt.Printf(formatVlc(player))
    default:
        os.Exit(0)
    }
}
