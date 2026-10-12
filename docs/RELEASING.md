# Releasing songsee

Publication runs only in `.github/workflows/release.yml`. Prepare a dated
`CHANGELOG.md` section on protected `main`, merge it, and wait for green push CI
on that exact commit. With explicit release authorization, push an annotated
SSH-signed `vX.Y.Z` tag whose signer is in `.github/release-allowed-signers`.
The workflow validates the signature and changelog, builds and verifies the
source archive, and creates a GitHub Release with the exact notes and SHA256SUMS.
Existing releases are refused; tags and assets are never replaced.

For a no-publish rehearsal, dispatch `release.yml` from `main` with
`dry_run=true` (the default). Leave `version` blank to validate the latest dated
changelog section against the current source. This builds and retains the same
archive, runs the checks, and never creates a tag or release. Dispatches cannot
publish: `dry_run=false` fails validation. Publication requires a signed tag push.

The Go module remains available at `github.com/steipete/songsee` from the same tag. After publication the release workflow explicitly calls the Homebrew handoff for `steipete/homebrew-tap`; it requires `HOMEBREW_TAP_TOKEN`. Docker CI only verifies local images; no registry image is shipped.
