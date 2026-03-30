package main

import (
	"context"
	"encoding/json"
	"strings"
	"time"
)

// GetTestFlight fetches beta groups and tester counts concurrently.
func (a *App) GetTestFlight(appID string) (TestFlightResponse, error) {
	if strings.TrimSpace(appID) == "" {
		return TestFlightResponse{Error: "app ID is required"}, nil
	}
	defer configGuard()()
	ascPath, err := a.resolveASCPath()
	if err != nil {
		return TestFlightResponse{Error: err.Error()}, nil
	}
	ctx, cancel := context.WithTimeout(a.contextOrBackground(), 30*time.Second)
	defer cancel()

	// 1. Fetch groups
	cmd := a.newASCCommand(ctx, ascPath, "testflight", "groups", "list", "--app", appID, "--output", "json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return TestFlightResponse{Error: strings.TrimSpace(string(out))}, nil
	}

	type rawGroup struct {
		ID         string `json:"id"`
		Attributes struct {
			Name            string `json:"name"`
			IsInternalGroup bool   `json:"isInternalGroup"`
			PublicLink      string `json:"publicLink"`
			FeedbackEnabled bool   `json:"feedbackEnabled"`
			CreatedDate     string `json:"createdDate"`
		} `json:"attributes"`
		Relationships struct {
			BetaTesters struct {
				Links struct {
					Related string `json:"related"`
				} `json:"links"`
			} `json:"betaTesters"`
		} `json:"relationships"`
	}
	var groupEnv struct {
		Data []rawGroup `json:"data"`
	}
	if json.Unmarshal(out, &groupEnv) != nil {
		return TestFlightResponse{Error: "failed to parse groups"}, nil
	}

	// 2. Fetch tester count per group concurrently (just need meta.paging.total)
	type countResult struct {
		idx   int
		count int
	}
	countCh := make(chan countResult, len(groupEnv.Data))
	for i, g := range groupEnv.Data {
		go func(idx int, groupID string) {
			cmd := a.newASCCommand(ctx, ascPath, "testflight", "testers", "list",
				"--group", groupID, "--limit", "1", "--output", "json")
			out, err := cmd.CombinedOutput()
			if err != nil {
				countCh <- countResult{idx: idx, count: 0}
				return
			}
			var env struct {
				Meta struct {
					Paging struct {
						Total int `json:"total"`
					} `json:"paging"`
				} `json:"meta"`
			}
			if json.Unmarshal(out, &env) == nil {
				countCh <- countResult{idx: idx, count: env.Meta.Paging.Total}
			} else {
				countCh <- countResult{idx: idx, count: 0}
			}
		}(i, g.ID)
	}

	groups := make([]BetaGroup, len(groupEnv.Data))
	for i, g := range groupEnv.Data {
		groups[i] = BetaGroup{
			ID:              g.ID,
			Name:            g.Attributes.Name,
			IsInternal:      g.Attributes.IsInternalGroup,
			PublicLink:      g.Attributes.PublicLink,
			FeedbackEnabled: g.Attributes.FeedbackEnabled,
			CreatedDate:     g.Attributes.CreatedDate,
		}
	}

	for range groupEnv.Data {
		r := <-countCh
		groups[r.idx].TesterCount = r.count
	}

	return TestFlightResponse{Groups: groups}, nil
}

// GetTestFlightTesters fetches ALL testers for a specific group (paginated).
func (a *App) GetTestFlightTesters(groupID string) (TestFlightResponse, error) {
	if strings.TrimSpace(groupID) == "" {
		return TestFlightResponse{Error: "group ID is required"}, nil
	}
	defer configGuard()()
	ascPath, err := a.resolveASCPath()
	if err != nil {
		return TestFlightResponse{Error: err.Error()}, nil
	}
	ctx, cancel := context.WithTimeout(a.contextOrBackground(), 120*time.Second)
	defer cancel()

	cmd := a.newASCCommand(ctx, ascPath, "testflight", "testers", "list",
		"--group", groupID, "--paginate", "--output", "json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return TestFlightResponse{Error: strings.TrimSpace(string(out))}, nil
	}

	type rawTester struct {
		Attributes struct {
			Email      string `json:"email"`
			FirstName  string `json:"firstName"`
			LastName   string `json:"lastName"`
			InviteType string `json:"inviteType"`
			State      string `json:"state"`
		} `json:"attributes"`
	}
	var env struct {
		Data []rawTester `json:"data"`
	}
	if json.Unmarshal(out, &env) != nil {
		return TestFlightResponse{Error: "failed to parse testers"}, nil
	}

	testers := make([]BetaTester, 0, len(env.Data))
	for _, t := range env.Data {
		testers = append(testers, BetaTester{
			Email:      t.Attributes.Email,
			FirstName:  t.Attributes.FirstName,
			LastName:   t.Attributes.LastName,
			InviteType: t.Attributes.InviteType,
			State:      t.Attributes.State,
		})
	}
	return TestFlightResponse{Testers: testers}, nil
}

// GetFeedback fetches TestFlight feedback list, then enriches each with detail view concurrently.
func (a *App) GetFeedback(appID string) (FeedbackResponse, error) {
	if strings.TrimSpace(appID) == "" {
		return FeedbackResponse{Error: "app ID is required"}, nil
	}
	defer configGuard()()
	ascPath, err := a.resolveASCPath()
	if err != nil {
		return FeedbackResponse{Error: err.Error()}, nil
	}
	ctx, cancel := context.WithTimeout(a.contextOrBackground(), 60*time.Second)
	defer cancel()

	// Fetch feedback list with screenshots
	cmd := a.newASCCommand(ctx, ascPath, "testflight", "feedback", "list",
		"--app", appID, "--include-screenshots", "--sort", "-createdDate", "--paginate", "--output", "json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return FeedbackResponse{Error: strings.TrimSpace(string(out))}, nil
	}

	type rawScreenshot struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	}
	type rawFeedback struct {
		ID         string `json:"id"`
		Attributes struct {
			Comment      string          `json:"comment"`
			Email        string          `json:"email"`
			DeviceModel  string          `json:"deviceModel"`
			OSVersion    string          `json:"osVersion"`
			AppPlatform  string          `json:"appPlatform"`
			CreatedDate  string          `json:"createdDate"`
			DeviceFamily string          `json:"deviceFamily"`
			Screenshots  []rawScreenshot `json:"screenshots"`
		} `json:"attributes"`
	}
	var listEnv struct {
		Data []rawFeedback `json:"data"`
		Meta struct {
			Paging struct {
				Total int `json:"total"`
			} `json:"paging"`
		} `json:"meta"`
	}
	if json.Unmarshal(out, &listEnv) != nil {
		return FeedbackResponse{Error: "failed to parse feedback list"}, nil
	}

	// Enrich each feedback item with detail view (concurrent, best-effort)
	type detailResult struct {
		idx    int
		locale string
		tz     string
		conn   string
		batt   int
	}
	ch := make(chan detailResult, len(listEnv.Data))
	for i, fb := range listEnv.Data {
		go func(idx int, fbID string) {
			cmd := a.newASCCommand(ctx, ascPath, "testflight", "feedback", "view",
				"--submission-id", fbID, "--output", "json")
			out, err := cmd.CombinedOutput()
			if err != nil {
				ch <- detailResult{idx: idx}
				return
			}
			var env struct {
				Data struct {
					Attributes struct {
						Locale         string `json:"locale"`
						TimeZone       string `json:"timeZone"`
						ConnectionType string `json:"connectionType"`
						Battery        int    `json:"batteryPercentage"`
					} `json:"attributes"`
				} `json:"data"`
			}
			if json.Unmarshal(out, &env) == nil {
				ch <- detailResult{
					idx:    idx,
					locale: env.Data.Attributes.Locale,
					tz:     env.Data.Attributes.TimeZone,
					conn:   env.Data.Attributes.ConnectionType,
					batt:   env.Data.Attributes.Battery,
				}
			} else {
				ch <- detailResult{idx: idx}
			}
		}(i, fb.ID)
	}

	items := make([]FeedbackItem, len(listEnv.Data))
	for i, fb := range listEnv.Data {
		var shots []FeedbackScreenshot
		for _, s := range fb.Attributes.Screenshots {
			shots = append(shots, FeedbackScreenshot{URL: s.URL, Width: s.Width, Height: s.Height})
		}
		items[i] = FeedbackItem{
			ID:           fb.ID,
			Comment:      fb.Attributes.Comment,
			Email:        fb.Attributes.Email,
			DeviceModel:  fb.Attributes.DeviceModel,
			DeviceFamily: fb.Attributes.DeviceFamily,
			OSVersion:    fb.Attributes.OSVersion,
			AppPlatform:  fb.Attributes.AppPlatform,
			CreatedDate:  fb.Attributes.CreatedDate,
			Screenshots:  shots,
		}
	}
	for range listEnv.Data {
		r := <-ch
		items[r.idx].Locale = r.locale
		items[r.idx].TimeZone = r.tz
		items[r.idx].ConnectionType = r.conn
		items[r.idx].Battery = r.batt
	}

	return FeedbackResponse{Feedback: items, Total: listEnv.Meta.Paging.Total}, nil
}
