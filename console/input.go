package console

import (
	"bufio"
	"os"
	"strings"
)

func ReadStr(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	BluePrintln(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSuffix(text, "\n")
}
