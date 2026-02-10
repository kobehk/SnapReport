package geo

import "testing"

func TestClassifyHighwayByCategory(t *testing.T) {
	tests := []struct {
		cat  string
		want bool
	}{
		{"motorway", true},
		{"trunk", true},
		{"motorway_link", true},
		{"trunk_link", true},
		{"primary", false},
		{"secondary", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := ClassifyHighway(tt.cat, ""); got != tt.want {
			t.Fatalf("cat %q => %v, want %v", tt.cat, got, tt.want)
		}
	}
}

func TestClassifyHighwayByRoadName(t *testing.T) {
	tests := []struct {
		road string
		want bool
	}{
		{"沪宁高速公路", true},
		{"广深沿江高速", true},
		{"某某Expressway", true},
		{"滨河大道", false},
		{"深南大道", false},
	}
	for _, tt := range tests {
		if got := ClassifyHighway("", tt.road); got != tt.want {
			t.Fatalf("road %q => %v, want %v", tt.road, got, tt.want)
		}
	}
}
