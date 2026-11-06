package suggest

type Suggestion struct {
  Label  string
  Detail string
  Value  string
}

type Context struct {
  Line   string
  Cursor int
  Cwd    string
}
