import DocsLayout from "@/components/docs/DocsLayout";
import CodeBlock from "@/components/docs/CodeBlock";
import { GitCommit, Layers, Clock, ShieldCheck } from "lucide-react";
import {
  commitInFlags as flags,
  commitInExitCodes as exitCodes,
  commitInAutoInit as autoInit,
  commitInProfileJson as profileJson,
} from "./commitInData";

const CommitInPage = () => (
  <DocsLayout>
    <div className="max-w-4xl space-y-10">
      <div>
        <div className="flex items-center gap-3 mb-2">
          <GitCommit className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold tracking-tight">commit-in</h1>
          <span className="font-mono text-xs px-2 py-1 rounded bg-primary/10 text-foreground border border-primary/20 dark:bg-primary/15 dark:text-primary dark:border-primary/40">
            alias: cin
          </span>
        </div>
        <p className="text-lg text-muted-foreground">
          Walk one or more SOURCE git repos in author-date order and APPEND each commit
          (preserving BOTH <code>AuthorDate</code> and <code>CommitterDate</code>) into a
          TARGET repo. Useful for stitching together project history that lives across forks,
          archives, or versioned siblings into a single canonical timeline — without ever
          rewriting an existing commit.
        </p>
        <p className="text-xs text-muted-foreground mt-2">
          Spec: <code>spec/03-commit-in/</code>
        </p>
      </div>

      <section>
        <h2 className="text-xl font-semibold mb-3">Overview</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[
            { icon: Clock, title: "Chronological replay", desc: "Inputs walked oldest -> newest by author date; both AuthorDate and CommitterDate preserved byte-for-byte." },
            { icon: Layers, title: "Multi-source", desc: "Comma-separated inputs, or use all / -N to pull every (or the latest N) versioned siblings." },
            { icon: ShieldCheck, title: "Idempotent", desc: "Dedupe via ShaMap means re-running never replays a commit twice across runs." },
          ].map((f) => (
            <div key={f.title} className="rounded-lg border border-border p-4 bg-card">
              <f.icon className="h-5 w-5 text-primary mb-2" />
              <h3 className="font-semibold text-sm mb-1">{f.title}</h3>
              <p className="text-xs text-muted-foreground">{f.desc}</p>
            </div>
          ))}
        </div>
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">Usage</h2>
        <CodeBlock code={`gitmap commit-in <source> <input1,input2,...> [flags]
gitmap cin       <source> all                    [flags]
gitmap cin       <source> -5                     [flags]`} />
        <p className="text-sm text-muted-foreground mt-3">
          <code>&lt;source&gt;</code> is the TARGET repo (the one receiving appended commits).
          Auto-init is fixed: URL → <code>git clone</code>; existing repo → reuse; existing
          non-repo folder → <code>git init</code> in place; missing path →{" "}
          <code>mkdir -p && git init</code>. No prompt, no flag.
        </p>
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">Flags</h2>
        <div className="overflow-x-auto">
          <table className="w-full text-sm border border-border rounded-lg">
            <thead>
              <tr className="bg-muted/50">
                <th className="text-left px-4 py-2 font-medium">Flag</th>
                <th className="text-left px-4 py-2 font-medium">Default</th>
                <th className="text-left px-4 py-2 font-medium">Description</th>
              </tr>
            </thead>
            <tbody>
              {flags.map((f) => (
                <tr key={f.flag} className="border-t border-border">
                  <td className="px-4 py-2 font-mono text-primary">{f.flag}</td>
                  <td className="px-4 py-2 font-mono text-muted-foreground">{f.def}</td>
                  <td className="px-4 py-2 text-muted-foreground">{f.desc}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">How &lt;source&gt; auto-init works</h2>
        <p className="text-sm text-muted-foreground mb-3">
          You never have to <code>git init</code> first. <code>commit-in</code> resolves
          <code> &lt;source&gt;</code> through a fixed dispatch table — no prompts, no flags,
          no surprises:
        </p>
        <div className="overflow-x-auto">
          <table className="w-full text-sm border border-border rounded-lg">
            <thead>
              <tr className="bg-muted/50">
                <th className="text-left px-4 py-2 font-medium">If &lt;source&gt; is…</th>
                <th className="text-left px-4 py-2 font-medium">commit-in does…</th>
              </tr>
            </thead>
            <tbody>
              {[
                { when: "An https:// or git@ URL", then: "git clone <url> into the derived folder name" },
                { when: "An existing path with .git/", then: "Reuse the repo in place — never re-init" },
                { when: "An existing folder, NO .git/", then: "git init in place (your files are kept untouched)" },
                { when: "A path that does not exist", then: "mkdir -p <path> && git init <path>" },
              ].map((row) => (
                <tr key={row.when} className="border-t border-border">
                  <td className="px-4 py-2 text-muted-foreground">{row.when}</td>
                  <td className="px-4 py-2 font-mono text-xs">{row.then}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">Examples</h2>

        <h3 className="font-semibold text-sm mt-4 mb-2 text-foreground">
          1 · Convert a plain folder of files into a git repo + replay history
        </h3>
        <p className="text-sm text-muted-foreground mb-2">
          You have <code>./my-project/</code> with code but no <code>.git/</code> yet.
          Point <code>commit-in</code> at it and pull history from a URL — the folder is
          auto-<code>git init</code>ed in place, your files stay where they are.
        </p>
        <CodeBlock language="bash" code={`# folder exists, no .git/ yet — commit-in will run \`git init\` for you
gitmap commit-in ./my-project https://github.com/me/my-project-archive.git`} />

        <h3 className="font-semibold text-sm mt-6 mb-2 text-foreground">
          2 · Mix a local folder + a remote URL as INPUTS into one canonical timeline
        </h3>
        <p className="text-sm text-muted-foreground mb-2">
          The first positional is the TARGET. The second is the comma-separated INPUTS to
          walk in author-date order. You can freely mix a local checkout with one or more
          remote URLs — each URL is shallow-cloned into{" "}
          <code>.gitmap/temp/&lt;runId&gt;/</code> and walked just like the local one.
        </p>
        <CodeBlock language="bash" code={`# target = ./canonical (auto-init if missing)
# inputs = local folder + 2 remote forks, walked oldest -> newest
gitmap cin ./canonical \\
    ./old-local-checkout,https://github.com/me/old-fork.git,git@github.com:me/new-fork.git`} />

        <h3 className="font-semibold text-sm mt-6 mb-2 text-foreground">
          3 · Brand-new target folder from scratch (mkdir + init + replay)
        </h3>
        <p className="text-sm text-muted-foreground mb-2">
          Pass a path that does not exist. <code>commit-in</code> creates the folder, runs
          <code> git init</code>, and starts appending — one command, zero setup.
        </p>
        <CodeBlock language="bash" code={`gitmap commit-in ./brand-new-canonical \\
    https://github.com/me/legacy-v1.git,https://github.com/me/legacy-v2.git`} />

        <h3 className="font-semibold text-sm mt-6 mb-2 text-foreground">
          4 · Replay every versioned sibling automatically
        </h3>
        <p className="text-sm text-muted-foreground mb-2">
          The <code>all</code> keyword expands to every <code>&lt;source&gt;-vN</code>{" "}
          sibling on disk. Use <code>-N</code> for the latest N only. Both work great with
          <code> --save-profile</code> so the next run is one word.
        </p>
        <CodeBlock language="bash" code={`# Every sibling, save the resolved settings as the default profile
gitmap commit-in ./canonical all --save-profile Default --set-default

# Just the last 3 siblings, dry-run, with per-language new-function intel
gitmap cin ./canonical -3 --dry-run --function-intel on --languages Go,TypeScript`} />

        <h3 className="font-semibold text-sm mt-6 mb-2 text-foreground">
          5 · Override author + scrub commit messages
        </h3>
        <CodeBlock language="bash" code={`gitmap cin git@github.com:me/canonical.git \\
    https://github.com/me/old-fork.git,https://github.com/me/new-fork.git \\
    --author-name "Jane Doe" --author-email jane@example.com \\
    --message-exclude "StartsWith:Signed-off-by:,Contains:[skip ci]" \\
    --title-suffix " — via gitmap"`} />

        <h3 className="font-semibold text-sm mt-6 mb-2 text-foreground">
          6 · Reuse a saved profile + only rewrite weak titles
        </h3>
        <CodeBlock language="bash" code={`gitmap cin ./canonical all --default \\
    --override-messages "Refine implementation,Improve module" \\
    --override-only-weak`} />

        <h3 className="font-semibold text-sm mt-6 mb-2 text-foreground">
          7 · Headless CI run (fail loudly on any unset value)
        </h3>
        <CodeBlock language="bash" code={`gitmap cin ./canonical all --profile CI --no-prompt`} />
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">Sample profile JSON</h2>
        <p className="text-sm text-muted-foreground mb-3">
          Drop this file at{" "}
          <code>.gitmap/commit-in/profiles/Default.json</code> (relative to your workspace
          root — the nearest ancestor containing <code>.gitmap/</code>) and load it with{" "}
          <code>--profile Default</code> or <code>--default</code>. Keys and enum values are
          <strong> PascalCase</strong>; the loader uses <em>strict</em> decoding, so unknown
          keys are an error. Edit anything you like — every field maps 1:1 to a CLI flag
          above.
        </p>
        <CodeBlock
          language="json"
          title=".gitmap/commit-in/profiles/Default.json"
          code={`{
  "Name": "Default",
  "SchemaVersion": 1,
  "SourceRepoPath": "/abs/path/to/canonical",
  "IsDefault": true,
  "ConflictMode": "ForceMerge",
  "Author": {
    "Name": "Jane Doe",
    "Email": "jane@example.com"
  },
  "Exclusions": [
    { "Kind": "PathFolder", "Value": "node_modules" },
    { "Kind": "PathFolder", "Value": "dist" },
    { "Kind": "PathFile",   "Value": "secrets.env" }
  ],
  "MessageRules": [
    { "Kind": "StartsWith", "Value": "Signed-off-by:" },
    { "Kind": "Contains",   "Value": "[skip ci]" },
    { "Kind": "EndsWith",   "Value": "(cherry picked from commit)" }
  ],
  "MessagePrefix":   ["chore:", "feat:", "fix:"],
  "MessageSuffix":   [],
  "TitlePrefix":     "",
  "TitleSuffix":     " — via gitmap",
  "OverrideMessages": ["Improve module", "Refine implementation"],
  "OverrideOnlyWeak": true,
  "WeakWords":        ["change", "update", "updates", "misc"],
  "FunctionIntel": {
    "IsEnabled": true,
    "Languages": ["Go", "TypeScript", "Python"]
  }
}`}
        />
        <p className="text-xs text-muted-foreground mt-3">
          <strong>Tip:</strong> let gitmap write the file for you the first time —
          <code> gitmap cin ./canonical all --save-profile Default --set-default</code> —
          then open the resulting JSON and tweak. Re-saving requires{" "}
          <code>--save-profile-overwrite</code>. Profiles bind by absolute symlink-resolved
          path, NOT by remote URL, so two clones of the same upstream can carry different
          policies.
        </p>
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">Exit Codes</h2>
        <div className="overflow-x-auto">
          <table className="w-full text-sm border border-border rounded-lg">
            <thead>
              <tr className="bg-muted/50">
                <th className="text-left px-4 py-2 font-medium">Code</th>
                <th className="text-left px-4 py-2 font-medium">Meaning</th>
              </tr>
            </thead>
            <tbody>
              {exitCodes.map((e) => (
                <tr key={e.code} className="border-t border-border">
                  <td className="px-4 py-2 font-mono text-primary">{e.code}</td>
                  <td className="px-4 py-2 text-muted-foreground">{e.meaning}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      <section>
        <h2 className="text-xl font-semibold mb-3">See Also</h2>
        <ul className="list-disc list-inside space-y-1 text-sm">
          <li><a href="/commit-left" className="text-primary hover:underline">commit-left</a> / <a href="/commit-right" className="text-primary hover:underline">commit-right</a> / <a href="/commit-both" className="text-primary hover:underline">commit-both</a></li>
          <li><a href="/merge-left" className="text-primary hover:underline">merge-left</a> / <a href="/merge-right" className="text-primary hover:underline">merge-right</a> / <a href="/merge-both" className="text-primary hover:underline">merge-both</a></li>
        </ul>
      </section>
    </div>
  </DocsLayout>
);

export default CommitInPage;
