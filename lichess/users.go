package lichess

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type UsersService service

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Online    bool   `json:"online"`
	CreatedAt int64  `json:"createdAt"`
	SeenAt    int64  `json:"seenAt"`
	Playtime  struct {
		Total int `json:"total"`
		Tv    int `json:"tv"`
	} `json:"playTime,omitempty"`
	Language       string   `json:"language"`
	URL            string   `json:"url"`
	CompletionRate int      `json:"completionRate,omitempty"`
	Profile        *Profile `json:"profile,omitempty"`
	Stat           *Stats   `json:"count,omitempty"`
}

type Profile struct {
	Country    string `json:"country,omitempty"`
	Location   string `json:"location,omitempty"`
	Bio        string `json:"bio,omitempty"`
	Firstname  string `json:"firstName,omitempty"`
	Lastname   string `json:"lastName,omitempty"`
	FideRating int    `json:"fideRating,omitempty"`
	UscfRating int    `json:"uscfRating,omitempty"`
	EcfRating  int    `json:"ecfRating,omitempty"`
	Links      string `json:"links,omitempty"`
}

type Stats struct {
	All      int `json:"all,omitempty"`
	Rated    int `json:"rated,omitempty"`
	Ai       int `json:"ai,omitempty"`
	Draw     int `json:"draw,omitempty"`
	DrawH    int `json:"drawH,omitempty"`
	Loss     int `json:"loss,omitempty"`
	LossH    int `json:"lossH,omitempty"`
	Win      int `json:"win,omitempty"`
	WinH     int `json:"winH,omitempty"`
	Bookmark int `json:"bookmark,omitempty"`
	Playing  int `json:"playing,omitempty"`
	Import   int `json:"import,omitempty"`
	Me       int `json:"me,omitempty"`
}

func (s *UsersService) Get(ctx context.Context, username string) (*User, *Response, error) {
	u := fmt.Sprintf("/api/user/%v", username)
	req, err := s.client.NewRequest("GET", u, nil)

	if err != nil {
		return nil, nil, errors.Wrap(err, "")
	}

	player := new(User)
	resp, err := s.client.Do(ctx, req, player)

	if err != nil {
		return nil, resp, err
	}

	return player, resp, nil
}
