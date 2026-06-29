# yze-layout

A [`yze`](https://github.com/gomatic/yze) analyzer (category `structure`) enforcing the **cross-package** correspondence of the gomatic three-tier CLI layout: every `internal/app/commands/<cmd>` package has a matching `internal/domain/<cmd>` package, and vice versa. A counterpart must be a real Go package — a directory containing Go source — not a bare or empty directory.

Each package checks its own counterpart on the filesystem, so both directions are reported without duplication. Pairs with the per-package [`yze-pkgstd`](https://github.com/gomatic/yze-pkgstd) analyzer (which enforces the standards *within* a command package).

- **Rule:** `yze/layout`
- **Library:** exports `Analyzer` and `Registration` for the [`yze`](https://github.com/gomatic/yze) aggregator and [`stickler`](https://github.com/gomatic/stickler) runner.
- **Binary:** `cmd/yze-layout` runs it standalone (`text`/`-json`, and as a `go vet -vettool`).

Built on the [`go-yze`](https://github.com/gomatic/go-yze) framework.
