package main

import "testing"

func TestRelHref(t *testing.T) {
	cases := []struct {
		from string
		to   string
		want string
	}{
		{"", "install/", "install/"},
		{"install/", "", "../"},
		{"zh/install/", "install/", "../../install/"},
		{"zh/", "zh/", "./"},
		{"", "", "./"},
		{"zh/", "", "../"},
		{"", "zh/", "zh/"},
		{"install/", "zh/install/", "../zh/install/"},
		{"zh/install/", "zh/", "../"},
	}
	for _, tc := range cases {
		if got := relHref(tc.from, tc.to); got != tc.want {
			t.Errorf("relHref(%q, %q) = %q, want %q", tc.from, tc.to, got, tc.want)
		}
	}
}

func TestRelAsset(t *testing.T) {
	if got := relAsset("zh/install/", "styles.css"); got != "../../styles.css" {
		t.Fatalf("got %q", got)
	}
	if got := relAsset("", "styles.css"); got != "styles.css" {
		t.Fatalf("got %q", got)
	}
}
