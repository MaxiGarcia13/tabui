package shellinit

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteZsh(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, "zsh"); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "tabui-complete") {
		t.Fatalf("missing widget: %q", out)
	}
	if !strings.Contains(out, "bindkey '^I'") {
		t.Fatalf("missing bindkey: %q", out)
	}
}

func TestWriteBash(t *testing.T) {
	var buf bytes.Buffer
	if err := Write(&buf, "bash"); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, `bind -x '"\t": tabui-complete'`) {
		t.Fatalf("missing bind: %q", out)
	}
}

func TestWriteUnknown(t *testing.T) {
	if err := Write(&bytes.Buffer{}, "fish"); err == nil {
		t.Fatal("expected error")
	}
}

func TestWriteMissing(t *testing.T) {
	if err := Write(&bytes.Buffer{}, ""); err == nil {
		t.Fatal("expected error")
	}
}
