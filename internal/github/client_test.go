package github

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/MakeNowJust/heredoc"
)

func Test_New(t *testing.T) {
	client := New(nil)
	if cli, ok := client.labels.(*Cli); !ok {
		t.Errorf("Client.labels is not expected type *Cli")
	} else if cli.Owner != ":owner" || cli.Repo != ":repo" {
		t.Errorf(`Client.labels is %v, want &{:owner :repo}`, cli)
	}
}

func Test_New_Mock(t *testing.T) {
	mock := &Mock{}
	client := New(mock)
	if mock != client.labels {
		t.Errorf("Client.labels is not expected type Mock")
	}
}

func Test_CreateLabel(t *testing.T) {
	tests := []struct {
		name   string
		stdout bytes.Buffer
		err    error
		want   Label
		wantE  bool
	}{
		{
			name:  "gh error",
			err:   errors.New("gh exited with code 1"),
			wantE: true,
		},
		{
			name:   "deserialization error",
			stdout: *bytes.NewBufferString("invalid JSON"),
			wantE:  true,
		},
		{
			name: "success",
			stdout: *bytes.NewBufferString(heredoc.Doc(`{
				"name": "test",
				"color": "112233",
				"description": "testing"
			}`)),
			want: Label{
				Name:        "test",
				Color:       "112233",
				Description: "testing",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := Mock{
				Stdout: tt.stdout,
				Err:    tt.err,
			}
			client := Client{
				labels: &mock,
			}
			label := Label{
				Name:        "test",
				Color:       "112233",
				Description: "testing",
			}
			if got, err := client.CreateLabel(label); (err != nil) != tt.wantE {
				t.Errorf("CreateLabel() error = %v, want: %v", err, tt.wantE)
			} else if got != tt.want {
				t.Errorf("CreateLabel() = %v, want: %v", got, tt.want)
			}
		})
	}
}

func Test_DeleteLabel(t *testing.T) {
	tests := []struct {
		name  string
		err   error
		wantE bool
	}{
		{
			name:  "gh error",
			err:   errors.New("gh exited with code 1"),
			wantE: true,
		},
		{
			name: "success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := Mock{
				Err: tt.err,
			}
			client := Client{
				labels: &mock,
			}
			if err := client.DeleteLabel("test"); (err != nil) != tt.wantE {
				t.Errorf("DeleteLabel() error = %v, want: %v", err, tt.wantE)
			}
		})
	}
}

func Test_ListLabels(t *testing.T) {
	tests := []struct {
		name   string
		stdout bytes.Buffer
		err    error
		want   Labels
		wantE  bool
	}{
		{
			name:  "gh error",
			err:   errors.New("gh exited with code 1"),
			wantE: true,
		},
		{
			name:   "deserialization error",
			stdout: *bytes.NewBufferString("invalid JSON"),
			wantE:  true,
		},
		{
			name: "success with single page",
			stdout: *bytes.NewBufferString(heredoc.Doc(`{
				"data": {
					"repository": {
						"labels": {
							"nodes": [
								{
									"name": "test",
									"color": "112233",
									"description": "testing"
								}
							],
							"pageInfo":{"hasNextPage":true,"endCursor":"abcd1234"}}}}}`)),
			want: Labels{
				{
					Name:        "test",
					Color:       "112233",
					Description: "testing",
				},
			},
		},
		// cSpell:ignore efgh5678
		{
			name: "success with multiple pages",
			stdout: *bytes.NewBufferString(heredoc.Doc(`{
				"data": {
					"repository": {
						"labels": {
							"nodes": [
								{
									"name": "test",
									"color": "112233",
									"description": "testing"
								}
							],
							"pageInfo":{"hasNextPage":true,"endCursor":"abcd1234"}}}}}{
				"data": {
					"repository": {
						"labels": {
							"nodes": [
								{
									"name": "test2",
									"color": "223344",
									"description": "testing again"
								}
							],
							"pageInfo":{"hasNextPage":false,"endCursor":"efgh5678"}}}}}`)),
			want: Labels{
				{
					Name:        "test",
					Color:       "112233",
					Description: "testing",
				},
				{
					Name:        "test2",
					Color:       "223344",
					Description: "testing again",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := Mock{
				Stdout: tt.stdout,
				Err:    tt.err,
			}
			client := Client{
				labels: &mock,
			}
			if got, err := client.ListLabels(""); (err != nil) != tt.wantE {
				t.Errorf("ListLabels() error = %v, want: %v", err, tt.wantE)
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListLabels() = %v, want: %v", got, tt.want)
			}
		})
	}
}

func Test_UpdateLabel(t *testing.T) {
	tests := []struct {
		name   string
		stdout bytes.Buffer
		err    error
		want   Label
		wantE  bool
	}{
		{
			name:  "gh error",
			err:   errors.New("gh exited with code 1"),
			wantE: true,
		},
		{
			name:   "deserialization error",
			stdout: *bytes.NewBufferString("invalid JSON"),
			wantE:  true,
		},
		{
			name: "success",
			stdout: *bytes.NewBufferString(heredoc.Doc(`{
				"name": "renamed",
				"color": "112233",
				"description": "testing"
			}`)),
			want: Label{
				Name:        "renamed",
				Color:       "112233",
				Description: "testing",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := Mock{
				Stdout: tt.stdout,
				Err:    tt.err,
			}
			client := Client{
				labels: &mock,
			}
			label := EditLabel{
				Label: Label{
					Name:        "test",
					Color:       "112233",
					Description: "testing",
				},
				NewName: "renamed",
			}
			if got, err := client.UpdateLabel(label); (err != nil) != tt.wantE {
				t.Errorf("UpdateLabel() error = %v, want: %v", err, tt.wantE)
			} else if got != tt.want {
				t.Errorf("UpdateLabel() = %v, want: %v", got, tt.want)
			}
		})
	}
}
