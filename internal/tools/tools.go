package tools

import (
	"context"
	"fmt"
	"slices"
)

type Permission string

const (
	PermRead      Permission = "read"
	PermWrite     Permission = "write"
	PermDangerous Permission = "dangerous"
)

type Result struct {
	Text     string
	ImageURL string
}

type Tool struct {
	Name        string
	Title       string
	Description string
	Permission  Permission
	Parameters  map[string]any
	Run         func(ctx context.Context, args string) (Result, error)
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

func (m *Manager) Run(ctx context.Context, name, args string) (Result, error) {
	t, ok := m.find(name)
	if !ok {
		return Result{}, fmt.Errorf("unknown or not allowed tool: %s", name)
	}
	return t.Run(ctx, args)
}

func (m *Manager) Title(name string) string {
	if t, ok := m.find(name); ok {
		return t.Title
	}
	return name
}

func (m *Manager) find(name string) (Tool, bool) {
	i := slices.IndexFunc(m.Tools, func(t Tool) bool { return t.Name == name })
	if i < 0 {
		return Tool{}, false
	}
	return m.Tools[i], true
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
