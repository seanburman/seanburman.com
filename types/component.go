package types

import "github.com/kitkitchen/fncmp"

type ComponentConfig struct {
	ID          string
	Class       string
	Style       string
	Children    []fncmp.Component
	Type        string
	Name        string
	Label       string
	For         string
	Value       string
	Placeholder string
	Options     []Option
	Required    string
}

type Option struct {
	Label string
	Value string
}
