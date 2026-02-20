package status

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
)

type includeSet struct {
	builds        bool
	testflight    bool
	appstore      bool
	submission    bool
	review        bool
	phasedRelease bool
	links         bool
}

type dashboardResponse struct {
	App           statusApp             `json:"app"`
	Builds        *buildsSection        `json:"builds,omitempty"`
	TestFlight    *testFlightSection    `json:"testflight,omitempty"`
	AppStore      *appStoreSection      `json:"appstore,omitempty"`
	Submission    *submissionSection    `json:"submission,omitempty"`
	Review        *reviewSection        `json:"review,omitempty"`
	PhasedRelease *phasedReleaseSection `json:"phasedRelease,omitempty"`
	Links         *linksSection         `json:"links,omitempty"`
}

type statusApp struct {
	ID       string `json:"id"`
	BundleID string `json:"bundleId"`
	Name     string `json:"name"`
}

type buildsSection struct {
	Latest *latestBuild `json:"latest,omitempty"`
}

type latestBuild struct {
	ID              string `json:"id"`
	Version         string `json:"version,omitempty"`
	BuildNumber     string `json:"buildNumber"`
	ProcessingState string `json:"processingState,omitempty"`
	UploadedDate    string `json:"uploadedDate,omitempty"`
	Platform        string `json:"platform,omitempty"`
}

type testFlightSection struct {
	LatestDistributedBuildID string `json:"latestDistributedBuildId,omitempty"`
	BetaReviewState          string `json:"betaReviewState,omitempty"`
	ExternalBuildState       string `json:"externalBuildState,omitempty"`
	SubmittedDate            string `json:"submittedDate,omitempty"`
}

type appStoreSection struct {
	VersionID   string `json:"versionId,omitempty"`
	Version     string `json:"version,omitempty"`
	State       string `json:"state,omitempty"`
	Platform    string `json:"platform,omitempty"`
	CreatedDate string `json:"createdDate,omitempty"`
}

type submissionSection struct {
	InFlight       bool     `json:"inFlight"`
	BlockingIssues []string `json:"blockingIssues"`
}

type reviewSection struct {
	LatestSubmissionID string `json:"latestSubmissionId,omitempty"`
	State              string `json:"state,omitempty"`
	SubmittedDate      string `json:"submittedDate,omitempty"`
	Platform           string `json:"platform,omitempty"`
}

type phasedReleaseSection struct {
	Configured         bool   `json:"configured"`
	ID                 string `json:"id,omitempty"`
	State              string `json:"state,omitempty"`
	StartDate          string `json:"startDate,omitempty"`
	CurrentDayNumber   int    `json:"currentDayNumber,omitempty"`
	TotalPauseDuration int    `json:"totalPauseDuration,omitempty"`
}

type linksSection struct {
	AppStoreConnect string `json:"appStoreConnect"`
	TestFlight      string `json:"testFlight"`
	Review          string `json:"review"`
}

type relationshipReference struct {
	Data asc.ResourceData `json:"data"`
}

type sectionTask struct {
	name string
	run  func() error
}

var allowedIncludes = []string{
	"builds",
	"testflight",
	"appstore",
	"submission",
	"review",
	"phased-release",
	"links",
}

// StatusCommand returns the root status dashboard command.
func StatusCommand() *ffcli.Command {
	fs := flag.NewFlagSet("status", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (required, or ASC_APP_ID env)")
	include := fs.String("include", "", "Comma-separated sections: builds,testflight,appstore,submission,review,phased-release,links")
	output := shared.BindOutputFlags(fs)

	return &ffcli.Command{
		Name:       "status",
		ShortUsage: "asc status [flags]",
		ShortHelp:  "Show a release pipeline dashboard for an app.",
		LongHelp: `Show a release pipeline dashboard for an app.

This command aggregates release signals into one deterministic payload for CI,
agents, and human review.

Examples:
  asc status --app "123456789"
  asc status --app "123456789" --include builds,testflight,submission
  asc status --app "123456789" --output table`,
		FlagSet:   fs,
		UsageFunc: shared.DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) > 0 {
				fmt.Fprintln(os.Stderr, "Error: status does not accept positional arguments")
				return flag.ErrHelp
			}

			resolvedAppID := shared.ResolveAppID(*appID)
			if resolvedAppID == "" {
				fmt.Fprintln(os.Stderr, "Error: --app is required (or set ASC_APP_ID)")
				return flag.ErrHelp
			}

			includes, err := parseInclude(*include)
			if err != nil {
				return shared.UsageError(err.Error())
			}

			client, err := shared.GetASCClient()
			if err != nil {
				return fmt.Errorf("status: %w", err)
			}

			requestCtx, cancel := shared.ContextWithTimeout(ctx)
			defer cancel()

			resp, err := collectDashboard(requestCtx, client, resolvedAppID, includes)
			if err != nil {
				return fmt.Errorf("status: %w", err)
			}

			return shared.PrintOutputWithRenderers(
				resp,
				*output.Output,
				*output.Pretty,
				func() error { renderTable(resp); return nil },
				func() error { renderMarkdown(resp); return nil },
			)
		},
	}
}

func parseInclude(value string) (includeSet, error) {
	parts := shared.SplitCSV(strings.ToLower(strings.TrimSpace(value)))
	if len(parts) == 0 {
		return includeSet{
			builds:        true,
			testflight:    true,
			appstore:      true,
			submission:    true,
			review:        true,
			phasedRelease: true,
			links:         true,
		}, nil
	}

	includes := includeSet{}
	for _, part := range parts {
		switch part {
		case "builds":
			includes.builds = true
		case "testflight":
			includes.testflight = true
		case "appstore":
			includes.appstore = true
		case "submission":
			includes.submission = true
		case "review":
			includes.review = true
		case "phased-release":
			includes.phasedRelease = true
		case "links":
			includes.links = true
		default:
			return includeSet{}, fmt.Errorf("--include contains unsupported section %q (allowed: %s)", part, strings.Join(allowedIncludes, ","))
		}
	}

	return includes, nil
}

func collectDashboard(ctx context.Context, client *asc.Client, appID string, includes includeSet) (*dashboardResponse, error) {
	appResp, err := client.GetApp(ctx, appID)
	if err != nil {
		return nil, err
	}

	resp := &dashboardResponse{
		App: statusApp{
			ID:       appResp.Data.ID,
			BundleID: appResp.Data.Attributes.BundleID,
			Name:     appResp.Data.Attributes.Name,
		},
	}

	if includes.links {
		resp.Links = &linksSection{
			AppStoreConnect: fmt.Sprintf("https://appstoreconnect.apple.com/apps/%s", appID),
			TestFlight:      fmt.Sprintf("https://appstoreconnect.apple.com/apps/%s/testflight/ios", appID),
			Review:          fmt.Sprintf("https://appstoreconnect.apple.com/apps/%s/appstore/review", appID),
		}
	}

	var tasks []sectionTask

	if includes.builds || includes.testflight {
		tasks = append(tasks, sectionTask{
			name: "builds/testflight",
			run: func() error {
				return fillBuildsAndTestFlight(ctx, client, appID, includes, resp)
			},
		})
	}
	if includes.appstore || includes.phasedRelease {
		tasks = append(tasks, sectionTask{
			name: "appstore/phased-release",
			run: func() error {
				return fillAppStoreAndPhasedRelease(ctx, client, appID, includes, resp)
			},
		})
	}
	if includes.submission || includes.review {
		tasks = append(tasks, sectionTask{
			name: "submission/review",
			run: func() error {
				return fillSubmissionAndReview(ctx, client, appID, includes, resp)
			},
		})
	}

	if err := runTasks(tasks, 3); err != nil {
		return nil, err
	}

	return resp, nil
}

func runTasks(tasks []sectionTask, limit int) error {
	if len(tasks) == 0 {
		return nil
	}

	if limit < 1 {
		limit = 1
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, limit)
	errCh := make(chan error, len(tasks))

	for _, task := range tasks {
		current := task
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			if err := current.run(); err != nil {
				errCh <- fmt.Errorf("%s: %w", current.name, err)
			}
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		return err
	}
	return nil
}

func fillBuildsAndTestFlight(ctx context.Context, client *asc.Client, appID string, includes includeSet, resp *dashboardResponse) error {
	buildsResp, err := client.GetBuilds(ctx, appID, asc.WithBuildsSort("-uploadedDate"), asc.WithBuildsLimit(50))
	if err != nil {
		return err
	}

	var latest *asc.Resource[asc.BuildAttributes]
	if len(buildsResp.Data) > 0 {
		latest = &buildsResp.Data[0]
	}

	if includes.builds {
		section := &buildsSection{}
		if latest != nil {
			entry := &latestBuild{
				ID:              latest.ID,
				BuildNumber:     latest.Attributes.Version,
				ProcessingState: latest.Attributes.ProcessingState,
				UploadedDate:    latest.Attributes.UploadedDate,
			}

			preRelease, preErr := client.GetBuildPreReleaseVersion(ctx, latest.ID)
			if preErr != nil {
				if !asc.IsNotFound(preErr) {
					return preErr
				}
			} else {
				entry.Version = preRelease.Data.Attributes.Version
				entry.Platform = string(preRelease.Data.Attributes.Platform)
			}
			section.Latest = entry
		}
		resp.Builds = section
	}

	if !includes.testflight {
		return nil
	}

	section := &testFlightSection{}
	if len(buildsResp.Data) == 0 {
		resp.TestFlight = section
		return nil
	}

	buildIDs := make([]string, 0, len(buildsResp.Data))
	for _, build := range buildsResp.Data {
		buildIDs = append(buildIDs, build.ID)
	}

	betaDetails, err := client.GetBuildBetaDetails(ctx,
		asc.WithBuildBetaDetailsBuildIDs(buildIDs),
		asc.WithBuildBetaDetailsLimit(200),
	)
	if err != nil {
		return err
	}

	externalStateByBuild := make(map[string]string, len(betaDetails.Data))
	for _, detail := range betaDetails.Data {
		buildID, relErr := relationshipResourceID(detail.Relationships, "build")
		if relErr != nil {
			return relErr
		}
		externalStateByBuild[buildID] = strings.TrimSpace(detail.Attributes.ExternalBuildState)
	}

	for _, build := range buildsResp.Data {
		state := strings.ToUpper(strings.TrimSpace(externalStateByBuild[build.ID]))
		if isDistributedState(state) {
			section.LatestDistributedBuildID = build.ID
			section.ExternalBuildState = state
			break
		}
	}

	reviewSubmissions, err := client.GetBetaAppReviewSubmissions(ctx,
		asc.WithBetaAppReviewSubmissionsBuildIDs(buildIDs),
		asc.WithBetaAppReviewSubmissionsLimit(200),
	)
	if err != nil {
		return err
	}
	latestReviewSubmission := selectLatestBetaReviewSubmission(reviewSubmissions.Data)
	if latestReviewSubmission != nil {
		section.BetaReviewState = latestReviewSubmission.Attributes.BetaReviewState
		section.SubmittedDate = latestReviewSubmission.Attributes.SubmittedDate
	}

	resp.TestFlight = section
	return nil
}

func fillAppStoreAndPhasedRelease(ctx context.Context, client *asc.Client, appID string, includes includeSet, resp *dashboardResponse) error {
	versions, err := client.GetAppStoreVersions(ctx, appID, asc.WithAppStoreVersionsLimit(200))
	if err != nil {
		return err
	}

	latestVersion := selectLatestAppStoreVersion(versions.Data)
	if includes.appstore {
		section := &appStoreSection{}
		if latestVersion != nil {
			section.VersionID = latestVersion.ID
			section.Version = latestVersion.Attributes.VersionString
			section.State = shared.ResolveAppStoreVersionState(latestVersion.Attributes)
			section.Platform = string(latestVersion.Attributes.Platform)
			section.CreatedDate = latestVersion.Attributes.CreatedDate
		}
		resp.AppStore = section
	}

	if !includes.phasedRelease {
		return nil
	}

	phased := &phasedReleaseSection{Configured: false}
	if latestVersion != nil {
		phaseResp, phaseErr := client.GetAppStoreVersionPhasedRelease(ctx, latestVersion.ID)
		if phaseErr != nil {
			if !asc.IsNotFound(phaseErr) {
				return phaseErr
			}
		} else {
			phased.Configured = true
			phased.ID = phaseResp.Data.ID
			phased.State = string(phaseResp.Data.Attributes.PhasedReleaseState)
			phased.StartDate = phaseResp.Data.Attributes.StartDate
			phased.CurrentDayNumber = phaseResp.Data.Attributes.CurrentDayNumber
			phased.TotalPauseDuration = phaseResp.Data.Attributes.TotalPauseDuration
		}
	}

	resp.PhasedRelease = phased
	return nil
}

func fillSubmissionAndReview(ctx context.Context, client *asc.Client, appID string, includes includeSet, resp *dashboardResponse) error {
	submissions, err := client.GetReviewSubmissions(ctx, appID, asc.WithReviewSubmissionsLimit(200))
	if err != nil {
		return err
	}

	if includes.submission {
		section := &submissionSection{
			InFlight:       false,
			BlockingIssues: []string{},
		}
		for _, submission := range submissions.Data {
			state := string(submission.Attributes.SubmissionState)
			if isInFlightSubmissionState(state) {
				section.InFlight = true
			}
			if strings.EqualFold(state, string(asc.ReviewSubmissionStateUnresolvedIssues)) {
				section.BlockingIssues = append(section.BlockingIssues, fmt.Sprintf("submission %s has unresolved issues", submission.ID))
			}
		}
		slices.Sort(section.BlockingIssues)
		resp.Submission = section
	}

	if includes.review {
		section := &reviewSection{}
		latest := selectLatestReviewSubmission(submissions.Data)
		if latest != nil {
			section.LatestSubmissionID = latest.ID
			section.State = string(latest.Attributes.SubmissionState)
			section.SubmittedDate = latest.Attributes.SubmittedDate
			section.Platform = string(latest.Attributes.Platform)
		}
		resp.Review = section
	}

	return nil
}

func relationshipResourceID(relationships json.RawMessage, key string) (string, error) {
	if len(relationships) == 0 {
		return "", fmt.Errorf("missing %s relationship", key)
	}

	var references map[string]relationshipReference
	if err := json.Unmarshal(relationships, &references); err != nil {
		return "", fmt.Errorf("parse relationships: %w", err)
	}

	reference, ok := references[key]
	if !ok {
		return "", fmt.Errorf("missing %s relationship", key)
	}

	id := strings.TrimSpace(reference.Data.ID)
	if id == "" {
		return "", fmt.Errorf("missing %s relationship id", key)
	}

	return id, nil
}

func selectLatestAppStoreVersion(versions []asc.Resource[asc.AppStoreVersionAttributes]) *asc.Resource[asc.AppStoreVersionAttributes] {
	if len(versions) == 0 {
		return nil
	}

	best := versions[0]
	for _, current := range versions[1:] {
		if current.Attributes.CreatedDate > best.Attributes.CreatedDate {
			best = current
			continue
		}
		if current.Attributes.CreatedDate == best.Attributes.CreatedDate && current.ID > best.ID {
			best = current
		}
	}
	return &best
}

func selectLatestReviewSubmission(submissions []asc.ReviewSubmissionResource) *asc.ReviewSubmissionResource {
	if len(submissions) == 0 {
		return nil
	}

	best := submissions[0]
	for _, current := range submissions[1:] {
		if current.Attributes.SubmittedDate > best.Attributes.SubmittedDate {
			best = current
			continue
		}
		if current.Attributes.SubmittedDate == best.Attributes.SubmittedDate && current.ID > best.ID {
			best = current
		}
	}
	return &best
}

func selectLatestBetaReviewSubmission(submissions []asc.Resource[asc.BetaAppReviewSubmissionAttributes]) *asc.Resource[asc.BetaAppReviewSubmissionAttributes] {
	if len(submissions) == 0 {
		return nil
	}

	best := submissions[0]
	for _, current := range submissions[1:] {
		if current.Attributes.SubmittedDate > best.Attributes.SubmittedDate {
			best = current
			continue
		}
		if current.Attributes.SubmittedDate == best.Attributes.SubmittedDate && current.ID > best.ID {
			best = current
		}
	}
	return &best
}

func isDistributedState(state string) bool {
	switch strings.ToUpper(strings.TrimSpace(state)) {
	case "IN_BETA_TESTING", "READY_FOR_TESTING":
		return true
	default:
		return false
	}
}

func isInFlightSubmissionState(state string) bool {
	switch strings.ToUpper(strings.TrimSpace(state)) {
	case string(asc.ReviewSubmissionStateReadyForReview),
		string(asc.ReviewSubmissionStateWaitingForReview),
		string(asc.ReviewSubmissionStateInReview),
		string(asc.ReviewSubmissionStateUnresolvedIssues),
		string(asc.ReviewSubmissionStateCanceling):
		return true
	default:
		return false
	}
}

func renderTable(resp *dashboardResponse) {
	renderDashboard(resp, false)
}

func renderMarkdown(resp *dashboardResponse) {
	renderDashboard(resp, true)
}

func renderDashboard(resp *dashboardResponse, markdown bool) {
	renderSection := func(title string, rows [][]string) {
		if markdown {
			fmt.Fprintf(os.Stdout, "### %s\n\n", title)
			asc.RenderMarkdown([]string{"field", "value"}, rows)
			fmt.Fprintln(os.Stdout)
			return
		}

		fmt.Fprintf(os.Stdout, "%s\n", strings.ToUpper(title))
		asc.RenderTable([]string{"field", "value"}, rows)
		fmt.Fprintln(os.Stdout)
	}

	renderSection("App", [][]string{
		{"id", resp.App.ID},
		{"name", resp.App.Name},
		{"bundleId", resp.App.BundleID},
	})

	if resp.Builds != nil {
		rows := [][]string{}
		if resp.Builds.Latest == nil {
			rows = append(rows, []string{"latest", "none"})
		} else {
			rows = append(rows,
				[]string{"latest.id", resp.Builds.Latest.ID},
				[]string{"latest.version", resp.Builds.Latest.Version},
				[]string{"latest.buildNumber", resp.Builds.Latest.BuildNumber},
				[]string{"latest.processingState", resp.Builds.Latest.ProcessingState},
				[]string{"latest.uploadedDate", resp.Builds.Latest.UploadedDate},
				[]string{"latest.platform", resp.Builds.Latest.Platform},
			)
		}
		renderSection("Builds", rows)
	}

	if resp.TestFlight != nil {
		renderSection("TestFlight", [][]string{
			{"latestDistributedBuildId", resp.TestFlight.LatestDistributedBuildID},
			{"betaReviewState", resp.TestFlight.BetaReviewState},
			{"externalBuildState", resp.TestFlight.ExternalBuildState},
			{"submittedDate", resp.TestFlight.SubmittedDate},
		})
	}

	if resp.AppStore != nil {
		renderSection("AppStore", [][]string{
			{"versionId", resp.AppStore.VersionID},
			{"version", resp.AppStore.Version},
			{"state", resp.AppStore.State},
			{"platform", resp.AppStore.Platform},
			{"createdDate", resp.AppStore.CreatedDate},
		})
	}

	if resp.Submission != nil {
		blocking := "none"
		if len(resp.Submission.BlockingIssues) > 0 {
			blocking = strings.Join(resp.Submission.BlockingIssues, "; ")
		}
		renderSection("Submission", [][]string{
			{"inFlight", fmt.Sprintf("%t", resp.Submission.InFlight)},
			{"blockingIssues", blocking},
		})
	}

	if resp.Review != nil {
		renderSection("Review", [][]string{
			{"latestSubmissionId", resp.Review.LatestSubmissionID},
			{"state", resp.Review.State},
			{"submittedDate", resp.Review.SubmittedDate},
			{"platform", resp.Review.Platform},
		})
	}

	if resp.PhasedRelease != nil {
		renderSection("PhasedRelease", [][]string{
			{"configured", fmt.Sprintf("%t", resp.PhasedRelease.Configured)},
			{"id", resp.PhasedRelease.ID},
			{"state", resp.PhasedRelease.State},
			{"startDate", resp.PhasedRelease.StartDate},
			{"currentDayNumber", fmt.Sprintf("%d", resp.PhasedRelease.CurrentDayNumber)},
			{"totalPauseDuration", fmt.Sprintf("%d", resp.PhasedRelease.TotalPauseDuration)},
		})
	}

	if resp.Links != nil {
		renderSection("Links", [][]string{
			{"appStoreConnect", resp.Links.AppStoreConnect},
			{"testFlight", resp.Links.TestFlight},
			{"review", resp.Links.Review},
		})
	}
}
