package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const baseURL = "http://unix"

type DockerClient struct {
	client *http.Client
}

func New(client *http.Client) *DockerClient {
	return &DockerClient{client}
}

func (dc *DockerClient) get(ctx context.Context, relativePath string, queryParams url.Values) (*http.Response, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to call get: %w", err)
	}

	u.Path = path.Join(u.Path, relativePath)
	u.RawQuery = queryParams.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call get: %w", err)
	}

	res, err := dc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call get: %w", err)
	}

	if res.StatusCode >= 400 {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to call get; failed to read body: %w", err)
		}
		return nil, fmt.Errorf("failed to call GET: status %d, body: %s", res.StatusCode, string(body))
	}

	return res, nil
}

func (dc *DockerClient) ContainerList(ctx context.Context, opts ContainerListOptions) ([]Container, error) {
	queryParams := url.Values{}

	if opts.All {
		queryParams.Set("all", "1")
	}

	if opts.Limit > 0 {
		queryParams.Set("limit", strconv.Itoa(opts.Limit))
	}

	if opts.Since != "" {
		queryParams.Set("since", opts.Since)
	}

	if opts.Before != "" {
		queryParams.Set("before", opts.Before)
	}

	if opts.Size {
		queryParams.Set("size", "1")
	}

	if opts.Filters.Len() > 0 {
		filterJSON, err := opts.Filters.Encode()
		if err != nil {
			return nil, err
		}

		queryParams.Set("filters", filterJSON)
	}

	url := "/containers/json"
	resp, err := dc.get(ctx, url, queryParams)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var containers []Container
	json.NewDecoder(resp.Body).Decode(&containers)

	return containers, nil
}

func (dc *DockerClient) ContainerLogs(ctx context.Context, containerID string, opts ContainerLogsOptions) (io.ReadCloser, error) {
	queryParams := url.Values{}
	if opts.Stdout {
		queryParams.Set("stdout", "1")
	}

	if opts.Stderr {
		queryParams.Set("stderr", "1")
	}

	if opts.Since != "" {
		ts, err := parseTimestamp(opts.Since)
		if err != nil {
			return nil, fmt.Errorf("failed to get container list: %w", err)
		}

		queryParams.Set("since", fmt.Sprintf("%d", ts))
	}

	if opts.Until != "" {
		ts, err := parseTimestamp(opts.Until)
		if err != nil {
			return nil, fmt.Errorf("failed to get container list: %w", err)
		}

		queryParams.Set("until", fmt.Sprintf("%d", ts))
	}

	if opts.Timestamp {
		queryParams.Set("timestamps", "1")
	}

	if opts.Follow {
		queryParams.Set("follow", "1")
	}

	if opts.Tail != "" {
		queryParams.Set("tail", opts.Tail)
	}

	url := fmt.Sprintf("/containers/%s/logs", containerID)
	resp, err := dc.get(ctx, url, queryParams)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (dc *DockerClient) Ping(ctx context.Context) (string, error) {
	resp, err := dc.get(ctx, "/_ping", url.Values{})
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (dc *DockerClient) Close() {
	if dc.client != nil {
		dc.client.CloseIdleConnections()
	}
}
