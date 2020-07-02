package V2BuildAssist

import (
	"context"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"net/http"
)

func CreateFileIfNotExist(accessToken, owner, repo, path, message string, content []byte) (string, error) {
	ctx := context.Background()
	client := getclient(accessToken, ctx)

	f, _, err := client.Repositories.CreateFile(ctx, owner, repo, path, &github.RepositoryContentFileOptions{Content: content, Message: &message})
	if err != nil {
		return "", err
	}
	url := f.GetURL()
	return url, nil
}

func getclient(accessToken string, ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return client
}

func GetReleaseFile(accessToken, owner, repo, tag, path string) ([]byte, int64, error) {
	ctx := context.Background()
	client := getclient(accessToken, ctx)
	rele, _, err := client.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
	if err != nil {
		return nil, 0, err
	}
	id := *rele.ID
	for _, v := range rele.Assets {
		if *v.Name == path {
			url := v.GetBrowserDownloadURL()
			resp, err := http.Get(url)
			if err != nil {
				return nil, 0, err
			}
			d, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, 0, err
			}
			return d, id, nil
		}
	}
	return nil, 0, io.EOF
}

func CreateCommentForCommit(accessToken, owner, repo, sha, message string) (string, error) {
	ctx := context.Background()
	client := getclient(accessToken, ctx)

	f, _, err := client.Repositories.CreateComment(ctx, owner, repo, sha, &github.RepositoryComment{Body: &message})
	if err != nil {
		return "", err
	}
	url := f.GetURL()
	return url, nil
}

func CreateCommentForPR(accessToken, owner, repo, message string, number int) (string, error) {
	ctx := context.Background()
	client := getclient(accessToken, ctx)
	f, _, err := client.Issues.CreateComment(ctx, owner, repo, number, &github.IssueComment{Body: &message})
	if err != nil {
		return "", err
	}
	url := f.GetURL()
	return url, nil
}
