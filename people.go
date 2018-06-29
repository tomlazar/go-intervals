package intervals

// PersonService is cool
type PersonService service

// PersonResponse the collection of people requtned from one request
type PersonResponse struct {
	*baseResponse
	People []Person `json:"person"`
}

// Person is cool
type Person struct {
	ID             string `json:"id"`
	LocalID        int    `json:"localid"`
	ClientID       int    `json:"clientid"`
	Title          string `json:"title"`
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	PrimaryAccount string `json:"primaryaccount"`
	Notes          string `json:"notes"`
	AllProjects    bool   `json:"allprojects"`
	Active         bool   `json:"active"`
	Private        bool   `json:"private"`
	Tips           string `json:"tips"`
	Username       string `json:"username"`
	GroupID        int    `json:"groupid"`
	Group          string `json:"group"`
	Client         string `json:"client"`
	NumLogins      int    `json:"numlogins"`
	LastLogin      string `json:"lastlogin"`
}

// PersonOptions is the options that can be passed into the API
type PersonOptions struct {
	Limit  int    `url:"limit,omitempty"`
	Search string `url:"search,omitempty"`
	Email  string `url:"email,omitempty"`
}

// List data
func (s *PersonService) List(opt *PersonOptions) ([]Person, error) {
	if opt == nil {
		opt = &PersonOptions{
			Limit: 100,
		}
	}

	u, err := addOptions("person", opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var personResponse *PersonResponse
	_, err = s.client.Do(req, &personResponse)

	return personResponse.People, err
}
