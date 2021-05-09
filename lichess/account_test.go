package lichess

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAccountService_GetMyEmail(t *testing.T) {
	client, mux, teardown := setUp()
	defer teardown()

	mux.HandleFunc("/api/account/email", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
				"email": "example@email.com"
			}`)
	})

	ctx := context.Background()
	email, _, err := client.Account.GetMyEmail(ctx)

	if err != nil {
		t.Errorf("Account.GetMyEmail returned error: %v", err)
	}

	want := "example@email.com"
	if !reflect.DeepEqual(email, want) {
		t.Errorf("Account.GetMyEmail returned %+v, want %+v", email, want)
	}
}

func TestAccountService_GetMyPreferences(t *testing.T) {
	client, mux, teardown := setUp()
	defer teardown()

	mux.HandleFunc("/api/account/preferences", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
				"prefs": {
					"dark": true,
					"transp": false,
					"bgImg": "//lichess1.org/assets/images/background/landscape.jpg",
					"is3d": false,
					"theme": "maple",
					"pieceSet": "cburnett",
					"rookCastle": 1
				}
			}`)
	})

	ctx := context.Background()
	pref, _, err := client.Account.GetMyPreferences(ctx)

	if err != nil {
		t.Errorf("Account.GetMyPreferences returned error: %v", err)
	}

	want := &Preferences{
		Dark:        true,
		Transparent: false,
		BgImg:       "//lichess1.org/assets/images/background/landscape.jpg",
		Is3D:        false,
		Theme:       "maple",
		PieceSet:    "cburnett",
		RookCastle:  1,
	}

	if diff := cmp.Diff(pref, want); diff != "" {
		t.Errorf("Account.GetMyPreferences returned %+v, want %+v", pref, want)
	}
}

func TestAccountService_GetMyProfile(t *testing.T) {
	client, mux, teardown := setUp()
	defer teardown()

	mux.HandleFunc("/api/account", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
				"id": "testID",
				"username": "Test",
				"online": false,
				"createdAt": 1617221900731,
				"seenAt": 1620387000153,
				"playTime": {
					"total": 121450,
					"tv": 0
				},
				"language": "en-US",
				"count": {
					"all": 176
				}
			}`)
	})

	ctx := context.Background()
	prof, _, err := client.Account.GetMyProfile(ctx)

	if err != nil {
		t.Errorf("Account.GetMyProfile returned error: %v", err)
	}

	want := &User{
		ID:        "testID",
		Username:  "Test",
		Online:    false,
		CreatedAt: 1617221900731,
		SeenAt:    1620387000153,
		Playtime: struct {
			Total int `json:"total"`
			Tv    int `json:"tv"`
		}{Total: 121450, Tv: 0},
		Language: "en-US",
		Stat:     &Stats{All: 176},
	}

	if diff := cmp.Diff(prof, want); diff != "" {
		t.Errorf("Account.GetMyProfile returned %+v, want %+v", prof, want)
	}
}
