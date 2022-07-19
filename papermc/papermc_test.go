//go:build integration

package papermc_test

import (
	"net/http"
	"testing"
	"time"

	"nw.codes/papermc-go/papermc"
)

func TestGetProjects(t *testing.T) {
	client := http.Client{Timeout: 10 * time.Second}

	pmc := papermc.NewPaperMCClient(&client)

	projects, err := pmc.GetProjects()
	if err != nil {
		t.Error(err)
	}

	if len(projects.Projects) == 0 {
		t.Error("no projects found")
	}

	t.Log(projects.Projects)
}

func TestGetProject(t *testing.T) {
	client := http.Client{Timeout: 10 * time.Second}

	pmc := papermc.NewPaperMCClient(&client)

	projectID := "paper"

	versions, err := pmc.GetProject(projectID)
	if err != nil {
		t.Error(err)
	}

	t.Log(versions)
}

func TestGetVersion(t *testing.T) {
	client := http.Client{Timeout: 10 * time.Second}

	pmc := papermc.NewPaperMCClient(&client)

	projectID := "paper"
	versionID := "1.19"

	builds, err := pmc.GetVersion(projectID, versionID)
	if err != nil {
		t.Error(err)
	}

	t.Log(builds)
}

func TestGetBuild(t *testing.T) {
	client := http.Client{Timeout: 10 * time.Second}

	pmc := papermc.NewPaperMCClient(&client)

	projectID := "paper"
	versionID := "1.19"
	buildID := 66

	build, err := pmc.GetBuild(projectID, versionID, buildID)
	if err != nil {
		t.Error(err)
	}

	t.Log(build)
}

func TestGetDownloadLink(t *testing.T) {
	projectID := "paper"
	versionID := "1.19"
	buildID := 66
	downloadID := "paper-1.19-66.jar"

	downloadLink := papermc.GetDownloadLink(projectID, versionID, buildID, downloadID)

	t.Log(downloadLink)
}

func TestDownload(t *testing.T) {
	url := "https://api.papermc.io/v2/projects/paper/versions/1.19/builds/66/downloads/paper-1.19-66.jar"
	hash := "b1e10c426f9fb70ff2ff382b4071df920614981100eed63b0b31ef4c73eaedca"
	filename := "paper-1.19-66.jar"

	client := http.Client{Timeout: 10 * time.Second}

	pmc := papermc.NewPaperMCClient(&client)

	err := pmc.Download(url, hash, filename)
	if err != nil {
		t.Error(err)
	}
}
