package intervals

import (
	"time"
)

// TimeService is a service that does the stuff for service
type TimeService service

// TimeResponse is the struct the stores the response
type TimeResponse struct {
	*baseResponse
	Time []Time `json:"time,omitempty"`
}

// Time is cool
type Time struct {
	ID           string `json:"id,omitempty"`
	ProjectID    string `json:"projectid,omitempty"`
	ModuleID     string `json:"moduleid,omitempty"`
	TaskID       string `json:"taskid,omitempty"`
	WorkTypeID   string `json:"worktypeid,omitempty"`
	PersonID     string `json:"personid,omitempty"`
	Date         string `json:"date,omitempty"`
	Time         string `json:"time,omitempty"`
	Description  string `json:"description,omitempty"`
	Billable     string `json:"billable,omitempty"`
	DateModified string `json:"datemodified,omitempty"`
	DateISO      string `json:"dateiso,omitempty"`
	Module       string `json:"module,omitempty"`
	Project      string `json:"project,omitempty"`
	WorkType     string `json:"worktype,omitempty"`
	FirstName    string `json:"firstname,omitempty"`
	LastName     string `json:"lastname,omitempty"`
	Active       string `json:"active,omitempty"`
	ClientID     string `json:"clientid,omitempty"`
	Client       string `json:"client,omitempty"`
	ClientActive string `json:"clientactive,omitempty"`
	StatusID     string `json:"statusid,omitempty"`
}

// TimeOptions set of options for Time
type TimeOptions struct {
	Limit       int       `url:"limit,omitempty"`
	DateBegin   time.Time `url:"date_begin,omitempty"`
	DateEnd     time.Time `url:"date_end,omitempty"`
	PersonID    int       `url:"person_id,omitempty"`
	ProjectID   int       `url:"project_id,omitempty"`
	ClientID    int       `url:"client_id,omitempty"`
	MilestoneID int       `url:"milestone_id,omitempty"`
}

// List all time for the options passed in
func (s *TimeService) List(opt *TimeOptions) ([]Time, error) {
	u, err := addOptions("time", opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var response *TimeResponse
	_, err = s.client.Do(req, &response)

	return response.Time, err
}
