package conversation

import (
	"log"
	"os"
)

func init() {
	// Configure logger to not print timestamps and other metadata
	log.SetFlags(0)
	// Set output to stdout for colored output support
	log.SetOutput(os.Stdout)
}

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
)

var colorList = []string{
	ColorRed,
	ColorGreen,
	ColorYellow,
	ColorBlue,
	ColorPurple,
	ColorCyan,
}

// GetParticipantColor returns a color for a participant based on their index
func GetParticipantColor(index int) string {
	return colorList[index%len(colorList)]
}
