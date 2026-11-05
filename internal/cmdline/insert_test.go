package cmdline

import "testing"

func TestInsertToken(t *testing.T) {
  tests := []struct {
    name  string
    left  string
    token string
    want  string
  }{
    {
      name:  "after tool space",
      left:  "pnpm ",
      token: "run",
      want:  "pnpm run",
    },
    {
      name:  "after run word",
      left:  "pnpm run",
      token: "build",
      want:  "pnpm run build",
    },
    {
      name:  "after run space",
      left:  "pnpm run ",
      token: "build",
      want:  "pnpm run build",
    },
    {
      name:  "replace partial script",
      left:  "pnpm run bu",
      token: "build",
      want:  "pnpm run build",
    },
    {
      name:  "replace partial command",
      left:  "npm in",
      token: "install",
      want:  "npm install",
    },
    {
      name:  "script auto run after tool",
      left:  "pnpm ",
      token: "run build",
      want:  "pnpm run build",
    },
    {
      name:  "script auto run replacing partial",
      left:  "pnpm bu",
      token: "run build",
      want:  "pnpm run build",
    },
  }

  for _, testCase := range tests {
    t.Run(testCase.name, func(t *testing.T) {
      got := InsertToken(testCase.left, testCase.token)
      if got != testCase.want {
        t.Fatalf(
          "InsertToken(%q, %q) = %q, want %q",
          testCase.left,
          testCase.token,
          got,
          testCase.want,
        )
      }
    })
  }
}
