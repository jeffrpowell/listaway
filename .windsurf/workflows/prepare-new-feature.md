---
description: Systematic process for preparing and documenting a new feature.
---

# New Feature Workflow

1. **Gather Feature Requirements**
    - Ask the user for the GitHub issue number, feature title, and a description of the feature (unless already provided).
    - Create a markdown file in `.windsurf/features/` named after the issue number and feature title (e.g., `issue3_feature-title.md`).
    - Record requirements in this file.

2. **Understand Project Context**
    - Read `README.md` and `architecture.md`.
    - Search for relevant code sections and review them.
    - Summarize any relevant context in the feature file. (ie. important classes, files, etc.)

3. **Draft Implementation Plan**
    - Present a tentative implementation plan to the user for feedback.
    - **Ask user:** “Does this plan look good, or would you like to adjust it?”
    - Once approved, add the plan to the feature file.

4. **Prepare for Implementation**
    - Check if you are on the `main` branch.
    - If yes: `git pull` to ensure up-to-date.
    - Create and checkout the new feature branch using the name of the feature file as the branch name.
    - If not on `main`: **Ask user:** “Should we switch to `main` or use the current branch?”

5. **Implement the Feature**
    - Use the checklist in the feature file to track progress.
    - Add any additional notes or open questions to the feature file as you go.

6. **Post-Implementation**
    - Update README.md and architecture.md, if appropriate.

---

## Example Feature File Template

```markdown
# Issue 3 - Add new feature

## Requirements
- Requirement 1
- Requirement 2

## Context / Background
- Key points from README or architecture docs

## Plan
1–2 sentence summary of plan
- [ ] Step 1
  - [ ] Substep 1
- [ ] Step 2
- [ ] Step 3

## Open Questions / Dependencies
- Question or dependency 1

## Additional Notes
- Note 1
- Note 2
```