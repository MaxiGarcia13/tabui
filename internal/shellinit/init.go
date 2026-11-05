package shellinit

import (
	"embed"
	"fmt"
	"io"
)

//go:embed tabui.zsh tabui.bash
var scripts embed.FS

func Write(w io.Writer, shell string) error {
	var name string
	switch shell {
	case "zsh":
		name = "tabui.zsh"
	case "bash":
		name = "tabui.bash"
	case "":
		return fmt.Errorf("usage: tabui init <zsh|bash>")
	default:
		return fmt.Errorf("unknown shell %q (want zsh or bash)", shell)
	}

	data, err := scripts.ReadFile(name)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
