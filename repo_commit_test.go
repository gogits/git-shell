package git

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_escapePath(t *testing.T) {
	tests := []struct {
		path    string
		expPath string
	}{
		{
			path:    "",
			expPath: "",
		},
		{
			path:    "normal",
			expPath: "normal",
		},
		{
			path:    ":normal",
			expPath: "\\:normal",
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			assert.Equal(t, test.expPath, escapePath(test.path))
		})
	}
}

func TestRepository_Log(t *testing.T) {
	tests := []struct {
		rev          string
		opt          LogOptions
		expCommitIDs []string
	}{
		{
			rev: "0eedd79eba4394bbef888c804e899731644367fe",
			opt: LogOptions{
				Since: time.Unix(1581250680, 0),
			},
			expCommitIDs: []string{
				"0eedd79eba4394bbef888c804e899731644367fe",
				"4e59b72440188e7c2578299fc28ea425fbe9aece",
			},
		},
		{
			rev: "0eedd79eba4394bbef888c804e899731644367fe",
			opt: LogOptions{
				Since: time.Now().AddDate(100, 0, 0),
			},
			expCommitIDs: []string{},
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			commits, err := testrepo.Log(test.rev, test.opt)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, test.expCommitIDs, commitsToIDs(commits))
		})
	}
}

func TestRepository_CommitByRevision(t *testing.T) {
	tests := []struct {
		rev    string
		opt    CommitByRevisionOptions
		expID  string
		expErr error
	}{
		{
			rev:    "4e59b72",
			expID:  "4e59b72440188e7c2578299fc28ea425fbe9aece",
			expErr: nil,
		},
		{
			rev:    "404",
			expID:  "",
			expErr: ErrRevisionNotExist,
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			c, err := testrepo.CommitByRevision(test.rev, test.opt)
			assert.Equal(t, test.expErr, err)
			if c != nil {
				assert.Equal(t, test.expID, c.ID().String())
			}
		})
	}
}

func TestRepository_CommitsSince(t *testing.T) {
	tests := []struct {
		rev          string
		since        time.Time
		opt          CommitsSinceOptions
		expCommitIDs []string
	}{
		{
			rev:   "0eedd79eba4394bbef888c804e899731644367fe",
			since: time.Unix(1581250680, 0),
			expCommitIDs: []string{
				"0eedd79eba4394bbef888c804e899731644367fe",
				"4e59b72440188e7c2578299fc28ea425fbe9aece",
			},
		},
		{
			rev:          "0eedd79eba4394bbef888c804e899731644367fe",
			since:        time.Now().AddDate(100, 0, 0),
			expCommitIDs: []string{},
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			commits, err := testrepo.CommitsSince(test.rev, test.since, test.opt)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, test.expCommitIDs, commitsToIDs(commits))
		})
	}
}

func TestRepository_DiffNameOnly(t *testing.T) {
	tests := []struct {
		base     string
		head     string
		opt      DiffNameOnlyOptions
		expFiles []string
	}{
		{
			base:     "ef7bebf8bdb1919d947afe46ab4b2fb4278039b3",
			head:     "978fb7f6388b49b532fbef8b856681cfa6fcaa0a",
			expFiles: []string{"fix.txt"},
		},
		{
			base: "45a30ea9afa413e226ca8614179c011d545ca883",
			head: "978fb7f6388b49b532fbef8b856681cfa6fcaa0a",
			opt: DiffNameOnlyOptions{
				NeedsMergeBase: true,
			},
			expFiles: []string{"fix.txt", "pom.xml", "src/test/java/com/github/AppTest.java"},
		},

		{
			base: "45a30ea9afa413e226ca8614179c011d545ca883",
			head: "978fb7f6388b49b532fbef8b856681cfa6fcaa0a",
			opt: DiffNameOnlyOptions{
				Path: "src",
			},
			expFiles: []string{"src/test/java/com/github/AppTest.java"},
		},
		{
			base: "45a30ea9afa413e226ca8614179c011d545ca883",
			head: "978fb7f6388b49b532fbef8b856681cfa6fcaa0a",
			opt: DiffNameOnlyOptions{
				Path: "resources",
			},
			expFiles: []string{},
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			files, err := testrepo.DiffNameOnly(test.base, test.head, test.opt)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, test.expFiles, files)
		})
	}
}

func TestRepository_RevListCount(t *testing.T) {
	tests := []struct {
		refspecs []string
		opt      RevListCountOptions
		expCount int64
		expErr   error
	}{
		{
			refspecs: []string{"755fd577edcfd9209d0ac072eed3b022cbe4d39b"},
			expCount: 1,
		},
		{
			refspecs: []string{"f5ed01959cffa4758ca0a49bf4c34b138d7eab0a"},
			expCount: 5,
		},
		{
			refspecs: []string{"978fb7f6388b49b532fbef8b856681cfa6fcaa0a"},
			expCount: 27,
		},

		{
			refspecs: []string{"7c5ee6478d137417ae602140c615e33aed91887c"},
			opt: RevListCountOptions{
				Path: "README.txt",
			},
			expCount: 3,
		},
		{
			refspecs: []string{"7c5ee6478d137417ae602140c615e33aed91887c"},
			opt: RevListCountOptions{
				Path: "resources",
			},
			expCount: 1,
		},

		{
			refspecs: []string{},
			expCount: 0,
			expErr:   errors.New("must have at least one refspec"),
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			count, err := testrepo.RevListCount(test.refspecs, test.opt)
			assert.Equal(t, test.expErr, err)
			assert.Equal(t, test.expCount, count)
		})
	}
}

func TestRepository_RevList(t *testing.T) {
	tests := []struct {
		refspecs     []string
		opt          RevListOptions
		expCommitIDs []string
		expErr       error
	}{
		{
			refspecs: []string{"45a30ea9afa413e226ca8614179c011d545ca883...978fb7f6388b49b532fbef8b856681cfa6fcaa0a"},
			expCommitIDs: []string{
				"978fb7f6388b49b532fbef8b856681cfa6fcaa0a",
				"ef7bebf8bdb1919d947afe46ab4b2fb4278039b3",
				"ebbbf773431ba07510251bb03f9525c7bab2b13a",
			},
		},
		{
			refspecs: []string{"45a30ea9afa413e226ca8614179c011d545ca883...978fb7f6388b49b532fbef8b856681cfa6fcaa0a"},
			opt: RevListOptions{
				Path: "src",
			},
			expCommitIDs: []string{
				"ebbbf773431ba07510251bb03f9525c7bab2b13a",
			},
		},

		{
			refspecs:     []string{},
			expCommitIDs: []string{},
			expErr:       errors.New("must have at least one refspec"),
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			commits, err := testrepo.RevList(test.refspecs, test.opt)
			assert.Equal(t, test.expErr, err)
			assert.Equal(t, test.expCommitIDs, commitsToIDs(commits))
		})
	}
}

func TestRepository_LatestCommitTime(t *testing.T) {
	tests := []struct {
		opt     LatestCommitTimeOptions
		expTime time.Time
	}{
		{
			opt: LatestCommitTimeOptions{
				Branch: "release-1.0",
			},
			expTime: time.Unix(1581256638, 0),
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			got, err := testrepo.LatestCommitTime(test.opt)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, test.expTime.Unix(), got.Unix())
		})
	}
}
