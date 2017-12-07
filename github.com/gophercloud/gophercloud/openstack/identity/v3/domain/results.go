package domain

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"

)

type domainResult struct {
	gophercloud.Result
}

// GetResult temporarily contains the response from the Get call.
type GetResult struct {
	domainResult
}

// CreateResult temporarily contains the reponse from the Create call.
type CreateResult struct {
	domainResult
}

// DeleteResult temporarily contains the response from the Delete call.
type DeleteResult struct {
	gophercloud.ErrResult
}

// UpdateResult temporarily contains the response from the Update call.
type UpdateResult struct {
	domainResult
}

// Project is a base unit of ownership.
type Domain struct {
	// Description is the description of the project.
	Description string `json:"description"`

	// Enabled is whether or not the project is enabled.
	Enabled bool `json:"enabled"`

	// ID is the unique ID of the project.
	ID string `json:"id"`

	// Name is the name of the project.
	Name string `json:"name"`

	// links
	Links map[string]string `json:"links"`

}

// ProjectPage is a single page of Project results.
type DomainPage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines whether or not a page of Projects contains any results.
func (r DomainPage) IsEmpty() (bool, error) {
	projects, err := ExtractDomains(r)
	return len(projects) == 0, err
}

// NextPageURL extracts the "next" link from the links section of the result.
func (r DomainPage) NextPageURL() (string, error) {
	var s struct {
		Links struct {
			Next     string `json:"next"`
			Previous string `json:"previous"`
		} `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return s.Links.Next, err
}

// ExtractProjects returns a slice of Projects contained in a single page of results.
func ExtractDomains(r pagination.Page) ([]Domain, error) {

	var s struct {
		Domains []Domain `json:"domains"`
	}
	err := (r.(DomainPage)).ExtractInto(&s)
	return s.Domains, err
}

// Extract interprets any projectResults as a Project.
func (r domainResult) Extract() (*Domain, error) {
	var s struct {
		Domain *Domain `json:"domain"`
	}
	err := r.ExtractInto(&s)
	return s.Domain, err
}
