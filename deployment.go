package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

//DeploymentRequest contains new deployment request data
type DeploymentRequest struct {
	Ref               string `json:"ref"`
	RefType           string `json:"ref_type"`
	Title             string `json:"title"`
	Summary           string `json:"summary"`
	Bypass            bool   `json:"bypass"`
	BypassAndStart    bool   `json:"bypass_and_start"`
	ScheduleStartUnix int    `json:"schedule_start_unix"`
	ScheduleEndUnix   int    `json:"schedule_end_unix"`
	Locked            bool   `json:"locked"`
}

//CreateDeploymentResponse contains response data on a new deployment request
type CreateDeploymentResponse struct {
	Data struct {
		Type       string `json:"type"`
		ID         string `json:"id"`
		Attributes struct {
			ID                int    `json:"id"`
			DateCreatedUnix   int    `json:"date_created_unix"`
			DateStartedUnix   int    `json:"date_started_unix"`
			DateRequestedUnix int    `json:"date_requested_unix"`
			DateApprovedUnix  int    `json:"date_approved_unix"`
			DateUpdatedUnix   int    `json:"date_updated_unix"`
			Title             string `json:"title"`
			Summary           string `json:"summary"`
			Changes           struct {
				CodeVersion struct {
					From string `json:"from"`
					To   string `json:"to"`
				} `json:"Code version"`
			} `json:"changes"`
			DeploymentType     string        `json:"deployment_type"`
			DeploymentEstimate string        `json:"deployment_estimate"`
			Sha                string        `json:"sha"`
			ShortSha           string        `json:"short_sha"`
			CommitSubject      string        `json:"commit_subject"`
			CommitMessage      string        `json:"commit_message"`
			CommitURL          string        `json:"commit_url"`
			Decision           string        `json:"decision"`
			Options            []interface{} `json:"options"`
			Messages           []struct {
				Text string `json:"text"`
				Code string `json:"code"`
			} `json:"messages"`
			Deployer struct {
				ID    int    `json:"id"`
				Email string `json:"email"`
				Role  string `json:"role"`
				Name  string `json:"name"`
			} `json:"deployer"`
			Approver interface{} `json:"approver"`
			Bypasser struct {
				ID    int    `json:"id"`
				Email string `json:"email"`
				Role  string `json:"role"`
				Name  string `json:"name"`
			} `json:"bypasser"`
			State          string `json:"state"`
			IsCurrentBuild bool   `json:"is_current_build"`
		} `json:"attributes"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	} `json:"data"`
}

type deploymentStatus struct {
	Data struct {
		Type       string `json:"type"`
		ID         int    `json:"id"`
		Attributes struct {
			ID                 int           `json:"id"`
			DateCreatedUnix    int           `json:"date_created_unix"`
			DateStartedUnix    int           `json:"date_started_unix"`
			DateRequestedUnix  int           `json:"date_requested_unix"`
			DateApprovedUnix   int           `json:"date_approved_unix"`
			DateUpdatedUnix    int           `json:"date_updated_unix"`
			Title              string        `json:"title"`
			Summary            string        `json:"summary"`
			DeploymentType     string        `json:"deployment_type"`
			DeploymentEstimate string        `json:"deployment_estimate"`
			Sha                string        `json:"sha"`
			ShortSha           string        `json:"short_sha"`
			CommitSubject      string        `json:"commit_subject"`
			CommitMessage      string        `json:"commit_message"`
			CommitURL          string        `json:"commit_url"`
			Decision           string        `json:"decision"`
			Options            []interface{} `json:"options"`
			Messages           []struct {
				Text string `json:"text"`
				Code string `json:"code"`
			} `json:"messages"`
			Deployer struct {
				ID    int    `json:"id"`
				Email string `json:"email"`
				Role  string `json:"role"`
				Name  string `json:"name"`
			} `json:"deployer"`
			Approver struct {
				ID    int    `json:"id"`
				Email string `json:"email"`
				Role  string `json:"role"`
				Name  string `json:"name"`
			} `json:"approver"`
			Bypasser       interface{} `json:"bypasser"`
			State          string      `json:"state"`
			IsCurrentBuild bool        `json:"is_current_build"`
		} `json:"attributes"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	} `json:"data"`
}

func (cdr *DeploymentRequest) toJSONBody() io.Reader {
	s, err := json.Marshal(cdr)
	check(err)
	return bytes.NewBuffer(s)
}

// Send new deployment request
func (cdr *DeploymentRequest) Send(c *Config, envID string, projectID string) {
	url := strings.Trim(c.BaseURL, "/") + "/project/" + projectID + "/environment/" + envID + "/deploys"
	fmt.Println("Create deployment...:", url)

	response, err := HttpCall(c.Email, c.APIToken, url, "POST", cdr.toJSONBody())
	check(err)

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	check(err)

	var res CreateDeploymentResponse
	err = json.Unmarshal(body, &res)
	check(err)

	timeout := time.Now().Unix() + 31*60

	if res.Data.ID != "0" {
		fmt.Println("Deployment created. Waiting for deployment finish...")
		for time.Now().Unix() < timeout {
			time.Sleep(5 * time.Second)
			if checkDeploymentProgress(c, &res, projectID) {
				fmt.Println("Deployment completed")
				return
			}
		}
		fmt.Println("Failed waiting for deployment finish. Timeout!")
		return
	}
	fmt.Println("Error deployment creation response. ID is invalid")
}

func checkDeploymentProgress(c *Config, res *CreateDeploymentResponse, projectID string) bool {
	if "0" == res.Data.ID {
		panic("Invalid ID")
	}
	fmt.Printf("\r%s Checking deployment...", time.Now().Format("15:04:05"))
	url := strings.Trim(c.BaseURL, "/") + "/project/" + projectID + "/environment/uat/deploys/" + res.Data.ID
	response, err := HttpCall(c.Email, c.APIToken, url, "GET", nil)
	check(err)

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	check(err)
	var result deploymentStatus
	json.Unmarshal(body, &result)
	if result.Data.Attributes.State == "Completed" {
		return true
	}
	return false
}
