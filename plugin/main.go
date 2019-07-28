package main

import (
	"context"
	"fmt"
)

type Endpoint func(context.Context, string) string

type Middleware interface {
	Wrap(Endpoint) Endpoint
}

type MiddlewareGroup struct {
	middlewares []Middleware
}

func NewMiddlewreGroup(mws ...Middleware) *MiddlewareGroup {
	return &MiddlewareGroup{mws}
}

func (g *MiddlewareGroup) Wrap(innerFunc Endpoint)  Endpoint {
	for _, middleware := range g.middlewares {
		innerFunc = middleware.Wrap(innerFunc)
	}
	return innerFunc
}

type Prefixer struct {
	prefix string
}

func NewPrefixer(prefix string) *Prefixer {
	return &Prefixer{prefix}
}

func (p *Prefixer) Wrap(innerFunc Endpoint)  Endpoint {
	return func(ctx context.Context, word string) string {
		var out = innerFunc(ctx, word)
		if GetState(ctx).IsSensitive() {
			return out
		}
		return p.prefix + out
	}
}

type Suffixer struct {
	suffix string
}

func NewSuffixer(suffix string) *Suffixer {
	return &Suffixer{suffix}
}

func (s *Suffixer) Wrap(innerFunc Endpoint)  Endpoint {
	return func(ctx context.Context, word string) string {
		var out = innerFunc(ctx, word)
		if GetState(ctx).IsSensitive() {
			return out
		}
		return out + s.suffix
	}
}

func Echo(ctx context.Context, word string) string {
	return word
}

type State struct {
	sensitive bool
}

func (s *State) MarkSensitive(){
	s.sensitive = true
}

func (s *State) IsSensitive() bool {
	return s.sensitive
}

func GetState(ctx context.Context) *State {
	return ctx.Value("state").(*State)
}

type SensitiveDetector struct {
	blackWords []string
}

func NewSensitiveDetector(blackWords []string) *SensitiveDetector {
	return &SensitiveDetector{blackWords}
}

func (d *SensitiveDetector) Wrap(innerFunc Endpoint) Endpoint {
	return func(ctx context.Context, word string) string {
		state := &State{false}
		ctx = context.WithValue(ctx, "state", state)
		for _, w := range d.blackWords {
			if word == w {
				state.MarkSensitive()
				return word
			}
		}
		return innerFunc(ctx, word)
	}
}

func main() {
	var group = NewMiddlewreGroup(
		NewPrefixer("hello "),
		NewSuffixer(" coders"),
		NewSensitiveDetector([]string{"sex", "politics"}))
	var endpoint = Echo
	endpoint = group.Wrap(endpoint)
	var ctx = context.Background()
	fmt.Println(endpoint(ctx, "sex"))
	fmt.Println(endpoint(ctx, "politics"))
	fmt.Println(endpoint(ctx, "bytedance"))
}
