package ccv2

import (
	"encoding/json"
	"net/url"

	"code.cloudfoundry.org/cli/api/cloudcontroller"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv2/internal"
)

// JobStatus is the current state of a job.
type JobStatus string

const (
	// JobStatusFailed is when the job is no longer running due to a failure.
	JobStatusFailed JobStatus = "failed"

	// JobStatusFinished is when the job is no longer and it was successful.
	JobStatusFinished JobStatus = "finished"

	// JobStatusQueued is when the job is waiting to be run.
	JobStatusQueued JobStatus = "queued"

	// JobStatusRunning is when the job is running.
	JobStatusRunning JobStatus = "running"
)

// Job represents a Cloud Controller Job.
type Job struct {
	Error  string
	GUID   string
	Status JobStatus
}

// UnmarshalJSON helps unmarshal a Cloud Controller Job response.
func (job *Job) UnmarshalJSON(data []byte) error {
	var ccJob struct {
		Entity struct {
			Error  string `json:"error"`
			GUID   string `json:"guid"`
			Status string `json:"status"`
		} `json:"entity"`
		Metadata internal.Metadata `json:"metadata"`
	}
	if err := json.Unmarshal(data, &ccJob); err != nil {
		return err
	}

	job.Error = ccJob.Entity.Error
	job.GUID = ccJob.Entity.GUID
	job.Status = JobStatus(ccJob.Entity.Status)
	return nil
}

func (job Job) Terminated() bool { // TODO: Add Test
	return job.Status == JobStatusFinished || job.Status == JobStatusFailed
}

func (job Job) Failed() bool { // TODO: Add Test
	return job.Status == JobStatusFailed
}

// GetJob returns a job for the provided GUID.
func (client *Client) GetJob(jobGUID string) (Job, Warnings, error) {
	request, err := client.newHTTPRequest(requestOptions{
		RequestName: internal.JobRequest,
		URIParams:   map[string]string{"job_guid": jobGUID},
	})
	if err != nil {
		return Job{}, nil, err
	}

	var job Job
	response := cloudcontroller.Response{
		Result: &job,
	}

	err = client.connection.Make(request, &response)
	return job, response.Warnings, err
}

// DeleteOrganization deletes the Organization associated with the provided
// GUID. It will return the Cloud Controller job that is assigned to the
// organization deletion.
func (client *Client) DeleteOrganization(orgGUID string) (Job, Warnings, error) {
	request, err := client.newHTTPRequest(requestOptions{
		RequestName: internal.DeleteOrganizationRequest,
		URIParams:   map[string]string{"organization_guid": orgGUID},
		Query: url.Values{
			"recursive": {"true"},
			"async":     {"true"},
		},
	})
	if err != nil {
		return Job{}, nil, err
	}

	var job Job
	response := cloudcontroller.Response{
		Result: &job,
	}

	err = client.connection.Make(request, &response)
	return job, response.Warnings, err
}
