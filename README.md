# agentskills-go

`agentskills-go` is a small Go implementation of the Agent Skills spec for
loading, validating, and discovering `SKILL.md` bundles.

It follows the current `agentskills` ecosystem conventions:

- `SKILL.md` frontmatter for metadata
- progressive disclosure through a registry/catalog
- deterministic discovery and listing

The module is intentionally small so other Go projects can use it as a reusable
skills core without taking a full agent framework dependency.
