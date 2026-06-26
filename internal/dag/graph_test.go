package dag

import (
	"fmt"
	"testing"

	"github.com/hyturing/compose-pulse/internal/compose"
)

func TestBuild_LinearChain(t *testing.T) {
	cfg := &compose.Config{
		Services: map[string]compose.Service{
			"postgres": {},
			"api":      {DependsOn: compose.DependsOn{"postgres": {}}},
			"frontend": {DependsOn: compose.DependsOn{"api": {}}},
		},
	}
	g, err := Build(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Roots) != 1 || g.Roots[0].Name != "postgres" {
		t.Errorf("expected single root 'postgres', got %v", g.Roots)
	}
	if g.ByName["api"].Level != 1 {
		t.Errorf("api should be level 1, got %d", g.ByName["api"].Level)
	}
	if g.ByName["frontend"].Level != 2 {
		t.Errorf("frontend should be level 2, got %d", g.ByName["frontend"].Level)
	}
}

func TestBuild_Diamond(t *testing.T) {
	// postgres, redis -> api, worker -> frontend
	cfg := &compose.Config{
		Services: map[string]compose.Service{
			"postgres": {},
			"redis":    {},
			"api":      {DependsOn: compose.DependsOn{"postgres": {}, "redis": {}}},
			"worker":   {DependsOn: compose.DependsOn{"postgres": {}, "redis": {}}},
			"frontend": {DependsOn: compose.DependsOn{"api": {}}},
		},
	}
	g, err := Build(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Roots) != 2 {
		t.Errorf("expected 2 roots, got %d", len(g.Roots))
	}
	if len(g.Ordered) != 5 {
		t.Errorf("expected 5 nodes ordered, got %d", len(g.Ordered))
	}
}

func TestBuild_CycleDetection(t *testing.T) {
	cfg := &compose.Config{
		Services: map[string]compose.Service{
			"a": {DependsOn: compose.DependsOn{"b": {}}},
			"b": {DependsOn: compose.DependsOn{"a": {}}},
		},
	}
	_, err := Build(cfg)
	if err == nil {
		t.Error("expected error for circular dependency, got nil")
	}
}

func TestBuild_IndependentRoots(t *testing.T) {
	cfg := &compose.Config{
		Services: map[string]compose.Service{
			"svc1": {},
			"svc2": {},
			"svc3": {},
		},
	}
	g, err := Build(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Roots) != 3 {
		t.Errorf("expected 3 roots, got %d: %v", len(g.Roots), func() []string {
			names := []string{}
			for _, r := range g.Roots {
				names = append(names, r.Name)
			}
			return names
		}())
	}
}

func TestBuild_Testdata(t *testing.T) {
	cfg, err := compose.Parse("../../testdata/docker-compose.yml")
	if err != nil {
		t.Skip("testdata not available:", err)
	}
	g, err := Build(cfg)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("\n  Topological order (%d services):\n", len(g.Ordered))
	for _, n := range g.Ordered {
		fmt.Printf("    level=%d  %-12s  deps=%v\n", n.Level, n.Name, n.Deps)
	}
}
