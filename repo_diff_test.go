// Copyright 2020 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_Diff(t *testing.T) {
	tests := []struct {
		rev          string
		maxFiles     int
		maxFileLines int
		maxLineChars int
		opt          DiffOptions
		expDiff      *Diff
	}{
		{
			rev: "978fb7f6388b49b532fbef8b856681cfa6fcaa0a",
			expDiff: &Diff{
				files: []*DiffFile{
					{
						name:         "fix.txt",
						typ:          DiffFileDelete,
						index:        "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
						sections:     nil,
						numAdditions: 0,
						numDeletions: 0,
						oldName:      "",
						isBinary:     false,
						isSubmodule:  false,
						isIncomplete: false,
					},
				},
				totalAdditions: 0,
				totalDeletions: 0,
				isIncomplete:   false,
			},
		},
		{
			rev: "755fd577edcfd9209d0ac072eed3b022cbe4d39b",
			expDiff: &Diff{
				files: []*DiffFile{
					{
						name:  "README.txt",
						typ:   DiffFileAdd,
						index: "1e24b564bf2298965d8037af42d3ae15ad7d225a",
						sections: []*DiffSection{
							{
								lines: []*DiffLine{
									{
										typ:     DiffLineSection,
										content: "@@ -0,0 +1,11 @@",
									},
									{
										typ:       DiffLineAdd,
										content:   "+This is a sample project students can use during Matthew's Git class.",
										leftLine:  0,
										rightLine: 1,
									},
									{
										typ:       DiffLineAdd,
										content:   "+",
										leftLine:  0,
										rightLine: 2,
									},
									{
										typ:       DiffLineAdd,
										content:   "+We can have a bit of fun with this repo, knowing that we can always reset it to a known good state.  We can apply labels, and branch, then add new code and merge it in to the master branch.",
										leftLine:  0,
										rightLine: 3,
									},
									{
										typ:       DiffLineAdd,
										content:   "+",
										leftLine:  0,
										rightLine: 4,
									},
									{
										typ:       DiffLineAdd,
										content:   "+As a quick reminder, this came from one of three locations in either SSH, Git, or HTTPS format:",
										leftLine:  0,
										rightLine: 5,
									},
									{
										typ:       DiffLineAdd,
										content:   "+",
										leftLine:  0,
										rightLine: 6,
									},
									{
										typ:       DiffLineAdd,
										content:   "+* git@github.com:matthewmccullough/hellogitworld.git",
										leftLine:  0,
										rightLine: 7,
									},
									{
										typ:       DiffLineAdd,
										content:   "+* git://github.com/matthewmccullough/hellogitworld.git",
										leftLine:  0,
										rightLine: 8,
									},
									{
										typ:       DiffLineAdd,
										content:   "+* https://matthewmccullough@github.com/matthewmccullough/hellogitworld.git",
										leftLine:  0,
										rightLine: 9,
									},
									{
										typ:       DiffLineAdd,
										content:   "+",
										leftLine:  0,
										rightLine: 10,
									},
									{
										typ:       DiffLineAdd,
										content:   "+We can, as an example effort, even modify this README and change it as if it were source code for the purposes of the class.",
										leftLine:  0,
										rightLine: 11,
									},
								},
								numAdditions: 11,
							},
						},
						numAdditions: 11,
						numDeletions: 0,
						oldName:      "",
						isBinary:     false,
						isSubmodule:  false,
						isIncomplete: false,
					},
					{
						name:  "resources/labels.properties",
						typ:   DiffFileAdd,
						index: "fbdcfef007c0c09061199e687087b18c3cf8e083",
						sections: []*DiffSection{
							{
								lines: []*DiffLine{
									{
										typ:     DiffLineSection,
										content: "@@ -0,0 +1,4 @@",
									},
									{
										typ:       DiffLineAdd,
										content:   "+app.title=Our App",
										leftLine:  0,
										rightLine: 1,
									},
									{
										typ:       DiffLineAdd,
										content:   "+app.welcome=Welcome to the application",
										leftLine:  0,
										rightLine: 2,
									},
									{
										typ:       DiffLineAdd,
										content:   "+app.goodbye=We hope you enjoyed using our application",
										leftLine:  0,
										rightLine: 3,
									},
									{
										typ:       DiffLineAdd,
										content:   "+cli.usage=This application doesn't use a command line interface",
										leftLine:  0,
										rightLine: 4,
									},
								},
								numAdditions: 4,
							},
						},
						numAdditions: 4,
						numDeletions: 0,
						oldName:      "",
						isBinary:     false,
						isSubmodule:  false,
						isIncomplete: false,
					},
					{
						name:  "src/Main.groovy",
						typ:   DiffFileAdd,
						index: "51680791956b43effdb2f16bccd2b4752d66078f",
						sections: []*DiffSection{
							{
								lines: []*DiffLine{
									{
										typ:     DiffLineSection,
										content: "@@ -0,0 +1,6 @@",
									},
									{
										typ:       DiffLineAdd,
										content:   "+def name = \"Matthew\"",
										leftLine:  0,
										rightLine: 1,
									},
									{
										typ:       DiffLineAdd,
										content:   "+",
										leftLine:  0,
										rightLine: 2,
									},
									{
										typ:       DiffLineAdd,
										content:   "+println \"Hello ${name}\"",
										leftLine:  0,
										rightLine: 3,
									},
									{
										typ:       DiffLineAdd,
										content:   "+",
										leftLine:  0,
										rightLine: 4,
									},
									{
										typ:       DiffLineAdd,
										content:   "+int programmingPoints = 10",
										leftLine:  0,
										rightLine: 5,
									},
									{
										typ:       DiffLineAdd,
										content:   "+println \"${name} has at least ${programmingPoints} programming points.\"",
										leftLine:  0,
										rightLine: 6,
									},
								},
								numAdditions: 6,
							},
						},
						numAdditions: 6,
						numDeletions: 0,
						oldName:      "",
						isBinary:     false,
						isSubmodule:  false,
						isIncomplete: false,
					},
				},
				totalAdditions: 21,
				totalDeletions: 0,
				isIncomplete:   false,
			},
		},
		{
			rev: "978fb7f6388b49b532fbef8b856681cfa6fcaa0a",
			opt: DiffOptions{
				Base: "ef7bebf8bdb1919d947afe46ab4b2fb4278039b3",
			},
			expDiff: &Diff{
				files: []*DiffFile{
					{
						name:         "fix.txt",
						typ:          DiffFileDelete,
						index:        "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
						sections:     nil,
						numAdditions: 0,
						numDeletions: 0,
						oldName:      "",
						isBinary:     false,
						isSubmodule:  false,
						isIncomplete: false,
					},
				},
				totalAdditions: 0,
				totalDeletions: 0,
				isIncomplete:   false,
			},
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			diff, err := testrepo.Diff(test.rev, test.maxFiles, test.maxFileLines, test.maxLineChars, test.opt)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, test.expDiff, diff)
		})
	}
}

func TestRepository_RawDiff(t *testing.T) {
	t.Run("invalid revision", func(t *testing.T) {
		err := testrepo.RawDiff("bad_revision", "bad_diff_type", nil)
		assert.Equal(t, ErrRevisionNotExist, err)
	})

	t.Run("invalid diffType", func(t *testing.T) {
		err := testrepo.RawDiff("978fb7f6388b49b532fbef8b856681cfa6fcaa0a", "bad_diff_type", nil)
		assert.Equal(t, errors.New("invalid diffType: bad_diff_type"), err)
	})

	tests := []struct {
		rev       string
		diffType  RawDiffFormat
		opt       RawDiffOptions
		expOutput string
	}{
		{
			rev:      "978fb7f6388b49b532fbef8b856681cfa6fcaa0a",
			diffType: RawDiffNormal,
			expOutput: `diff --git a/fix.txt b/fix.txt
deleted file mode 100644
index e69de29..0000000
`,
		},
		{
			rev:      "978fb7f6388b49b532fbef8b856681cfa6fcaa0a",
			diffType: RawDiffPatch,
			expOutput: `Date: Sun, 9 Feb 2020 17:22:24 +0800
Subject: [PATCH] Delete fix.txt

---
 fix.txt | 0
 1 file changed, 0 insertions(+), 0 deletions(-)
 delete mode 100644 fix.txt

diff --git a/fix.txt b/fix.txt
deleted file mode 100644
index e69de29..0000000
`,
		},
		{
			rev:      "755fd577edcfd9209d0ac072eed3b022cbe4d39b",
			diffType: RawDiffNormal,
			expOutput: `commit 755fd577edcfd9209d0ac072eed3b022cbe4d39b
Author: Matthew McCullough <matthewm@ambientideas.com>
Date:   Mon Nov 24 21:22:01 2008 -0700

    Addition of the README and basic Groovy source samples.
    
    - Addition of the README.txt file explaining what this repository is all about.
    - Addition of Groovy sample source.
    - Addition of sample resource Properties file.

diff --git a/README.txt b/README.txt
new file mode 100644
index 0000000..1e24b56
--- /dev/null
+++ b/README.txt
@@ -0,0 +1,11 @@
+This is a sample project students can use during Matthew's Git class.
+
+We can have a bit of fun with this repo, knowing that we can always reset it to a known good state.  We can apply labels, and branch, then add new code and merge it in to the master branch.
+
+As a quick reminder, this came from one of three locations in either SSH, Git, or HTTPS format:
+
+* git@github.com:matthewmccullough/hellogitworld.git
+* git://github.com/matthewmccullough/hellogitworld.git
+* https://matthewmccullough@github.com/matthewmccullough/hellogitworld.git
+
+We can, as an example effort, even modify this README and change it as if it were source code for the purposes of the class.
\ No newline at end of file
diff --git a/resources/labels.properties b/resources/labels.properties
new file mode 100644
index 0000000..fbdcfef
--- /dev/null
+++ b/resources/labels.properties
@@ -0,0 +1,4 @@
+app.title=Our App
+app.welcome=Welcome to the application
+app.goodbye=We hope you enjoyed using our application
+cli.usage=This application doesn't use a command line interface
diff --git a/src/Main.groovy b/src/Main.groovy
new file mode 100644
index 0000000..5168079
--- /dev/null
+++ b/src/Main.groovy
@@ -0,0 +1,6 @@
+def name = "Matthew"
+
+println "Hello ${name}"
+
+int programmingPoints = 10
+println "${name} has at least ${programmingPoints} programming points."
\ No newline at end of file
`,
		},
		{
			rev:      "755fd577edcfd9209d0ac072eed3b022cbe4d39b",
			diffType: RawDiffPatch,
			expOutput: `Date: Mon, 24 Nov 2008 21:22:01 -0700
Subject: [PATCH] Addition of the README and basic Groovy source samples.

- Addition of the README.txt file explaining what this repository is all about.
- Addition of Groovy sample source.
- Addition of sample resource Properties file.
---
 README.txt                  | 11 +++++++++++
 resources/labels.properties |  4 ++++
 src/Main.groovy             |  6 ++++++
 3 files changed, 21 insertions(+)
 create mode 100644 README.txt
 create mode 100644 resources/labels.properties
 create mode 100644 src/Main.groovy

diff --git a/README.txt b/README.txt
new file mode 100644
index 0000000..1e24b56
--- /dev/null
+++ b/README.txt
@@ -0,0 +1,11 @@
+This is a sample project students can use during Matthew's Git class.
+
+We can have a bit of fun with this repo, knowing that we can always reset it to a known good state.  We can apply labels, and branch, then add new code and merge it in to the master branch.
+
+As a quick reminder, this came from one of three locations in either SSH, Git, or HTTPS format:
+
+* git@github.com:matthewmccullough/hellogitworld.git
+* git://github.com/matthewmccullough/hellogitworld.git
+* https://matthewmccullough@github.com/matthewmccullough/hellogitworld.git
+
+We can, as an example effort, even modify this README and change it as if it were source code for the purposes of the class.
\ No newline at end of file
diff --git a/resources/labels.properties b/resources/labels.properties
new file mode 100644
index 0000000..fbdcfef
--- /dev/null
+++ b/resources/labels.properties
@@ -0,0 +1,4 @@
+app.title=Our App
+app.welcome=Welcome to the application
+app.goodbye=We hope you enjoyed using our application
+cli.usage=This application doesn't use a command line interface
diff --git a/src/Main.groovy b/src/Main.groovy
new file mode 100644
index 0000000..5168079
--- /dev/null
+++ b/src/Main.groovy
@@ -0,0 +1,6 @@
+def name = "Matthew"
+
+println "Hello ${name}"
+
+int programmingPoints = 10
+println "${name} has at least ${programmingPoints} programming points."
\ No newline at end of file
`,
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := testrepo.RawDiff(test.rev, test.diffType, buf, test.opt)
			if err != nil {
				t.Fatal(err)
			}
			output := buf.String()

			// Only check the content after "Date:" line, which is deterministic.
			i := strings.Index(output, "Date:")
			if i > 0 && test.diffType == RawDiffPatch {
				output = output[i:]
			}

			assert.Equal(t, test.expOutput, output)
		})
	}
}

func TestRepository_DiffBinary(t *testing.T) {
	tests := []struct {
		base      string
		head      string
		opt       DiffBinaryOptions
		expOutput string
	}{
		{
			base: "4eaa8d4b05e731e950e2eaf9e8b92f522303ab41",
			head: "4e59b72440188e7c2578299fc28ea425fbe9aece",
			expOutput: `diff --git a/.gitmodules b/.gitmodules
new file mode 100644
index 0000000..6abde17
--- /dev/null
+++ b/.gitmodules
@@ -0,0 +1,3 @@
+[submodule "gogs/docs-api"]
+	path = gogs/docs-api
+	url = https://github.com/gogs/docs-api.git
diff --git a/gogs/docs-api b/gogs/docs-api
new file mode 160000
index 0000000..6b08f76
--- /dev/null
+++ b/gogs/docs-api
@@ -0,0 +1 @@
+Subproject commit 6b08f76a5313fa3d26859515b30aa17a5faa2807
`,
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			p, err := testrepo.DiffBinary(test.base, test.head, test.opt)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, test.expOutput, string(p))
		})
	}
}
