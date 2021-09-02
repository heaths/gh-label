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
	Url         string `json:"url,omitempty"`
}

type Labels []Label

type Client struct {
	labels LabelsService
}

type LabelsService interface {
	CreateLabel(label Label) (bytes.Buffer, error)
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

func (c *Client) CreateLabel(label Label) (Label, error) {
	buf, err := c.labels.CreateLabel(label)
	if err != nil {
		return Label{}, err
	}

	label = Label{}
	if err = json.Unmarshal(buf.Bytes(), &label); err != nil {
		return Label{}, fmt.Errorf("failed to read label; error: %w, data: %s", err, buf.String())
	}

	return label, nil
}

func (c *Client) ListLabels(substr string) (Labels, error) {
	buf, err := c.labels.ListLabels(substr)
	if err != nil {
		return nil, err
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

func (m *Mock) CreateLabel(label Label) (bytes.Buffer, error) {
	return m.Stdout, m.Err
}

func (m *Mock) ListLabels(substr string) (bytes.Buffer, error) {
	return m.Stdout, m.Err
}

func (m *Mock) DeleteLabel(name string) error {
	return m.Err
}
