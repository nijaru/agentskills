# agentskills

`agentskills` is a small Go implementation of the [Agent Skills
spec](https://agentskills.io/specification). It loads `SKILL.md` bundles,
validates skill names, and gives you a deterministic registry for discovery and
progressive disclosure.

If you want the broader ecosystem context, start here:

- Spec: [agentskills.io/specification](https://agentskills.io/specification)
- Implementor guide: [agentskills.io/client-implementation/adding-skills-support](https://agentskills.io/client-implementation/adding-skills-support)
- Canonical reference repo: [github.com/agentskills/agentskills](https://github.com/agentskills/agentskills)

## Use It

```go
package main

import (
	"fmt"

	"github.com/nijaru/agentskills"
)

func main() {
	reg := agentskills.NewRegistry("skills")
	if err := reg.Discover(); err != nil {
		panic(err)
	}

	for _, skill := range reg.List() {
		fmt.Println(skill.Summary())
	}
}
```

## What It Exposes

- `Skill` for parsed `SKILL.md` metadata and body text
- `Load` for reading a single `SKILL.md`
- `Registry` for discovery, lookup, and listing
- `ValidateName` for the skill-name rules used by the runtime tools

The module stays intentionally small so other Go projects can depend on it
without pulling in a full agent framework.
