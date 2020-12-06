package twitch

import (
	"context"
	"net/http"
	"strings"

	"github.com/Adeithe/go-twitch/api"
	"github.com/shurcooL/graphql"
)

type Tripper struct {
	http.RoundTripper
}

func (t Tripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("Client-ID", api.Official.ID)
	return t.RoundTripper.RoundTrip(r)
}

type QueryUsers struct {
	Users []User `graphql:"users(ids: $ids)"`
}

type User struct {
	ID          graphql.ID
	Login       graphql.String
	DisplayName graphql.String
	Stream      struct {
		ID           graphql.String
		Title        graphql.String
		ViewersCount graphql.Int
		Game         struct {
			ID   graphql.ID
			Name graphql.String
		}
		CreatedAt graphql.String
	}
}

// GraphQL stores client for GraphQL requests
var GraphQL *graphql.Client = graphql.NewClient("https://gql.twitch.tv/gql", &http.Client{Transport: Tripper{http.DefaultTransport}})

func ToChannelName(s string) string {
	return strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(strings.TrimSuffix(strings.ToLower(s), ","), "#"), "@"))
}

func GetUsers(ids ...string) ([]User, error) {
	users := QueryUsers{}
	var gids []graphql.ID
	for _, id := range ids {
		gids = append(gids, id)
	}
	vars := map[string]interface{}{
		"ids": gids,
	}
	if err := GraphQL.Query(context.Background(), &users, vars); err != nil {
		return []User{}, err
	}
	return users.Users, nil
}
