package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	intranet "github.com/hujia-team/intranet-sdk"
	"github.com/hujia-team/intranet-sdk/models"
)

const (
	baseURL         = "https://intranet.minieye.tech/sys-api"
	accessKeyID     = "aiplorer"
	accessKeySecret = "0f3ba4a3396aafdaaaccb283de362166"
	rootCommitHash  = "45dc8b12ce68b4e7a1e27418465e2636"
	rootType        = "app"
)

func main() {
	client, err := intranet.NewClient(
		intranet.WithBaseURL(baseURL),
		intranet.WithAccessKeyID(accessKeyID),
		intranet.WithAccessKeySecret(accessKeySecret),
	)
	if err != nil {
		log.Fatalf("init sdk failed: %v", err)
	}

	lookup := &models.ArtifactLookupOptions{
		ArtifactType: rootType,
	}

	rootArtifact, err := client.Artifact.GetArtifactByCommitHash(rootCommitHash, lookup)
	if err != nil {
		log.Fatalf("load root artifact failed: %v", err)
	}

	downloadRoot := filepath.Join("downloads", rootCommitHash, "root")
	if rootArtifact.ID == nil {
		log.Fatal("root artifact id is empty")
	}
	rootPlan, err := client.Artifact.DownloadByArtifactID(*rootArtifact.ID, downloadRoot)
	if err != nil {
		log.Fatalf("download root artifact failed: %v", err)
	}

	msgChild, err := findMsgChild(client, rootArtifact)
	if err != nil {
		log.Fatalf("find msg child failed: %v", err)
	}
	if msgChild.ID == nil {
		log.Fatal("msg child id is empty")
	}

	downloadMsg := filepath.Join("downloads", rootCommitHash, "msg")
	msgPlan, err := client.Artifact.DownloadByArtifactID(*msgChild.ID, downloadMsg)
	if err != nil {
		log.Fatalf("download msg artifact failed: %v", err)
	}

	fmt.Printf("root artifact downloaded: %s\n", rootPlan.TargetPath)
	fmt.Printf("msg artifact name: %s\n", value(msgChild.Name))
	fmt.Printf("msg artifact commit hash: %s\n", value(msgChild.CommitHash))
	fmt.Printf("msg artifact downloaded: %s\n", msgPlan.TargetPath)
}

func findMsgChild(client *intranet.Client, root *models.ArtifactInfo) (*models.ArtifactInfo, error) {
	if root == nil {
		return nil, fmt.Errorf("root artifact is nil")
	}
	if root.ProjectName == nil || strings.TrimSpace(*root.ProjectName) == "" {
		return nil, fmt.Errorf("root project name is empty")
	}

	wantName := strings.TrimSpace(*root.ProjectName) + "-msg"
	wantPlatform := strings.TrimSpace(value(root.Platform))

	for i := range root.Dependencies {
		dep := root.Dependencies[i]
		if dep.ID == nil {
			continue
		}
		if strings.TrimSpace(value(dep.Name)) != wantName {
			continue
		}
		if strings.TrimSpace(value(dep.Type)) != "app" {
			continue
		}
		if strings.TrimSpace(value(dep.ModulePath)) != "" {
			continue
		}
		if dep.IsVirtual != nil && *dep.IsVirtual {
			continue
		}
		if wantPlatform != "" && strings.TrimSpace(value(dep.Platform)) != wantPlatform {
			continue
		}

		detail, err := client.Artifact.GetArtifactByID(*dep.ID)
		if err != nil {
			return nil, fmt.Errorf("load msg dependency detail by id=%d: %w", *dep.ID, err)
		}
		ok, err := isStableMsgArtifact(detail, strings.TrimSpace(*root.ProjectName), wantPlatform)
		if err != nil {
			return nil, fmt.Errorf("validate msg dependency detail by id=%d: %w", *dep.ID, err)
		}
		if !ok {
			continue
		}
		return detail, nil
	}
	return nil, fmt.Errorf("stable msg child not found")
}

func value(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func isStableMsgArtifact(artifact *models.ArtifactInfo, projectName, platform string) (bool, error) {
	if artifact == nil {
		return false, nil
	}

	extra := strings.TrimSpace(value(artifact.Extra))
	if extra == "" {
		return false, nil
	}

	var payload map[string]string
	if err := json.Unmarshal([]byte(extra), &payload); err != nil {
		return false, err
	}

	if strings.TrimSpace(payload["artifact_type"]) != "msg" {
		return false, nil
	}
	if projectName != "" && strings.TrimSpace(payload["project_name"]) != projectName {
		return false, nil
	}
	if platform != "" && strings.TrimSpace(payload["platform"]) != platform {
		return false, nil
	}
	return true, nil
}

func init() {
	if err := os.MkdirAll(filepath.Join("downloads", rootCommitHash), 0o755); err != nil {
		log.Fatalf("prepare download directory failed: %v", err)
	}
}
