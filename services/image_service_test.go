package services

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/hujia-team/intranet-sdk/client"
	"github.com/hujia-team/intranet-sdk/models"
	"github.com/hujia-team/intranet-sdk/utils"
)

func newImageTestService(t *testing.T, handler http.HandlerFunc) *imageService {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	httpClient, err := client.NewHTTPClient(&client.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatalf("new http client: %v", err)
	}
	return NewImageService(httpClient).(*imageService)
}

func imageStringPtr(value string) *string { return &value }
func uint64Ptr(value uint64) *uint64      { return &value }

func TestListImages(t *testing.T) {
	service := newImageTestService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.URL.Path != "/api/v1/images" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Body == nil {
			t.Fatal("expected request body to be readable")
		}
		body := make([]byte, 1)
		if n, _ := r.Body.Read(body); n != 0 {
			t.Fatalf("GET request must not contain a body: %q", body[:n])
		}

		query := r.URL.Query()
		expected := url.Values{
			"page_num":              {"2"},
			"page_size":             {"20"},
			"git_repo":              {"group/project"},
			"git_branch":            {"main/release"},
			"img_fullname":          {"registry.example.com/project/image:tag"},
			"digest":                {"sha256:abc+123"},
			"git_commit_short_hash": {"abc&123"},
			"pipeline_id":           {"10"},
			"job_id":                {"20"},
		}
		if query.Encode() != expected.Encode() {
			t.Fatalf("unexpected query: got %s, want %s", query.Encode(), expected.Encode())
		}

		_, _ = w.Write([]byte(`{"code":0,"msg":"success","data":{"total":1,"items":[{"id":7,"created_at":"2026-09-05T00:00:00Z","updated_at":"2026-09-05T01:00:00Z","img_registry":"registry.example.com","img_project":"project","img_repo":"image","digest":"sha256:abc+123","os":"linux","arch":"amd64","img_tag":"tag","job_url":"https://ci/job/20","pipeline_url":"https://ci/pipeline/10","git_repo":"group/project","git_branch":"main/release","git_commit_hash":"abcdef","git_commit_at":"2026-09-04T00:00:00Z"}]}}`))
	})

	result, err := service.ListImages(&models.ImageListReq{
		PageNum:            uint64Ptr(2),
		PageSize:           uint64Ptr(20),
		GitRepo:            imageStringPtr("group/project"),
		GitBranch:          imageStringPtr("main/release"),
		ImgFullname:        imageStringPtr("registry.example.com/project/image:tag"),
		Digest:             imageStringPtr("sha256:abc+123"),
		GitCommitShortHash: imageStringPtr("abc&123"),
		PipelineID:         uint64Ptr(10),
		JobID:              uint64Ptr(20),
	})
	if err != nil {
		t.Fatalf("ListImages error: %v", err)
	}
	if result.Total != 1 || len(result.Items) != 1 {
		t.Fatalf("unexpected result: %#v", result)
	}
	item := result.Items[0]
	if item.ID != 7 || item.ImgRepo != "image" || item.Digest != "sha256:abc+123" {
		t.Fatalf("unexpected image item: %#v", item)
	}
	if item.GitCommitHash == nil || *item.GitCommitHash != "abcdef" {
		t.Fatalf("unexpected commit hash: %#v", item.GitCommitHash)
	}
}

func TestListImagesOmitsUnsetQuery(t *testing.T) {
	service := newImageTestService(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Encode(); got != "page_num=1" {
			t.Fatalf("unexpected query: %s", got)
		}
		_, _ = w.Write([]byte(`{"code":0,"data":{"total":0,"items":[]}}`))
	})

	_, err := service.ListImages(&models.ImageListReq{PageNum: uint64Ptr(1)})
	if err != nil {
		t.Fatalf("ListImages error: %v", err)
	}
}

func TestListImagesNilRequest(t *testing.T) {
	service := newImageTestService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"code":0,"data":null}`))
	})

	result, err := service.ListImages(nil)
	if err != nil {
		t.Fatalf("ListImages error: %v", err)
	}
	if result == nil || result.Total != 0 || result.Items == nil {
		t.Fatalf("expected empty typed result, got %#v", result)
	}
}

func TestListImagesSendsExplicitZeroAndEmptyValues(t *testing.T) {
	service := newImageTestService(t, func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if value, ok := query["page_num"]; !ok || len(value) != 1 || value[0] != "0" {
			t.Fatalf("unexpected page_num: %#v", value)
		}
		if value, ok := query["git_repo"]; !ok || len(value) != 1 || value[0] != "" {
			t.Fatalf("unexpected git_repo: %#v", value)
		}
		_, _ = w.Write([]byte(`{"code":0,"data":{"total":0,"items":[]}}`))
	})

	_, err := service.ListImages(&models.ImageListReq{
		PageNum: uint64Ptr(0),
		GitRepo: imageStringPtr(""),
	})
	if err != nil {
		t.Fatalf("ListImages error: %v", err)
	}
}

func TestListImagesBusinessError(t *testing.T) {
	service := newImageTestService(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"code":4,"msg":"image not found"}`))
	})

	_, err := service.ListImages(nil)
	if err == nil || !strings.Contains(err.Error(), "image not found") {
		t.Fatalf("expected business error, got %v", err)
	}
	var sdkErr *utils.SDKError
	if !errors.As(err, &sdkErr) || sdkErr.Code != utils.ErrCodeAPIError {
		t.Fatalf("unexpected error type: %T %v", err, err)
	}
}

func TestListImagesHTTPError(t *testing.T) {
	service := newImageTestService(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	})

	_, err := service.ListImages(nil)
	if err == nil {
		t.Fatal("expected HTTP error")
	}
	var outer *utils.SDKError
	if !errors.As(err, &outer) || outer.Code != utils.ErrCodeAPIError {
		t.Fatalf("unexpected outer error: %T %v", err, err)
	}
	var inner *utils.SDKError
	if !errors.As(outer.Unwrap(), &inner) || inner.Code != utils.ErrCodeUnauthorized {
		t.Fatalf("expected unauthorized cause, got %v", outer.Unwrap())
	}
}
