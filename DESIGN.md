---
name: jobtrack
description: Self-hosted job pipeline UI — departure-board rhythm on a dark rail shell
colors:
  surface: "#f4f1ea"
  surface-muted: "#e8e4da"
  surface-elevated: "#fffffc"
  foreground: "#161a20"
  foreground-muted: "#525862"
  border: "#c6c0b2"
  accent: "#b46212"
  accent-foreground: "#fffcf5"
  rail: "#1c2028"
  rail-foreground: "#f0ece4"
  status-applied: "#2962a8"
  status-offer: "#168050"
  status-rejected: "#b93030"
typography:
  ui:
    fontFamily: '"Source Sans 3", system-ui, sans-serif'
    fontWeight: 400
  data:
    fontFamily: '"JetBrains Mono", ui-monospace, monospace'
    fontWeight: 500
rounded:
  sm: "2px"
  md: "6px"
  lg: "8px"
spacing:
  page-x: "2rem"
  page-y: "2rem"
components:
  button-primary:
    backgroundColor: "{colors.accent}"
    textColor: "{colors.accent-foreground}"
    rounded: "{rounded.md}"
    padding: "0 1rem"
  panel:
    backgroundColor: "{colors.surface-elevated}"
    rounded: "{rounded.md}"
---

## Overview

Operate-mode web app for job seekers self-hosting **jobtrack**. Visual world: airport departure board (FIDS)—dark navigation rail, amber primary actions, monospace timestamps and counts, hairline grid on the work surface. Signature: horizontal scan animation on the top bar while Gmail sync runs.

## Colors

Light mode uses warm paper surfaces; dark mode deep charcoal with amber accent. Semantic status hues map to pipeline stages (applied, replied, interview, offer, rejected). State colors are never reused as decoration-only accents.

## Typography

**Source Sans 3** for UI copy and labels. **JetBrains Mono** (`.font-data`) for times, counts, sync status, and table numerics. Status badges use uppercase tracked labels at 11px.

## Layout

Fixed **52px** top bar with search + sync; **208px** dark sidebar on `md+`; bottom tab bar on mobile. Main content max readable width on settings/detail; board scrolls horizontally. Page titles use `.page-title` / `.page-lead`—no kickers.

## Elevation & Depth

`.panel` and `shadow-board` / `shadow-board-lg`: 1px border plus soft vertical offset shadow—no glow halos or glass decoration.

## Shapes

Most controls use **6px** radius (`rounded-md`). Badges use **2px** (`rounded-sm`). Avoid nested card-in-card layouts.

## Components

- **Button**: primary (amber), secondary (bordered panel), ghost, danger.
- **Badge**: uppercase status chips with semantic border/background tints.
- **EmptyState**: dashed-free panel with centered copy.
- **Nav**: rail with amber active bar; mobile bottom tabs.

## Do's and Don'ts

- Do use monospace for measurable fields (dates, counts, sync clock).
- Do keep Sync as the single saturated action in the top bar.
- Don't add gradient text, hero metric grids of four identical cards, or decorative blur.
- Don't use monospace for body paragraphs.
