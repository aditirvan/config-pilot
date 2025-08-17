package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Commit represents a GitHub commit
type Commit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Message string `json:"message"`
	} `json:"commit"`
	HTMLURL string `json:"html_url"`
}

// Client handles GitHub API interactions
type Client struct {
	token   string
	baseURL string
	client  *http.Client
	owner   string
	repo    string
}

// NewClient creates a new GitHub client
func NewClient(token, owner, repo string) *Client {
	return &Client{
		token:   token,
		baseURL: "https://api.github.com",
		client:  &http.Client{Timeout: 10 * time.Second},
		owner:   owner,
		repo:    repo,
	}
}

// GetLatestCommit fetches the latest commit from the repository
func (c *Client) GetLatestCommit(path string) (*Commit, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/commits", c.baseURL, c.owner, c.repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	if path != "" {
		q.Add("path", path)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	var commits []Commit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(commits) == 0 {
		return nil, fmt.Errorf("no commits found")
	}

	return &commits[0], nil
}
