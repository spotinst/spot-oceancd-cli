package utils

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/mbndr/figlet4go"
)

func WaitSpinner() {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond) // Build our new spinner
	s.Start()                                                   // Start the spinner
	time.Sleep(4 * time.Second)                                 // Run for some time to simulate work
	s.Stop()
}

func MessageWithProgress(msg string, sec int) {

	s := spinner.New(spinner.CharSets[22], time.Duration(100)*time.Millisecond) // Build our new spinner
	//s.Prefix = " "
	s.Prefix = msg + " "
	s.FinalMSG = color.GreenString("\x08\x08v\n")

	s.Start()                                    // Start the spinner
	time.Sleep(time.Duration(sec) * time.Second) // Run for some time to simulate work

	s.Stop()

}

func RenderTerminalString(message string) string {
	ascii := figlet4go.NewAsciiRender()

	// The underscore would be an error
	renderStr, _ := ascii.Render(message)
	return renderStr
}
