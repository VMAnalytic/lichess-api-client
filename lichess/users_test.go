package lichess

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUsersService_Get(t *testing.T) {
	client, mux, teardown := setUp()
	defer teardown()

	mux.HandleFunc("/user/vmyroslav", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"id": "vmyroslav",
			"username": "VMyroslav",
			"online": false,
			"createdAt": 1617221900731,
			"seenAt": 1619639888472,
			"playTime": {
				"total": 96338,
				"tv": 0
			},
			"language": "en-US",
			"url": "https://lichess.org/@/VMyroslav",
			"nbFollowing": 0,
			"nbFollowers": 0,
			"completionRate": 97,
			"count": {
				"all": 145,
				"rated": 143,
				"ai": 2,
				"draw": 13,
				"drawH": 13,
				"loss": 56,
				"lossH": 56,
				"win": 76,
				"winH": 74,
				"bookmark": 0,
				"playing": 0,
				"import": 0,
				"me": 0
			}
		}`)
	})

	ctx := context.Background()
	user, _, err := client.Users.Get(ctx, "vmyroslav")

	if err != nil {
		t.Errorf("Account.GetMyEmail returned error: %v", err)
	}

	want := &User{
		ID:        "vmyroslav",
		Username:  "VMyroslav",
		Online:    false,
		CreatedAt: 1617221900731,
		SeenAt:    1619639888472,
		Playtime: struct {
			Total int `json:"total"`
			Tv    int `json:"tv"`
		}{96338, 0},
		Language:       "en-US",
		URL:            "https://lichess.org/@/VMyroslav",
		CompletionRate: 97,
		Profile:        nil,
		Stat: &Stats{
			All:      145,
			Rated:    143,
			Ai:       2,
			Draw:     13,
			DrawH:    13,
			Loss:     56,
			LossH:    56,
			Win:      76,
			WinH:     74,
			Bookmark: 0,
			Playing:  0,
			Import:   0,
			Me:       0,
		},
	}

	if diff := cmp.Diff(user, want); diff != "" {
		t.Errorf("Responses do not match. Diff: %+v", diff)
	}
}
