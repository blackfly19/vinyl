# vinyl

> A small version control system, built from scratch in Go — designed to be **easy for beginners**. Command: `vin`.

`vinyl` records snapshots of your files over time so you can see what changed and go back to
an earlier version — the same core job as Git, but with a deliberately smaller, friendlier
mental model. Like a vinyl record, it keeps a faithful recording of your project that you can
play back to any point. It's built from the ground up (content-addressed storage, Merkle
trees, a commit history) as a **glass-box tool for understanding how version control actually
works under the hood**.

> **Status: early development.** The core workflow — initialize a repo, commit snapshots,
> and view history — works today. `revert` (restoring an old snapshot) is in progress. See
> [Command status](#command-status) for exactly what is and isn't implemented.

---

## Why vinyl?

Git is powerful, but its mental model (the staging area/index, detached HEAD, the reflog,
rebase…) is famously hard for newcomers. vinyl's guiding question for every feature is:
*"would a first-time user understand this?"* Concretely, that led to one notable design
decision:

**vinyl has no staging area and no `add` step.** In Git you `add` files to an index before
committing — a two-phase model that's one of the top things beginners trip on. vinyl drops
it: `commit` snapshots your working directory directly (what you see is what you commit).
Selecting *what not to* commit is handled by an ignore file rather than an opt-in staging
step — designed for the common case ("commit everything") instead of taxing it. *(The
ignore layer is on the roadmap; see below.)*

---

## Core concepts (in plain language)

If you've never used a VCS before, here's the whole idea in three terms:

| Term | What it means in vinyl |
| --- | --- |
| **Working directory** | The normal folder you edit your files in. |
| **Commit** | A saved snapshot of your project at one point in time, with a message describing it. Created with `vin commit`. |
| **Store** (`.vinyl/`) | A hidden folder vinyl creates to remember all your history. You never edit it by hand. |

A typical session:

```sh
vin init                       # start tracking this folder
vin commit -m "First version"  # save a snapshot of everything
vin log                        # see your history
```

---

## How it works under the hood

You don't need any of this to *use* vinyl — but it's the interesting part, and it's
documented here for the curious. vinyl borrows the proven ideas behind Git.

### Content-addressed storage

Every file's contents are hashed, and the file is stored as an **object** named by that hash
under `.vinyl/objects/`. Two identical files (or an unchanged file across many commits) map
to the same hash, so they're **stored only once** (deduplication). Because the name *is* the
hash of the content, storage is tamper-evident: change the bytes and the name would have to
change too.

### Merkle-tree snapshots

When you commit, vinyl walks your project and builds a **Merkle tree** that mirrors your
folders and files, where:

- each **file node** carries the hash of that file's contents, and
- each **directory node** carries a hash derived from its children's hashes.

Because every parent's hash depends on its children, the **root hash is a fingerprint of the
entire project** at that moment. Change one byte in one file and the change "bubbles up" to a
new root hash. Two directory nodes with the same hash are identical underneath — which is
what makes comparing two snapshots efficient (equal subtrees can be skipped).

### Commits and history

Each commit is a small record pointing at a snapshot's root tree plus metadata (message,
timestamp, parent commit). Commits form a **chain**: each points to its parent, and a single
`HEAD` file points to the latest commit. That's what powers `vin log` and moving between
versions. `HEAD` is the one mutable pointer; everything under `.vinyl/objects/` is immutable.

---

## Installation

Requires **Go 1.25+**.

```sh
git clone https://github.com/blackfly19/vcs.git vinyl
cd vinyl
go build -o vin .           # produces the `vin` binary
# optionally move it onto your PATH:
# mv vin /usr/local/bin/
```

---

## Commands

### `vin init`
Initializes a new store in the current folder. Creates the hidden `.vinyl/` directory and the
structures vinyl uses to track history. Run this once per project.

### `vin commit -m "message"`
Snapshots your entire working directory and saves it as a commit.

```sh
vin commit -m "Add project outline"
vin commit -m "Milestone build" -c    # -c marks this commit as a checkpoint
```

- `-m, --message` — the message describing this commit.
- `-c, --checkpoint` — mark the commit as a **checkpoint** (reserved for guaranteed
  full-snapshot "safe points" to return to).

### `vin log`
Displays your commit history — id, date, and message — newest first, walking the parent
chain from `HEAD`.

### `vin status` *(in progress)*
Intended to list files that have changed since your last commit. Being reworked to derive
its answer directly from the latest commit's tree (rather than a separate cached index).

### `vin revert <commitID>` *(in progress)*
Will restore your working directory to the snapshot captured by an earlier commit.

---

## Command status

| Command | Status |
| --- | --- |
| `init` | ✅ Working |
| `commit` | ✅ Working |
| `log` | ✅ Working |
| `status` | 🚧 In progress (being reworked) |
| `revert` | 🚧 In progress |

---

## The `.vinyl/` store layout

After `vin init`, your project contains a hidden `.vinyl/` folder. **You never edit this by
hand** — it's how vinyl remembers your history:

```
.vinyl/
├── objects/    # blob and tree objects — file contents and directory nodes,
│               #   each named by the hash of its content
├── commits/    # commit records, each named by its commit id
└── HEAD        # a pointer to the current (latest) commit
```

---

## Project layout (for contributors)

```
src/
├── cmd/          # the CLI commands (cobra): init, commit, log, status, revert
├── vinyl/        # core VCS logic: Merkle tree, snapshots, commits, persistence
├── constants/    # shared paths and names
└── utils/        # hashing and encoding helpers
```

---

## A note on the name

A **vinyl** is a record — and recording your project's history is exactly what this tool
does. Each commit lays down another groove; `log` reads back the track list; and (soon)
`revert` drops the needle on any earlier point. The command is shortened to **`vin`** so it
stays quick to type.
