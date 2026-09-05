---
summary: "Release checklist for songsee (GitHub release + Homebrew tap + Pages)"
---

# Releasing songsee

Follow these steps for each release. Title GitHub releases as `songsee <version>`.

## Checklist

- Obtain explicit release authorization before creating or pushing a tag or publishing a GitHub release.
- Finalize `CHANGELOG.md` from all user-visible changes since the previous tag; replace `Unreleased` with the release date without changing earlier releases. Commit and push the release preparation.
- Verify the exact release commit has green `ci` and `docker` workflows. Locally run `go test ./... -cover`, `golangci-lint run`, `node --test scripts/build-docs-site.test.mjs`, and `make docs-site`.
- Build with `make VERSION=v<version>` and run `bin/songsee --version`, then render `testdata/sine.mp3` and a non-WAV/MP3 input through real ffmpeg. Check the output images, including a render with all nine modes.
- Tag the verified release commit: `git tag -a v<version> -m "Release <version>"`, then push that tag.
- Create the GitHub release for `v<version>` with title `songsee <version>` and a body matching the finalized changelog section. Verify the published tag, source archive, and release notes.
- Verify the `Update Homebrew Tap` workflow triggered by the tag. It dispatches `update-formula.yml` in `steipete/homebrew-tap`, which updates the formula from `https://github.com/openclaw/songsee/archive/refs/tags/v<version>.tar.gz`. The dispatcher requires the repository's `HOMEBREW_TAP_TOKEN` secret and waits for the tap run to complete. If it fails, resolve the reported cause before rerunning; do not substitute a separately generated archive checksum.
- Verify Homebrew install: `brew update && brew reinstall steipete/tap/songsee && songsee --version`, then render the smoke-test fixture with the installed binary.
- Verify `go install github.com/steipete/songsee/cmd/songsee@v<version>` reports the release version from the newly installed binary. The module path retains the original owner name.
- Verify the latest applicable `pages` workflow deployed the current docs to `songsee.sh`.
- Reopen the changelog with the next patch version's `Unreleased` section after the release; preserve the dated section and published release notes.

The release channels are GitHub source archives, the Go module tag, and Homebrew. The Dockerfile supports local image builds; this repository has no container-registry publishing workflow. npm publication and macOS app notarization do not apply.
