# Learning Log

`learning-log.jsonl` is the compact curriculum and review history for this repository. Each line is one complete JSON object for one exercise.

## What it records

- exercise number, folder, difficulty, and completion state
- the branch containing the submitted solution
- concepts practiced and demonstrated
- concrete mistakes or gaps found during review
- changes made during review
- useful next topics
- how the solution was validated

## Workflow

1. When a new exercise is published, add a record with `"status":"assigned"`.
2. A pushed branch is a submission, not automatic proof of completion.
3. Review the solution against its README and tests.
4. Update that exercise's record to `"completed"` or `"needs_revision"`.
5. Use `next_focus`, recent topics, and demonstrated difficulty when selecting the next exercise.
6. Read the most recent exercise source and tests when authoring or reviewing. The log is an index, not a replacement for checking code.

Keep one line per exercise and keep every line valid JSON so scripts and assistants can parse the file without a custom database.
