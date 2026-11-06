package clihelp

import (
	"errors"
	"testing"
)

func TestParseTarget(t *testing.T) {
  previous := lookPath
  lookPath = func(file string) (string, error) {
    if file == "docker" {
      return "/bin/docker", nil
    }
    return "", errors.New("not found")
  }
  t.Cleanup(func() { lookPath = previous })

  tests := []struct {
    name         string
    line         string
    wantOK       bool
    wantHelpArgs []string
    wantPartial  string
  }{
    {
      name:         "root space",
      line:         "docker ",
      wantOK:       true,
      wantHelpArgs: nil,
      wantPartial:  "",
    },
    {
      name:         "partial subcommand",
      line:         "docker c",
      wantOK:       true,
      wantHelpArgs: nil,
      wantPartial:  "c",
    },
    {
      name:         "compose flag partial",
      line:         "docker compose -",
      wantOK:       true,
      wantHelpArgs: []string{"compose"},
      wantPartial:  "-",
    },
    {
      name:         "strip complete flags",
      line:         "docker run -d ",
      wantOK:       true,
      wantHelpArgs: []string{"run"},
      wantPartial:  "",
    },
    {
      name:         "flag as partial keeps parent",
      line:         "docker run -d",
      wantOK:       true,
      wantHelpArgs: []string{"run"},
      wantPartial:  "-d",
    },
  }

  for _, testCase := range tests {
    t.Run(testCase.name, func(t *testing.T) {
      got := parseTarget(testCase.line)
      if got.ok != testCase.wantOK {
        t.Fatalf("ok = %v, want %v", got.ok, testCase.wantOK)
      }
      if got.partial != testCase.wantPartial {
        t.Fatalf("partial = %q, want %q", got.partial, testCase.wantPartial)
      }
      if len(got.helpArgs) != len(testCase.wantHelpArgs) {
        t.Fatalf(
          "helpArgs = %#v, want %#v",
          got.helpArgs,
          testCase.wantHelpArgs,
        )
      }
      for index, arg := range testCase.wantHelpArgs {
        if got.helpArgs[index] != arg {
          t.Fatalf(
            "helpArgs = %#v, want %#v",
            got.helpArgs,
            testCase.wantHelpArgs,
          )
        }
      }
    })
  }
}
