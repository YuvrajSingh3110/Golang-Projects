package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
	"github.com/gen2brain/beeep"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

const (
	markName  = "CLI_Reminder_Tool"
	markvalue = "1"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <hh:mm> <Text Message>\n", os.Args[0])
		os.Exit(1)
	}

	now := time.Now()
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	t, err := w.Parse(os.Args[1], now)
	if err != nil {
		fmt.Println("Error while parsing: ",err)
		os.Exit(2)
	}
	if t == nil {
		fmt.Println("Unable to parse time")
		os.Exit(2)
	}
	if now.After(t.Time) {
		fmt.Println("Set a time that is in future")
		os.Exit(3)
	}

	diff := t.Time.Sub(now)
	if os.Getenv(markName) == markvalue {
		time.Sleep(diff)
		err := beeep.Alert("Reminder", strings.Join(os.Args[2:], " "), "assets/information.png")
		if err != nil {
			fmt.Println("Error in beeep: ", err)
			os.Exit(4)
		}
	} else {
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Env = append(os.Environ(), fmt.Sprintf("%s = %s", markName, markvalue))
		if err := cmd.Start(); err != nil {
			fmt.Println("Error while starting cmd: ", err)
			os.Exit(5)
		}
		fmt.Println("Remider will be shown after: ", diff.Round(time.Second))
		os.Exit(0)
	}
}
