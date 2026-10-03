# vinyl

> A small version control system, built from scratch in Go and designed to be **easy for beginners**. Command: `vin`.

`vinyl` records snapshots of your files over time so you can see what changed and go back to
an earlier version. It does the same core job as Git, but with a deliberately smaller,
friendlier mental model. Like a vinyl record, it keeps a faithful recording of your project
that you can play back to any point. It is built from the ground up (content-addressed
storage, Merkle trees, a commit history) as a **glass-box tool for understanding how version
control actually works under the hood**.

```sh
vin init                        # start tracking this folder
vin commit -m "First version"   # save a snapshot of everything
vin status                      # see what's changed since the last commit
vin log                         # read back your history
vin restore notes.txt           # put a file back the way it was committed
```

---

## Why vinyl?

vinyl takes its philosophy from Go, the language it is written in: **stay simple, and give one
obvious way to do each thing.** Go keeps its feature set small on purpose; there is a single
loop keyword, a single tool to build and format code, and very little overlap between
features. vinyl aims for the same feeling. Each command does one job, there is one way to do
it, and there are no flags or modes you need to learn before you can be productive.

That principle drives the one decision newcomers notice most: **vinyl has no staging area and
no `add` step.** In Git you `add` files to an index before committing, a two-phase model that
is one of the top things beginners trip on. vinyl drops it. `commit` snapshots your working
directory directly, so what you see is what you commit. One less concept to learn, and the
common case (commit everything) is the default instead of something you opt into file by file.

---

## Core concepts (in plain language)

If you have never used a VCS before, here is the whole idea in three terms:

| Term | What it means in vinyl |
| --- | --- |
| **Working directory** | The normal folder you edit your files in. |
| **Commit** | A saved snapshot of your project at one point in time, with a message describing it. Created with `vin commit`. |
| **Store** (`.vinyl/`) | A hidden folder vinyl creates to remember all your history. You never edit it by hand. |

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

## Getting started

```sh
mkdir my-project && cd my-project
vin init                          # create the .vinyl/ store

echo "hello" > notes.txt
vin commit -m "Start the project" # snapshot everything

echo "hello, world" > notes.txt
vin status                        # notes.txt shows up as modified

vin restore notes.txt             # back to "hello"
vin log                           # see the commit you made
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
vin commit -m "Milestone build" -c    # mark this commit as a checkpoint
```

- `-m, --message` sets the message describing this commit.
- `-c, --checkpoint` marks the commit as a **checkpoint**, a guaranteed full-snapshot safe
  point you can always return to.

### `vin status`
Shows what has changed in your working directory since the last commit, grouped into
**added**, **modified**, and **deleted** files. It compares your files against the latest
commit's tree directly, so there is no separate index to keep in sync. In a brand-new repo
with no commits yet, every file is listed as new.

### `vin log`
Displays your commit history (id, date, and message), newest first, walking the parent chain
from `HEAD`.

### `vin restore [-c <commitID>] <path>`
Puts a file back to the way it was in a commit, discarding your working changes to it. By
default it restores from the latest commit; pass `-c` to restore from a specific one.

```sh
vin restore notes.txt                            # restore from the latest commit
vin restore -c <commitID> notes.txt              # restore from an earlier commit (id from `vin log`)
```

If the file (or its folder) no longer exists in your working directory, vinyl recreates it.

---

## How it works under the hood

You do not need any of this to *use* vinyl, but it is the interesting part. vinyl borrows the
proven ideas behind Git.

### Content-addressed storage

Every file's contents are hashed, and the file is stored as an **object** named by that hash
under `.vinyl/objects/`. Two identical files, or an unchanged file across many commits, map to
the same hash, so they are **stored only once** (deduplication). Because the name *is* the
hash of the content, storage is tamper-evident: change the bytes and the name would have to
change too.

### Merkle-tree snapshots

When you commit, vinyl walks your project and builds a **Merkle tree** that mirrors your
folders and files, where:

- each **file node** carries the hash of that file's contents, and
- each **directory node** carries a hash derived from its children's hashes.

Because every parent's hash depends on its children, the **root hash is a fingerprint of the
entire project** at that moment. Change one byte in one file and the change bubbles up to a
new root hash. Two directory nodes with the same hash are identical underneath, which is what
lets `vin status` compare two snapshots efficiently: equal subtrees are skipped instead of
re-examined file by file.

### Commits and history

Each commit is a small record pointing at a snapshot's root tree plus metadata (message,
timestamp, parent commit). Commits form a **chain**: each points to its parent, and a single
`HEAD` pointer tracks the latest commit. That is what powers `vin log` and restoring from
older versions. `HEAD` is the one mutable pointer; everything under `.vinyl/objects/` is
immutable.

---

## The `.vinyl/` store layout

After `vin init`, your project contains a hidden `.vinyl/` folder. **You never edit this by
hand.** It is how vinyl remembers your history:

```
.vinyl/
├── objects/    # blob and tree objects: file contents and directory nodes,
│               #   each named by the hash of its content
├── commits/    # commit records, each named by its commit id
└── HEAD        # a pointer to the current (latest) commit
```

---

## Project layout (for contributors)

```
src/
├── cmd/          # the CLI commands (cobra): init, commit, log, status, restore
├── vinyl/        # core VCS logic: Merkle tree, snapshots, commits, persistence
├── constants/    # shared paths and names
└── utils/        # hashing and encoding helpers
```

---

## A note on the name

A **vinyl** is a record, and recording your project's history is exactly what this tool does.
Each commit lays down another groove, `log` reads back the track list, and `restore` drops the
needle on any earlier point. The command is shortened to **`vin`** so it stays quick to type.
