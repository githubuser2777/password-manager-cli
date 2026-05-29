# Agents System for Antigravity (Spec-Driven Development)

This document outlines how Agents in the **Antigravity** system interact, coordinate, and execute tasks based on the **Spec-Driven Development (SDD)** philosophy (similar to GSD / GitHub Spec Kit).

Core objective: Prevent "Context Rot" (loss of context when chats get too long) and "Vibe Coding" (coding based on feelings without a plan).

---

## 1. Agent Routing Mechanism
Agents are not tied to cliché personas. Their roles are assigned based on the project's **Phase** and the current **Context**.

### A. Spec Agent
- **Behavior:** Collects input data from the user (ideas, old files, bug reports) and asks questions to strictly lock down the scope.
- **Responsibility:** Focuses 100% on "What" (What are we building?) and "Why" (Why does it need to be built?). DOES NOT write code or design architecture.
- **Output:** Creates/Updates `docs/spec.md`.

### B. Planning Agent
- **Behavior:** Consumes `docs/spec.md` to draft the technical design (System Blueprint).
- **Responsibility:** Makes decisions on "How" (How to build it?). Defines the stack, folder structure, API schemas, Data models, and libraries to be used.
- **Output:** Creates/Updates `docs/plan.md`.

### C. Task Agent
- **Behavior:** Converts `docs/plan.md` into an independently executable task checklist.
- **Responsibility:** Breaks down the workload. If a task takes too many tokens/steps to process, it must be broken down further.
- **Output:** Creates/Updates `docs/tasks.md`.

### D. Code Agent / Editor
- **Behavior:** Reads `docs/tasks.md` and executes each task sequentially.
- **Responsibility:** Opens files, writes code, fixes bugs, runs tests. After completing a task, automatically checks `[x]` in `docs/tasks.md` before moving to a new task. Avoids reading the entire source code unless necessary.

---

## 2. Communication Protocol
- **State Handoff:** When an Agent finishes its Phase, it must explicitly name the output file for the next Agent to take over (e.g., "Spec completed at `docs/spec.md`. Ready for Planning Phase").
- **Tool Activation:** Encouraged to use Command Line Tools (bash/zsh/go/npx) to automate linting, testing, or compiling directly in Antigravity's terminal.
