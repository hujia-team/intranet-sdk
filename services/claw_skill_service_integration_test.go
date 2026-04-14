package services

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hujia-team/intranet-sdk/client"
	"github.com/hujia-team/intranet-sdk/models"
)

func TestClawSkillReportPrivateSkillHubEventAgainstLocalService(t *testing.T) {
	if os.Getenv("RUN_LOCAL_PRIVATE_SKILL_HUB_INTEGRATION") != "1" {
		t.Skip("set RUN_LOCAL_PRIVATE_SKILL_HUB_INTEGRATION=1 to run against local service")
	}

	accessKeyID := strings.TrimSpace(os.Getenv("INTRANET_ACCESS_KEY_ID"))
	accessKeySecret := strings.TrimSpace(os.Getenv("INTRANET_ACCESS_KEY_SECRET"))
	if accessKeyID == "" || accessKeySecret == "" {
		t.Fatal("INTRANET_ACCESS_KEY_ID and INTRANET_ACCESS_KEY_SECRET are required")
	}

	baseURL := strings.TrimSpace(os.Getenv("LOCAL_CORE_API_BASE_URL"))
	if baseURL == "" {
		baseURL = "http://127.0.0.1:9100"
	}

	httpClient, err := client.NewHTTPClient(&client.Config{
		BaseURL:         baseURL,
		AccessKeyID:     accessKeyID,
		AccessKeySecret: accessKeySecret,
	})
	if err != nil {
		t.Fatalf("new http client: %v", err)
	}

	service := NewClawSkillService(httpClient).(*clawSkillService)
	skillName := "sdk-local-verify-" + time.Now().Format("20060102150405")
	clientName := "intranet-sdk-integration"

	reportResult, err := service.ReportPrivateSkillHubEvent(
		&models.PrivateSkillHubEventReportRequest{
			SkillName:  skillName,
			Action:     "install",
			Success:    true,
			ClientName: &clientName,
		},
		nil,
	)
	if err != nil {
		t.Fatalf("report event: %v", err)
	}
	if reportResult.StatusCode == 0 {
		t.Fatalf("unexpected report result: %#v", reportResult)
	}

	var statsResp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Total uint64 `json:"total"`
			Items []struct {
				SkillName    *string `json:"skillName"`
				InstallCount *uint64 `json:"installCount"`
			} `json:"items"`
		} `json:"data"`
	}
	keyword := skillName
	if err := httpClient.Post("/claw/skill/private-hub/stats/list", map[string]any{
		"page":     1,
		"pageSize": 10,
		"keyword":  &keyword,
	}, &statsResp); err != nil {
		t.Fatalf("query stats: %v", err)
	}

	found := false
	for _, item := range statsResp.Data.Items {
		if item.SkillName != nil && *item.SkillName == skillName {
			found = true
			if item.InstallCount == nil || *item.InstallCount == 0 {
				t.Fatalf("expected install count for %s, got %#v", skillName, item.InstallCount)
			}
			break
		}
	}
	if !found {
		t.Fatalf("skill %s not found in stats response: %#v", skillName, statsResp.Data.Items)
	}
}
