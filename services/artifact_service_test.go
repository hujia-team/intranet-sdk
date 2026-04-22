package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hujia-team/intranet-sdk/client"
	"github.com/hujia-team/intranet-sdk/models"
)

func newArtifactTestService(t *testing.T, handler http.HandlerFunc) *artifactService {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	httpClient, err := client.NewHTTPClient(&client.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatalf("new http client: %v", err)
	}
	return NewArtifactService(httpClient).(*artifactService)
}

func decodeBody(t *testing.T, r *http.Request) map[string]any {
	t.Helper()
	defer r.Body.Close()
	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	return payload
}

func TestListArtifacts(t *testing.T) {
	service := newArtifactTestService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/aiplorer/artifact/list" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		payload := decodeBody(t, r)
		if payload["page"].(float64) != 1 || payload["pageSize"].(float64) != 20 {
			t.Fatalf("unexpected payload: %#v", payload)
		}
		_, _ = w.Write([]byte(`{"code":0,"data":{"total":1,"data":[{"id":12,"name":"artifact-a"}]}}`))
	})

	result, err := service.ListArtifacts(&models.ArtifactListReq{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("ListArtifacts error: %v", err)
	}
	if result.Total != 1 || result.Data[0].Name == nil || *result.Data[0].Name != "artifact-a" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestGetArtifactByName(t *testing.T) {
	call := 0
	service := newArtifactTestService(t, func(w http.ResponseWriter, r *http.Request) {
		call++
		switch r.URL.Path {
		case "/aiplorer/artifact/list":
			payload := decodeBody(t, r)
			if payload["name"].(string) != "artifact-a" {
				t.Fatalf("unexpected name payload: %#v", payload)
			}
			_, _ = w.Write([]byte(`{"code":0,"data":{"total":1,"data":[{"id":12,"name":"artifact-a"}]}}`))
		case "/aiplorer/artifact":
			payload := decodeBody(t, r)
			if payload["id"].(float64) != 12 {
				t.Fatalf("unexpected id payload: %#v", payload)
			}
			_, _ = w.Write([]byte(`{"code":0,"data":{"id":12,"name":"artifact-a","projectName":"proj-a"}}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	})

	result, err := service.GetArtifactByName("artifact-a", nil)
	if err != nil {
		t.Fatalf("GetArtifactByName error: %v", err)
	}
	if call != 2 || result.ID == nil || *result.ID != 12 {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestGetArtifactByCommitHashAndChildHashes(t *testing.T) {
	call := 0
	service := newArtifactTestService(t, func(w http.ResponseWriter, r *http.Request) {
		call++
		switch r.URL.Path {
		case "/aiplorer/artifact/by-commit-hash":
			payload := decodeBody(t, r)
			if payload["commitHash"].(string) != "root-hash" {
				t.Fatalf("unexpected commitHash payload: %#v", payload)
			}
			if payload["type"].(string) != "pkg" {
				t.Fatalf("unexpected type payload: %#v", payload)
			}
			_, _ = w.Write([]byte(`{"code":0,"data":{"id":12,"name":"artifact-a","type":"pkg","commitHash":"root-hash","dependencies":[{"id":21,"name":"child-a","type":"pkg","commitHash":"child-hash-a","parentId":12},{"id":22,"name":"child-b","type":"mcu","commitHash":"child-hash-b","parentId":21}]}}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	})

	artifact, err := service.GetArtifactByCommitHash("root-hash", &models.ArtifactLookupOptions{
		ArtifactType: "pkg",
	})
	if err != nil {
		t.Fatalf("GetArtifactByCommitHash error: %v", err)
	}
	if artifact.ID == nil || *artifact.ID != 12 {
		t.Fatalf("unexpected artifact: %#v", artifact)
	}

	childHashes, err := service.GetChildArtifactHashesByCommitHash("root-hash", &models.ArtifactLookupOptions{
		ArtifactType: "pkg",
	})
	if err != nil {
		t.Fatalf("GetChildArtifactHashesByCommitHash error: %v", err)
	}
	if len(childHashes.ChildHashes) != 2 {
		t.Fatalf("unexpected child hashes: %#v", childHashes)
	}
	if childHashes.ChildHashes[0].CommitHash == nil || *childHashes.ChildHashes[0].CommitHash != "child-hash-a" {
		t.Fatalf("unexpected first child hash: %#v", childHashes.ChildHashes[0])
	}
	if call != 2 {
		t.Fatalf("unexpected call count: %d", call)
	}
}

func TestGetArtifactByCommitHashAllowsExplicitEmptyPlatform(t *testing.T) {
	service := newArtifactTestService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/aiplorer/artifact/by-commit-hash" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		payload := decodeBody(t, r)
		if payload["commitHash"].(string) != "root-hash" {
			t.Fatalf("unexpected commitHash payload: %#v", payload)
		}
		platform, ok := payload["platform"]
		if !ok {
			t.Fatalf("expected platform field to be present, payload: %#v", payload)
		}
		if platform.(string) != "" {
			t.Fatalf("expected empty platform payload, got %#v", platform)
		}
		_, _ = w.Write([]byte(`{"code":0,"data":{"id":12,"name":"artifact-a","type":"app","commitHash":"root-hash","platform":""}}`))
	})

	emptyPlatform := ""
	artifact, err := service.GetArtifactByCommitHash("root-hash", &models.ArtifactLookupOptions{
		ArtifactType: "app",
		Platform:     &emptyPlatform,
	})
	if err != nil {
		t.Fatalf("GetArtifactByCommitHash error: %v", err)
	}
	if artifact.ID == nil || *artifact.ID != 12 {
		t.Fatalf("unexpected artifact: %#v", artifact)
	}
}

func TestCheckExistsByCommitHashAllowsExplicitEmptyPlatform(t *testing.T) {
	service := newArtifactTestService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/aiplorer/artifact/by-commit-hash" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		payload := decodeBody(t, r)
		if payload["commitHash"].(string) != "root-hash" {
			t.Fatalf("unexpected commitHash payload: %#v", payload)
		}
		platform, ok := payload["platform"]
		if !ok {
			t.Fatalf("expected platform field to be present, payload: %#v", payload)
		}
		if platform.(string) != "" {
			t.Fatalf("expected empty platform payload, got %#v", platform)
		}
		_, _ = w.Write([]byte(`{"code":0,"data":{"id":12,"name":"artifact-a","type":"app","commitHash":"root-hash","platform":"","fullPath":"repo/path/artifact.tar.gz","fileHash":"abc123"}}`))
	})

	emptyPlatform := ""
	exists, err := service.CheckExistsByCommitHash("root-hash", &models.ArtifactLookupOptions{
		ArtifactType: "app",
		Platform:     &emptyPlatform,
	})
	if err != nil {
		t.Fatalf("CheckExistsByCommitHash error: %v", err)
	}
	if !exists {
		t.Fatal("expected artifact to exist")
	}
}

func TestCheckExistsByCommitHashReturnsFalseForVirtualArtifact(t *testing.T) {
	service := newArtifactTestService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/aiplorer/artifact/by-commit-hash" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"code":0,"data":{"id":12,"name":"artifact-a","type":"app","commitHash":"root-hash"}}`))
	})

	exists, err := service.CheckExistsByCommitHash("root-hash", &models.ArtifactLookupOptions{
		ArtifactType: "app",
	})
	if err != nil {
		t.Fatalf("CheckExistsByCommitHash error: %v", err)
	}
	if exists {
		t.Fatal("expected virtual artifact without fullPath to be treated as non-existent")
	}
}

func TestCheckExistsByCommitHashReturnsFalseWithoutFileHash(t *testing.T) {
	service := newArtifactTestService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/aiplorer/artifact/by-commit-hash" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"code":0,"data":{"id":12,"name":"artifact-a","type":"app","commitHash":"root-hash","fullPath":"repo/path/artifact.tar.gz"}}`))
	})

	exists, err := service.CheckExistsByCommitHash("root-hash", &models.ArtifactLookupOptions{
		ArtifactType: "app",
	})
	if err != nil {
		t.Fatalf("CheckExistsByCommitHash error: %v", err)
	}
	if exists {
		t.Fatal("expected artifact without fileHash to be treated as non-existent")
	}
}

func TestCheckExistsByCommitHashReturnsFalseWhenArtifactMissing(t *testing.T) {
	service := newArtifactTestService(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/aiplorer/artifact/by-commit-hash" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"code":500,"msg":"artifact not found by commit hash: root-hash"}`))
	})

	exists, err := service.CheckExistsByCommitHash("root-hash", &models.ArtifactLookupOptions{
		ArtifactType: "app",
	})
	if err != nil {
		t.Fatalf("CheckExistsByCommitHash error: %v", err)
	}
	if exists {
		t.Fatal("expected missing artifact to return false")
	}
}

func TestCheckExistsPrepareDownloadAndVersionMetadata(t *testing.T) {
	service := newArtifactTestService(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/aiplorer/artifact":
			payload := decodeBody(t, r)
			if payload["id"].(float64) != 12 {
				t.Fatalf("unexpected artifact id payload: %#v", payload)
			}
			_, _ = w.Write([]byte(`{"code":0,"data":{"id":12,"name":"artifact-a","type":"pkg","commitHash":"root-hash","projectName":"proj-a","fileHash":"abc123","fullPath":"repo/path/artifact.zip"}}`))
		case "/aiplorer/artifact/by-commit-hash":
			_, _ = w.Write([]byte(`{"code":0,"data":{"id":12,"name":"artifact-a","type":"pkg","commitHash":"root-hash","projectName":"proj-a","fileHash":"abc123","fullPath":"repo/path/artifact.zip"}}`))
		case "/aiplorer/jfrog/token":
			_, _ = w.Write([]byte(`{"code":0,"data":{"token_id":"tid","access_token":"token","expires_in":3600,"token_type":"Bearer","scope":"scope","url":"https://jfrog.example.com"}}`))
		case "/aiplorer/artifact/download-url":
			_, _ = w.Write([]byte(`{"code":0,"data":{"downloadUrl":"https://jfrog.example.com/file","expireTime":"2026-03-18 12:00:00","fileName":"artifact.zip","filePath":"repo/path/artifact.zip"}}`))
		case "/aiplorer/artifact/version-metadata":
			_, _ = w.Write([]byte(`{"code":0,"data":{"artifactId":12,"commitHash":"root-hash","name":"artifact-a","metadataPath":"repo/path/version.json","metadataFileName":"version.json","rawContent":"{\"version\":\"1.2.3\"}"}}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	})
	service.downloadArtifact = func(token *models.JfrogTokenInfo, filePath, targetDir string) error {
		if token.AccessToken != "token" || filePath != "repo/path/artifact.zip" {
			t.Fatalf("unexpected download args: %#v %s %s", token, filePath, targetDir)
		}
		return nil
	}

	exists, err := service.CheckExistsByCommitHash("root-hash", &models.ArtifactLookupOptions{ArtifactType: "pkg"})
	if err != nil || !exists {
		t.Fatalf("CheckExistsByCommitHash = %v, %v", exists, err)
	}

	plan, err := service.PrepareDownloadByCommitHash("root-hash", &models.ArtifactLookupOptions{ArtifactType: "pkg"}, "downloads/")
	if err != nil {
		t.Fatalf("PrepareDownloadByCommitHash error: %v", err)
	}
	if plan.TargetPath != "downloads/artifact.zip" {
		t.Fatalf("unexpected target path: %#v", plan)
	}

	plan, err = service.PrepareDownloadByArtifactID(12, "downloads/")
	if err != nil {
		t.Fatalf("PrepareDownloadByArtifactID error: %v", err)
	}
	if plan.TargetPath != "downloads/artifact.zip" {
		t.Fatalf("unexpected artifact-id target path: %#v", plan)
	}

	plan, err = service.DownloadByCommitHash("root-hash", &models.ArtifactLookupOptions{ArtifactType: "pkg"}, "downloads/")
	if err != nil {
		t.Fatalf("DownloadByCommitHash error: %v", err)
	}
	if plan.DownloadURL == nil || plan.DownloadURL.FileName != "artifact.zip" {
		t.Fatalf("unexpected download plan: %#v", plan)
	}

	plan, err = service.DownloadByArtifactID(12, "downloads/")
	if err != nil {
		t.Fatalf("DownloadByArtifactID error: %v", err)
	}
	if plan.DownloadURL == nil || plan.DownloadURL.FileName != "artifact.zip" {
		t.Fatalf("unexpected artifact-id download plan: %#v", plan)
	}

	metadata, err := service.GetVersionMetadataByCommitHash("root-hash", &models.ArtifactLookupOptions{ArtifactType: "pkg"})
	if err != nil {
		t.Fatalf("GetVersionMetadataByCommitHash error: %v", err)
	}
	if metadata.Parsed["version"].(string) != "1.2.3" {
		t.Fatalf("unexpected parsed metadata: %#v", metadata.Parsed)
	}
}

func TestDownloadByCommitHashSkipsExistingFileWithMatchingHash(t *testing.T) {
	service := newArtifactTestService(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/aiplorer/artifact/by-commit-hash":
			_, _ = w.Write([]byte(`{"code":0,"data":{"id":12,"name":"artifact-a","type":"pkg","commitHash":"root-hash","projectName":"proj-a","fileHash":"5d41402abc4b2a76b9719d911017c592","fullPath":"repo/path/artifact.zip"}}`))
		case "/aiplorer/jfrog/token":
			_, _ = w.Write([]byte(`{"code":0,"data":{"token_id":"tid","access_token":"token","expires_in":3600,"token_type":"Bearer","scope":"scope","url":"https://jfrog.example.com"}}`))
		case "/aiplorer/artifact/download-url":
			_, _ = w.Write([]byte(`{"code":0,"data":{"downloadUrl":"https://jfrog.example.com/file","expireTime":"2026-03-18 12:00:00","fileName":"artifact.zip","filePath":"repo/path/artifact.zip"}}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	})
	called := false
	service.downloadArtifact = func(token *models.JfrogTokenInfo, filePath, targetDir string) error {
		called = true
		return nil
	}

	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "artifact.zip")
	if err := os.WriteFile(targetPath, []byte("hello"), 0o644); err != nil {
		t.Fatalf("write existing file: %v", err)
	}

	plan, err := service.DownloadByCommitHash("root-hash", &models.ArtifactLookupOptions{ArtifactType: "pkg"}, tempDir)
	if err != nil {
		t.Fatalf("DownloadByCommitHash error: %v", err)
	}
	if called {
		t.Fatal("expected downloadArtifact not to be called")
	}
	if !plan.SkippedExisting {
		t.Fatalf("expected skipped existing plan: %#v", plan)
	}
}

func TestGetArtifactTagSchemaJSONAndUpdateTags(t *testing.T) {
	service := newArtifactTestService(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/aiplorer/artifact/tag-schema":
			_, _ = w.Write([]byte(`{"code":0,"data":{"version":"0.2.0","content":"{\"version\":\"0.2.0\",\"type\":\"object\"}"}}`))
		case "/aiplorer/artifact/update":
			payload := decodeBody(t, r)
			if payload["id"].(float64) != 12 {
				t.Fatalf("unexpected update payload: %#v", payload)
			}
			if !strings.Contains(payload["tags"].(string), `"schema_version":"0.2.0"`) {
				t.Fatalf("unexpected tags payload: %#v", payload)
			}
			_, _ = w.Write([]byte(`{"code":0,"msg":"updated"}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	})

	schema, err := service.GetArtifactTagSchemaJSON("0.2.0")
	if err != nil {
		t.Fatalf("GetArtifactTagSchemaJSON error: %v", err)
	}
	if schema["version"].(string) != "0.2.0" {
		t.Fatalf("unexpected schema: %#v", schema)
	}

	result, err := service.UpdateArtifactTags(12, map[string]any{
		"schema_version": "0.2.0",
		"decision":       "pass",
	}, "")
	if err != nil {
		t.Fatalf("UpdateArtifactTags error: %v", err)
	}
	if result.Msg != "updated" {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestGetJfrogTokenAndDownloadURL(t *testing.T) {
	service := newArtifactTestService(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/aiplorer/jfrog/token":
			_, _ = w.Write([]byte(`{"code":0,"data":{"token_id":"tid","access_token":"token","expires_in":3600,"token_type":"Bearer","scope":"scope","url":"https://jfrog.example.com"}}`))
		case "/aiplorer/artifact/download-url":
			payload := decodeBody(t, r)
			if payload["downloadType"].(string) != "artifact" {
				t.Fatalf("unexpected download payload: %#v", payload)
			}
			_, _ = w.Write([]byte(`{"code":0,"data":{"downloadUrl":"https://jfrog.example.com/file","expireTime":"2026-03-18 12:00:00","fileName":"artifact.zip","filePath":"repo/path/artifact.zip"}}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	})

	token, err := service.GetJfrogToken("proj-a")
	if err != nil {
		t.Fatalf("GetJfrogToken error: %v", err)
	}
	if token.AccessToken != "token" {
		t.Fatalf("unexpected token: %#v", token)
	}

	download, err := service.GetArtifactDownloadURL(12, "")
	if err != nil {
		t.Fatalf("GetArtifactDownloadURL error: %v", err)
	}
	if download.FileName != "artifact.zip" {
		t.Fatalf("unexpected download: %#v", download)
	}
}

func TestDownloadByNameUsesResolvedArtifactID(t *testing.T) {
	service := newArtifactTestService(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/aiplorer/artifact/list":
			_, _ = w.Write([]byte(`{"code":0,"data":{"total":1,"data":[{"id":12,"name":"artifact-a"}]}}`))
		case "/aiplorer/artifact":
			payload := decodeBody(t, r)
			if payload["id"].(float64) != 12 {
				t.Fatalf("unexpected artifact id payload: %#v", payload)
			}
			_, _ = w.Write([]byte(`{"code":0,"data":{"id":12,"name":"artifact-a","projectName":"proj-a","fileHash":"abc123","fullPath":"repo/path/artifact.zip"}}`))
		case "/aiplorer/jfrog/token":
			_, _ = w.Write([]byte(`{"code":0,"data":{"token_id":"tid","access_token":"token","expires_in":3600,"token_type":"Bearer","scope":"scope","url":"https://jfrog.example.com"}}`))
		case "/aiplorer/artifact/download-url":
			payload := decodeBody(t, r)
			if payload["artifactId"].(float64) != 12 {
				t.Fatalf("unexpected download payload: %#v", payload)
			}
			_, _ = w.Write([]byte(`{"code":0,"data":{"downloadUrl":"https://jfrog.example.com/file","expireTime":"2026-03-18 12:00:00","fileName":"artifact.zip","filePath":"repo/path/artifact.zip"}}`))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	})
	service.downloadArtifact = func(token *models.JfrogTokenInfo, filePath, targetDir string) error {
		if token.AccessToken != "token" || filePath != "repo/path/artifact.zip" {
			t.Fatalf("unexpected download args: %#v %s %s", token, filePath, targetDir)
		}
		return nil
	}

	plan, err := service.DownloadByName("artifact-a", nil, "downloads/")
	if err != nil {
		t.Fatalf("DownloadByName error: %v", err)
	}
	if plan.DownloadURL == nil || plan.DownloadURL.FileName != "artifact.zip" {
		t.Fatalf("unexpected download plan: %#v", plan)
	}
}
