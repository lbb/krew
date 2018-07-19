// Copyright © 2018 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gitutil

import (
	"fmt"
	"os"

	"github.com/golang/glog"
	"gopkg.in/src-d/go-git.v4"
)

// EnsureCloned will clone into the destination path, otherwise will return no error.
func EnsureCloned(uri, destinationPath string) error {
	co := &git.CloneOptions{URL: uri}
	if glog.V(2) {
		co.Progress = os.Stderr
	}
	if _, err := git.PlainClone(destinationPath, false, co); err == git.ErrRepositoryAlreadyExists {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to clone repo, err: %v", err)
	}

	return nil
}

// IsGitCloned will test if the path is a git dir.
func IsGitCloned(gitPath string) (bool, error) {
	_, err := git.PlainOpen(gitPath)
	if err == git.ErrRepositoryNotExists {
		return false, nil
	}
	return err == nil, err
}

// update will fetch origin and set HEAD to origin/HEAD.
func update(destinationPath string) error {
	g, err := git.PlainOpen(destinationPath)
	if err != nil {
		return fmt.Errorf("failed to open git repo, err: %v", err)
	}

	w, err := g.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get the git worktree, err: %v", err)
	}

	po := &git.PullOptions{RemoteName: "origin"}
	if glog.V(2) {
		po.Progress = os.Stderr
	}

	if err = w.Pull(po); err == git.NoErrAlreadyUpToDate {
		glog.V(2).Infof("Already Up To Date")
	} else if err != nil {
		return fmt.Errorf("failed to fetch the origin, err: %v", err)
	}
	return nil
}

// EnsureUpdated will ensure the destination path exsists and is up to date.
func EnsureUpdated(uri, destinationPath string) error {
	if err := EnsureCloned(uri, destinationPath); err != nil {
		return err
	}
	return update(destinationPath)
}
