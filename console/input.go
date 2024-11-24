package console

import (
	"bufio"
	"os"
)

func ReadStr(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	BluePrintln(prompt)
	text, _ := reader.ReadString('\n')
	return text
}
