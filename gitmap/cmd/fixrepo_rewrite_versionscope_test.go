package cmd

// v5.39.0: regression matrix locking down the bare-base rewrite scope
// across current versions v1..v4. Bare `{base}` substitution must ONLY
// fire on the v1→v2 transition (current==2 with v1 in targets). Every
// other (current, targets) shape must preserve standalone `{base}`.

import "testing"

func TestApplyAllTargets_VersionScopeMatrix(t *testing.T) {
	const base = "gitmap"
	type tc struct {
		name    string
		current int
		targets []int
		in      string
		want    string
	}
	cases := []tc{
		{
			// current=1 is a no-op floor: nothing to bump to.
			name:    "v1_no_rewrite",
			current: 1,
			targets: []int{},
			in:      "gitmap and gitmap-v1 stay put",
			want:    "gitmap and gitmap-v1 stay put",
		},
		{
			// v1→v2: the ONLY case where bare base is rewritten.
			name:    "v2_bare_base_rewritten",
			current: 2,
			targets: []int{1},
			in:      "url=https://github.com/x/gitmap plus gitmap-v1 token",
			want:    "url=https://github.com/x/gitmap-v2 plus gitmap-v2 token",
		},
		{
			// v3: bare base preserved even with v1 in targets.
			name:    "v3_bare_base_preserved",
			current: 3,
			targets: []int{1, 2},
			in:      "gitmap binary and gitmap-v1 and gitmap-v2",
			want:    "gitmap binary and gitmap-v3 and gitmap-v3",
		},
		{
			// v4: bare base preserved across full target sweep.
			name:    "v4_bare_base_preserved",
			current: 4,
			targets: []int{1, 2, 3},
			in:      "gitmap binary, gitmap-v1, gitmap-v2, gitmap-v3",
			want:    "gitmap binary, gitmap-v4, gitmap-v4, gitmap-v4",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, _ := applyAllTargets(c.in, base, c.current, c.targets)
			if got != c.want {
				t.Fatalf("scope mismatch (current=%d).\n got:  %q\n want: %q", c.current, got, c.want)
			}
		})
	}
}
