# genapi docs site — UI audit

Date: 2026-11-05. Surface: English and Chinese home, install, examples, and changelog. Benchmark: Awwwards, Webby, and FWA product pages, where type, space, and one motion idea carry the page. The wordmark, favicons, URLs, and facts stay.

Chosen direction after the audit: **ink and signal**. Warm paper and near-black ink, with the Go-blue dot as the only chromatic signal. Code is set as a document, not a fake operating-system window.

## Typography

Gap. Inter at a single display size, with feature pills and identical card titles, flattens the page. Body measure is unconstrained on subpages, so lines run the full content width. Install captions are file paths and line numbers (`README.md (line 241)`), which is source control, not a type hierarchy.

Response. Keep Inter for UI and JetBrains Mono for code and the wordmark. Tighten the display size steps, cap body measure near 68ch, and replace path captions on the install page with instruction titles. Headlines stay `text-wrap: balance`.

## Whitespace

Gap. Sections stack at a constant padding. The four “why” cards and the getting-started steps use the same gap, so groups and separations feel identical. The hero terminal sits in a generic card with little air around the headline.

Response. More space above section headings than below them. The reasons become a definition list with a wide row gap. Subpages keep the home container inset.

## Visual hierarchy

Gap. Hero, four equal cards, tabs, and steps all compete. Primary and secondary buttons, pills, and the blog link have similar weight. Subpage `h1` is only slightly larger than the following `h2`.

Response. The headline and the interface listing lead. Reasons are a list, not a card grid. Install is a short sequence of commands. Changelog releases stay as separate articles, with the latest mark as the only badge.

## Color

Gap. Indigo (`#3949ab` / `#7986cb`) is a default SaaS palette and fights the Go-blue dot in the wordmark. Secondary gray (`#6c757d`) is the category gray, not a tint of the paper. Dark mode is a neutral `#121212` with the same indigo.

Response. Paper `#f4f1ea` and ink `#161411`. Links and focus use a darkened Go blue that clears 4.5:1. The literal `#00ADD8` stays on the wordmark and as a hairline, not as small text. Dark theme inverts the same paper and ink. Buttons are ink, not indigo.

## Motion

Gap. Almost none, aside from color transitions and a mobile drawer. Nothing marks the code as the arrival. There is no `prefers-reduced-motion` guard on scroll smoothing beyond what the browser does.

Response. One moment: the hairline under the hero code sheet draws in with an ease-out. It is off when `prefers-reduced-motion: reduce`. Scroll smoothing and hover shifts respect the same query. No scroll-triggered fade-in of content.

## Micro-interaction

Gap. Tabs switch with no feedback beyond a border. Code cannot be copied. Focus rings exist but theme toggle, selection, and scrollbars are browser defaults. The fake red/yellow/green window dots are decoration, not a control.

Response. Copy control on each listing. Tab and nav underlines move. Selection, caret, and scrollbars use the paper/ink palette. Window dots are removed.

## Responsive

Gap. The layout holds (inset, no page-level sideways scroll, balanced English headline). The hero code still looks like a desktop window squeezed onto a phone. Feature pills wrap as loose chips. Install code cards are full-bleed within the column but their captions are long filenames.

Response. Keep the container inset and internal code scroll. Stack the hero. Let the code sheet be full width of the column with a short file label. Touch targets on tabs and the theme control stay at least 40px.

## Originality

Gap. The page is a familiar generator landing: gradient-free but still indigo, four icon-less clones of the same card, macOS terminal chrome, and pill tags. Nothing belongs specifically to “interface in, client out” except the wordmark.

Response. Drop the terminal costume and the equal cards. The interface listing is the picture. The only color signal is the Go-blue hairline, matching the dot in the wordmark.
