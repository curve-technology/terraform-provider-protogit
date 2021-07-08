package store

import (
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
)

const MainBranch = "master"

type GitRepo struct {
	repo         *git.Repository
	protosAbsLoc string
}

func NewGitRepo(url string, tagVersion string, rootGitFolderPath string, protosLocation string) (*GitRepo, error) {
	protoAbsoluteLocation := path.Join(rootGitFolderPath, protosLocation)
	gitRepo := &GitRepo{protosAbsLoc: protoAbsoluteLocation}

	if _, err := os.Stat(rootGitFolderPath); os.IsNotExist(err) {
		gitRepo.repo, err = git.PlainClone(rootGitFolderPath, false, &git.CloneOptions{
			URL: url,
			// Progress: os.Stdout,
		})
		if err != nil {
			return nil, err
		}
	} else {
		gitRepo.repo, err = git.PlainOpen(rootGitFolderPath)
		if err != nil {
			return nil, err
		}
	}

	worktree, err := gitRepo.repo.Worktree()
	if err != nil {
		return nil, err
	}

	refName := plumbing.NewBranchReferenceName(MainBranch)

	err = worktree.Checkout(&git.CheckoutOptions{Branch: refName, Force: true})
	if err != nil {
		return nil, err
	}

	err = worktree.Pull(&git.PullOptions{ReferenceName: refName, Force: true})
	if err != git.NoErrAlreadyUpToDate && err != nil {
		return nil, err
	}

	tagRefName := plumbing.NewTagReferenceName(tagVersion)

	err = worktree.Checkout(&git.CheckoutOptions{Branch: tagRefName, Force: true})
	if err != nil {
		return nil, err
	}

	return gitRepo, nil
}

// Function that returns proto content given a filepath.
// filepath is the relative path from protosLocation specified in the constructor.
func (repo *GitRepo) GetFileDescriptor(filepath string) (*desc.FileDescriptor, error) {
	parser := protoparse.Parser{ImportPaths: []string{repo.protosAbsLoc}}

	descriptors, err := parser.ParseFiles(filepath)
	if err != nil {
		return nil, err
	}

	fileDescriptor := descriptors[0]
	return fileDescriptor, nil
}
