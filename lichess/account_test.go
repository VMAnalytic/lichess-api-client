package lichess

import (
	"context"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"reflect"
	"testing"
)

func TestAccountService_GetMyEmail(t *testing.T) {
	client, mux, _, teardown := setUp()
	defer teardown()

	mux.HandleFunc("/account/email", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
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
	client, mux, _, teardown := setUp()
	defer teardown()

	mux.HandleFunc("/account/preferences", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
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
		Dark:          true,
		Transparent:   false,
		BgImg:         "//lichess1.org/assets/images/background/landscape.jpg",
		Is3D:          false,
		Theme:         "maple",
		PieceSet:      "cburnett",
		RookCastle:    1,
	}

	if diff := cmp.Diff(pref, want); diff != "" {
		t.Errorf("Account.GetMyPreferences returned %+v, want %+v", pref, want)
	}
}
