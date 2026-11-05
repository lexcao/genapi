#!/usr/bin/env python3
"""Render favicon and Open Graph PNGs from the brand SVGs.

Source SVGs in website/brand/ are used as given. This script does not edit
them. It shells out to rsvg-convert and composites with Pillow.

Requires: rsvg-convert (librsvg2-bin), Pillow, JetBrains Mono, and a CJK font
(WenQuanYi Micro Hei or Noto Sans CJK).

    python3 website/scripts/build-assets.py
"""

import io
import subprocess
import sys
from pathlib import Path

from PIL import Image, ImageDraw, ImageFont

ROOT = Path(__file__).resolve().parents[1]
BRAND = ROOT / "brand"
OUT = ROOT / "assets"

WORDMARK = BRAND / "wordmark-light.svg"
FAVICON = BRAND / "favicon-light.svg"
FAVICON_DARK = BRAND / "favicon-dark.svg"
FAVICON_16 = BRAND / "favicon-16-light.svg"
FAVICON_16_DARK = BRAND / "favicon-16-dark.svg"

BG = (244, 247, 251, 255)
INK = (17, 17, 17, 255)
ACCENT = (0, 173, 216, 255)
APPLE_BG = (255, 255, 255, 255)

TAGLINE_EN = "Declarative HTTP clients for Go —\ninterface in, client out."
TAGLINE_ZH = "用 interface 声明，生成 Go HTTP 客户端"

LATIN_FONTS = [
    "/usr/share/fonts/truetype/jetbrains-mono/JetBrainsMono-Medium.ttf",
    "/usr/share/fonts/truetype/jetbrains-mono/JetBrainsMono-Regular.ttf",
    "/usr/share/fonts/truetype/macos/JetBrainsMono-Regular.ttf",
]
CJK_FONTS = [
    "/usr/share/fonts/truetype/wqy/wqy-microhei.ttc",
    "/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc",
    "/usr/share/fonts/truetype/noto/NotoSansCJK-Regular.ttc",
]


def die(message):
    print("build-assets: " + message, file=sys.stderr)
    sys.exit(1)


def first_existing(paths):
    for path in paths:
        if Path(path).is_file():
            return path
    die("none of these fonts exist: " + ", ".join(paths))


def render_svg(src, width, height=None):
    if not src.is_file():
        die("missing SVG " + str(src))
    cmd = ["rsvg-convert", "-w", str(width)]
    if height is not None:
        cmd += ["-h", str(height)]
    cmd.append(str(src))
    try:
        png = subprocess.check_output(cmd)
    except FileNotFoundError:
        die("rsvg-convert is not installed (package librsvg2-bin)")
    except subprocess.CalledProcessError as exc:
        die("rsvg-convert failed for " + src.name + ": " + str(exc))
    image = Image.open(io.BytesIO(png)).convert("RGBA")
    if height is not None and image.size != (width, height):
        die("%s rendered %s, want %dx%d" % (src.name, image.size, width, height))
    return image


def font(path, size):
    return ImageFont.truetype(path, size)


def is_latin(ch):
    return ord(ch) < 128


def text_width(text, latin, cjk):
    width = 0
    for ch in text:
        face = latin if is_latin(ch) else cjk
        width += face.getlength(ch)
    return width


def wrap_text(text, latin, cjk, max_width):
    if "\n" in text:
        lines = []
        for part in text.split("\n"):
            lines.extend(wrap_paragraph(part, latin, cjk, max_width))
        return lines
    return wrap_paragraph(text, latin, cjk, max_width)


def wrap_paragraph(text, latin, cjk, max_width):
    lines = []
    current = ""
    for ch in text:
        trial = current + ch
        if text_width(trial, latin, cjk) <= max_width or current.strip() == "":
            current = trial
            continue
        break_at = current.rfind(" ")
        if break_at > 0:
            lines.append(current[:break_at].rstrip())
            current = (current[break_at + 1 :] + ch).lstrip()
        else:
            lines.append(current.rstrip())
            current = ch if ch != " " else ""
    if current.strip():
        lines.append(current.rstrip())
    return lines


def draw_line(draw, xy, text, latin, cjk, fill):
    x, y = xy
    for ch in text:
        face = latin if is_latin(ch) else cjk
        draw.text((x, y), ch, font=face, fill=fill, anchor="ls")
        x += face.getlength(ch)
    return x


def intersects(a, b, gap=0):
    ax, ay, aw, ah = a
    bx, by, bw, bh = b
    return not (
        ax + aw + gap <= bx
        or bx + bw + gap <= ax
        or ay + ah + gap <= by
        or by + bh + gap <= ay
    )


def build_og(path, tagline):
    latin_path = first_existing(LATIN_FONTS)
    cjk_path = first_existing(CJK_FONTS)
    canvas = Image.new("RGBA", (1200, 630), BG)
    draw = ImageDraw.Draw(canvas)
    draw.rectangle((0, 0, 10, 630), fill=ACCENT)

    wordmark = render_svg(WORDMARK, 520)
    margin_x = 88
    text_width_limit = 1200 - margin_x - 72

    chosen = None
    for size in (48, 44, 40, 36, 34, 32):
        latin = font(latin_path, size)
        cjk = font(cjk_path, size)
        lines = wrap_text(tagline, latin, cjk, text_width_limit)
        if lines and len(lines) <= 3 and all(text_width(line, latin, cjk) <= text_width_limit for line in lines):
            chosen = (latin, cjk, lines, size)
            # Prefer the largest size that keeps the tagline on one line when it fits.
            if len(lines) == 1 or "\n" in tagline:
                break
    if chosen is None:
        die("tagline does not fit on the OG image: " + tagline)
    latin, cjk, lines, size = chosen
    line_gap = int(size * 1.45)
    block_h = wordmark.height + 40 + line_gap * len(lines)
    top = max(48, (630 - block_h) // 2)
    wm_box = (margin_x, top, wordmark.width, wordmark.height)
    canvas.alpha_composite(wordmark, (wm_box[0], wm_box[1]))

    text_top = top + wordmark.height + 40
    for index, line in enumerate(lines):
        baseline = text_top + int(size * 0.82) + index * line_gap
        width = text_width(line, latin, cjk)
        text_box = (margin_x, baseline - size, int(width) + 2, size + 8)
        if intersects(wm_box, text_box, gap=24):
            die("OG text overlaps the wordmark")
        if text_box[0] < 72 or text_box[0] + text_box[2] > 1160:
            die("OG text collides with the canvas edge: " + line)
        if baseline > 590:
            die("OG text runs off the bottom")
        draw_line(draw, (margin_x, baseline), line, latin, cjk, INK)

    OUT.mkdir(parents=True, exist_ok=True)
    canvas.convert("RGB").save(path, "PNG")


def build_apple(path):
    icon = render_svg(FAVICON, 112, 112)
    canvas = Image.new("RGBA", (180, 180), APPLE_BG)
    offset = (180 - 112) // 2
    canvas.alpha_composite(icon, (offset, offset))
    # The favicon SVG is transparent around the glyph. Flatten onto the solid fill.
    flat = Image.new("RGB", (180, 180), APPLE_BG[:3])
    flat.paste(canvas, mask=canvas.getchannel("A"))
    flat.save(path, "PNG")


def build_ico(path, png16, png32):
    image16 = Image.open(png16).convert("RGBA")
    image32 = Image.open(png32).convert("RGBA")
    image32.save(path, format="ICO", sizes=[(16, 16), (32, 32)], append_images=[image16])


def main():
    OUT.mkdir(parents=True, exist_ok=True)
    pairs = [
        (FAVICON_16, OUT / "favicon-16.png", 16),
        (FAVICON_16_DARK, OUT / "favicon-16-dark.png", 16),
        (FAVICON, OUT / "favicon-32.png", 32),
        (FAVICON_DARK, OUT / "favicon-32-dark.png", 32),
    ]
    for src, dest, size in pairs:
        render_svg(src, size, size).save(dest, "PNG")
    build_apple(OUT / "apple-touch-icon.png")
    build_ico(OUT / "favicon.ico", OUT / "favicon-16.png", OUT / "favicon-32.png")
    build_og(OUT / "og-en.png", TAGLINE_EN)
    build_og(OUT / "og-zh.png", TAGLINE_ZH)

    expected = {
        "favicon-16.png": (16, 16),
        "favicon-16-dark.png": (16, 16),
        "favicon-32.png": (32, 32),
        "favicon-32-dark.png": (32, 32),
        "apple-touch-icon.png": (180, 180),
        "og-en.png": (1200, 630),
        "og-zh.png": (1200, 630),
    }
    for name, size in expected.items():
        image = Image.open(OUT / name)
        if image.size != size:
            die("%s is %s, want %s" % (name, image.size, size))
    print("wrote", ", ".join(sorted(p.name for p in OUT.iterdir())))


if __name__ == "__main__":
    main()
