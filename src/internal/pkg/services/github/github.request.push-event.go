package github

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/nobuenhombre/suikat/pkg/ge"
	configgithub "go-github-webhook-cicd/src/internal/pkg/services/github/config"
	"io"
	"net/http"
	"strings"
	"time"
)

type PushEventRequestHeaders struct {
	Event                      string
	Delivery                   string
	HookID                     string
	HookInstallationTargetID   string
	HookInstallationTargetType string
	Signature                  string
	Signature256               string
}

func NewPushEventRequestHeaders(r *http.Request) *PushEventRequestHeaders {
	headers := &PushEventRequestHeaders{}
	headers.Event = r.Header.Get("X-GitHub-Event")
	headers.Delivery = r.Header.Get("X-GitHub-Delivery")
	headers.HookID = r.Header.Get("X-GitHub-Hook-ID")
	headers.HookInstallationTargetID = r.Header.Get("X-GitHub-Hook-Installation-Target-ID")
	headers.HookInstallationTargetType = r.Header.Get("X-GitHub-Hook-Installation-Target-Type")
	headers.Signature = r.Header.Get("X-Hub-Signature")
	headers.Signature256 = r.Header.Get("X-Hub-Signature-256")

	return headers
}

func (prh *PushEventRequestHeaders) Validate() error {
	if prh.Event != "push" {
		return ge.Pin(ge.New("Event header is not [push]"))
	}

	if prh.Delivery == "" {
		return ge.Pin(ge.New("Delivery header is empty"))
	}

	if prh.HookID == "" {
		return ge.Pin(ge.New("HookID header is empty"))
	}

	if prh.HookInstallationTargetID == "" {
		return ge.Pin(ge.New("HookInstallationTargetID header is empty"))
	}

	if prh.HookInstallationTargetType == "" {
		return ge.Pin(ge.New("HookInstallationTargetType header is empty"))
	}

	if prh.Signature == "" {
		return ge.Pin(ge.New("Signature header is empty"))
	}

	if prh.Signature256 == "" {
		return ge.Pin(ge.New("Signature256 header is empty"))
	}

	return nil
}

type PushEventRequest struct {
	Headers *PushEventRequestHeaders
	Body    []byte
	Data    *PushEventData
}

func NewPushEventRequest(r *http.Request) (*PushEventRequest, error) {
	request := &PushEventRequest{}

	request.Headers = NewPushEventRequestHeaders(r)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, ge.Pin(err)
	}

	request.Body = body

	request.Data = new(PushEventData)

	err = json.Unmarshal(request.Body, request.Data)
	if err != nil {
		return nil, ge.Pin(err)
	}

	return request, nil
}

const SignatureTypeSHA1 = "sha1"

func (pr *PushEventRequest) Validate(project *configgithub.GitHubProjectConfig) error {
	err := pr.Headers.Validate()
	if err != nil {
		return ge.Pin(err)
	}

	signatureParts := strings.SplitN(pr.Headers.Signature, "=", 2)
	if len(signatureParts) != 2 {
		return ge.Pin(ge.New("Signature header does not contain (hash type and hash)", ge.Params{"prh.Signature": pr.Headers.Signature}))
	}

	if signatureParts[0] != SignatureTypeSHA1 {
		return ge.Pin(&ge.MismatchError{
			ComparedItems: "signature type",
			Expected:      SignatureTypeSHA1,
			Actual:        signatureParts[0],
		})
	}

	hm := hmac.New(sha1.New, []byte(project.Secret))
	hm.Write(pr.Body)

	hash := fmt.Sprintf("%x", hm.Sum(nil))

	if !hmac.Equal([]byte(hash), []byte(signatureParts[1])) {
		return ge.Pin(ge.New("Signature is invalid"))
	}

	if pr.Data.Repository.FullName != project.Repository {
		return ge.Pin(&ge.MismatchError{
			ComparedItems: "Repository",
			Expected:      project.Repository,
			Actual:        pr.Data.Repository.FullName,
		})
	}

	ref := fmt.Sprintf("refs/heads/%v", project.Branch)
	if pr.Data.Ref != ref {
		return ge.Pin(&ge.MismatchError{
			ComparedItems: "Ref",
			Expected:      ref,
			Actual:        pr.Data.Ref,
		})
	}

	return nil
}

type Owner struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	Login             string `json:"login"`
	Id                int    `json:"id"`
	NodeId            string `json:"node_id"`
	AvatarUrl         string `json:"avatar_url"`
	GravatarId        string `json:"gravatar_id"`
	Url               string `json:"url"`
	HtmlUrl           string `json:"html_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReposUrl          string `json:"repos_url"`
	EventsUrl         string `json:"events_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Pusher struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Sender struct {
	Login             string `json:"login"`
	Id                int    `json:"id"`
	NodeId            string `json:"node_id"`
	AvatarUrl         string `json:"avatar_url"`
	GravatarId        string `json:"gravatar_id"`
	Url               string `json:"url"`
	HtmlUrl           string `json:"html_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReposUrl          string `json:"repos_url"`
	EventsUrl         string `json:"events_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type Commit struct {
	Id        string    `json:"id"`
	TreeId    string    `json:"tree_id"`
	Distinct  bool      `json:"distinct"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Url       string    `json:"url"`
	Author    User      `json:"author"`
	Committer User      `json:"committer"`
	Added     []string  `json:"added"`
	Removed   []string  `json:"removed"`
	Modified  []string  `json:"modified"`
}

type Repository struct {
	Id                       int           `json:"id"`
	NodeId                   string        `json:"node_id"`
	Name                     string        `json:"name"`
	FullName                 string        `json:"full_name"`
	Private                  bool          `json:"private"`
	Owner                    Owner         `json:"owner"`
	HtmlUrl                  string        `json:"html_url"`
	Description              string        `json:"description"`
	Fork                     bool          `json:"fork"`
	Url                      string        `json:"url"`
	ForksUrl                 string        `json:"forks_url"`
	KeysUrl                  string        `json:"keys_url"`
	CollaboratorsUrl         string        `json:"collaborators_url"`
	TeamsUrl                 string        `json:"teams_url"`
	HooksUrl                 string        `json:"hooks_url"`
	IssueEventsUrl           string        `json:"issue_events_url"`
	EventsUrl                string        `json:"events_url"`
	AssigneesUrl             string        `json:"assignees_url"`
	BranchesUrl              string        `json:"branches_url"`
	TagsUrl                  string        `json:"tags_url"`
	BlobsUrl                 string        `json:"blobs_url"`
	GitTagsUrl               string        `json:"git_tags_url"`
	GitRefsUrl               string        `json:"git_refs_url"`
	TreesUrl                 string        `json:"trees_url"`
	StatusesUrl              string        `json:"statuses_url"`
	LanguagesUrl             string        `json:"languages_url"`
	StargazersUrl            string        `json:"stargazers_url"`
	ContributorsUrl          string        `json:"contributors_url"`
	SubscribersUrl           string        `json:"subscribers_url"`
	SubscriptionUrl          string        `json:"subscription_url"`
	CommitsUrl               string        `json:"commits_url"`
	GitCommitsUrl            string        `json:"git_commits_url"`
	CommentsUrl              string        `json:"comments_url"`
	IssueCommentUrl          string        `json:"issue_comment_url"`
	ContentsUrl              string        `json:"contents_url"`
	CompareUrl               string        `json:"compare_url"`
	MergesUrl                string        `json:"merges_url"`
	ArchiveUrl               string        `json:"archive_url"`
	DownloadsUrl             string        `json:"downloads_url"`
	IssuesUrl                string        `json:"issues_url"`
	PullsUrl                 string        `json:"pulls_url"`
	MilestonesUrl            string        `json:"milestones_url"`
	NotificationsUrl         string        `json:"notifications_url"`
	LabelsUrl                string        `json:"labels_url"`
	ReleasesUrl              string        `json:"releases_url"`
	DeploymentsUrl           string        `json:"deployments_url"`
	CreatedAt                int           `json:"created_at"`
	UpdatedAt                time.Time     `json:"updated_at"`
	PushedAt                 int           `json:"pushed_at"`
	GitUrl                   string        `json:"git_url"`
	SshUrl                   string        `json:"ssh_url"`
	CloneUrl                 string        `json:"clone_url"`
	SvnUrl                   string        `json:"svn_url"`
	Homepage                 interface{}   `json:"homepage"`
	Size                     int           `json:"size"`
	StargazersCount          int           `json:"stargazers_count"`
	WatchersCount            int           `json:"watchers_count"`
	Language                 interface{}   `json:"language"`
	HasIssues                bool          `json:"has_issues"`
	HasProjects              bool          `json:"has_projects"`
	HasDownloads             bool          `json:"has_downloads"`
	HasWiki                  bool          `json:"has_wiki"`
	HasPages                 bool          `json:"has_pages"`
	HasDiscussions           bool          `json:"has_discussions"`
	ForksCount               int           `json:"forks_count"`
	MirrorUrl                interface{}   `json:"mirror_url"`
	Archived                 bool          `json:"archived"`
	Disabled                 bool          `json:"disabled"`
	OpenIssuesCount          int           `json:"open_issues_count"`
	License                  interface{}   `json:"license"`
	AllowForking             bool          `json:"allow_forking"`
	IsTemplate               bool          `json:"is_template"`
	WebCommitSignoffRequired bool          `json:"web_commit_signoff_required"`
	Topics                   []interface{} `json:"topics"`
	Visibility               string        `json:"visibility"`
	Forks                    int           `json:"forks"`
	OpenIssues               int           `json:"open_issues"`
	Watchers                 int           `json:"watchers"`
	DefaultBranch            string        `json:"default_branch"`
	Stargazers               int           `json:"stargazers"`
	MasterBranch             string        `json:"master_branch"`
}

type PushEventData struct {
	Ref        string      `json:"ref"`
	Before     string      `json:"before"`
	After      string      `json:"after"`
	Repository Repository  `json:"repository"`
	Pusher     Pusher      `json:"pusher"`
	Sender     Sender      `json:"sender"`
	Created    bool        `json:"created"`
	Deleted    bool        `json:"deleted"`
	Forced     bool        `json:"forced"`
	BaseRef    interface{} `json:"base_ref"`
	Compare    string      `json:"compare"`
	Commits    []Commit    `json:"commits"`
	HeadCommit Commit      `json:"head_commit"`
}
