package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type User struct {
	ID          string `json:"id,omitempty"`
	Email       string `json:"email,omitempty"`
	Name        string `json:"name,omitempty"`
	SiteRole    string `json:"siteRole,omitempty"`
	AuthSetting string `json:"authSetting,omitempty"`
}

type UserRequest struct {
	User User `json:"user"`
}

type UserResponse struct {
	User User `json:"user"`
}

type UserListResponse struct {
	Users []User `json:"user"`
}

type GetUserResponse struct {
	Users      UserListResponse `json:"users"`
	Pagination Pagination       `json:"pagination"`
}

func (c *TableauClient) CreateUser(email string, siteRole string, authSetting string) (*User, error) {
	newUser := User{
		Email:       email,
		Name:        email,
		SiteRole:    siteRole,
		AuthSetting: authSetting,
	}
	userRequest := UserRequest{
		User: newUser,
	}

	payload, err := json.Marshal(userRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users", c.ApiUrl), strings.NewReader(string(payload)))
	if err != nil {
		return nil, err
	}

	body, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	resp := UserResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.User, nil
}

func (c *TableauClient) GetUser(userID string) (*User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/%s/", c.ApiUrl, userID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	resp := UserResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.User, nil
}

func (c *TableauClient) GetUserByEmail(userEmail string) (*User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users?filter=name:eq:%s", c.ApiUrl, url.QueryEscape(userEmail)), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	resp := GetUserResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	for _, user := range resp.Users.Users {
		if user.Email == userEmail {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("unable to find user with email '%s'", userEmail)
}

func (c *TableauClient) UpdateUser(userID string, email string, siteRole string, authSetting string) (*User, error) {
	updatedUser := User{
		Email:       email,
		Name:        email,
		SiteRole:    siteRole,
		AuthSetting: authSetting,
	}
	userRequest := UserRequest{
		User: updatedUser,
	}

	payload, err := json.Marshal(userRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/users/%s", c.ApiUrl, userID), strings.NewReader(string(payload)))
	if err != nil {
		return nil, err
	}

	body, err := c.sendRequest(req)
	if err != nil {
		return nil, err
	}

	resp := UserResponse{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.User, nil
}

func (c *TableauClient) DeleteUser(userID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/users/%s", c.ApiUrl, userID), nil)
	if err != nil {
		return err
	}

	_, err = c.sendRequest(req)
	if err != nil {
		return err
	}

	return nil
}
