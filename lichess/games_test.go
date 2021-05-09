package lichess

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGamesService_Get(t *testing.T) {
	client, mux, teardown := setUp()
	defer teardown()

	mux.HandleFunc("/game/export/12345678", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mediaTypeEnableNDJson)
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `
		{
		  "id": "12345678",
		  "rated": true,
		  "variant": "standard",
		  "speed": "blitz",
		  "status": "draw",
		  "opening": {
			"eco": "D31",
			"name": "Semi-Slav Defense: Marshall Gambit",
			"ply": 7
		  },
		  "clock": {
			"initial": 300,
			"increment": 3,
			"totalTime": 420
		  }
		}
	`)
	})

	ctx := context.Background()
	game, _, err := client.Games.Get(ctx, "12345678")

	if err != nil {
		t.Errorf("Games.Get returned error: %v", err)
	}

	want := &Game{
		ID:      "12345678",
		Rated:   true,
		Variant: "standard",
		Speed:   "blitz",
		Status:  "draw",
		Clock: struct {
			Initial   int `json:"initial"`
			Increment int `json:"increment"`
			TotalTime int `json:"totalTime"`
		}{Initial: 300, Increment: 3, TotalTime: 420},
		Opening: struct {
			Eco  string `json:"eco"`
			Name string `json:"name"`
			Ply  int64  `json:"ply"`
		}(struct {
			Eco  string
			Name string
			Ply  int64
		}{Eco: "D31", Name: "Semi-Slav Defense: Marshall Gambit", Ply: 7}),
	}

	if diff := cmp.Diff(game, want); diff != "" {
		t.Errorf("Responses do not match. Diff: %+v", diff)
	}
}

func TestGamesService_List(t *testing.T) {
	client, mux, teardown := setUp()
	defer teardown()

	mux.HandleFunc("/api/games/user/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mediaTypeEnableNDJson)
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `
		{
			"id": "id_1",
			"rated": true,
			"variant": "standard",
			"perf": "rapid",
			"createdAt": 1620384484273,
			"status": "resign"
		}
		{
			"id": "id_2",
			"rated": false,
			"variant": "standard",
			"perf": "rapid",
			"createdAt": 1620381701704,
			"status": "mate"
		}
		`)
	})

	ctx := context.Background()
	games, _, err := client.Games.List(ctx, "test", ListOptions{})

	if err != nil {
		t.Errorf("Games.List returned error: %v", err)
	}

	want := []*Game{{
		ID:        "id_1",
		Rated:     true,
		Variant:   "standard",
		Perf:      "rapid",
		CreatedAt: 1620384484273,
		Status:    "resign",
	}, {ID: "id_2",
		Rated:     false,
		Variant:   "standard",
		Perf:      "rapid",
		CreatedAt: 1620381701704,
		Status:    "mate",
	}}

	if diff := cmp.Diff(games, want); diff != "" {
		t.Errorf("Responses do not match. Diff: %+v", diff)
	}
}
