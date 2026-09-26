<!-- mlc-dochub:begin — auto-managed, do not edit between these markers -->
## MLC Doc Hub — mlcssystray

Structured documentation lives in `.mlcai/`, maintained through the
`mlc-dochub` MCP server.

- **Project ID:** `mlcssystray` — MLC SSH Tray Agent
- **Source directory:** `/mnt/data2tb/mlcssystray` (read source files with native Read —
  the MCP tools only touch `.mlcai/`)

### Existing docs

Read any of these directly (native Read is fine) to gather context before you work:

_(no docs created yet)_

### Recommended, not yet created

- `.mlcai/INTEGRATION.md` — How this project fits into the larger system
- `.mlcai/TECH_STACK.md` — Stack & dependencies
- `.mlcai/API_CONTRACT.md` — API contract / endpoints

### Tickets, worklog, release notes

`BACKLOG.md`, `WORKLOG.md` and `RELEASE_NOTES.md` are indexes the server
generates — never write them. Each entry is its own file, readable with
`get_doc` (or native Read): `backlog/<ID>.md`, `worklog/<ID>.md`,
`releases/<version>.md`. Tickets: `*_backlog_item` tools; close with
`update_backlog_item status=done` and the commit hash (`resolution=wontfix`
etc. when it was not fixed). A fix that still needs verifying:
`state=retest`; an accepted known issue: `state=known` (+ `user_facing=true`
if users notice it — it then belongs in the release notes). A missing or stale doc
is a doc ticket (`type=doc`, `doc_type=X.md`). At the end of a session write
**your** work strand: `update_worklog` with `entry_id` = the ticket you worked
on (or the `W-…` id you got back); it never touches another agent's strand.

### Working with `.mlcai/`

`.mlcai/` must be tracked by git — preferably as a submodule with its own
private repo, so the docs reach the Doc Hub and every other checkout. Each
write reports its git result: if it says a commit/push **failed** or the docs
are **NOT VERSIONED**, stop and tell the user instead of carrying on.

**Never write a `.mlcai/` file with a native editor.** Every create / update /
delete goes through the `mlc-dochub` MCP tools — they stamp the `## 📋 Meta`
footer, append to the activity log and guard against concurrent edits. Reading
with a native Read is fine and usually cheaper.

**If the tools are not available to you, read but do not write** — and point the
user at `task install-all` in the mlcintegration checkout (https://github.com/mlc911/mlcintegration).

The server states its full operating rules on connect (`author=`, `base_modified`,
which doc serves which purpose). Clients that drop server-level instructions —
Antigravity does, verified 29.08.2026 — get the same rules from the global
`~/.gemini/GEMINI.md`, section 6.
<!-- mlc-dochub:end -->
