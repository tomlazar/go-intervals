package intervals

type ProjectService service

type Project struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Datestart       string `json:"datestart"`
	Dateend         string `json:"dateend"`
	AlertPercent    string `json:"alert_percent"`
	AlertDate       string `json:"alert_date"`
	Active          string `json:"active"`
	Billable        string `json:"billable"`
	Budget          string `json:"budget"`
	Clientid        string `json:"clientid"`
	Client          string `json:"client"`
	Clientlocalid   string `json:"clientlocalid"`
	Localidunpadded string `json:"localidunpadded"`
	Localid         string `json:"localid"`
	Manager         string `json:"manager"`
	Managerid       string `json:"managerid"`
}

type ProjectResult struct {
	*baseResponse
	Projects []Project `json:"project"`
}

func (ps *ProjectService) List() ([]Project, error) {
	u, err := addOptions("project", nil)
	if err != nil {
		return nil, err
	}

	req, err := ps.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var personResponse *ProjectResult
	_, err = ps.client.Do(req, &personResponse)

	return personResponse.Projects, err
}
