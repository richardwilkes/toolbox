// Package git provides simple git repository access.
package git

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
)

var (
	tagsRefListRegex   = regexp.MustCompile(`(?m-s)(?:tags)/(\S+)$`)
	originRefListRegex = regexp.MustCompile(`(?m-s)(?:origin)/(\S+)$`)
)

// Repo provides access to a git repository.
type Repo struct {
	remote string
	local  string
}

// NewRepo creates a new git repository access object.
func NewRepo(remote, local string) (*Repo, error) {
	if _, err := exec.LookPath("git"); err != nil {
		return nil, errs.New(i18n.Text("git is not installed"))
	}
	repo := &Repo{
		remote: remote,
		local:  local,
	}
	if repo.CheckLocal() {
		out, err := repo.runFromDir("git", "config", "--get", "remote.origin.url")
		if err != nil {
			return nil, errs.NewWithCause(i18n.Text("Unable to retrieve local repository information"), err)
		}
		localRemote := strings.TrimSpace(string(out))
		if remote != "" && localRemote != remote {
			return nil, errs.Newf(i18n.Text("Existing remote (%) does not match requested remote (%s)"), localRemote, remote)
		}
		if remote == "" && localRemote != "" {
			repo.remote = localRemote
		}
	}
	return repo, nil
}

// CheckLocal verifies the local location is a Git repo.
func (repo *Repo) CheckLocal() bool {
	_, err := os.Stat(repo.local + "/.git")
	return err == nil
}

// Init initializes a git repository at the local location.
func (repo *Repo) Init() error {
	if _, err := exec.Command("git", "init", repo.local).CombinedOutput(); err != nil { //nolint:gosec
		return errs.NewWithCause(i18n.Text("Unable to initialize repository"), err)
	}
	return nil
}

// Clone a repository.
func (repo *Repo) Clone() error {
	if _, err := exec.Command("git", "clone", repo.remote, repo.local).CombinedOutput(); err != nil { //nolint:gosec
		return errs.NewWithCause(i18n.Text("Unable to clone repository"), err)
	}
	return nil
}

// Checkout a revision, branch or tag.
func (repo *Repo) Checkout(revisionBranchOrTag string) error {
	if _, err := repo.runFromDir("git", "checkout", revisionBranchOrTag); err != nil {
		return errs.NewWithCausef(err, i18n.Text("Unable to check out '%s'"), revisionBranchOrTag)
	}
	return nil
}

// Fetch a repository.
func (repo *Repo) Fetch() error {
	if _, err := repo.runFromDir("git", "fetch", "--tags"); err != nil {
		return errs.NewWithCause(i18n.Text("Unable to fetch"), err)
	}
	return nil
}

// Pull a repository.
func (repo *Repo) Pull() error {
	if _, err := repo.runFromDir("git", "pull"); err != nil {
		return errs.NewWithCause(i18n.Text("Unable to pull"), err)
	}
	return nil
}

// HasDetachedHead returns true if the repo is currently in a "detached head"
// state.
func (repo *Repo) HasDetachedHead() bool {
	contents, err := ioutil.ReadFile(filepath.Join(repo.local, ".git", "HEAD"))
	return err != nil && !bytes.HasPrefix(bytes.TrimSpace(contents), []byte("ref: "))
}

// Date retrieves the date on the latest commit.
func (repo *Repo) Date() (time.Time, error) {
	out, err := repo.runFromDir("git", "log", "-1", "--date=iso", "--pretty=format:%cd")
	if err != nil {
		return time.Time{}, errs.NewWithCause(i18n.Text("Unable to retrieve revision date"), err)
	}
	t, err := time.Parse("2006-01-02 15:04:05 -0700", string(out))
	if err != nil {
		return time.Time{}, errs.NewWithCause(i18n.Text("Unable to retrieve revision date"), err)
	}
	return t, nil
}

// Branches returns a list of available branches.
func (repo *Repo) Branches() ([]string, error) {
	out, err := repo.runFromDir("git", "show-ref")
	if err != nil {
		return []string{}, errs.NewWithCause(i18n.Text("Unable to retrieve branches"), err)
	}
	return repo.referenceList(string(out), originRefListRegex), nil
}

// Revision retrieves the current revision.
func (repo *Repo) Revision() (string, error) {
	out, err := repo.runFromDir("git", "rev-parse", "HEAD")
	if err != nil {
		return "", errs.NewWithCause(i18n.Text("Unable to retrieve checked out revision"), err)
	}
	return strings.TrimSpace(string(out)), nil
}

// Current returns the current branch/tag/revision.
// * Branch name if on the tip of the branch
// * Tag if on a tag
// * Otherwise a revision id
func (repo *Repo) Current() (string, error) {
	if out, err := repo.runFromDir("git", "symbolic-ref", "HEAD"); err == nil {
		return string(bytes.TrimSpace(bytes.TrimPrefix(out, []byte("refs/heads/")))), nil
	}
	rev, err := repo.Revision()
	if err != nil {
		return "", err
	}
	tags, err := repo.TagsFromCommit(rev)
	if err != nil {
		return "", err
	}
	if len(tags) > 0 {
		return tags[0], nil
	}
	return rev, nil
}

// Tags returns a list of available tags.
func (repo *Repo) Tags() ([]string, error) {
	out, err := repo.runFromDir("git", "show-ref")
	if err != nil {
		return []string{}, errs.NewWithCause(i18n.Text("Unable to retrieve tags"), err)
	}
	return repo.referenceList(string(out), tagsRefListRegex), nil
}

// TagsFromCommit retrieves the tags from a revision.
func (repo *Repo) TagsFromCommit(rev string) ([]string, error) {
	out, err := repo.runFromDir("git", "show-ref", "-d")
	if err != nil {
		return nil, errs.NewWithCause(i18n.Text("Unable to retrieve tags"), err)
	}
	lines := strings.Split(string(out), "\n")
	list := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, rev) {
			list = append(list, line)
		}
	}
	tags := repo.referenceList(strings.Join(list, "\n"), tagsRefListRegex)
	result := make([]string, 0, len(tags))
	for _, t := range tags {
		result = append(result, strings.TrimSuffix(t, "^{}"))
	}
	return result, nil
}

func (repo *Repo) referenceList(content string, re *regexp.Regexp) []string {
	var out []string //nolint:prealloc
	for _, m := range re.FindAllStringSubmatch(content, -1) {
		out = append(out, m[1])
	}
	return out
}

// HasChanges returns true if changes are present.
func (repo *Repo) HasChanges() bool {
	out, err := repo.runFromDir("git", "status", "--porcelain")
	return err != nil || len(out) != 0
}

func (repo *Repo) runFromDir(cmd string, args ...string) ([]byte, error) {
	c := exec.Command(cmd, args...)
	c.Dir = repo.local
	c.Env = mergeEnvLists([]string{"PWD=" + c.Dir}, os.Environ())
	return c.CombinedOutput()
}

func mergeEnvLists(in, out []string) []string {
NextVar:
	for _, inkv := range in {
		k := strings.SplitAfterN(inkv, "=", 2)[0] + "="
		for i, outkv := range out {
			if strings.HasPrefix(outkv, k) {
				out[i] = inkv
				continue NextVar
			}
		}
		out = append(out, inkv)
	}
	return out
}
