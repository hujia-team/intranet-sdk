package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hujia-team/intranet-sdk/client"
	"github.com/hujia-team/intranet-sdk/models"
)

func newArtifactTestService(t *testing.T, handler http.HandlerFunc) ArtifactService {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	httpClient, err := client.NewHTTPClient(&client.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatalf("new http client: %v", err)
	}
	return NewArtifactService(httpClient)
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
