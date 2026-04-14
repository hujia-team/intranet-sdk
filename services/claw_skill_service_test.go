package services

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hujia-team/intranet-sdk/client"
	"github.com/hujia-team/intranet-sdk/models"
)

func newClawSkillTestService(t *testing.T, transport http.RoundTripper) *clawSkillService {
	t.Helper()
	httpClient, err := client.NewHTTPClient(&client.Config{
		BaseURL:         "https://intranet.minieye.tech/sys-api",
		AccessKeyID:     "alice",
		AccessKeySecret: "secret-123",
		HTTPClient: &http.Client{
			Transport: transport,
		},
	})
	if err != nil {
		t.Fatalf("new http client: %v", err)
	}
	return NewClawSkillService(httpClient).(*clawSkillService)
}

func TestClawSkillUploadLocalSkill(t *testing.T) {
	service := newClawSkillTestService(t, roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://upload.example.com/claw/skill/local/upload" {
			t.Fatalf("unexpected url: %s", r.URL.String())
		}
		if got := r.Header.Get("x-sts-uid"); got != "alice" {
			t.Fatalf("unexpected sts uid: %s", got)
		}
		if got := r.Header.Get("X-Skill-Token"); got != "claw_skill_demo" {
			t.Fatalf("unexpected upload token header: %s", got)
		}
		if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Fatalf("unexpected content type: %s", r.Header.Get("Content-Type"))
		}
		if err := r.ParseMultipartForm(8 << 20); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		if got := r.FormValue("version"); got != "1.2.3" {
			t.Fatalf("unexpected version field: %s", got)
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("form file: %v", err)
		}
		defer file.Close()
		content, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("read uploaded file: %v", err)
		}
		if string(content) != "zip-bytes" {
			t.Fatalf("unexpected file content: %q", string(content))
		}
		return &http.Response{
			StatusCode: http.StatusCreated,
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(`{
				"code": 0,
				"msg": "success",
				"data": {
					"created": true,
					"skipped": false,
					"uploadToken": "claw_skill_demo",
					"skill": {
						"slug": "demo-skill",
						"version": "1.2.3"
					}
				}
			}`)),
		}, nil
	}))

	result, err := service.UploadLocalSkill(
		"https://upload.example.com/claw/skill/local/upload",
		"demo-skill.zip",
		[]byte("zip-bytes"),
		"1.2.3",
		"claw_skill_demo",
		nil,
	)
	if err != nil {
		t.Fatalf("UploadLocalSkill error: %v", err)
	}
	if result.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected status code: %d", result.StatusCode)
	}
	if result.Parsed == nil || result.Parsed.Data.Skill.Slug == nil || *result.Parsed.Data.Skill.Slug != "demo-skill" {
		t.Fatalf("unexpected parsed result: %#v", result)
	}
}

func TestClawSkillResetLocalSkillUploadToken(t *testing.T) {
	service := newClawSkillTestService(t, roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://upload.example.com/claw/skill/local/token/reset" {
			t.Fatalf("unexpected url: %s", r.URL.String())
		}
		if got := r.Header.Get("x-sts-uid"); got != "alice" {
			t.Fatalf("unexpected sts uid: %s", got)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if !strings.Contains(string(body), `"slug":"demo-skill"`) {
			t.Fatalf("unexpected body: %s", string(body))
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(`{
				"code": 0,
				"msg": "success",
				"data": {
					"uploadToken": "claw_skill_reset",
					"skill": {
						"slug": "demo-skill"
					}
				}
			}`)),
		}, nil
	}))

	result, err := service.ResetLocalSkillUploadToken(
		"https://upload.example.com/claw/skill/local/token/reset",
		"demo-skill",
		nil,
	)
	if err != nil {
		t.Fatalf("ResetLocalSkillUploadToken error: %v", err)
	}
	if result.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", result.StatusCode)
	}
	if result.Parsed == nil || result.Parsed.Data.UploadToken != "claw_skill_reset" {
		t.Fatalf("unexpected parsed result: %#v", result)
	}
}

func TestClawSkillReportPrivateSkillHubEvent(t *testing.T) {
	service := newClawSkillTestService(t, roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() != "https://upload.example.com/claw/skill/private-hub/event/report" {
			t.Fatalf("unexpected url: %s", r.URL.String())
		}
		if got := r.Header.Get("x-sts-uid"); got != "alice" {
			t.Fatalf("unexpected sts uid: %s", got)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if !strings.Contains(string(body), `"skillName":"common.skill-hub.publisher"`) {
			t.Fatalf("unexpected body: %s", string(body))
		}
		if !strings.Contains(string(body), `"action":"install"`) {
			t.Fatalf("unexpected body: %s", string(body))
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"code":0,"msg":"success"}`)),
		}, nil
	}))

	result, err := service.ReportPrivateSkillHubEvent(
		"https://upload.example.com/claw/skill/private-hub/event/report",
		&models.PrivateSkillHubEventReportRequest{
			SkillName: "common.skill-hub.publisher",
			Action:    "install",
			Success:   true,
		},
		nil,
	)
	if err != nil {
		t.Fatalf("ReportPrivateSkillHubEvent error: %v", err)
	}
	if result.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", result.StatusCode)
	}
	if result.Parsed == nil || result.Parsed.Msg != "success" {
		t.Fatalf("unexpected parsed result: %#v", result)
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
