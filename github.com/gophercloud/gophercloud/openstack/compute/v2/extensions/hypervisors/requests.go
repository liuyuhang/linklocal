package hypervisors

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)


type ListOptsBuilder interface {
	ToHypervisorListQuery() (string, error)
}

type ListOpts struct {
	HypervisorHostnamePattern string `json:"hypervisor_hostname_pattern"`
}

// ToServerListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToHypervisorListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}


// List makes a request against the API to list hypervisors.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := hypervisorsListDetailURL(client)
	if opts != nil {
		query, err := opts.ToHypervisorListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return HypervisorPage{pagination.SinglePageBase(r)}
	})
}


func Get(client *gophercloud.ServiceClient,id string)(r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 203},
	})
	return
}