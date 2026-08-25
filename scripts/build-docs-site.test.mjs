import assert from "node:assert/strict";
import { execFileSync } from "node:child_process";
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import test from "node:test";
import { fileURLToPath } from "node:url";

const sourceRoot = path.dirname(path.dirname(fileURLToPath(import.meta.url)));

test("builds safe TOC text from formatted headings in source order", (t) => {
  const root = fs.mkdtempSync(path.join(os.tmpdir(), "songsee-docs-"));
  t.after(() => fs.rmSync(root, { recursive: true, force: true }));
  fs.cpSync(path.join(sourceRoot, "scripts"), path.join(root, "scripts"), {
    recursive: true,
  });
  fs.mkdirSync(path.join(root, "docs"));
  fs.writeFileSync(path.join(root, "docs", "index.md"), "# Fixture docs\n");
  fs.writeFileSync(path.join(root, "docs", "quickstart.md"), "# Quickstart\n");
  fs.writeFileSync(
    path.join(root, "docs", "install.md"),
    [
      "---",
      "title: Fixture",
      "---",
      "# Fixture",
      "## **Bold** and *emphasis*",
      "> ### `code` and [linked](#target)",
      '## Tag-shaped <script>alert("x")</script>',
      "### Target",
      "",
    ].join("\n"),
  );

  execFileSync(process.execPath, ["scripts/build-docs-site.mjs"], {
    cwd: root,
  });
  const html = fs.readFileSync(
    path.join(root, "dist", "docs-site", "install.html"),
    "utf8",
  );
  const tocStart = html.indexOf('<nav class="toc"');
  const toc = html.slice(
    tocStart,
    html.indexOf("</nav>", tocStart) + "</nav>".length,
  );
  const links = [
    '<a class="toc-l2" href="#bold-and-emphasis">Bold and emphasis</a>',
    '<a class="toc-l3" href="#code-and-linked-target">code and linked</a>',
    '<a class="toc-l2" href="#tag-shaped-script-alert-x-script">Tag-shaped &lt;script&gt;alert(&quot;x&quot;)&lt;/script&gt;</a>',
    '<a class="toc-l3" href="#target">Target</a>',
  ];

  let previous = -1;
  for (const link of links) {
    const index = toc.indexOf(link);
    assert.ok(index > previous, `missing or out-of-order TOC link: ${link}`);
    previous = index;
  }
  assert.match(
    html,
    /<h3 id="code-and-linked-target"><a class="anchor" href="#code-and-linked-target"/,
  );
  assert.doesNotMatch(toc, /&amp;lt;|<script>/i);
  assert.doesNotMatch(html, /<script>alert\("x"\)<\/script>/);
});
