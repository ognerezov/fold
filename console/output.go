package console

import "fmt"

var (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)

func ColorPrintln(str string, color string) {
	fmt.Println(color + str + Reset)
}

func BluePrintln(str string) {
	ColorPrintln(str, Blue)
}

func MagentaPrintln(str string) {
	ColorPrintln(str, Magenta)
}

func CyanPrintln(str string) {
	ColorPrintln(str, Cyan)
}

func GreenPrintln(str string) {
	ColorPrintln(str, Green)
}

func RedPrintln(str string) {
	ColorPrintln(str, Red)
}
