package lichess

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type AccountService service

type Preferences struct {
	Dark          bool   `json:"dark"`
	Transparent   bool   `json:"transp"`
	BgImg         string `json:"bgImg"`
	Is3D          bool   `json:"is3d"`
	Theme         string `json:"theme"`
	PieceSet      string `json:"pieceSet"`
	Theme3D       string `json:"theme3d"`
	PieceSet3D    string `json:"pieceSet3d"`
	SoundSet      string `json:"soundSet"`
	Blindfold     int    `json:"blindfold"`
	AutoQueen     int    `json:"autoQueen"`
	AutoThreefold int    `json:"autoThreefold"`
	TakeBack      int    `json:"takeback"`
	MoreTime      int    `json:"moretime"`
	ClockTenths   int    `json:"clockTenths"`
	ClockBar      bool   `json:"clockBar"`
	ClockSound    bool   `json:"clockSound"`
	PreMove       bool   `json:"premove"`
	Animation     int    `json:"animation"`
	Captured      bool   `json:"captured"`
	Follow        bool   `json:"follow"`
	Highlight     bool   `json:"highlight"`
	Destination   bool   `json:"destination"`
	Coords        int    `json:"coords"`
	Replay        int    `json:"replay"`
	Challenge     int    `json:"challenge"`
	Message       int    `json:"message"`
	CoordColor    int    `json:"coordColor"`
	SubmitMove    int    `json:"submitMove"`
	ConfirmResign int    `json:"confirmResign"`
	InsightShare  int    `json:"insightShare"`
	KeyboardMove  int    `json:"keyboardMove"`
	Zen           int    `json:"zen"`
	MoveEvent     int    `json:"moveEvent"`
	RookCastle    int    `json:"rookCastle"`
}

func (s *AccountService) GetMyProfile(ctx context.Context) (*User, *Response, error) {
	u := fmt.Sprint("/api/account")
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

func (s *AccountService) GetMyEmail(ctx context.Context) (string, *Response, error) {
	u := fmt.Sprint("/api/account/email")
	req, err := s.client.NewRequest("GET", u, nil)

	if err != nil {
		return "", nil, errors.Wrap(err, "")
	}

	var e = struct {
		Email string `json:"email"`
	}{}

	resp, err := s.client.Do(ctx, req, &e)

	if err != nil {
		return "", resp, err
	}

	return e.Email, resp, nil
}

func (s *AccountService) GetMyPreferences(ctx context.Context) (*Preferences, *Response, error) {
	type prefResp struct {
		Pref *Preferences `json:"prefs"`
	}

	u := fmt.Sprint("/api/account/preferences")
	req, err := s.client.NewRequest("GET", u, nil)

	if err != nil {
		return nil, nil, errors.Wrap(err, "")
	}

	pref := new(prefResp)

	resp, err := s.client.Do(ctx, req, pref)

	if err != nil {
		return nil, resp, err
	}

	return pref.Pref, resp, nil
}
