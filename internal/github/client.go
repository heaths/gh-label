package github

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Label struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description,omitempty"`
}

type Labels []Label

type Client struct {
	labels LabelsService
}

type LabelsService interface {
	CreateOrUpdateLabel(label Label) error
	ListLabels(substr string) (bytes.Buffer, error)
	DeleteLabel(name string) error
}

func New(labels LabelsService) *Client {
	if labels == nil {
		labels = &Cli{
			Owner: ":owner",
			Repo:  ":repo",
		}
	}

	return &Client{
		labels,
	}
}

func (c *Client) ListLabels(substr string) (Labels, error) {
	buf, err := c.labels.ListLabels(substr)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize labels, error: %w", err)
	}

	type response struct {
		Data struct {
			Repository struct {
				Labels struct {
					Nodes Labels
				}
			}
		}
	}

	var labels Labels

	// Work around https://github.com/cli/cli/issues/1268 by splitting responses after cursor info.
	for _, data := range bytes.SplitAfter(buf.Bytes(), []byte("}}}}}")) {
		if len(data) == 0 {
			break
		}

		var resp response
		if err = json.Unmarshal(data, &resp); err != nil {
			return nil, fmt.Errorf("failed to read labels; error: %w, data: %s", err, data)
		}

		labels = append(labels, resp.Data.Repository.Labels.Nodes...)
	}

	return labels, nil
}

type Mock struct {
	Stdout bytes.Buffer
	Err    error
}

func (m *Mock) CreateOrUpdateLabel(label Label) error {
	return m.Err
}

func (m *Mock) ListLabels(substr string) (bytes.Buffer, error) {
	return m.Stdout, m.Err
}

func (m *Mock) DeleteLabel(name string) error {
	return m.Err
}
