#!/usr/bin/env python3
from __future__ import annotations

import json
from pathlib import Path
from urllib.parse import urlparse

START_MARKER = "<!-- WALL-OF-APPS:START -->"
END_MARKER = "<!-- WALL-OF-APPS:END -->"
ALLOWED_PLATFORMS = ("IOS", "MAC_OS", "TV_OS", "VISION_OS")


def read_entries(source_path: Path) -> list[dict[str, str | list[str]]]:
    if not source_path.exists():
        raise SystemExit(f"Missing source file: {source_path}")

    raw = source_path.read_text(encoding="utf-8").strip()
    if not raw:
        raise SystemExit(f"Source file is empty: {source_path}")

    try:
        parsed = json.loads(raw)
    except json.JSONDecodeError as exc:
        raise SystemExit(f"Invalid JSON in {source_path}: {exc}") from exc

    if not isinstance(parsed, list):
        raise SystemExit(f"Expected a JSON array in {source_path}")

    entries = [normalize_entry(item, idx + 1) for idx, item in enumerate(parsed)]
    entries.sort(key=lambda entry: (entry["app"].lower(), entry["link"].lower()))
    return entries


def normalize_entry(entry: object, index: int) -> dict[str, str | list[str]]:
    if not isinstance(entry, dict):
        raise SystemExit(f"Entry #{index}: expected object")

    app = str(entry.get("app", "")).strip()
    link = str(entry.get("link", "")).strip()
    creator = str(entry.get("creator", "")).strip()
    platforms_raw = entry.get("platform")

    if app == "":
        raise SystemExit(f"Entry #{index}: 'app' is required")
    if link == "":
        raise SystemExit(f"Entry #{index}: 'link' is required")
    if creator == "":
        raise SystemExit(f"Entry #{index}: 'creator' is required")

    parsed_url = urlparse(link)
    if parsed_url.scheme not in {"http", "https"} or parsed_url.netloc == "":
        raise SystemExit(f"Entry #{index}: 'link' must be a valid http/https URL")

    if not isinstance(platforms_raw, list) or len(platforms_raw) == 0:
        raise SystemExit(f"Entry #{index}: 'platform' must be a non-empty array")

    platforms: list[str] = []
    for value in platforms_raw:
        platform = str(value).strip().upper()
        if platform not in ALLOWED_PLATFORMS:
            allowed = ", ".join(ALLOWED_PLATFORMS)
            raise SystemExit(
                f"Entry #{index}: invalid platform {platform!r} (allowed: {allowed})"
            )
        if platform not in platforms:
            platforms.append(platform)

    return {
        "app": app,
        "link": link,
        "creator": creator,
        "platform": platforms,
    }


def build_snippet(entries: list[dict[str, str | list[str]]]) -> str:
    lines = [
        "## Wall of Apps",
        "",
        "Apps shipping with asc-cli. [Add yours via PR](https://github.com/rudrankriyam/App-Store-Connect-CLI/pulls)!",
        "",
        "| App | Link | Creator | Platform |",
        "|:----|:-----|:--------|:---------|",
    ]
    for entry in entries:
        app = escape_cell(str(entry["app"]))
        link = str(entry["link"])
        creator = escape_cell(str(entry["creator"]))
        platforms = ", ".join(str(value) for value in entry["platform"])
        lines.append(
            f"| {app} | [Open]({link}) | {creator} | {escape_cell(platforms)} |"
        )
    return "\n".join(lines) + "\n"


def escape_cell(value: str) -> str:
    return value.replace("|", "\\|").replace("\n", " ").strip()


def write_generated_doc(snippet: str, generated_path: Path) -> None:
    generated_path.parent.mkdir(parents=True, exist_ok=True)
    header = (
        "<!-- Generated from docs/wall-of-apps.json by scripts/update-wall-of-apps.py. -->\n\n"
    )
    generated_path.write_text(header + snippet, encoding="utf-8")


def sync_readme(snippet: str, readme_path: Path) -> None:
    if not readme_path.exists():
        raise SystemExit(f"Missing README file: {readme_path}")

    content = readme_path.read_text(encoding="utf-8")
    if START_MARKER not in content or END_MARKER not in content:
        raise SystemExit(
            "README markers not found. Expected WALL-OF-APPS markers in README.md."
        )

    before, remainder = content.split(START_MARKER, 1)
    _, after = remainder.split(END_MARKER, 1)

    updated = f"{before}{START_MARKER}\n{snippet}{END_MARKER}{after}"
    readme_path.write_text(updated, encoding="utf-8")


def main() -> None:
    repo_root = Path(__file__).resolve().parents[1]
    source_path = repo_root / "docs" / "wall-of-apps.json"
    generated_path = repo_root / "docs" / "generated" / "app-wall.md"
    readme_path = repo_root / "README.md"

    entries = read_entries(source_path)
    snippet = build_snippet(entries)
    write_generated_doc(snippet, generated_path)
    sync_readme(snippet, readme_path)

    print(f"Updated {generated_path}")
    print(f"Synced snippet markers in {readme_path}")


if __name__ == "__main__":
    main()
