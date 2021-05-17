package lichess

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type GamesService service

type Game struct {
	ID         string `json:"id"`
	Rated      bool   `json:"rated"`
	Variant    string `json:"variant"`
	Speed      string `json:"speed"`
	Perf       string `json:"perf"`
	CreatedAt  int64  `json:"createdAt"`
	LastMoveAt int64  `json:"lastMoveAt"`
	Status     string `json:"status"`
	Players    struct {
		White struct {
			User struct {
				Name string `json:"name"`
				ID   string `json:"id"`
			} `json:"user"`
			Rating     int       `json:"rating"`
			RatingDiff int       `json:"ratingDiff"`
			Analysis   *Analysis `json:"analysis"`
		} `json:"white"`
		Black struct {
			User struct {
				Name string `json:"name"`
				ID   string `json:"id"`
			} `json:"user"`
			Rating     int       `json:"rating"`
			RatingDiff int       `json:"ratingDiff"`
			Analysis   *Analysis `json:"analysis"`
		} `json:"black"`
	} `json:"players"`
	Winner  string `json:"winner"`
	Moves   string `json:"moves"`
	Pgn     string `json:"pgn"`
	Opening struct {
		Eco  string `json:"eco"`
		Name string `json:"name"`
		Ply  int64  `json:"ply"`
	} `json:"opening"`
	Clock struct {
		Initial   int `json:"initial"`
		Increment int `json:"increment"`
		TotalTime int `json:"totalTime"`
	} `json:"clock"`
}

type Analysis struct {
	Inaccuracy uint8 `json:"inaccuracy"`
	Mistake    uint8 `json:"mistake"`
	Blunder    uint8 `json:"blunder"`
	ACPL       uint8 `json:"acpl"`
}

type ListOptions struct {
	Since int64 `url:"since,omitempty"`
}

func (s *GamesService) Get(ctx context.Context, ID string) (*Game, *Response, error) {
	u := fmt.Sprintf("/game/export/%v?pgnInJson=true", ID)
	req, err := s.client.NewRequest("GET", u, nil)

	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	game := new(Game)

	resp, err := s.client.Do(ctx, req, game)

	if err != nil {
		return nil, resp, err
	}

	return game, resp, nil
}

func (s *GamesService) List(ctx context.Context, username string, opts ListOptions) ([]*Game, *Response, error) {
	u := fmt.Sprintf("/api/games/user/%v?pgnInJson=true&since=%v&opening=true&cloacks=true", username, opts.Since)
	req, err := s.client.NewRequest("GET", u, nil)

	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	req.Header.Set("Accept", mediaTypeEnableNDJson)

	var games []*Game

	resp, err := s.client.Do(ctx, req, &games)

	if err != nil {
		return nil, resp, errors.WithStack(err)
	}

	return games, resp, nil
}

func (s *GamesService) All(ctx context.Context, username string) (<-chan *Game, <-chan error) {
	max := 50
	since := 0
	gch := make(chan *Game, max)
	errCh := make(chan error)

	defer func() {
		close(gch)
		close(errCh)
	}()

	go func() {
		for {
			u := fmt.Sprintf("/api/games/user/%v?pgnInJson=true&since=%v&max=%v", username, since, max)
			req, err := s.client.NewRequest("GET", u, nil)

			if err != nil {
				errCh <- err
				break
			}

			req.Header.Set("Accept", mediaTypeEnableNDJson)

			var games []*Game

			_, err = s.client.Do(ctx, req, &games)

			if err != nil {
				errCh <- err
				break
			}

			for _, g := range games {
				gch <- g
			}

			if len(games) < max {
				close(gch)
				close(errCh)

				break
			}

			since = int(games[max-1].CreatedAt)
		}
	}()

	return gch, errCh
}
