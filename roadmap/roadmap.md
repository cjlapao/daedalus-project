# Project Daedalus

*Codename: Daedalus — the architect who builds systems that run themselves.*

A code-driven orchestration engine for AI agents. Deterministic control flow on the
outside, intelligence on the inside, with context scoped tightly enough that decay
can't happen.

---

## 1. The Problem

Current agentic systems put the LLM in charge of orchestration. The orchestration
agent holds the harness in its context, and as work accumulates that context fills,
the original instructions slip out of the window, and the orchestrator quietly stops
orchestrating and starts doing the job itself. It becomes the very thing it was meant
to manage. This is context decay, and no amount of model intelligence fixes it —
it's an architectural problem, not a capability one.

These systems also demand constant human babysitting, because the moment the agent
loses the plot, a person has to step in and re-anchor it.

## 2. The Thesis

Stop asking the LLM to *be* the orchestrator. Make it a *worker the orchestrator
calls*.

- **Control flow lives in code.** Deterministic, inspectable, never forgotten.
- **Intelligence lives in stateless, single-shot agent calls.** Each call gets exactly
  the context it needs, does one small job, returns a result, and forgets everything.
- **State lives in a database.** The event-sourced ledger is the memory the LLM was
  failing to hold. The code remembers what the model can't.

Context limits stop mattering because no single call ever approaches them. The big
picture lives in code; the small picture lives, briefly, in the worker.

> Determinism on the outside, intelligence on the inside.

## 3. Core Objectives

1. **Eliminate context decay** by never asking one LLM to hold an entire job.
2. **Full automation** of the happy path — the orchestrator drives end to end, pausing
   only at human checkpoints that the roadmap itself declares.
3. **Durable, resumable execution** — a crash mid-run loses nothing; the system reboots,
   reads the ledger, and knows exactly where it stood.
4. **Honest reconciliation** — retries are never blind replays. The system inspects
   reality before re-attempting, and reality always wins over recorded state.
5. **Self-hostable, homelab-native** — no cloud dependency, no cluster you don't want
   to run, one coherent state model, zero vendor lock-in.
6. **Live observability** — watch the pipeline flow through its swim lanes in real time.

## 4. Architecture Overview

```
Roadmap template (per issue type, versioned, static)
        │  selected at kickoff, generated into ↓
Roadmap instance (this job's tasks + checkpoints, frozen on green light)
        │  traversed by ↓
Orchestrator (code: traverse · delegate · verify · reconcile · gate)
        │  reads/writes every step ↓
Event store (append-only log + state projection)  ← trusted record
        │                                              ↑ audited against
        │  delegates via Executor interface ↓          │
Agent = SKILL.md (user-written) + contract (Daedalus-bound)
        │  one slice, exact context, forgets after     │
        │  inspects / mutates ↓                         │
Real world (local filesystem)  ─────────────────────────┘ source of truth
```

The event store sits dead center because every arrow eventually touches it. That's
the tell that it's the spine, not a side concern. It is simultaneously the
crash-recovery ledger **and** the live feed powering the dashboard — one mechanism,
seen twice.

### 4.1 The Layers

**Development map (template).** A structured, declarative artifact defining *how a
type of work gets done here* — the decomposition rules, ordering/dependency rules,
where human checkpoints sit, and the definition of "done" for each phase. Stricter
than a prompt. It is not advice the model might forget; it is a spec the code reads
and enforces. Prompts are suggestions. This is law. Different issue types (feature,
bugfix, refactor) load different maps, and the library of maps grows over time into
durable institutional memory.

**Planner (one-shot AI call).** Reads the template for the issue type plus the actual
request, and generates the concrete task list — turning "implement feature X" into
capability-sized tasks with explicit declared dependencies. Runs once, bounded, up
front. Its output is the roadmap instance. The creativity happens here, exactly once.

**Roadmap instance.** The concrete task graph for this specific job. Frozen the moment
the human gives the green light. After freezing it is a *data structure code traverses*,
not a plan an AI re-litigates at every step.

**Orchestrator (deterministic code).** A chat-facing entry point. It does not *think*
about how to run a job — it *executes* a known procedure: traverse the roadmap, delegate
each task, verify each result, reconcile on failure, halt at checkpoints, write every
step to the event store. It never loses the plot because the plot isn't in an LLM's
context — it's in the program counter.

**Executor (pluggable interface).** The orchestrator never knows *how* a task runs. It
hands an executor-agnostic task contract to an `Executor` and gets a result back. Two
v1 implementations: Claude Code (CLI agent) and a direct LiteLLM call. The contract is
described in terms of *what capability must exist and how to verify it* — never "run
this specific command."

**Worker / reviewer / reconciliation agents (stateless singletons).** Each call is
fire-and-forget context. The builder builds; a *separate* reviewer judges (never the
agent that wrote the code); the reconciliation agent inspects reality. No agent carries
the weight of the whole project.

### 4.2 Agents as Declarative Skills

An agent is a **declarative artifact, not a hardcoded integration.** A user writes a
`SKILL.md` — the same plain format already used across the homelab — describing what an
agent knows and does (e.g. a "Go engineer agent"). Daedalus *reads* that file and brings
the agent to life. Adding a new agent never requires touching orchestrator code: drop in
a skill file and Daedalus creates it. This is MCP-like pluggability, but lighter — a skill
file is enough, no server required, though an MCP server can be one of the backing shapes.

This is the project's own philosophy repeated one level down: **the user writes the
intelligence (the skill); Daedalus writes the determinism (the contract).** The skill
says *what the agent is*; Daedalus binds the orchestration wrapper that makes it a
well-behaved, observable, callable tool in the harness. In Daedalus terms an agent is
therefore **skill + contract**, never the skill alone.

What Daedalus adds on top of the user's skill:

- **How to respond** — a structured response the orchestrator can read as data, not free
  prose it has to guess at.
- **What to call back when done** — a defined completion signal, so the agent doesn't
  merely stop; it reports back through a known channel and the orchestrator collects the
  result.
- **Liveness / health** — is this agent online, reachable, responding? Code owns this.
- **Harness plumbing** — parsing the response, routing it into the event store, handing
  it to the done-check, feeding reconciliation.

This is the same seam as the `Executor` interface, viewed from the other side. Whatever
backs an agent — Claude Code, a raw LiteLLM call, an MCP server — it must honor the
contract Daedalus imposes. Because every agent speaks the same response/callback contract,
the differences between executor backends get **normalized at this boundary**: the
orchestrator does not care what is behind an agent, only that the agent honors the
contract.

**Response contract — v1: universal.** For v1 Daedalus injects a single fixed response
envelope for *every* agent (status, completion signal, errors — the fields the
orchestrator must have to function). One response shape to parse, one callback path, one
validation rule; every agent looks identical from the orchestrator's seat. This proves
the whole skill-to-agent pipeline before any per-agent complexity is added.

**Response contract — v2: hybrid (planned).** The mandatory universal envelope stays, and
a skill *may additionally declare* a task-specific payload inside it (a Go engineer agent
returning `files_changed`; a reviewer returning `verdict` + `issues`) that the done-check
knows how to read. Because v2 only *adds* an optional payload within the v1 envelope,
nothing built in v1 is discarded — it is a clean, additive upgrade.

### 4.3 Key Design Properties

**Capability-sized tasks.** Not as granular as "create this file," not as coarse as
"implement an authenticated API." A task is one coherent capability that either works
or doesn't — "implement bearer token auth," "implement CORS support," "implement config
environment variables." The line is drawn by *capability boundary*, because that is the
granularity at which "is this done?" has a crisp, unambiguous answer. The agent owns
file-level decisions (the *how*); the roadmap owns the capability (the *what*).

**The done-check, used twice.** Every task carries both an execution instruction *and*
a verification contract. Two gates, deterministic first:

1. **Deterministic gate** — tests pass, it compiles, the linter is happy. Objective,
   cheap, certain. If this is red, the task failed; no tokens are spent going further.
2. **Semantic gate** — a separate reviewer agent judges whether the capability is
   actually implemented and sane. Only runs once gate 1 is green.

The same done-check that decides "is this task complete?" is what the reconciliation
agent re-runs to decide "was this already done before we retried?" One mechanism, two
uses.

**Reconciliation, not replay.** When a task fails mid-flight, the orchestrator never
assumes a clean slate. Before re-delegating, a reconciliation agent inspects the actual
filesystem to learn what is *really* there. The orchestrator diffs that against what the
ledger expected: if the world agrees, resume from the gap; if the world disagrees, the
world wins and the ledger is corrected. Idempotency by construction — every task must be
able to answer "am I already done?"

**Event sourcing.** State is an append-only log of immutable facts
("task 4 attempt 2 failed with X," "task 4 attempt 3 succeeded"); current state is a
projection derived from those events. This gives the audit trail for free, makes retries
trivial to reason about, turns "what's already been done?" into a query, and means live
visualization and historical replay are the same data read two ways.

## 5. v1 Scope (the first buildable cut)

A Go orchestrator on a single host, fronted by a chat interface and a live React
dashboard. You ask it to implement a software feature. It loads a per-issue-type
roadmap template, uses a one-shot planner call to decompose the request into
capability-sized tasks with explicit dependencies, freezes that on your green light,
then walks it **sequentially** — delegating each task through a pluggable executor
(Claude Code or LiteLLM) to the **local filesystem**. Each result faces the two-gate
done-check. On failure, a reconciliation agent inspects the actual filesystem before
retrying. Every step is an immutable event that serves as both crash-recovery ledger
and live dashboard feed.

### In scope for v1

| Area | v1 decision |
|---|---|
| First pipeline | Software feature implementation (the API example) |
| Codebase mode | **Greenfield-first** (existing-codebase supported by the abstractions, exercised in v1.5) |
| Execution model | Sequential (topological order over a dependency graph) |
| Executor | Pluggable interface; Claude Code + LiteLLM implementations |
| Agents | Declarative — defined by user-written `SKILL.md`, loaded by Daedalus |
| Agent contract | Universal response/callback envelope for every agent (hybrid in v2) |
| Location | Local filesystem, same host as orchestrator |
| Real world | Filesystem as-is; reconciliation agent inspects what exists |
| Done-check | Deterministic gate (tests/build/lint) **then** separate reviewer agent |
| Interaction | Chat (to the orchestrator) + live dashboard, together |
| Visualization | Real-time / live, streamed from the event log |
| Concurrency | Single user, single run at a time |
| State store | SQLite to start; Postgres if/when needed; Redis only if a real need appears |

### Explicitly NOT in v1 (deliberately deferred)

- Parallel task execution (designed for via the dependency graph, executed later)
- Existing-codebase discovery mode (abstractions ready; not the v1 target)
- Remote / SSH / Proxmox execution and per-task containers
- Multiple concurrent runs, run isolation, multi-user
- Historical replay UI (falls out of event sourcing later, nearly free)
- Heavy durable-execution machinery (Temporal-style replay) — not needed at
  sequential, single-run scale

## 6. Build-vs-Buy Stance

The advantage is **not** reimplementing solved plumbing. It is owning the layers that
make Daedalus *Daedalus*, while keeping it self-hostable and free of lock-in.

**Build (this is the product):** roadmap template system, planner contract + frozen-roadmap
discipline, reconciliation loop, task schema + done-check, the orchestration policy,
the visualization / swim-lane view.

**Thin self-built plumbing for v1 (genuinely small at this scale):** durable execution
for a single sequential cursor is "read the ledger on startup, resume from the last
committed task, keep tasks idempotent" — a few hundred lines, fully understood, zero
lock-in. The topological walk over your own task table is textbook code.

**Keep swappable:** the orchestrator talks to "the ledger" and "the executor" through
interfaces narrow enough that heavy machinery (Temporal for durable execution, LangGraph
for graph execution) can slot in *behind* those boundaries later — when parallel,
long-running, multi-day workflows actually justify it — without touching the policy layer
that is the real product. Build it yourself, but build it swappable. "Build" today does
not foreclose "buy" tomorrow.

## 7. Technology Choices

- **Backend / orchestrator:** Go.
- **State store:** SQLite first (single host, single run — perfect fit). Postgres as a
  later swap if scale or concurrency demands it. Redis only if a concrete need surfaces
  (it is not assumed).
- **Frontend / dashboard:** TypeScript + React.
- **Live transport:** event stream (websocket or SSE) tapping the same append-only log
  that provides durability.
- **Agent execution:** via the pluggable `Executor` interface — Claude Code and the
  existing LiteLLM routing layer as the two initial backends.

## 8. High-Level Feature List

- **Issue-type roadmap templates** — selectable at kickoff; a growing, versioned library.
- **Declarative agents (SKILL.md)** — users add agents by writing a skill file; Daedalus
  reads it and creates the agent. No orchestrator code changes to add capability.
- **Agent contract layer** — Daedalus binds a response/callback contract around each skill
  (universal envelope in v1, hybrid with optional declared payload in v2).
- **Agent health / liveness** — code tracks whether each agent is online and responding.
- **AI planner** — one-shot decomposition of a request into a capability-sized,
  dependency-aware task graph.
- **Green-light gate** — human approval that freezes the roadmap before execution.
- **In-roadmap checkpoints** — human gates declared as node types in the roadmap, halting
  the generic orchestrator wherever the map says.
- **Sequential orchestrator** — topological traversal, delegation, verification, gating.
- **Pluggable executor interface** — Claude Code + LiteLLM backends, swappable and
  extensible.
- **Two-gate done-check** — deterministic checks, then independent reviewer agent.
- **Reconciliation engine** — inspect-reality-before-retry; world beats ledger.
- **Event-sourced state store** — durable, resumable, auditable; the system's spine.
- **Chat interface** — converse with the deterministic orchestrator (not an agent).
- **Live dashboard** — real-time swim-lane view streamed from the event log.

## 9. Phased Roadmap

**Phase 0 — Foundations & contracts.** Define the task schema (execution instruction +
verification contract + dependency edges), the event model, the `Executor` interface, and
the **universal agent response/callback contract** plus the `SKILL.md` shape Daedalus
reads to create an agent. Get the boundaries right; everything hangs off these.

**Phase 1 — The spine.** Event store (SQLite) with append + projection. Orchestrator
skeleton that can persist and resume from the ledger. No AI yet — drive a hand-written
roadmap of trivial tasks to prove traverse / record / resume.

**Phase 2 — Execution & agents.** Implement the `Executor` interface with one backend, the
**skill loader** (read a `SKILL.md`, bind the universal contract, expose it as a callable
agent), and the **health/liveness** check. Wire in the deterministic gate (run
tests/build/lint, read the result). Prove a real task can be delegated to a
skill-defined agent, run, be checked, and recorded.

**Phase 3 — Verification & reconciliation.** Add the reviewer agent (semantic gate) and
the reconciliation agent (inspect filesystem, diff against ledger, world wins). Prove the
retry loop reconciles rather than replays.

**Phase 4 — Planning.** Add the planner: template + request → frozen task graph. Add the
green-light gate and in-roadmap checkpoints. Now the full happy path runs from a chat
request to a built feature.

**Phase 5 — Surfaces.** Chat front-end to the orchestrator; live React dashboard streaming
the event log into swim lanes. Make the running machine watchable.

**Phase 6 — Hardening & second executor.** Add the second executor backend, exercise the
greenfield feature pipeline end to end on something real, tighten failure handling.

*Beyond v1:* hybrid agent response contract (optional declared payload inside the
universal envelope), existing-codebase discovery mode (v1.5), parallel execution over the
graph, historical replay UI, remote/containerized execution, heavier durable-execution
machinery behind the existing interfaces if and when justified.

## 10. Guiding Principles (the things not to drift from)

1. The orchestrator is code. The agents are tools the code calls. Never the reverse.
2. No LLM ever holds the whole job. Context scoped per call, forgotten after.
3. The ledger is the trusted record; reality is the ultimate truth; reality wins.
4. Roadmaps are frozen after green light. The planner runs once, not continuously.
5. Every task knows how to check whether it is already done.
6. Build the parts that are yours; keep the plumbing thin and swappable.
7. Agents are declarative — users write skills, Daedalus binds the contract. Adding an
   agent never means editing orchestrator code.
8. Self-hostable, one coherent state model, no lock-in.