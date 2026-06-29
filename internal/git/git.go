package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Client struct {
	repoPath string
}

func New(repoPath string) *Client {
	return &Client{repoPath: repoPath}
}

func (c *Client) IsRepo() bool {
	err := exec.Command("git", "-C", c.repoPath, "rev-parse", "--git-dir").Run()
	return err == nil
}

func (c *Client) run(args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", c.repoPath}, args...)...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, stderr.String())
	}
	return strings.TrimRight(stdout.String(), "\n"), nil
}

func (c *Client) CurrentBranch() (string, error) {
	return c.run("rev-parse", "--abbrev-ref", "HEAD")
}

func (c *Client) LatestCommitHash() (string, error) {
	return c.run("rev-parse", "HEAD")
}

func (c *Client) LatestCommitMessage() (string, error) {
	return c.run("log", "-1", "--pretty=format:%s")
}

func (c *Client) Diff() (string, error) {
	return c.run("diff", "--no-color")
}

func (c *Client) DiffStaged() (string, error) {
	return c.run("diff", "--cached", "--no-color")
}

func (c *Client) DiffWithHEAD() (string, error) {
	return c.run("diff", "HEAD", "--no-color")
}

func (c *Client) DiffBetween(from, to string) (string, error) {
	return c.run("diff", "--no-color", from, to)
}

func (c *Client) ChangedFiles() ([]string, error) {
	output, err := c.run("diff", "--name-only", "--no-color")
	if err != nil {
		return nil, err
	}
	if output == "" {
		return nil, nil
	}
	return strings.Split(output, "\n"), nil
}

func (c *Client) ChangedFilesWithStatus() (map[string]string, error) {
	output, err := c.run("status", "--porcelain")
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if len(line) < 3 {
			continue
		}
		status := strings.TrimSpace(line[:2])
		file := strings.TrimSpace(line[2:])
		result[file] = status
	}
	return result, nil
}

func (c *Client) RecentCommits(n int) ([]Commit, error) {
	output, err := c.run("log", "--oneline", fmt.Sprintf("-%d", n), "--no-color")
	if err != nil {
		return nil, err
	}

	if output == "" {
		return nil, nil
	}

	var commits []Commit
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			commits = append(commits, Commit{
				Hash:    parts[0],
				Message: parts[1],
			})
		}
	}
	return commits, nil
}

func (c *Client) FileDiff(filename string) (string, error) {
	return c.run("diff", "--no-color", "--", filename)
}

func (c *Client) IsMerging() bool {
	_, err := c.run("rev-parse", "--verify", "MERGE_HEAD")
	return err == nil
}

func (c *Client) IsRebasing() bool {
	_, err := c.run("rev-parse", "--verify", "REBASE_HEAD")
	return err == nil
}

func (c *Client) StashList() ([]string, error) {
	output, err := c.run("stash", "list")
	if err != nil {
		return nil, nil
	}
	if output == "" {
		return nil, nil
	}
	return strings.Split(output, "\n"), nil
}

func (c *Client) HasUncommittedChanges() bool {
	output, err := c.run("status", "--porcelain")
	return err == nil && output != ""
}

type Commit struct {
	Hash    string
	Message string
}

type ChangeInfo struct {
	Branch        string
	CommitHash    string
	CommitMessage string
	ChangedFiles  map[string]string
	Diff          string
	Uncommitted   bool
	RecentCommits []Commit
}
