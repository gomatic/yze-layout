# yze-go-layout

A [`yze`](https://github.com/gomatic/yze) analyzer (group `go`, category `structure`) enforcing the **cross-package** correspondence of the gomatic three-tier CLI layout: every `internal/app/commands/<cmd>` package has a matching `internal/domain/<cmd>` package, and vice versa.

Each package checks its own counterpart directory on the filesystem, so both directions are reported without duplication. Pairs with the per-package [`yze-go-pkgstd`](https://github.com/gomatic/yze-go-pkgstd) analyzer (which enforces the standards *within* a command package).

- **Rule:** `yze/go/layout`
- **Binary:** `cmd/yze-go-layout` runs it standalone.

Built on the [`go-yze`](https://github.com/gomatic/go-yze) framework.
