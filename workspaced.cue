package workspaced

// Real upstream trees for assimilate demos. Lock pins the commit.
// Place under testdata/projects/ (gitignored):
//   mise run projects:lock && mise run projects:sync

#project: {}

#project: [string]: {
	from:    string
	version: string | *"HEAD"
	origin:  string | *"."
	dest:    string
}

#project: csmith: {
	from:    "github:csmith-project/csmith"
	version: "csmith-2.3.0"
	dest:    "testdata/projects/csmith"
}

#project: rhai: {
	from:    "github:rhaiscript/rhai"
	version: "v1.25.0"
	dest:    "testdata/projects/rhai"
}

workspaced: {
	inputs: {
		self: {from: "self"}
		for name, p in #project {
			(name): {
				from:    p.from
				version: p.version
			}
		}
	}
	modules: {
		for name, p in #project {
			(name): {
				from: "core:place"
				config: {
					items: {
						"\(p.dest)": "\(name):\(p.origin)"
					}
				}
			}
		}
	}
}
