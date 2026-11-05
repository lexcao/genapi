package main

import (
	"encoding/xml"
	"html"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestStaticSite(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	dir := filepath.Dir(file)
	readme, err := os.ReadFile(filepath.Join(dir, "..", "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(readme), "favorate") {
		t.Fatal("README still contains favorate")
	}
	if !strings.Contains(string(readme), "favorite") {
		t.Fatal("README missing favorite")
	}

	pages := []struct {
		path      string
		lang      string
		canonical string
		title     string
		home      bool
		blog      string
	}{
		{"index.html", "en", "https://lexcao.io/genapi/", "genapi — Declarative HTTP clients for Go — interface in, client out.", true, "https://lexcao.io/posts/genapi/"},
		{"zh/index.html", "zh-CN", "https://lexcao.io/genapi/zh/", "genapi — 用 interface 声明，生成 Go HTTP 客户端", true, "https://lexcao.io/zh/posts/genapi/"},
		{"install/index.html", "en", "https://lexcao.io/genapi/install/", "Install — genapi", false, "https://lexcao.io/posts/genapi/"},
		{"zh/install/index.html", "zh-CN", "https://lexcao.io/genapi/zh/install/", "安装 — genapi", false, "https://lexcao.io/zh/posts/genapi/"},
		{"examples/index.html", "en", "https://lexcao.io/genapi/examples/", "Examples — genapi", false, "https://lexcao.io/posts/genapi/"},
		{"zh/examples/index.html", "zh-CN", "https://lexcao.io/genapi/zh/examples/", "示例 — genapi", false, "https://lexcao.io/zh/posts/genapi/"},
		{"changelog/index.html", "en", "https://lexcao.io/genapi/changelog/", "Changelog — genapi", false, "https://lexcao.io/posts/genapi/"},
		{"zh/changelog/index.html", "zh-CN", "https://lexcao.io/genapi/zh/changelog/", "更新日志 — genapi", false, "https://lexcao.io/zh/posts/genapi/"},
	}

	var bodies []string
	for _, page := range pages {
		body := string(mustRead(t, filepath.Join(dir, page.path)))
		bodies = append(bodies, body)
		assertContains(t, page.path, body, `<html lang="`+page.lang+`"`)
		assertContains(t, page.path, body, "<title>"+page.title+"</title>")
		assertContains(t, page.path, body, `<meta name="description" content="`)
		assertContains(t, page.path, body, `<link rel="canonical" href="`+page.canonical+`">`)
		assertContains(t, page.path, body, `hreflang="en"`)
		assertContains(t, page.path, body, `hreflang="zh-CN"`)
		assertContains(t, page.path, body, `hreflang="x-default"`)
		assertContains(t, page.path, body, `<meta property="og:title"`)
		assertContains(t, page.path, body, `<meta property="og:description"`)
		assertContains(t, page.path, body, `<meta property="og:url" content="`+page.canonical+`">`)
		assertContains(t, page.path, body, `<meta property="og:type" content="website">`)
		assertContains(t, page.path, body, `<meta property="og:image" content="https://lexcao.io/genapi/assets/og-`)
		assertContains(t, page.path, body, `<meta property="og:image:width" content="1200">`)
		assertContains(t, page.path, body, `<meta property="og:image:height" content="630">`)
		assertContains(t, page.path, body, `<meta property="og:image:alt"`)
		assertContains(t, page.path, body, `<meta name="twitter:card" content="summary_large_image">`)
		assertContains(t, page.path, body, `"@type": "SoftwareSourceCode"`)
		assertContains(t, page.path, body, `"codeRepository": "https://github.com/lexcao/genapi"`)
		assertContains(t, page.path, body, `"programmingLanguage": "Go"`)
		assertContains(t, page.path, body, `id="intro"`)
		assertContains(t, page.path, body, page.blog)
		assertContains(t, page.path, body, "https://github.com/lexcao/genapi")
		assertContains(t, page.path, body, "https://pkg.go.dev/github.com/lexcao/genapi")
		assertContains(t, page.path, body, `class="wordmark wordmark-light"`)
		assertContains(t, page.path, body, `class="wordmark wordmark-dark"`)
		assertContains(t, page.path, body, `alt="genapi"`)
		for _, icon := range []string{
			"favicon-light.svg?v=1",
			"favicon-dark.svg?v=1",
			"favicon-16.png?v=1",
			"favicon-16-dark.png?v=1",
			"favicon-32.png?v=1",
			"favicon-32-dark.png?v=1",
			"apple-touch-icon.png?v=1",
			"favicon.ico?v=1",
		} {
			assertContains(t, page.path, body, icon)
		}
		if strings.Contains(body, "TODO") {
			t.Errorf("%s still has a TODO placeholder", page.path)
		}
		if strings.Contains(body, "og:image\" content=\"\"") || strings.Contains(body, "rel=\"icon\" href=\"\"") {
			t.Errorf("%s has an empty icon or image", page.path)
		}
		if page.home {
			if strings.Contains(body, "BreadcrumbList") {
				t.Errorf("%s should not include BreadcrumbList", page.path)
			}
		} else if !strings.Contains(body, "BreadcrumbList") {
			t.Errorf("%s missing BreadcrumbList", page.path)
		}
		if page.lang == "en" && strings.Contains(body, "https://lexcao.io/zh/posts/genapi/") {
			t.Errorf("%s links to the Chinese blog", page.path)
		}
		if page.lang == "zh-CN" && strings.Contains(body, "https://lexcao.io/posts/genapi/") {
			t.Errorf("%s links to the English blog", page.path)
		}
	}

	examples := html.UnescapeString(bodies[4])
	api := strings.TrimRight(string(mustRead(t, filepath.Join(dir, "..", "examples", "github", "api", "api.go"))), "\n")
	if !strings.Contains(examples, api) {
		t.Fatal("English examples page does not contain examples/github/api/api.go")
	}
	if !strings.Contains(bodies[6], "First release!") {
		t.Fatal("changelog missing v0.0.1 release text")
	}
	if !strings.Contains(bodies[6], "https://github.com/lexcao/genapi/compare/v0.0.2...v0.0.3") {
		t.Fatal("changelog missing v0.0.3 compare URL")
	}
	if !strings.Contains(bodies[2], "go get github.com/lexcao/genapi") {
		t.Fatal("install page missing go get")
	}
	if !strings.Contains(bodies[2], "There is no go install command") {
		t.Fatal("install page missing the go install note")
	}
	if strings.Contains(bodies[2], "README.md (line") {
		t.Fatal("install page still cites README line numbers")
	}

	sitemap := string(mustRead(t, filepath.Join(dir, "sitemap.xml")))
	dec := xml.NewDecoder(strings.NewReader(sitemap))
	for {
		_, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}
	}
	if strings.Count(sitemap, "<loc>") != 8 {
		t.Fatalf("sitemap loc count %d", strings.Count(sitemap, "<loc>"))
	}
	for _, loc := range []string{
		"https://lexcao.io/genapi/",
		"https://lexcao.io/genapi/zh/",
		"https://lexcao.io/genapi/install/",
		"https://lexcao.io/genapi/zh/install/",
		"https://lexcao.io/genapi/examples/",
		"https://lexcao.io/genapi/zh/examples/",
		"https://lexcao.io/genapi/changelog/",
		"https://lexcao.io/genapi/zh/changelog/",
	} {
		assertContains(t, "sitemap.xml", sitemap, "<loc>"+loc+"</loc>")
	}
	if strings.Count(sitemap, `hreflang="en"`) != 8 || strings.Count(sitemap, `hreflang="zh-CN"`) != 8 || strings.Count(sitemap, `hreflang="x-default"`) != 8 {
		t.Fatal("sitemap hreflang counts")
	}
	if _, err := os.Stat(filepath.Join(dir, "robots.txt")); !os.IsNotExist(err) {
		t.Fatal("website/robots.txt must not be published")
	}
	for i, page := range pages {
		body := bodies[i]
		const marker = `<link rel="canonical" href="`
		start := strings.Index(body, marker)
		if start < 0 {
			t.Fatalf("%s missing canonical", page.path)
		}
		rest := body[start+len(marker):]
		end := strings.Index(rest, `">`)
		if end < 0 {
			t.Fatalf("%s canonical unclosed", page.path)
		}
		canonical := rest[:end]
		if !strings.HasSuffix(canonical, "/") {
			t.Fatalf("%s canonical %q has no trailing slash", page.path, canonical)
		}
		if canonical != page.canonical {
			t.Fatalf("%s canonical %q, want %q", page.path, canonical, page.canonical)
		}
		if !strings.Contains(sitemap, "<loc>"+canonical+"</loc>") {
			t.Fatalf("sitemap missing loc for canonical %s", canonical)
		}
	}
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func assertContains(t *testing.T, name, body, needle string) {
	t.Helper()
	if !strings.Contains(body, needle) {
		t.Errorf("%s missing %q", name, needle)
	}
}
