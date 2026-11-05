// Command buildsite generates the static genapi documentation site.
//
// It reads install commands from the README and go.mod files, example
// sources from the repository, and changelog entries from `gh release view`
// plus git tags. Output under website/ is generated; do not edit it by hand.
//
//	go run ./scripts/buildsite
package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	siteOrigin = "https://lexcao.io/genapi/"
	repoURL    = "https://github.com/lexcao/genapi"
	pkgURL     = "https://pkg.go.dev/github.com/lexcao/genapi"
	blogEN     = "https://lexcao.io/posts/genapi/"
	blogZH     = "https://lexcao.io/zh/posts/genapi/"
	licenseURL = "https://github.com/lexcao/genapi/blob/main/LICENSE"
)

func main() {
	root, err := repoRoot()
	if err != nil {
		fatal(err)
	}
	if err := generate(root); err != nil {
		fatal(err)
	}
}

func generate(root string) error {
	readme := mustRead(filepath.Join(root, "README.md"))
	restyReadme := mustRead(filepath.Join(root, "pkg/clients/resty/README.md"))
	assertNoGoInstall(
		filepath.Join(root, "README.md"),
		filepath.Join(root, "pkg/clients/resty/README.md"),
		filepath.Join(root, "go.mod"),
		filepath.Join(root, "pkg/clients/resty/go.mod"),
	)

	releases, tags, err := loadReleases(root)
	if err != nil {
		return err
	}
	if err := writeReleaseFile(root, tags, releases); err != nil {
		return err
	}

	snippets := loadSnippets(root, readme, restyReadme)
	pages := buildPages(snippets, releases)

	site := filepath.Join(root, "website")
	for _, page := range pages {
		if err := writePage(site, page); err != nil {
			return err
		}
	}
	if err := writeSitemap(site, pages); err != nil {
		return err
	}
	if err := writeRobots(site); err != nil {
		return err
	}
	fmt.Printf("generated %d HTML pages, sitemap.xml, and robots.txt\n", len(pages))
	return nil
}

func repoRoot() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func mustRead(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		fatal(fmt.Errorf("read %s: %w", path, err))
	}
	return string(b)
}

func assertNoGoInstall(paths ...string) {
	const needle = "go install"
	for _, p := range paths {
		body := mustRead(p)
		if strings.Contains(body, needle) {
			fatal(fmt.Errorf("%s contains %q; quote that command on the install page instead of reporting it absent", p, needle))
		}
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "buildsite: %v\n", err)
	os.Exit(1)
}

type ghView struct {
	TagName      string `json:"tagName"`
	Name         string `json:"name"`
	PublishedAt  string `json:"publishedAt"`
	Body         string `json:"body"`
	URL          string `json:"url"`
	IsDraft      bool   `json:"isDraft"`
	IsPrerelease bool   `json:"isPrerelease"`
}

type ghListItem struct {
	TagName      string `json:"tagName"`
	Name         string `json:"name"`
	PublishedAt  string `json:"publishedAt"`
	IsDraft      bool   `json:"isDraft"`
	IsPrerelease bool   `json:"isPrerelease"`
	IsLatest     bool   `json:"isLatest"`
}

type releaseRecord struct {
	TagName          string `json:"tagName"`
	Name             string `json:"name"`
	PublishedAt      string `json:"publishedAt"`
	PublishedBeijing string `json:"publishedAtBeijing"`
	Body             string `json:"body"`
	URL              string `json:"url"`
	IsDraft          bool   `json:"isDraft"`
	IsPrerelease     bool   `json:"isPrerelease"`
	IsLatest         bool   `json:"isLatest"`
}

type releaseFile struct {
	GeneratedBy string          `json:"generatedBy"`
	Source      string          `json:"source"`
	GitTags     []string        `json:"gitTags"`
	Releases    []releaseRecord `json:"releases"`
}

func loadReleases(root string) ([]releaseRecord, []string, error) {
	tagOut, err := git(root, "tag", "-l")
	if err != nil {
		return nil, nil, err
	}
	var tags []string
	for _, line := range strings.Split(strings.TrimSpace(tagOut), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			tags = append(tags, line)
		}
	}
	sort.Strings(tags)

	listOut, err := gh(root, "release", "list", "--limit", "100", "--json", "tagName,name,publishedAt,isDraft,isPrerelease,isLatest")
	if err != nil {
		return nil, nil, err
	}
	var listed []ghListItem
	if err := json.Unmarshal([]byte(listOut), &listed); err != nil {
		return nil, nil, fmt.Errorf("parse gh release list: %w", err)
	}
	latest := map[string]bool{}
	listedTags := map[string]bool{}
	for _, item := range listed {
		listedTags[item.TagName] = true
		latest[item.TagName] = item.IsLatest
	}

	tagSet := map[string]bool{}
	for _, tag := range tags {
		tagSet[tag] = true
	}
	if len(tagSet) != len(listedTags) || !sameSet(tagSet, listedTags) {
		return nil, nil, fmt.Errorf("git tags %v and GitHub releases %v differ", tags, keys(listedTags))
	}

	var releases []releaseRecord
	for _, tag := range tags {
		viewOut, err := gh(root, "release", "view", tag, "--json", "tagName,name,publishedAt,body,url,isDraft,isPrerelease")
		if err != nil {
			return nil, nil, err
		}
		var view ghView
		if err := json.Unmarshal([]byte(viewOut), &view); err != nil {
			return nil, nil, fmt.Errorf("parse gh release view %s: %w", tag, err)
		}
		if view.TagName != tag {
			return nil, nil, fmt.Errorf("gh release view %s returned tag %s", tag, view.TagName)
		}
		if strings.TrimSpace(view.Body) == "" || view.URL == "" || view.PublishedAt == "" {
			return nil, nil, fmt.Errorf("release %s is missing body, url, or publishedAt", tag)
		}
		published, err := time.Parse(time.RFC3339, view.PublishedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("release %s publishedAt: %w", tag, err)
		}
		beijing := published.In(time.FixedZone("UTC+8", 8*60*60)).Format("2006-01-02 15:04:05")
		releases = append(releases, releaseRecord{
			TagName:          view.TagName,
			Name:             view.Name,
			PublishedAt:      view.PublishedAt,
			PublishedBeijing: beijing,
			Body:             view.Body,
			URL:              view.URL,
			IsDraft:          view.IsDraft,
			IsPrerelease:     view.IsPrerelease,
			IsLatest:         latest[tag],
		})
	}
	sort.Slice(releases, func(i, j int) bool {
		return releases[i].PublishedAt > releases[j].PublishedAt
	})
	return releases, tags, nil
}

func sameSet(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func git(root string, args ...string) (string, error) {
	return command(root, "git", args...)
}

func gh(root string, args ...string) (string, error) {
	return command(root, "gh", args...)
}

func command(root, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = root
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s %s: %w: %s", name, strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return string(out), nil
}

func writeReleaseFile(root string, tags []string, releases []releaseRecord) error {
	payload := releaseFile{
		GeneratedBy: "scripts/buildsite",
		Source:      "git tag -l and gh release view",
		GitTags:     tags,
		Releases:    releases,
	}
	raw, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	dir := filepath.Join(root, "scripts", "buildsite", "data")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "releases.json"), raw, 0o644)
}

type Snippet struct {
	Caption string
	Body    string
	Lang    string
	BlobURL string
}

// Snippets are verbatim excerpts. Field names are exported for the templates.
type Snippets struct {
	APIGo        Snippet
	APIGen       Snippet
	MainGo       Snippet
	RestyMain    Snippet
	BaseAnn      Snippet
	MethodsGH    Snippet
	MethodsTodo  Snippet
	Alias        Snippet
	Query        Snippet
	BodyIface    Snippet
	BodyGen      Snippet
	ResultBind   Snippet
	GoMod        Snippet
	RestyMod     Snippet
	GoGet        Snippet
	GoGenerate   Snippet
	GoDirective  Snippet
	RestyGet     Snippet
	RestyGetBare Snippet
	MakeGen      Snippet
	Clone        Snippet
	GoVersionDoc Snippet
}

func loadSnippets(root, readme, restyReadme string) Snippets {
	return Snippets{
		APIGo:        snippetFull(root, "examples/github/api/api.go", "go"),
		APIGen:       snippetFull(root, "examples/github/api/api.gen.go", "go"),
		MainGo:       snippetFull(root, "examples/github/main.go", "go"),
		RestyMain:    snippetFull(root, "pkg/clients/resty/example/main.go", "go"),
		BaseAnn:      snippetLines(root, "examples/github/api/api.go", 11, 21, "go", `@BaseURL("https://api.github.com")`),
		MethodsGH:    snippetLines(root, "examples/github/api/api.go", 16, 20, "go", `@GET("/repos/{owner}/{repo}/contributors")`),
		MethodsTodo:  snippetLines(root, "website/api/todo.go", 24, 28, "go", `@DELETE("/todos/{id}")`),
		Alias:        snippetLines(root, "internal/build/parser/annotation/declarations.go", 137, 149, "go", "http.MethodTrace"),
		Query:        snippetLines(root, "test/e2e/httpbin.go", 12, 18, "go", `@Query("value", "{value}")`),
		BodyIface:    snippetLines(root, "examples/github/api/api.go", 19, 34, "go", "CreateIssue"),
		BodyGen:      snippetLines(root, "examples/github/api/api.gen.go", 32, 44, "go", "Body:    issue"),
		ResultBind:   snippetLines(root, "internal/build/binder/result.go", 24, 47, "go", "MustHandleResponse"),
		GoMod:        snippetLines(root, "go.mod", 1, 3, "text", "go 1.18"),
		RestyMod:     snippetLines(root, "pkg/clients/resty/go.mod", 1, 5, "text", "toolchain go1.24.5"),
		GoGet:        lineSnippet(readme, "README.md", "go get github.com/lexcao/genapi", "bash"),
		GoGenerate:   lineSnippet(readme, "README.md", "$ go generate ./api", "bash"),
		GoDirective:  lineSnippet(readme, "README.md", "//go:generate go run github.com/lexcao/genapi/cmd/genapi  -file $GOFILE", "go"),
		RestyGet:     lineSnippet(readme, "README.md", "$ go get github.com/lexcao/genapi/pkg/clients/resty", "bash"),
		RestyGetBare: lineSnippet(restyReadme, "pkg/clients/resty/README.md", "go get github.com/lexcao/genapi/pkg/clients/resty", "bash"),
		MakeGen:      lineSnippetBlock(mustRead(filepath.Join(root, "Makefile")), "Makefile", "generate:", "\tgo generate ./...", "text"),
		Clone:        fenceContaining(readme, "README.md", "git clone https://github.com/lexcao/genapi.git", "bash"),
		GoVersionDoc: lineSnippet(readme, "README.md", "- Go 1.18 or higher", "text"),
	}
}

func snippetFull(root, rel, lang string) Snippet {
	text := mustRead(filepath.Join(root, rel))
	return Snippet{
		Caption: rel,
		Body:    strings.TrimRight(text, "\n"),
		Lang:    lang,
		BlobURL: blob(rel, 0, 0),
	}
}

func snippetLines(root, rel string, start, end int, lang, mustContain string) Snippet {
	text := mustRead(filepath.Join(root, rel))
	lines := strings.Split(text, "\n")
	// Drop the trailing empty line produced by a final newline so indexes match file lines.
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if start < 1 || end > len(lines) || start > end {
		fatal(fmt.Errorf("%s lines %d-%d out of range (file has %d lines)", rel, start, end, len(lines)))
	}
	body := strings.Join(lines[start-1:end], "\n")
	if !strings.Contains(body, mustContain) {
		fatal(fmt.Errorf("%s lines %d-%d do not contain %q", rel, start, end, mustContain))
	}
	return Snippet{
		Caption: fmt.Sprintf("%s (lines %d–%d)", rel, start, end),
		Body:    body,
		Lang:    lang,
		BlobURL: blob(rel, start, end),
	}
}

func lineSnippet(fileText, rel, exact, lang string) Snippet {
	n := findExactLine(fileText, exact)
	if n == 0 {
		fatal(fmt.Errorf("%s does not contain exact line %q", rel, exact))
	}
	return Snippet{
		Caption: fmt.Sprintf("%s (line %d)", rel, n),
		Body:    exact,
		Lang:    lang,
		BlobURL: blob(rel, n, n),
	}
}

func lineSnippetBlock(fileText, rel, first, second, lang string) Snippet {
	n1 := findExactLine(fileText, first)
	n2 := findExactLine(fileText, second)
	if n1 == 0 || n2 == 0 || n2 != n1+1 {
		fatal(fmt.Errorf("%s does not contain adjacent lines %q and %q", rel, first, second))
	}
	return Snippet{
		Caption: fmt.Sprintf("%s (lines %d–%d)", rel, n1, n2),
		Body:    first + "\n" + second,
		Lang:    lang,
		BlobURL: blob(rel, n1, n2),
	}
}

func fenceContaining(fileText, rel, needle, lang string) Snippet {
	parts := strings.Split(fileText, "```")
	for i := 1; i < len(parts); i += 2 {
		block := parts[i]
		nl := strings.IndexByte(block, '\n')
		if nl == -1 {
			continue
		}
		body := strings.Trim(block[nl+1:], "\n")
		if !strings.Contains(body, needle) {
			continue
		}
		// Locate the fence body in the original file for a blob link.
		idx := strings.Index(fileText, body)
		if idx < 0 {
			fatal(fmt.Errorf("fence body not found in %s", rel))
		}
		start := strings.Count(fileText[:idx], "\n") + 1
		end := start + strings.Count(body, "\n")
		return Snippet{
			Caption: fmt.Sprintf("%s (lines %d–%d)", rel, start, end),
			Body:    body,
			Lang:    lang,
			BlobURL: blob(rel, start, end),
		}
	}
	fatal(fmt.Errorf("%s has no fenced block containing %q", rel, needle))
	return Snippet{}
}

func findExactLine(text, exact string) int {
	for i, line := range strings.Split(text, "\n") {
		if line == exact {
			return i + 1
		}
	}
	return 0
}

func blob(rel string, start, end int) string {
	u := repoURL + "/blob/main/" + path.Clean(rel)
	if start <= 0 {
		return u
	}
	if end <= start {
		return fmt.Sprintf("%s#L%d", u, start)
	}
	return fmt.Sprintf("%s#L%d-L%d", u, start, end)
}

func relPrefix(from string) string {
	from = strings.Trim(from, "/")
	if from == "" {
		return ""
	}
	return strings.Repeat("../", strings.Count(from, "/")+1)
}

func relHref(from, to string) string {
	if from == to {
		return "./"
	}
	var fromParts, toParts []string
	if trimmed := strings.Trim(from, "/"); trimmed != "" {
		fromParts = strings.Split(trimmed, "/")
	}
	if trimmed := strings.Trim(to, "/"); trimmed != "" {
		toParts = strings.Split(trimmed, "/")
	}
	i := 0
	for i < len(fromParts) && i < len(toParts) && fromParts[i] == toParts[i] {
		i++
	}
	out := make([]string, 0, len(fromParts)-i+len(toParts)-i)
	for n := len(fromParts) - i; n > 0; n-- {
		out = append(out, "..")
	}
	out = append(out, toParts[i:]...)
	if len(out) == 0 {
		return "./"
	}
	rel := strings.Join(out, "/")
	if !strings.HasSuffix(rel, "/") {
		rel += "/"
	}
	return rel
}

func relAsset(from, name string) string {
	return relPrefix(from) + name
}

func absolute(sitePath string) string {
	return siteOrigin + sitePath
}

type crumb struct {
	Name string
	Href string
	URL  string
}

type page struct {
	ID                string
	LangID            string
	LangCode          string
	OGLocale          string
	OGLocaleAlt       string
	SitePath          string
	Title             string
	Description       string
	H1                string
	Intro             string
	Canonical         string
	HrefCSS           string
	HrefJS            string
	HrefHome          string
	HrefInstall       string
	HrefExamples      string
	HrefChangelog     string
	HrefReference     string
	HrefWhy           string
	HrefStart         string
	HrefEN            string
	HrefZH            string
	BlogURL           string
	RepoURL           string
	PkgURL            string
	LicenseURL        string
	IssuesURL         string
	DiscussURL        string
	ContribURL        string
	ReleasesURL       string
	GoURL             string
	JSONLD            template.JS
	OGImage           string
	OGAlt             string
	HrefWordmarkLight string
	HrefWordmarkDark  string
	HrefFaviconLight  string
	HrefFaviconDark   string
	HrefFavicon16     string
	HrefFavicon16Dark string
	HrefFavicon32     string
	HrefFavicon32Dark string
	HrefApple         string
	HrefIco           string
	Crumbs            []crumb
	Snippets          Snippets
	Releases          []releaseRecord
}

func buildPages(sn Snippets, releases []releaseRecord) []page {
	type spec struct {
		id string
		en string
		zh string
	}
	specs := []spec{
		{"home", "", "zh/"},
		{"install", "install/", "zh/install/"},
		{"examples", "examples/", "zh/examples/"},
		{"changelog", "changelog/", "zh/changelog/"},
	}
	var pages []page
	for _, spec := range specs {
		pages = append(pages, newPage(spec.id, "en", spec.en, spec.zh, sn, releases))
		pages = append(pages, newPage(spec.id, "zh", spec.zh, spec.en, sn, releases))
	}
	return pages
}

func newPage(id, lang, sitePath, altPath string, sn Snippets, releases []releaseRecord) page {
	enPath, zhPath := sitePath, altPath
	if lang == "zh" {
		enPath, zhPath = altPath, sitePath
	}
	home := ""
	if lang == "zh" {
		home = "zh/"
	}
	p := page{
		ID:                id,
		LangID:            lang,
		SitePath:          sitePath,
		Canonical:         absolute(sitePath),
		HrefCSS:           relAsset(sitePath, "styles.css"),
		HrefJS:            relAsset(sitePath, "script.js"),
		HrefHome:          relHref(sitePath, home),
		HrefInstall:       relHref(sitePath, joinLang(lang, "install/")),
		HrefExamples:      relHref(sitePath, joinLang(lang, "examples/")),
		HrefChangelog:     relHref(sitePath, joinLang(lang, "changelog/")),
		HrefReference:     relHref(sitePath, home) + "#documentation",
		HrefWhy:           relHref(sitePath, home) + "#why-genapi",
		HrefStart:         relHref(sitePath, home) + "#getting-started",
		HrefEN:            relHref(sitePath, enPath),
		HrefZH:            relHref(sitePath, zhPath),
		RepoURL:           repoURL,
		PkgURL:            pkgURL,
		LicenseURL:        licenseURL,
		IssuesURL:         repoURL + "/issues",
		DiscussURL:        repoURL + "/discussions",
		ContribURL:        repoURL + "/blob/main/CONTRIBUTING.md",
		ReleasesURL:       repoURL + "/releases",
		GoURL:             "https://go.dev",
		HrefWordmarkLight: relAsset(sitePath, "brand/wordmark-light.svg"),
		HrefWordmarkDark:  relAsset(sitePath, "brand/wordmark-dark.svg"),
		HrefFaviconLight:  relAsset(sitePath, "brand/favicon-light.svg"),
		HrefFaviconDark:   relAsset(sitePath, "brand/favicon-dark.svg"),
		HrefFavicon16:     relAsset(sitePath, "assets/favicon-16.png"),
		HrefFavicon16Dark: relAsset(sitePath, "assets/favicon-16-dark.png"),
		HrefFavicon32:     relAsset(sitePath, "assets/favicon-32.png"),
		HrefFavicon32Dark: relAsset(sitePath, "assets/favicon-32-dark.png"),
		HrefApple:         relAsset(sitePath, "assets/apple-touch-icon.png"),
		HrefIco:           relAsset(sitePath, "assets/favicon.ico"),
		Snippets:          sn,
		Releases:          releases,
	}
	if lang == "zh" {
		p.LangCode = "zh-CN"
		p.OGLocale = "zh_CN"
		p.OGLocaleAlt = "en_US"
		p.BlogURL = blogZH
		p.OGImage = siteOrigin + "assets/og-zh.png"
		fillZH(&p, releases)
	} else {
		p.LangCode = "en"
		p.OGLocale = "en_US"
		p.OGLocaleAlt = "zh_CN"
		p.BlogURL = blogEN
		p.OGImage = siteOrigin + "assets/og-en.png"
		fillEN(&p, releases)
	}
	if id != "home" {
		p.Crumbs = []crumb{
			{Name: "genapi", Href: p.HrefHome, URL: absolute(home)},
			{Name: p.H1, Href: "./", URL: p.Canonical},
		}
	}
	p.JSONLD = template.JS(mustJSONLD(p))
	return p
}

func joinLang(lang, leaf string) string {
	if lang == "zh" {
		return "zh/" + leaf
	}
	return leaf
}

func fillEN(p *page, releases []releaseRecord) {
	const tagline = "Declarative HTTP clients for Go — interface in, client out."
	p.OGAlt = "genapi wordmark. " + tagline
	switch p.ID {
	case "home":
		p.Title = "genapi — " + tagline
		p.Description = "Declarative HTTP clients for Go — interface in, client out. Declare a Go interface, run go generate, and call the client."
		p.H1 = tagline
		p.Intro = "Write a Go interface, annotate each call, and run go generate. The generated file implements that interface, and your code calls it with ordinary Go types."
	case "install":
		p.Title = "Install — genapi"
		p.H1 = "Install"
		p.Description = "Install github.com/lexcao/genapi with the go get command from the README. The module requires Go 1.18. The repository does not document a go install command."
		p.Intro = "Commands on this page are copied from README.md, pkg/clients/resty/README.md, the Makefile, and the go.mod files."
	case "examples":
		p.Title = "Examples — genapi"
		p.H1 = "Examples"
		p.Description = "GitHub client example copied from examples/github in the genapi repository: the interface, the generated client, and main."
		p.Intro = "Each listing is a file from examples/github, copied verbatim from this repository."
	case "changelog":
		p.Title = "Changelog — genapi"
		p.H1 = "Changelog"
		p.Description = changelogDescription("en", releases)
		p.Intro = "Each entry is a GitHub Release. The note is the release body, copied verbatim. Dates show the publishedAt timestamp from GitHub and the same instant in Beijing time (UTC+8)."
	}
}

func fillZH(p *page, releases []releaseRecord) {
	const tagline = "用 interface 声明，生成 Go HTTP 客户端"
	p.OGAlt = "genapi 字标。" + tagline
	switch p.ID {
	case "home":
		p.Title = "genapi — " + tagline
		p.Description = "用 interface 声明，生成 Go HTTP 客户端。声明 Go interface，运行 go generate，然后直接调用。"
		p.H1 = tagline
		p.Intro = "写一个 Go interface，给每次调用加上注解，然后运行 go generate。生成的文件实现这个 interface，业务代码用普通的 Go 类型去调用。"
	case "install":
		p.Title = "安装 — genapi"
		p.H1 = "安装"
		p.Description = "用 README 里的 go get 命令安装 github.com/lexcao/genapi。模块要求 Go 1.18。仓库没有写 go install 命令。"
		p.Intro = "本页命令摘自 README.md、pkg/clients/resty/README.md、Makefile 和 go.mod，原文照录。"
	case "examples":
		p.Title = "示例 — genapi"
		p.H1 = "示例"
		p.Description = "摘自仓库 examples/github 的 GitHub 客户端：interface、生成的客户端，以及 main。"
		p.Intro = "下面每一段都是 examples/github 里的文件，按仓库内容原文贴出。"
	case "changelog":
		p.Title = "更新日志 — genapi"
		p.H1 = "更新日志"
		p.Description = changelogDescription("zh", releases)
		p.Intro = "每一节对应一个 GitHub Release。说明文字是 Release 正文的原文。时间同时给出 GitHub 的 publishedAt，以及同一时刻的北京时间（UTC+8）。"
	}
}

func changelogDescription(lang string, releases []releaseRecord) string {
	tags := make([]string, len(releases))
	for i, rel := range releases {
		tags[i] = rel.TagName
	}
	// releases are newest first; mention them in that order.
	list := strings.Join(tags, ", ")
	if lang == "zh" {
		return "genapi 更新日志，内容来自 GitHub Releases：" + list + "。"
	}
	return "genapi changelog taken from GitHub Releases: " + list + "."
}

func mustJSONLD(p page) string {
	software := map[string]interface{}{
		"@type":               "SoftwareSourceCode",
		"name":                "genapi",
		"description":         p.Description,
		"url":                 p.Canonical,
		"codeRepository":      repoURL,
		"programmingLanguage": "Go",
		"license":             licenseURL,
		"inLanguage":          []string{"en", "zh-CN"},
	}
	var payload interface{}
	if p.ID == "home" {
		software["@context"] = "https://schema.org"
		payload = software
	} else {
		crumbs := make([]map[string]interface{}, len(p.Crumbs))
		for i, c := range p.Crumbs {
			crumbs[i] = map[string]interface{}{
				"@type":    "ListItem",
				"position": i + 1,
				"name":     c.Name,
				"item":     c.URL,
			}
		}
		payload = map[string]interface{}{
			"@context": "https://schema.org",
			"@graph": []interface{}{
				software,
				map[string]interface{}{
					"@type":           "BreadcrumbList",
					"itemListElement": crumbs,
				},
			},
		}
	}
	raw, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		fatal(err)
	}
	return string(raw)
}

func writePage(site string, p page) error {
	var buf bytes.Buffer
	if err := render(&buf, p); err != nil {
		return fmt.Errorf("render %s %s: %w", p.LangID, p.ID, err)
	}
	out := filepath.Join(site, filepath.FromSlash(p.SitePath), "index.html")
	if p.SitePath == "" {
		out = filepath.Join(site, "index.html")
	}
	if err := ensureInside(site, out); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}
	return os.WriteFile(out, buf.Bytes(), 0o644)
}

func ensureInside(root, target string) error {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(absRoot, absTarget)
	if err != nil {
		return err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("refusing to write outside website: %s", target)
	}
	return nil
}

func writeSitemap(site string, pages []page) error {
	type pair struct{ en, zh string }
	seen := []pair{}
	index := map[string]int{}
	for _, p := range pages {
		if p.LangID != "en" {
			continue
		}
		zh := ""
		for _, other := range pages {
			if other.ID == p.ID && other.LangID == "zh" {
				zh = other.Canonical
			}
		}
		seen = append(seen, pair{en: p.Canonical, zh: zh})
		index[p.ID] = len(seen) - 1
	}
	// Stable order: home, install, examples, changelog.
	order := []string{"home", "install", "examples", "changelog"}
	var b strings.Builder
	b.WriteString(xml.Header)
	b.WriteString("<!-- Generated by scripts/buildsite. Do not edit by hand. -->\n")
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">` + "\n")
	for _, id := range order {
		item := seen[index[id]]
		for _, loc := range []string{item.en, item.zh} {
			fmt.Fprintf(&b, "  <url>\n    <loc>%s</loc>\n", xmlEscape(loc))
			for _, alt := range []struct{ lang, href string }{
				{"en", item.en},
				{"zh-CN", item.zh},
				{"x-default", item.en},
			} {
				fmt.Fprintf(&b, "    <xhtml:link rel=\"alternate\" hreflang=\"%s\" href=\"%s\"/>\n", alt.lang, xmlEscape(alt.href))
			}
			b.WriteString("  </url>\n")
		}
	}
	b.WriteString("</urlset>\n")
	return os.WriteFile(filepath.Join(site, "sitemap.xml"), []byte(b.String()), 0o644)
}

func xmlEscape(s string) string {
	var buf bytes.Buffer
	_ = xml.EscapeText(&buf, []byte(s))
	return buf.String()
}

func writeRobots(site string) error {
	body := "# Generated by scripts/buildsite. Do not edit by hand.\n" +
		"User-agent: *\n" +
		"Allow: /\n\n" +
		"Sitemap: " + siteOrigin + "sitemap.xml\n"
	return os.WriteFile(filepath.Join(site, "robots.txt"), []byte(body), 0o644)
}
