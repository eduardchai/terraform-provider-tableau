package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type GroupMembershipEmailList struct {
	GroupID    string
	UserEmails []string
}

type GroupMembershipRequest struct {
	User User `json:"user"`
}

type GroupMembershipResponse struct {
	User User `json:"user"`
}

type GetGroupMembershipResponse struct {
	Users      UserListResponse `json:"users"`
	Pagination Pagination       `json:"pagination"`
}

func (c *TableauClient) CreateGroupMembershipByUserID(groupID string, userID string) error {
	// Create request object
	groupMembershipRequest := GroupMembershipRequest{
		User: User{
			ID: userID,
		},
	}

	// Create JSON payload
	payload, err := json.Marshal(groupMembershipRequest)
	if err != nil {
		return err
	}

	// Create request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/groups/%s/users", c.ApiUrl, groupID), strings.NewReader(string(payload)))
	if err != nil {
		return err
	}

	// Send request
	body, err := c.sendRequest(req)
	if err != nil {
		return err
	}

	// Unmarshal response
	resp := GroupMembershipResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return err
	}

	return nil
}

func (c *TableauClient) CreateGroupMembershipByUserEmail(groupID string, userEmail string) error {
	// Get user by email
	user, err := c.GetUserByEmail(userEmail)
	if err != nil {
		return err
	}

	// Create group membership
	err = c.CreateGroupMembershipByUserID(groupID, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (c *TableauClient) GetGroupMembership(groupID string) (*GroupMembershipEmailList, error) {
	// Create request
	// TODO: need to handle pagination
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/groups/%s/users?pageSize=1000", c.ApiUrl, groupID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	resp := GetGroupMembershipResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	var userEmails []string
	for _, user := range resp.Users.Users {
		userEmails = append(userEmails, user.Email)
	}

	groupMembershipEmailList := GroupMembershipEmailList{
		GroupID:    groupID,
		UserEmails: userEmails,
	}

	return &groupMembershipEmailList, nil
}

func (c *TableauClient) DeleteGroupMembershipByUserID(groupID string, userID string) error {
	// Create delete request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/groups/%s/users/%s", c.ApiUrl, groupID, userID), nil)
	if err != nil {
		return err
	}

	// Send request
	_, err = c.sendRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *TableauClient) DeleteGroupMembershipByUserEmail(groupID string, userEmail string) error {
	// Get user by email
	user, err := c.GetUserByEmail(userEmail)
	if err != nil {
		return err
	}

	// Delete group membership
	err = c.DeleteGroupMembershipByUserID(groupID, user.ID)
	if err != nil {
		return err
	}

	return nil
}
