package main

import (
	"encoding/json"
	"time"

	"github.com/trustero/api/go/receptor_sdk"
	"github.com/trustero/api/go/receptor_sdk/cmd"
	"github.com/xanzy/go-gitlab"
)

type GitLabUser struct {
	Username         string     `trustero:"id:;display:Username;order:1"`
	Name             string     `trustero:"display:Name;order:2"`
	IsAdmin          bool       `trustero:"display:Admin;order:3"`
	CreatedAt        *time.Time `trustero:"display:Created On;order:4"`
	TwoFactorEnabled bool       `trustero:"display:MFA Enabled;order:5"`
	LastActivityOn   *time.Time `trustero:"display:Last Activity On;order:6"`
}

type Receptor struct {
	Token   string
	GroupID string
}

func (r *Receptor) GetReceptorType() (receptorType string) {
	receptorType = "example_gitlab"
	return
}

func (r *Receptor) UnmarshalCredentials(credentials string) (obj interface{}, err error) {
	obj, err = receptor_sdk.UnmarshalCredentials(credentials, r)
	return
}

func (r *Receptor) Verify(credentials interface{}) (ok bool, err error) {
	c := credentials.(*Receptor)
	ok = false
	var git *gitlab.Client
	if git, err = gitlab.NewClient(c.Token); err == nil {
		if _, _, err = git.Groups.ListGroupMembers(c.GroupID, &gitlab.ListGroupMembersOptions{}); err != nil {
			return
		}
	}
	if err == nil {
		ok = true
	}
	return
}

func (r *Receptor) Discover(credentials interface{}) (services []*receptor_sdk.Service, err error) {
	c := credentials.(*Receptor)
	var git *gitlab.Client
	if git, err = gitlab.NewClient(c.Token); err == nil {
		services = []*receptor_sdk.Service{}

		// Get group members.  "User" is the Service.Name and the user name is the Service.InstanceId
		var members []*gitlab.GroupMember
		members, _, err = git.Groups.ListGroupMembers(c.GroupID, &gitlab.ListGroupMembersOptions{})
		for _, member := range members {
			services = append(services, receptor_sdk.NewService("User", member.Username))
		}
	}
	return
}

func (r *Receptor) Report(credentials interface{}) (evidences []*receptor_sdk.Evidence, err error) {
	c := credentials.(*Receptor)
	var git *gitlab.Client
	if git, err = gitlab.NewClient(c.Token); err == nil {
		var user *gitlab.User
		var members []*gitlab.GroupMember
		members, _, err = git.Groups.ListGroupMembers(c.GroupID, &gitlab.ListGroupMembersOptions{})
		evidence := receptor_sdk.NewEvidence("User", "GitLab Users", "")
		for _, member := range members {
			user, _, err = git.Users.GetUser(member.ID, gitlab.GetUsersOptions{})
			evidence.AddSource("git.Users.GetUser(member.ID, gitlab.GetUsersOptions{})", user)
			evidence.AddRow(*newGitLabUser(user))
		}
		evidences = append(evidences, evidence)
	}
	return
}

func newGitLabUser(user *gitlab.User) (guser *GitLabUser) {
	guser = &GitLabUser{
		Username:         user.Username,
		Name:             user.Name,
		IsAdmin:          user.IsAdmin,
		CreatedAt:        user.CreatedAt,
		TwoFactorEnabled: user.TwoFactorEnabled,
		LastActivityOn:   (*time.Time)(user.LastActivityOn),
	}

	return
}

func main() {
	var token string
	var groupId string

	// Add convenience token
	cmd.RootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "GitLab user access token")
	cmd.RootCmd.PersistentFlags().StringVarP(&groupId, "gid", "g", "", "GitLab group id")

	// Get credentials from flags
	receptor_sdk.CredentialsFromFlags = func() (j string) {
		j = ""
		if len(token) > 0 {
			b, err := json.Marshal(&Receptor{Token: token, GroupID: groupId})
			if err == nil {
				j = string(b)
			}
		}
		return
	}

	cmd.Execute(&Receptor{})
}
