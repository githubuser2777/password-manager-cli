# Core Rules

This ruleset acts as the "Constitution" for all Agents operating within the Antigravity environment.

## 1. Chat & Formatting Behavior (No Fluff)
- **NO CLICHÉ PATTERNS:** Strictly forbidden to start prompts or answers with "You are an expert software engineer...", "I am a world-class AI...", "As an AI...". Get straight to solving the problem.
- **Short & Concise:** Answer directly. If the user asks for code, output only logic, code, and changes. Limit apologies or lengthy explanations.
- **Flexibility:** Rules are made to guide, not to slow down progress. If a change is extremely small (e.g., fixing a typo, changing 1 line of CSS), it is allowed to skip the Spec/Plan process and do a "Fast Track".

## 2. Error Handling & Assumptions (Anti-Hallucination)
- If the user provides a vague request that has a major impact on the architecture, **the Agent MUST stop and ask questions**.
- Do not arbitrarily generate mock APIs or non-existent libraries. Base actions on the specs in `docs/plan.md`.
- If an error occurs during execution (Error Logs), always read the log thoroughly, point out the anomaly before suggesting a fix. Avoid blindly spamming fix code.

## 3. File & Context Management (Context Limits)
- To avoid context overflow (Context Rot), do not open/read more than 5 files at once unless performing a search operation (grep/search).
- Write highly modularized code.

## 4. Coding Standards
- Clarity is prioritized over brevity. Use variable names that accurately describe their function.
- Adhere to DRY (Don't Repeat Yourself) reasonably; don't force DRY if it makes the code overly abstracted and hard to read.
- Leave comments in functions with complex business logic, branching algorithms, or code segments that are "workarounds".

## 5. Git & Changelog Process
- **Git Commit & Push**: The Agent must automatically commit and push after completing significant changes (e.g., finishing a phase).
- **Review before Commit**: DO NOT use `git add .` blindly. Check which files have changed (`git status`) and only add the files actually related to that commit.
- **Changelog**: Automatically create and update the `docs/changelog.md` file for every major change. Each change must be clearly separated (e.g., using Date headers, versions) but kept in the same file.
