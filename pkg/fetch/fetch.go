package fetch

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"os"
	"strings"
)

// Fetcher fetching data from remote
type Fetcher interface {
	Fetch() (*object.Commit, error)
}

type fetcher struct {
	url       string
	directory string
	branch    string
	auth      transport.AuthMethod
}

// NewFetcher creates new Fetcher
func NewFetcher(url string, directory string, branch string, auth transport.AuthMethod) Fetcher {
	fetcher := &fetcher{
		url:       url,
		directory: directory,
		branch:    branch,
		auth:      auth,
	}
	return fetcher
}

// Fetch from remote
func (f *fetcher) Fetch() (*object.Commit, error) {
	var r *git.Repository
	var err error

	if r, err = git.PlainOpen(f.directory); err != nil {
		log.Infof("Repository not found in '%s' cloning...", f.directory)
		r, err = git.PlainClone(f.directory, false, &git.CloneOptions{
			URL:           f.url,
			Auth:          f.auth,
			Depth:         1,
			SingleBranch:  true,
			ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", f.branch)),
		})
		if err != nil {
			return nil, err
		}
	} else {
		log.Infof("Repository found in '%s' opening...", f.directory)
	}

	if branch, err := r.Branch(f.branch); branch == nil || err != nil {
		log.Infof("Branch switched to '%s'", f.branch)
		os.RemoveAll(f.directory)
		r, err = git.PlainClone(f.directory, false, &git.CloneOptions{
			URL:           f.url,
			Auth:          f.auth,
			Depth:         1,
			SingleBranch:  true,
			ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", f.branch)),
		})
		if err != nil {
			return nil, err
		}
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	log.Println("Pulling changes")
	err = w.Pull(&git.PullOptions{
		Auth:          f.auth,
		RemoteName:    "origin",
		Force:         true,
		SingleBranch:  true,
		Depth:         1,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", f.branch)),
	})
	if err != nil && err.Error() != "already up-to-date" {
		return nil, err
	}

	ref, err := r.Head()
	if err != nil {
		return nil, err
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	log.Infof("Pulled ref '%s'", ref.Hash())

	return commit, nil
}

// NewAuth creates new AuthMethod based on URI
func NewAuth(git string) (transport.AuthMethod, error) {
	var auth transport.AuthMethod

	ep, err := transport.NewEndpoint(git)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(ep.Protocol, "http") && ep.User != "" && ep.Password != "" {
		auth = &http.BasicAuth{
			Username: ep.User,
			Password: ep.Password,
		}
	}

	return auth, nil
}