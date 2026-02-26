package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"
)

/*
	This is the adapter that turns raw Loki log lines into structured audit history.
*/

// AuditEvent represents a single parsed audit log entry returned to the user.
type AuditEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Details   string    `json:"details,omitempty"`
}

// ServiceEvent represents a parsed service status event emitted by the watcher.
type ServiceEvent struct {
	Timestamp      time.Time `json:"timestamp"`
	Resource       string    `json:"resource"`
	Event          string    `json:"event"`
	PreviousStatus string    `json:"previous_status,omitempty"`
	CurrentStatus  string    `json:"current_status,omitempty"`
}

type LokiService struct {
	lokiURL string
	client  *http.Client
}

func NewLokiService(lokiURL string) *LokiService {
	return &LokiService{
		lokiURL: lokiURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

type lokiQueryResponse struct {
	Data struct {
		Result []struct {
			Values [][]string `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

func (s *LokiService) GetAuditLogs(ctx context.Context, databaseName string, since time.Duration) ([]AuditEvent, error) {
	query := fmt.Sprintf(
		`{namespace="paas-control-plane"} | json | log_type="audit" | resource="%s"`,
		databaseName,
	)
	return s.queryAuditLogs(ctx, query, since)
}

func (s *LokiService) GetGlobalAuditLogs(ctx context.Context, since time.Duration) ([]AuditEvent, error) {
	query := `{namespace="paas-control-plane"} | json | log_type="audit"`
	return s.queryAuditLogs(ctx, query, since)
}

func (s *LokiService) GetServiceLogs(ctx context.Context, databaseName string, since time.Duration) ([]ServiceEvent, error) {
	query := fmt.Sprintf(
		`{namespace="paas-control-plane"} | json | log_type="service" | resource="%s"`,
		databaseName,
	)
	return s.queryServiceLogs(ctx, query, since)
}

func (s *LokiService) queryAuditLogs(ctx context.Context, query string, since time.Duration) ([]AuditEvent, error) {
	// Build the HTTP request to Loki's query_range endpoint
	params := url.Values{}
	params.Set("query", query)
	params.Set("start", fmt.Sprintf("%d", time.Now().Add(-since).UnixNano()))
	params.Set("end", fmt.Sprintf("%d", time.Now().UnixNano()))
	params.Set("limit", "100")
	params.Set("direction", "backward") // newest first

	reqURL := fmt.Sprintf("%s/loki/api/v1/query_range?%s", s.lokiURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("building loki request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("querying loki: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("loki returned status %d", resp.StatusCode)
	}

	var lokiResp lokiQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&lokiResp); err != nil {
		return nil, fmt.Errorf("decoding loki response: %w", err)
	}

	var events []AuditEvent
	for _, result := range lokiResp.Data.Result {
		for _, value := range result.Values {
			if len(value) < 2 {
				continue
			}

			// value[0] is unix nanosecond timestamp as string
			// value[1] is the raw log line (JSON from zap)
			var logLine map[string]interface{}
			if err := json.Unmarshal([]byte(value[1]), &logLine); err != nil {
				continue // skip malformed lines
			}

			ts := time.Now()
			if tsFloat, ok := logLine["ts"].(float64); ok {
				ts = time.Unix(int64(tsFloat), 0)
			}

			event := AuditEvent{
				Timestamp: ts,
				Action:    stringField(logLine, "action"),
				Resource:  stringField(logLine, "resource"),
			}

			details := ""
			if instances, ok := logLine["instances"]; ok {
				details += fmt.Sprintf("instances=%v ", instances)
			}
			if storage, ok := logLine["storage"]; ok {
				details += fmt.Sprintf("storage=%v", storage)
			}
			event.Details = details

			events = append(events, event)
		}
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.After(events[j].Timestamp)
	})

	return events, nil
}

func (s *LokiService) queryServiceLogs(ctx context.Context, query string, since time.Duration) ([]ServiceEvent, error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("start", fmt.Sprintf("%d", time.Now().Add(-since).UnixNano()))
	params.Set("end", fmt.Sprintf("%d", time.Now().UnixNano()))
	params.Set("limit", "100")
	params.Set("direction", "backward") // newest first

	reqURL := fmt.Sprintf("%s/loki/api/v1/query_range?%s", s.lokiURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("building loki request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("querying loki: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("loki returned status %d", resp.StatusCode)
	}

	var lokiResp lokiQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&lokiResp); err != nil {
		return nil, fmt.Errorf("decoding loki response: %w", err)
	}

	var events []ServiceEvent
	for _, result := range lokiResp.Data.Result {
		for _, value := range result.Values {
			if len(value) < 2 {
				continue
			}

			var logLine map[string]interface{}
			if err := json.Unmarshal([]byte(value[1]), &logLine); err != nil {
				continue
			}

			ts := time.Now()
			if tsFloat, ok := logLine["ts"].(float64); ok {
				ts = time.Unix(int64(tsFloat), 0)
			}

			event := ServiceEvent{
				Timestamp:      ts,
				Resource:       stringField(logLine, "resource"),
				Event:          stringField(logLine, "event"),
				PreviousStatus: stringField(logLine, "previous_status"),
				CurrentStatus:  stringField(logLine, "current_status"),
			}

			events = append(events, event)
		}
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.After(events[j].Timestamp)
	})

	return events, nil
}

func stringField(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
