package papermc

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	PaperMCProjectsURL      = "https://api.papermc.io/v2/projects"
	PaperMCVersionsURL      = "https://api.papermc.io/v2/projects/%s"
	PaperMCBuildsURL        = "https://api.papermc.io/v2/projects/%s/versions/%s"
	PaperMCBuildURL         = "https://api.papermc.io/v2/projects/%s/versions/%s/builds/%d"
	PaperMCBuildDownloadURL = "https://api.papermc.io/v2/projects/%s/versions/%s/builds/%d/downloads/%s"
)

type PaperMCClient struct {
	client *http.Client
}

func NewPaperMCClient(c *http.Client) *PaperMCClient {
	return &PaperMCClient{
		client: c,
	}
}

type Projects struct {
	Projects []string `json:"projects"`
}

func (pmcc *PaperMCClient) GetProjects() (*Projects, error) {
	req, err := http.NewRequest("GET", PaperMCProjectsURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := pmcc.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var projects Projects
	err = json.NewDecoder(resp.Body).Decode(&projects)
	if err != nil {
		return nil, err
	}

	return &projects, nil
}

type Project struct {
	ProjectID     string   `json:"project_id"`
	ProjectName   string   `json:"project_name"`
	VersionGroups []string `json:"version_groups"`
	Versions      []string `json:"versions"`
}

func (pmcc *PaperMCClient) GetProject(project string) (*Project, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(PaperMCVersionsURL, project), nil)
	if err != nil {
		return nil, err
	}

	resp, err := pmcc.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var p Project
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil

}

type Version struct {
	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
	Version     string `json:"version"`
	Builds      []int  `json:"builds"`
}

func (pmcc *PaperMCClient) GetVersion(project string, version string) (*Version, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(PaperMCBuildsURL, project, version), nil)
	if err != nil {
		return nil, err
	}

	resp, err := pmcc.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var v Version
	err = json.NewDecoder(resp.Body).Decode(&v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

type Build struct {
	ProjectID   string    `json:"project_id"`
	ProjectName string    `json:"project_name"`
	Version     string    `json:"version"`
	Build       int       `json:"build"`
	Time        time.Time `json:"time"`
	Channel     string    `json:"channel"`
	Promoted    bool      `json:"promoted"`
	Changes     []struct {
		Commit  string `json:"commit"`
		Summary string `json:"summary"`
		Message string `json:"message"`
	} `json:"changes"`
	Downloads struct {
		Application struct {
			Name   string `json:"name"`
			Sha256 string `json:"sha256"`
		} `json:"application"`
		MojangMappings struct {
			Name   string `json:"name"`
			Sha256 string `json:"sha256"`
		} `json:"mojang-mappings"`
	} `json:"downloads"`
}

func (pmcc *PaperMCClient) GetBuild(project string, version string, build int) (*Build, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(PaperMCBuildURL, project, version, build), nil)
	if err != nil {
		return nil, err
	}

	resp, err := pmcc.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var b Build
	err = json.NewDecoder(resp.Body).Decode(&b)
	if err != nil {
		return nil, err
	}

	return &b, nil
}

func GetDownloadLink(project string, version string, build int, name string) string {
	return fmt.Sprintf(PaperMCBuildDownloadURL, project, version, build, name)
}

func (pmcc *PaperMCClient) Download(link string, hash string, fileLoc string) error {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return err
	}

	resp, err := pmcc.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	downloadBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	shaHash := sha256.Sum256(downloadBytes)
	shaString := hex.EncodeToString(shaHash[:])
	if shaString != hash {
		return fmt.Errorf("SHA256 hash mismatch")
	}

	err = os.WriteFile(fileLoc, downloadBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
