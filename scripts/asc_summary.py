#!/usr/bin/env python3
import argparse
import json
from datetime import datetime
from pathlib import Path


def parse_iso(value):
    if not value:
        return None
    try:
        if value.endswith("Z"):
            value = value[:-1] + "+00:00"
        return datetime.fromisoformat(value)
    except ValueError:
        return None


def sort_key(item):
    attrs = item.get("attributes", {})
    created = attrs.get("createdDate") or ""
    parsed = parse_iso(created)
    return parsed or datetime.min


def detect_type(items):
    for item in items:
        attrs = item.get("attributes", {})
        if "rating" in attrs:
            return "reviews"
        if "comment" in attrs:
            return "feedback"
    return "unknown"


def summarize_reviews(items):
    ratings = []
    for item in items:
        rating = item.get("attributes", {}).get("rating")
        if isinstance(rating, int):
            ratings.append(rating)
    avg_rating = sum(ratings) / len(ratings) if ratings else None
    latest = sorted(items, key=sort_key, reverse=True)[:5]
    latest_entries = []
    for item in latest:
        attrs = item.get("attributes", {})
        latest_entries.append(
            {
                "id": item.get("id", ""),
                "createdDate": attrs.get("createdDate", ""),
                "rating": attrs.get("rating"),
                "title": attrs.get("title", ""),
                "body": attrs.get("body", ""),
                "territory": attrs.get("territory", ""),
            }
        )
    return {
        "count": len(items),
        "averageRating": avg_rating,
        "latest": latest_entries,
    }


def summarize_feedback(items):
    latest = sorted(items, key=sort_key, reverse=True)[:5]
    latest_entries = []
    for item in latest:
        attrs = item.get("attributes", {})
        latest_entries.append(
            {
                "id": item.get("id", ""),
                "createdDate": attrs.get("createdDate", ""),
                "comment": attrs.get("comment", ""),
                "deviceModel": attrs.get("deviceModel", ""),
                "devicePlatform": attrs.get("devicePlatform", ""),
                "osVersion": attrs.get("osVersion", ""),
            }
        )
    return {
        "count": len(items),
        "latest": latest_entries,
    }


def summarize_file(path, forced_type):
    data = json.loads(Path(path).read_text())
    items = data.get("data")
    if not isinstance(items, list):
        items = []

    data_type = forced_type or detect_type(items)
    if data_type == "reviews":
        summary = summarize_reviews(items)
    elif data_type == "feedback":
        summary = summarize_feedback(items)
    else:
        summary = {"count": len(items), "latest": []}

    return {
        "file": str(path),
        "type": data_type,
        **summary,
    }


def print_text(summary):
    print(f"file: {summary['file']}")
    print(f"type: {summary['type']}")
    print(f"count: {summary['count']}")
    if summary["type"] == "reviews":
        avg = summary.get("averageRating")
        print(f"averageRating: {avg:.2f}" if avg is not None else "averageRating: n/a")
    print("latest:")
    for entry in summary.get("latest", []):
        created = entry.get("createdDate", "")
        if summary["type"] == "reviews":
            rating = entry.get("rating")
            territory = entry.get("territory", "")
            title = entry.get("title", "")
            body = entry.get("body", "")
            print(f"- {created} | {rating} | {territory} | {title} | {body}")
        else:
            comment = entry.get("comment", "")
            device = f"{entry.get('devicePlatform','')} {entry.get('deviceModel','')}".strip()
            os_version = entry.get("osVersion", "")
            print(f"- {created} | {device} {os_version} | {comment}")
    print("---")


def main():
    parser = argparse.ArgumentParser(description="Summarize ASC JSON output files.")
    parser.add_argument("files", nargs="+", help="Path(s) to JSON files")
    parser.add_argument(
        "--type",
        choices=["reviews", "feedback", "auto"],
        default="auto",
        help="Force summary type or auto-detect",
    )
    parser.add_argument(
        "--format",
        choices=["json", "text"],
        default="json",
        help="Output format",
    )
    parser.add_argument("--pretty", action="store_true", help="Pretty-print JSON output")
    args = parser.parse_args()

    forced = None if args.type == "auto" else args.type
    summaries = [summarize_file(path, forced) for path in args.files]

    if args.format == "text":
        for summary in summaries:
            print_text(summary)
    else:
        output = {"summaries": summaries}
        if args.pretty:
            print(json.dumps(output, indent=2))
        else:
            print(json.dumps(output, separators=(",", ":")))


if __name__ == "__main__":
    main()
