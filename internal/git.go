package internal

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitRepo interface {
	Checkout(branch string) error
	// Fetch eg: git fetch
	Fetch() error
	// ResolveRevision eg: git rev-parse origin/master
	ResolveRevision(rev string) (commitID string, err error)
	// HardReset eg: git reset --hard {commitID}
	HardReset(commitID string) error
}

type gitRepoImpl struct {
	repo *git.Repository
}

func CloneRepo(destPath string, srcURL string) (GitRepo, error) {
	option := &git.CloneOptions{
		URL:               srcURL,
		Progress:          nil,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}
	r, err := git.PlainClone(destPath, false, option)
	if err != nil {
		return nil, err
	}

	return &gitRepoImpl{r}, nil
}

func (g gitRepoImpl) Checkout(branchName string) error {
	// refs/heads/<localBranchName>
	localBranchReferenceName := plumbing.NewBranchReferenceName(branchName)
	// refs/remotes/origin/<remoteBranchName>
	remoteReferenceName := plumbing.NewRemoteReferenceName("origin", branchName)

	if err := g.repo.CreateBranch(&config.Branch{Name: branchName, Remote: "origin", Merge: localBranchReferenceName}); err != nil {
		return err
	}
	newReference := plumbing.NewSymbolicReference(localBranchReferenceName, remoteReferenceName)
	if err := g.repo.Storer.SetReference(newReference); err != nil {
		return err
	}
	return nil
}

func (g gitRepoImpl) Fetch() error {
	o := &git.FetchOptions{
		RemoteName: "origin",
	}
	err := g.repo.Fetch(o)
	if err != nil {
		return err
	}

	return nil
}

func (g gitRepoImpl) ResolveRevision(rev string) (string, error) {
	v, err := g.repo.ResolveRevision(plumbing.Revision(rev))
	if err != nil {
		return "", err
	}

	return v.String(), nil
}

func (g gitRepoImpl) HardReset(commitID string) error {
	w, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	o := &git.ResetOptions{
		Commit: plumbing.NewHash(commitID),
		Mode:   git.HardReset,
	}
	if err := w.Reset(o); err != nil {
		return err
	}

	return nil
}
