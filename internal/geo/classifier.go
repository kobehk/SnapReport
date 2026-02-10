package geo

func ClassifyHighway(category string, road string) bool {
	if category == "motorway" || category == "trunk" || category == "motorway_link" || category == "trunk_link" {
		return true
	}
	if road != "" {
		if containsAny(road, []string{"高速", "高速公路", "Expressway"}) {
			return true
		}
	}
	return false
}

func containsAny(s string, subs []string) bool {
	for _, sub := range subs {
		if contains(s, sub) {
			return true
		}
	}
	return false
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (func() bool { return stringIndex(s, sub) >= 0 })()
}

func stringIndex(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
