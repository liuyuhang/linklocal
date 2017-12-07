package domain

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/ccwings/log"
)

// ListOptsBuilder allows extensions to add additional parameters to
// the List request
type ListOptsBuilder interface {
	ToDomainListQuery() (string, error)
}

// ListOpts allows you to query the List method.
type ListOpts struct {
	// Enabled filters the response by enabled projects.
	Enabled *bool `q:"enabled"`

	// Name filters the response by project name.
	Name string `q:"name"`
}

// ToDomainListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToDomainListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List enumerats the Projects to which the current token has access.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToDomainListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return DomainPage{pagination.LinkedPageBase{PageResult: r}}
	})
}
/*

// GetOptsBuilder allows extensions to add additional parameters to
// the Get request.
type GetOptsBuilder interface {
	ToProjectGetQuery() (string, error)
}

// GetOpts allows you to modify the details included in the Get request.
type GetOpts struct{}

// ToProjectGetQuery formats a GetOpts into a query string.
func (opts GetOpts) ToProjectGetQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}
*/
// Get retrieves details on a single project, by ID.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	url := getURL(client, id)
	_, r.Err = client.Get(url, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateOptsBuilder interface {
	ToDomainCreateMap() (map[string]interface{}, error)
}


// CreateOpts allows you to modify the details included in the Create request.
type CreateOpts struct {
	// Enabled sets the project status to enabled or disabled.
	Enabled *bool `json:"enabled,omitempty"`
	// Name is the name of the project.
	Name string `json:"name,required"`
	// Description is the description of the project.
	Description string `json:"description,omitempty"`
}

// ToProjectCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToDomainCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "domain")
}

// Create creates a new Project.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToDomainCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), &b, &r.Body, nil)
	return
}

type UpdateOptsBuilder interface {
	ToDomainUpdateMap() (map[string]interface{}, error)
}

type UpdateOpts struct {
	// Enabled sets the project status to enabled or disabled.
	Enabled *bool `json:"enabled,omitempty"`
	// Name is the name of the project.
	Name string `json:"name,omitempty"`
	// Description is the description of the project.
	Description string `json:"description,omitempty"`
}

func (opts UpdateOpts) ToDomainUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "domain")
}

func Update(client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToDomainUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Patch(updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	log.Debug(deleteURL(client, id))
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}


/*
// Delete deletes a project.
func Delete(client *gophercloud.ServiceClient, projectID string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, projectID), nil)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateOptsBuilder interface {
	ToProjectUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts allows you to modify the details included in the Update request.
type UpdateOpts struct {
	// DomainID is the ID this project will belong under.
	DomainID string `json:"domain_id,omitempty"`

	// Enabled sets the project status to enabled or disabled.
	Enabled *bool `json:"enabled,omitempty"`

	// IsDomain indicates if this project is a domain.
	IsDomain *bool `json:"is_domain,omitempty"`

	// Name is the name of the project.
	Name string `json:"name,omitempty"`

	// ParentID specifies the parent project of this new project.
	ParentID string `json:"parent_id,omitempty"`

	// Description is the description of the project.
	Description string `json:"description,omitempty"`
}

// ToUpdateCreateMap formats a UpdateOpts into an update request.
func (opts UpdateOpts) ToProjectUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "project")
}

// Update modifies the attributes of a project.
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToProjectUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Patch(updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}



// add by ccwings
//List Projects By User ID
func ListByUserId(client *gophercloud.ServiceClient, user_id string) pagination.Pager {
	url := listByUserIdURL(client,user_id)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ProjectPage{pagination.LinkedPageBase{PageResult: r}}
	})
}*/
