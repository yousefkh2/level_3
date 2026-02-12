package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	cnpgv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1"
	"github.com/labstack/echo/v4"
	"github.com/yousefkh2/level_3/week_4/api/app"
	"github.com/yousefkh2/level_3/week_4/api/models"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type MockDatabaseService struct {
	CreateFunc func(ctx context.Context, name string, instances int, storage string) (*cnpgv1.Cluster, error)
	ListFunc   func(ctx context.Context) ([]cnpgv1.Cluster, error)
	GetFunc    func(ctx context.Context, name string) (*cnpgv1.Cluster, error)
	DeleteFunc func(ctx context.Context, name string) error
}

// boiler plate code that satisfies the interface - doesn't do any "work" by itself; it just checks if
// i've provided a custom behavior
func (m *MockDatabaseService) CreateDatabase(ctx context.Context, name string, instances int, storage string) (*cnpgv1.Cluster, error) {
	if m.CreateFunc != nil {
		// if a custom func was provided, run it and return its result
		return m.CreateFunc(ctx, name, instances, storage)
	}
	return nil, nil
}

func (m *MockDatabaseService) ListDatabases(ctx context.Context) ([]cnpgv1.Cluster, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return nil, nil
}

func (m *MockDatabaseService) GetDatabase(ctx context.Context, name string) (*cnpgv1.Cluster, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, name)
	}
	return nil, nil
}

func (m *MockDatabaseService) DeleteDatabase(ctx context.Context, name string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, name)
	}
	return nil
}

func TestCreateDatabase(t *testing.T) {
	// each test case swaps in a different CreateFunc to simulate a different service outcome
	tests := []struct {
		name         string // sub-test name
		body         string // JSON request body
		createFunc   func(ctx context.Context, name string, instances int, storage string) (*cnpgv1.Cluster, error)
		wantStatus   int                      // expected HTTP status code
		wantErrKey   string                   // for error responses (if non-empty, expect this key in a map[string]string response)
		wantResponse *models.DatabaseResponse // for success responses (if non-nil, decode into DatabaseResponse and check fields)
	}{
		{
			name: "Success",
			body: `{"name":"test-db","instances":3,"storage":"1Gi"}`,
			createFunc: func(ctx context.Context, name string, instances int, storage string) (*cnpgv1.Cluster, error) {
				return &cnpgv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: name},
					Spec: cnpgv1.ClusterSpec{
						Instances:            instances,
						StorageConfiguration: cnpgv1.StorageConfiguration{Size: storage},
					},
				}, nil
			},
			wantStatus: http.StatusCreated,
			wantResponse: &models.DatabaseResponse{
				Name:   "test-db",
				Spec:   models.DatabaseSpec{Instances: 3, Storage: "1Gi"},
				Status: "creating",
			},
		},
		{
			name: "AlreadyExists",
			body: `{"name":"existing-db","instances":3,"storage":"1Gi"}`,
			createFunc: func(ctx context.Context, name string, instances int, storage string) (*cnpgv1.Cluster, error) {
				return nil, apierrors.NewAlreadyExists(
					schema.GroupResource{Group: "postgresql.cnpg.io", Resource: "clusters"}, name,
				)
			},
			wantStatus: http.StatusConflict,
			wantErrKey: "error",
		},
		{
			name:       "InvalidBody",
			body:       `not json at all`,
			createFunc: nil, // should never be called
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name: "InvalidConfig",
			body: `{"name":"bad-db","instances":-1,"storage":"1Gi"}`,
			createFunc: func(ctx context.Context, name string, instances int, storage string) (*cnpgv1.Cluster, error) {
				return nil, apierrors.NewInvalid(
					schema.GroupKind{Group: "postgresql.cnpg.io", Kind: "Cluster"},
					name, nil,
				)
			},
			wantStatus: http.StatusBadRequest,
			wantErrKey: "error",
		},
		{
			name: "Unauthorized",
			body: `{"name":"secret-db","instances":1,"storage":"1Gi"}`,
			createFunc: func(ctx context.Context, name string, instances int, storage string) (*cnpgv1.Cluster, error) {
				return nil, apierrors.NewUnauthorized("not allowed")
			},
			wantStatus: http.StatusUnauthorized,
			wantErrKey: "error",
		},
		{
			name: "InternalServerError",
			body: `{"name":"fail-db","instances":1,"storage":"1Gi"}`,
			createFunc: func(ctx context.Context, name string, instances int, storage string) (*cnpgv1.Cluster, error) {
				return nil, fmt.Errorf("something broke")
			},
			wantStatus: http.StatusInternalServerError,
			wantErrKey: "error",
		},
	}
	// _ is the index (i don't care about it)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// 1. build request / recorder
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/databases", strings.NewReader(tc.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// 2. inject mock
			c.Set("app", &app.App{
				DBService: &MockDatabaseService{CreateFunc: tc.createFunc},
			})

			// 3. execute
			err := CreateDatabase(c)
			if err != nil {
				t.Fatalf("expected no error from handler, got %v", err)
			}

			// 4. check status code
			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}

			// 5. check response body
			if tc.wantResponse != nil {
				var got models.DatabaseResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				if got.Name != tc.wantResponse.Name {
					t.Errorf("name = %q, want %q", got.Name, tc.wantResponse.Name)
				}
				if got.Spec.Instances != tc.wantResponse.Spec.Instances {
					t.Errorf("instances = %d, want %d", got.Spec.Instances, tc.wantResponse.Spec.Instances)
				}
				if got.Spec.Storage != tc.wantResponse.Spec.Storage {
					t.Errorf("storage = %q, want %q", got.Spec.Storage, tc.wantResponse.Spec.Storage)
				}
				if got.Status != tc.wantResponse.Status {
					t.Errorf("status = %q, want %q", got.Status, tc.wantResponse.Status)
				}
			}
			if tc.wantErrKey != "" {
				var got map[string]string
				if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
					t.Fatalf("failed to parse error response: %v", err)
				}
				if got[tc.wantErrKey] == "" {
					t.Errorf("expected %q key in response, got %v", tc.wantErrKey, got)
				}
			}
		})
	}
}

func TestListDatabases(t *testing.T) {
	tests := []struct {
		name       string
		listFunc   func(ctx context.Context) ([]cnpgv1.Cluster, error)
		wantStatus int
		wantCount  int    // expected number of items (-1 to skip check)
		wantErrKey string // if non-empty, expect error response
	}{
		{
			name: "Success",
			listFunc: func(ctx context.Context) ([]cnpgv1.Cluster, error) {
				return []cnpgv1.Cluster{
					{
						ObjectMeta: metav1.ObjectMeta{Name: "db1"},
						Spec: cnpgv1.ClusterSpec{
							Instances:            3,
							StorageConfiguration: cnpgv1.StorageConfiguration{Size: "1Gi"},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{Name: "db2"},
						Spec: cnpgv1.ClusterSpec{
							Instances:            5,
							StorageConfiguration: cnpgv1.StorageConfiguration{Size: "2Gi"},
						},
					},
				}, nil
			},
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
		{
			name: "EmptyList",
			listFunc: func(ctx context.Context) ([]cnpgv1.Cluster, error) {
				return []cnpgv1.Cluster{}, nil
			},
			wantStatus: http.StatusOK,
			wantCount:  0,
		},
		{
			name: "InternalServerError",
			listFunc: func(ctx context.Context) ([]cnpgv1.Cluster, error) {
				return nil, fmt.Errorf("k8s unreachable")
			},
			wantStatus: http.StatusInternalServerError,
			wantCount:  -1,
			wantErrKey: "error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/databases", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.Set("app", &app.App{
				DBService: &MockDatabaseService{ListFunc: tc.listFunc},
			})

			err := ListDatabases(c)
			if err != nil {
				t.Fatalf("expected no error from handler, got %v", err)
			}

			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}

			if tc.wantCount >= 0 {
				var got []models.DatabaseResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				if len(got) != tc.wantCount {
					t.Errorf("count = %d, want %d", len(got), tc.wantCount)
				}
			}
			if tc.wantErrKey != "" {
				var got map[string]string
				if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
					t.Fatalf("failed to parse error response: %v", err)
				}
				if got[tc.wantErrKey] == "" {
					t.Errorf("expected %q key in response, got %v", tc.wantErrKey, got)
				}
			}
		})
	}
}

func TestGetDatabase(t *testing.T) {
	tests := []struct {
		name       string
		paramName  string // URL param :name
		getFunc    func(ctx context.Context, name string) (*cnpgv1.Cluster, error)
		wantStatus int
		wantErrKey string                         // if non-empty, expect error response
		wantResp   *models.DatabaseDetailResponse // if non-nil, check fields
	}{
		{
			name:      "Success",
			paramName: "my-db",
			getFunc: func(ctx context.Context, name string) (*cnpgv1.Cluster, error) {
				return &cnpgv1.Cluster{
					ObjectMeta: metav1.ObjectMeta{Name: name},
					Spec: cnpgv1.ClusterSpec{
						Instances:            3,
						StorageConfiguration: cnpgv1.StorageConfiguration{Size: "1Gi"},
					},
				}, nil
			},
			wantStatus: http.StatusOK,
			wantResp: &models.DatabaseDetailResponse{
				Name: "my-db",
				Spec: models.DatabaseSpec{Instances: 3, Storage: "1Gi"},
				Connection: models.ConnectionInfo{
					Host:     "my-db-rw.default.svc.cluster.local",
					Port:     5432,
					Database: "app",
					Username: "app",
				},
			},
		},
		{
			name:      "NotFound",
			paramName: "ghost-db",
			getFunc: func(ctx context.Context, name string) (*cnpgv1.Cluster, error) {
				return nil, apierrors.NewNotFound(
					schema.GroupResource{Group: "postgresql.cnpg.io", Resource: "clusters"}, name,
				)
			},
			wantStatus: http.StatusNotFound,
			wantErrKey: "error",
		},
		{
			name:      "InternalServerError",
			paramName: "fail-db",
			getFunc: func(ctx context.Context, name string) (*cnpgv1.Cluster, error) {
				return nil, fmt.Errorf("k8s blew up")
			},
			wantStatus: http.StatusInternalServerError,
			wantErrKey: "error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/databases/"+tc.paramName, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			// set the :name URL param so c.Param("name") works inside the handler
			c.SetParamNames("name")
			c.SetParamValues(tc.paramName)

			c.Set("app", &app.App{
				DBService: &MockDatabaseService{GetFunc: tc.getFunc},
			})

			err := GetDatabase(c)
			if err != nil {
				t.Fatalf("expected no error from handler, got %v", err)
			}

			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}

			if tc.wantResp != nil {
				var got models.DatabaseDetailResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
					t.Fatalf("failed to parse response: %v", err)
				}
				if got.Name != tc.wantResp.Name {
					t.Errorf("name = %q, want %q", got.Name, tc.wantResp.Name)
				}
				if got.Spec.Instances != tc.wantResp.Spec.Instances {
					t.Errorf("instances = %d, want %d", got.Spec.Instances, tc.wantResp.Spec.Instances)
				}
				if got.Spec.Storage != tc.wantResp.Spec.Storage {
					t.Errorf("storage = %q, want %q", got.Spec.Storage, tc.wantResp.Spec.Storage)
				}
				if got.Connection.Host != tc.wantResp.Connection.Host {
					t.Errorf("host = %q, want %q", got.Connection.Host, tc.wantResp.Connection.Host)
				}
				if got.Connection.Port != tc.wantResp.Connection.Port {
					t.Errorf("port = %d, want %d", got.Connection.Port, tc.wantResp.Connection.Port)
				}
			}
			if tc.wantErrKey != "" {
				var got map[string]string
				if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
					t.Fatalf("failed to parse error response: %v", err)
				}
				if got[tc.wantErrKey] == "" {
					t.Errorf("expected %q key in response, got %v", tc.wantErrKey, got)
				}
			}
		})
	}
}

func TestDeleteDatabase(t *testing.T) {
	tests := []struct {
		name       string
		paramName  string // URL param :name
		deleteFunc func(ctx context.Context, name string) error
		wantStatus int
		wantErrKey string // if non-empty, expect JSON error body
	}{
		{
			name:      "Success",
			paramName: "my-db",
			deleteFunc: func(ctx context.Context, name string) error {
				return nil
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:      "NotFound_Idempotent",
			paramName: "gone-db",
			deleteFunc: func(ctx context.Context, name string) error {
				return apierrors.NewNotFound(
					schema.GroupResource{Group: "postgresql.cnpg.io", Resource: "clusters"}, name,
				)
			},
			wantStatus: http.StatusNoContent, // idempotent - already gone
		},
		{
			name:      "InternalServerError",
			paramName: "fail-db",
			deleteFunc: func(ctx context.Context, name string) error {
				return fmt.Errorf("something broke")
			},
			wantStatus: http.StatusInternalServerError,
			wantErrKey: "error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/databases/"+tc.paramName, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("name")
			c.SetParamValues(tc.paramName)

			c.Set("app", &app.App{
				DBService: &MockDatabaseService{DeleteFunc: tc.deleteFunc},
			})

			err := DeleteDatabase(c)
			if err != nil {
				t.Fatalf("expected no error from handler, got %v", err)
			}

			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}

			if tc.wantErrKey != "" {
				var got map[string]string
				if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
					t.Fatalf("failed to parse error response: %v", err)
				}
				if got[tc.wantErrKey] == "" {
					t.Errorf("expected %q key in response, got %v", tc.wantErrKey, got)
				}
			}
		})
	}
}
