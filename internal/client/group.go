package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Group struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type GroupRequest struct {
	Group Group `json:"group"`
}

type GroupResponse struct {
	Group Group `json:"group"`
}

type GroupListResponse struct {
	Groups []Group `json:"group"`
}

type GetGroupResponse struct {
	Groups     GroupListResponse `json:"groups"`
	Pagination Pagination        `json:"pagination"`
}

func (c *TableauClient) CreateGroup(name string) (*Group, error) {
	newGroup := Group{
		Name: name,
	}
	groupRequest := GroupRequest{
		Group: newGroup,
	}

	payload, err := json.Marshal(groupRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/groups", c.ApiUrl), strings.NewReader(string(payload)))
	if err != nil {
		return nil, err
	}

	body, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	resp := GroupResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Group, nil
}

func (c *TableauClient) GetGroupByName(groupName string) (*Group, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/groups?filter=name:eq:%s", c.ApiUrl, url.QueryEscape(groupName)), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	resp := GetGroupResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	for _, group := range resp.Groups.Groups {
		if group.Name == groupName {
			return &group, nil
		}
	}

	return nil, fmt.Errorf("unable to find group with name '%s'", groupName)
}

func (c *TableauClient) GetGroupByID(groupID string) (*Group, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/groups?pageSize=1000", c.ApiUrl), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	resp := GetGroupResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	for _, group := range resp.Groups.Groups {
		if group.ID == groupID {
			return &group, nil
		}
	}

	return nil, fmt.Errorf("unable to find group with id %s", groupID)
}

func (c *TableauClient) UpdateGroup(groupID string, name string) (*Group, error) {
	updatedGroup := Group{
		Name: name,
	}
	groupRequest := GroupRequest{
		Group: updatedGroup,
	}

	payload, err := json.Marshal(groupRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/groups/%s", c.ApiUrl, groupID), strings.NewReader(string(payload)))
	if err != nil {
		return nil, err
	}

	body, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	resp := GroupResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Group, nil
}

func (c *TableauClient) DeleteGroup(groupID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/groups/%s", c.ApiUrl, groupID), nil)
	if err != nil {
		return err
	}

	_, err = c.sendRequest(req)
	if err != nil {
		return err
	}

	return nil
}
