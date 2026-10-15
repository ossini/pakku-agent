package tools

import (
	"fmt"
	"slices"
)

type Permission string

const (
	PermRead      Permission = "read"
	PermWrite     Permission = "write"
	PermDangerous Permission = "dangerous"
)

type Tool struct {
	Name        string
	Title       string
	Description string
	Permission  Permission
	Parameters  map[string]any
	Run         func(args string) (string, error)
}

var allTools []Tool

func register(t Tool) {
	allTools = append(allTools, t)
}

type Manager struct {
	Tools []Tool
}

func NewManager(allowed ...Permission) *Manager {
	m := &Manager{}
	for _, t := range allTools {
		if slices.Contains(allowed, t.Permission) {
			m.Tools = append(m.Tools, t)
		}
	}
	return m
}

func (m *Manager) Run(name, args string) (string, error) {
	for _, t := range m.Tools {
		if t.Name == name {
			return t.Run(args)
		}
	}
	return "", fmt.Errorf("unknown or not allowed tool: %s", name)
}

func (m *Manager) Definitions() []map[string]any {
	var defs []map[string]any
	for _, t := range m.Tools {
		defs = append(defs, map[string]any{
			"type":        "function",
			"name":        t.Name,
			"description": t.Description,
			"parameters":  t.Parameters,
		})
	}
	return defs
}
