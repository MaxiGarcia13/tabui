package suggest

type Provider interface {
  Match(Context) bool
  Suggest(Context) []Suggestion
}

type Funcs struct {
  MatchFn   func(Context) bool
  SuggestFn func(Context) []Suggestion
}

func (funcs Funcs) Match(ctx Context) bool {
  if funcs.MatchFn == nil {
    return false
  }
  return funcs.MatchFn(ctx)
}

func (funcs Funcs) Suggest(ctx Context) []Suggestion {
  if funcs.SuggestFn == nil {
    return nil
  }
  return funcs.SuggestFn(ctx)
}

func FirstMatch(ctx Context, providers ...Provider) []Suggestion {
  for _, provider := range providers {
    if provider.Match(ctx) {
      return provider.Suggest(ctx)
    }
  }
  return nil
}
