package asc

import (
	"net/url"
	"strings"
)

// GameCenterAchievementAttributes represents a Game Center achievement resource.
type GameCenterAchievementAttributes struct {
	ReferenceName    string `json:"referenceName"`
	VendorIdentifier string `json:"vendorIdentifier"`
	Points           int    `json:"points"`
	ShowBeforeEarned bool   `json:"showBeforeEarned"`
	Repeatable       bool   `json:"repeatable"`
	Archived         bool   `json:"archived,omitempty"`
}

// GameCenterAchievementCreateAttributes describes attributes for creating an achievement.
type GameCenterAchievementCreateAttributes struct {
	ReferenceName    string `json:"referenceName"`
	VendorIdentifier string `json:"vendorIdentifier"`
	Points           int    `json:"points"`
	ShowBeforeEarned bool   `json:"showBeforeEarned"`
	Repeatable       bool   `json:"repeatable"`
}

// GameCenterAchievementUpdateAttributes describes attributes for updating an achievement.
type GameCenterAchievementUpdateAttributes struct {
	ReferenceName    *string `json:"referenceName,omitempty"`
	Points           *int    `json:"points,omitempty"`
	ShowBeforeEarned *bool   `json:"showBeforeEarned,omitempty"`
	Repeatable       *bool   `json:"repeatable,omitempty"`
	Archived         *bool   `json:"archived,omitempty"`
}

// GameCenterAchievementRelationships describes relationships for achievements.
type GameCenterAchievementRelationships struct {
	GameCenterDetail *Relationship `json:"gameCenterDetail"`
}

// GameCenterAchievementCreateData is the data portion of an achievement create request.
type GameCenterAchievementCreateData struct {
	Type          ResourceType                          `json:"type"`
	Attributes    GameCenterAchievementCreateAttributes `json:"attributes"`
	Relationships *GameCenterAchievementRelationships   `json:"relationships,omitempty"`
}

// GameCenterAchievementCreateRequest is a request to create an achievement.
type GameCenterAchievementCreateRequest struct {
	Data GameCenterAchievementCreateData `json:"data"`
}

// GameCenterAchievementUpdateData is the data portion of an achievement update request.
type GameCenterAchievementUpdateData struct {
	Type       ResourceType                           `json:"type"`
	ID         string                                 `json:"id"`
	Attributes *GameCenterAchievementUpdateAttributes `json:"attributes,omitempty"`
}

// GameCenterAchievementUpdateRequest is a request to update an achievement.
type GameCenterAchievementUpdateRequest struct {
	Data GameCenterAchievementUpdateData `json:"data"`
}

// GameCenterDetailAttributes represents a Game Center detail resource.
type GameCenterDetailAttributes struct {
	ArcadeEnabled                      bool `json:"arcadeEnabled,omitempty"`
	ChallengeEnabled                   bool `json:"challengeEnabled,omitempty"`
	LeaderboardSetEnabled              bool `json:"leaderboardSetEnabled,omitempty"`
	LeaderboardEnabled                 bool `json:"leaderboardEnabled,omitempty"`
	AchievementEnabled                 bool `json:"achievementEnabled,omitempty"`
	MultiplayerSessionEnabled          bool `json:"multiplayerSessionEnabled,omitempty"`
	MultiplayerTurnBasedSessionEnabled bool `json:"multiplayerTurnBasedSessionEnabled,omitempty"`
}

// GameCenterAchievementsResponse is the response from achievement list endpoints.
type GameCenterAchievementsResponse = Response[GameCenterAchievementAttributes]

// GameCenterAchievementResponse is the response from achievement detail endpoints.
type GameCenterAchievementResponse = SingleResponse[GameCenterAchievementAttributes]

// GameCenterDetailResponse is the response from Game Center detail endpoints.
type GameCenterDetailResponse = SingleResponse[GameCenterDetailAttributes]

// GameCenterAchievementDeleteResult represents CLI output for achievement deletions.
type GameCenterAchievementDeleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// GCAchievementsOption is a functional option for GetGameCenterAchievements.
type GCAchievementsOption func(*gcAchievementsQuery)

type gcAchievementsQuery struct {
	listQuery
}

// WithGCAchievementsLimit sets the max number of achievements to return.
func WithGCAchievementsLimit(limit int) GCAchievementsOption {
	return func(q *gcAchievementsQuery) {
		if limit > 0 {
			q.limit = limit
		}
	}
}

// WithGCAchievementsNextURL uses a next page URL directly.
func WithGCAchievementsNextURL(next string) GCAchievementsOption {
	return func(q *gcAchievementsQuery) {
		if strings.TrimSpace(next) != "" {
			q.nextURL = strings.TrimSpace(next)
		}
	}
}

func buildGCAchievementsQuery(query *gcAchievementsQuery) string {
	values := url.Values{}
	addLimit(values, query.limit)
	return values.Encode()
}

// GameCenterLeaderboardAttributes represents a Game Center leaderboard resource.
type GameCenterLeaderboardAttributes struct {
	ReferenceName       string `json:"referenceName"`
	VendorIdentifier    string `json:"vendorIdentifier"`
	DefaultFormatter    string `json:"defaultFormatter"`
	ScoreSortType       string `json:"scoreSortType"`
	ScoreRangeStart     string `json:"scoreRangeStart,omitempty"`
	ScoreRangeEnd       string `json:"scoreRangeEnd,omitempty"`
	RecurrenceStartDate string `json:"recurrenceStartDate,omitempty"`
	RecurrenceDuration  string `json:"recurrenceDuration,omitempty"`
	RecurrenceRule      string `json:"recurrenceRule,omitempty"`
	SubmissionType      string `json:"submissionType"`
	Archived            bool   `json:"archived,omitempty"`
}

// GameCenterLeaderboardCreateAttributes describes attributes for creating a leaderboard.
type GameCenterLeaderboardCreateAttributes struct {
	ReferenceName       string `json:"referenceName"`
	VendorIdentifier    string `json:"vendorIdentifier"`
	DefaultFormatter    string `json:"defaultFormatter"`
	ScoreSortType       string `json:"scoreSortType"`
	ScoreRangeStart     string `json:"scoreRangeStart,omitempty"`
	ScoreRangeEnd       string `json:"scoreRangeEnd,omitempty"`
	RecurrenceStartDate string `json:"recurrenceStartDate,omitempty"`
	RecurrenceDuration  string `json:"recurrenceDuration,omitempty"`
	RecurrenceRule      string `json:"recurrenceRule,omitempty"`
	SubmissionType      string `json:"submissionType"`
}

// GameCenterLeaderboardUpdateAttributes describes attributes for updating a leaderboard.
type GameCenterLeaderboardUpdateAttributes struct {
	ReferenceName       *string `json:"referenceName,omitempty"`
	DefaultFormatter    *string `json:"defaultFormatter,omitempty"`
	ScoreSortType       *string `json:"scoreSortType,omitempty"`
	ScoreRangeStart     *string `json:"scoreRangeStart,omitempty"`
	ScoreRangeEnd       *string `json:"scoreRangeEnd,omitempty"`
	RecurrenceStartDate *string `json:"recurrenceStartDate,omitempty"`
	RecurrenceDuration  *string `json:"recurrenceDuration,omitempty"`
	RecurrenceRule      *string `json:"recurrenceRule,omitempty"`
	SubmissionType      *string `json:"submissionType,omitempty"`
	Archived            *bool   `json:"archived,omitempty"`
}

// GameCenterLeaderboardRelationships describes relationships for leaderboards.
type GameCenterLeaderboardRelationships struct {
	GameCenterDetail *Relationship `json:"gameCenterDetail"`
}

// GameCenterLeaderboardCreateData is the data portion of a leaderboard create request.
type GameCenterLeaderboardCreateData struct {
	Type          ResourceType                          `json:"type"`
	Attributes    GameCenterLeaderboardCreateAttributes `json:"attributes"`
	Relationships *GameCenterLeaderboardRelationships   `json:"relationships,omitempty"`
}

// GameCenterLeaderboardCreateRequest is a request to create a leaderboard.
type GameCenterLeaderboardCreateRequest struct {
	Data GameCenterLeaderboardCreateData `json:"data"`
}

// GameCenterLeaderboardUpdateData is the data portion of a leaderboard update request.
type GameCenterLeaderboardUpdateData struct {
	Type       ResourceType                           `json:"type"`
	ID         string                                 `json:"id"`
	Attributes *GameCenterLeaderboardUpdateAttributes `json:"attributes,omitempty"`
}

// GameCenterLeaderboardUpdateRequest is a request to update a leaderboard.
type GameCenterLeaderboardUpdateRequest struct {
	Data GameCenterLeaderboardUpdateData `json:"data"`
}

// GameCenterLeaderboardsResponse is the response from leaderboard list endpoints.
type GameCenterLeaderboardsResponse = Response[GameCenterLeaderboardAttributes]

// GameCenterLeaderboardResponse is the response from leaderboard detail endpoints.
type GameCenterLeaderboardResponse = SingleResponse[GameCenterLeaderboardAttributes]

// GameCenterLeaderboardDeleteResult represents CLI output for leaderboard deletions.
type GameCenterLeaderboardDeleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// Valid leaderboard formatters.
var ValidLeaderboardFormatters = []string{
	"INTEGER",
	"DECIMAL_POINT_1_PLACE",
	"DECIMAL_POINT_2_PLACE",
	"DECIMAL_POINT_3_PLACE",
	"ELAPSED_TIME_MILLISECOND",
	"ELAPSED_TIME_SECOND",
	"ELAPSED_TIME_MINUTE",
	"MONEY_WHOLE",
	"MONEY_POINT_2_PLACE",
}

// Valid leaderboard score sort types.
var ValidScoreSortTypes = []string{
	"ASC",
	"DESC",
}

// Valid leaderboard submission types.
var ValidSubmissionTypes = []string{
	"BEST_SCORE",
	"MOST_RECENT_SCORE",
}

// GCLeaderboardsOption is a functional option for GetGameCenterLeaderboards.
type GCLeaderboardsOption func(*gcLeaderboardsQuery)

type gcLeaderboardsQuery struct {
	listQuery
}

// WithGCLeaderboardsLimit sets the max number of leaderboards to return.
func WithGCLeaderboardsLimit(limit int) GCLeaderboardsOption {
	return func(q *gcLeaderboardsQuery) {
		if limit > 0 {
			q.limit = limit
		}
	}
}

// WithGCLeaderboardsNextURL uses a next page URL directly.
func WithGCLeaderboardsNextURL(next string) GCLeaderboardsOption {
	return func(q *gcLeaderboardsQuery) {
		if strings.TrimSpace(next) != "" {
			q.nextURL = strings.TrimSpace(next)
		}
	}
}

func buildGCLeaderboardsQuery(query *gcLeaderboardsQuery) string {
	values := url.Values{}
	addLimit(values, query.limit)
	return values.Encode()
}

// GameCenterLeaderboardSetAttributes represents a Game Center leaderboard set resource.
type GameCenterLeaderboardSetAttributes struct {
	ReferenceName    string `json:"referenceName"`
	VendorIdentifier string `json:"vendorIdentifier"`
}

// GameCenterLeaderboardSetCreateAttributes describes attributes for creating a leaderboard set.
type GameCenterLeaderboardSetCreateAttributes struct {
	ReferenceName    string `json:"referenceName"`
	VendorIdentifier string `json:"vendorIdentifier"`
}

// GameCenterLeaderboardSetUpdateAttributes describes attributes for updating a leaderboard set.
type GameCenterLeaderboardSetUpdateAttributes struct {
	ReferenceName *string `json:"referenceName,omitempty"`
}

// GameCenterLeaderboardSetRelationships describes relationships for leaderboard sets.
type GameCenterLeaderboardSetRelationships struct {
	GameCenterDetail *Relationship `json:"gameCenterDetail"`
}

// GameCenterLeaderboardSetCreateData is the data portion of a leaderboard set create request.
type GameCenterLeaderboardSetCreateData struct {
	Type          ResourceType                             `json:"type"`
	Attributes    GameCenterLeaderboardSetCreateAttributes `json:"attributes"`
	Relationships *GameCenterLeaderboardSetRelationships   `json:"relationships,omitempty"`
}

// GameCenterLeaderboardSetCreateRequest is a request to create a leaderboard set.
type GameCenterLeaderboardSetCreateRequest struct {
	Data GameCenterLeaderboardSetCreateData `json:"data"`
}

// GameCenterLeaderboardSetUpdateData is the data portion of a leaderboard set update request.
type GameCenterLeaderboardSetUpdateData struct {
	Type       ResourceType                              `json:"type"`
	ID         string                                    `json:"id"`
	Attributes *GameCenterLeaderboardSetUpdateAttributes `json:"attributes,omitempty"`
}

// GameCenterLeaderboardSetUpdateRequest is a request to update a leaderboard set.
type GameCenterLeaderboardSetUpdateRequest struct {
	Data GameCenterLeaderboardSetUpdateData `json:"data"`
}

// GameCenterLeaderboardSetsResponse is the response from leaderboard set list endpoints.
type GameCenterLeaderboardSetsResponse = Response[GameCenterLeaderboardSetAttributes]

// GameCenterLeaderboardSetResponse is the response from leaderboard set detail endpoints.
type GameCenterLeaderboardSetResponse = SingleResponse[GameCenterLeaderboardSetAttributes]

// GameCenterLeaderboardSetDeleteResult represents CLI output for leaderboard set deletions.
type GameCenterLeaderboardSetDeleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// GCLeaderboardSetsOption is a functional option for GetGameCenterLeaderboardSets.
type GCLeaderboardSetsOption func(*gcLeaderboardSetsQuery)

type gcLeaderboardSetsQuery struct {
	listQuery
}

// WithGCLeaderboardSetsLimit sets the max number of leaderboard sets to return.
func WithGCLeaderboardSetsLimit(limit int) GCLeaderboardSetsOption {
	return func(q *gcLeaderboardSetsQuery) {
		if limit > 0 {
			q.limit = limit
		}
	}
}

// WithGCLeaderboardSetsNextURL uses a next page URL directly.
func WithGCLeaderboardSetsNextURL(next string) GCLeaderboardSetsOption {
	return func(q *gcLeaderboardSetsQuery) {
		if strings.TrimSpace(next) != "" {
			q.nextURL = strings.TrimSpace(next)
		}
	}
}

func buildGCLeaderboardSetsQuery(query *gcLeaderboardSetsQuery) string {
	values := url.Values{}
	addLimit(values, query.limit)
	return values.Encode()
}

// GameCenterLeaderboardSetLocalizationAttributes represents a Game Center leaderboard set localization resource.
type GameCenterLeaderboardSetLocalizationAttributes struct {
	Locale string `json:"locale"`
	Name   string `json:"name,omitempty"`
}

// GameCenterLeaderboardSetLocalizationCreateAttributes describes attributes for creating a leaderboard set localization.
type GameCenterLeaderboardSetLocalizationCreateAttributes struct {
	Locale string `json:"locale"`
	Name   string `json:"name,omitempty"`
}

// GameCenterLeaderboardSetLocalizationUpdateAttributes describes attributes for updating a leaderboard set localization.
type GameCenterLeaderboardSetLocalizationUpdateAttributes struct {
	Name *string `json:"name,omitempty"`
}

// GameCenterLeaderboardSetLocalizationRelationships describes relationships for leaderboard set localizations.
type GameCenterLeaderboardSetLocalizationRelationships struct {
	GameCenterLeaderboardSet *Relationship `json:"gameCenterLeaderboardSet"`
}

// GameCenterLeaderboardSetLocalizationCreateData is the data portion of a leaderboard set localization create request.
type GameCenterLeaderboardSetLocalizationCreateData struct {
	Type          ResourceType                                         `json:"type"`
	Attributes    GameCenterLeaderboardSetLocalizationCreateAttributes `json:"attributes"`
	Relationships *GameCenterLeaderboardSetLocalizationRelationships   `json:"relationships,omitempty"`
}

// GameCenterLeaderboardSetLocalizationCreateRequest is a request to create a leaderboard set localization.
type GameCenterLeaderboardSetLocalizationCreateRequest struct {
	Data GameCenterLeaderboardSetLocalizationCreateData `json:"data"`
}

// GameCenterLeaderboardSetLocalizationUpdateData is the data portion of a leaderboard set localization update request.
type GameCenterLeaderboardSetLocalizationUpdateData struct {
	Type       ResourceType                                          `json:"type"`
	ID         string                                                `json:"id"`
	Attributes *GameCenterLeaderboardSetLocalizationUpdateAttributes `json:"attributes,omitempty"`
}

// GameCenterLeaderboardSetLocalizationUpdateRequest is a request to update a leaderboard set localization.
type GameCenterLeaderboardSetLocalizationUpdateRequest struct {
	Data GameCenterLeaderboardSetLocalizationUpdateData `json:"data"`
}

// GameCenterLeaderboardSetLocalizationsResponse is the response from leaderboard set localization list endpoints.
type GameCenterLeaderboardSetLocalizationsResponse = Response[GameCenterLeaderboardSetLocalizationAttributes]

// GameCenterLeaderboardSetLocalizationResponse is the response from leaderboard set localization detail endpoints.
type GameCenterLeaderboardSetLocalizationResponse = SingleResponse[GameCenterLeaderboardSetLocalizationAttributes]

// GameCenterLeaderboardSetLocalizationDeleteResult represents CLI output for leaderboard set localization deletions.
type GameCenterLeaderboardSetLocalizationDeleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// GCLeaderboardSetLocalizationsOption is a functional option for GetGameCenterLeaderboardSetLocalizations.
type GCLeaderboardSetLocalizationsOption func(*gcLeaderboardSetLocalizationsQuery)

type gcLeaderboardSetLocalizationsQuery struct {
	listQuery
}

// WithGCLeaderboardSetLocalizationsLimit sets the max number of leaderboard set localizations to return.
func WithGCLeaderboardSetLocalizationsLimit(limit int) GCLeaderboardSetLocalizationsOption {
	return func(q *gcLeaderboardSetLocalizationsQuery) {
		if limit > 0 {
			q.limit = limit
		}
	}
}

// WithGCLeaderboardSetLocalizationsNextURL uses a next page URL directly.
func WithGCLeaderboardSetLocalizationsNextURL(next string) GCLeaderboardSetLocalizationsOption {
	return func(q *gcLeaderboardSetLocalizationsQuery) {
		if strings.TrimSpace(next) != "" {
			q.nextURL = strings.TrimSpace(next)
		}
	}
}

func buildGCLeaderboardSetLocalizationsQuery(query *gcLeaderboardSetLocalizationsQuery) string {
	values := url.Values{}
	addLimit(values, query.limit)
	return values.Encode()
}

// GameCenterLeaderboardLocalizationAttributes represents a Game Center leaderboard localization resource.
type GameCenterLeaderboardLocalizationAttributes struct {
	Locale                  string `json:"locale"`
	Name                    string `json:"name"`
	FormatterOverride       string `json:"formatterOverride,omitempty"`
	FormatterSuffix         string `json:"formatterSuffix,omitempty"`
	FormatterSuffixSingular string `json:"formatterSuffixSingular,omitempty"`
}

// GameCenterLeaderboardLocalizationCreateAttributes describes attributes for creating a localization.
type GameCenterLeaderboardLocalizationCreateAttributes struct {
	Locale                  string `json:"locale"`
	Name                    string `json:"name"`
	FormatterOverride       string `json:"formatterOverride,omitempty"`
	FormatterSuffix         string `json:"formatterSuffix,omitempty"`
	FormatterSuffixSingular string `json:"formatterSuffixSingular,omitempty"`
}

// GameCenterLeaderboardLocalizationUpdateAttributes describes attributes for updating a localization.
type GameCenterLeaderboardLocalizationUpdateAttributes struct {
	Name                    *string `json:"name,omitempty"`
	FormatterOverride       *string `json:"formatterOverride,omitempty"`
	FormatterSuffix         *string `json:"formatterSuffix,omitempty"`
	FormatterSuffixSingular *string `json:"formatterSuffixSingular,omitempty"`
}

// GameCenterLeaderboardLocalizationRelationships describes relationships for leaderboard localizations.
type GameCenterLeaderboardLocalizationRelationships struct {
	GameCenterLeaderboard *Relationship `json:"gameCenterLeaderboard"`
}

// GameCenterLeaderboardLocalizationCreateData is the data portion of a localization create request.
type GameCenterLeaderboardLocalizationCreateData struct {
	Type          ResourceType                                      `json:"type"`
	Attributes    GameCenterLeaderboardLocalizationCreateAttributes `json:"attributes"`
	Relationships *GameCenterLeaderboardLocalizationRelationships   `json:"relationships,omitempty"`
}

// GameCenterLeaderboardLocalizationCreateRequest is a request to create a localization.
type GameCenterLeaderboardLocalizationCreateRequest struct {
	Data GameCenterLeaderboardLocalizationCreateData `json:"data"`
}

// GameCenterLeaderboardLocalizationUpdateData is the data portion of a localization update request.
type GameCenterLeaderboardLocalizationUpdateData struct {
	Type       ResourceType                                       `json:"type"`
	ID         string                                             `json:"id"`
	Attributes *GameCenterLeaderboardLocalizationUpdateAttributes `json:"attributes,omitempty"`
}

// GameCenterLeaderboardLocalizationUpdateRequest is a request to update a localization.
type GameCenterLeaderboardLocalizationUpdateRequest struct {
	Data GameCenterLeaderboardLocalizationUpdateData `json:"data"`
}

// GameCenterLeaderboardLocalizationsResponse is the response from leaderboard localization list endpoints.
type GameCenterLeaderboardLocalizationsResponse = Response[GameCenterLeaderboardLocalizationAttributes]

// GameCenterLeaderboardLocalizationResponse is the response from leaderboard localization detail endpoints.
type GameCenterLeaderboardLocalizationResponse = SingleResponse[GameCenterLeaderboardLocalizationAttributes]

// GameCenterLeaderboardLocalizationDeleteResult represents CLI output for localization deletions.
type GameCenterLeaderboardLocalizationDeleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// GCLeaderboardLocalizationsOption is a functional option for GetGameCenterLeaderboardLocalizations.
type GCLeaderboardLocalizationsOption func(*gcLeaderboardLocalizationsQuery)

type gcLeaderboardLocalizationsQuery struct {
	listQuery
}

// WithGCLeaderboardLocalizationsLimit sets the max number of localizations to return.
func WithGCLeaderboardLocalizationsLimit(limit int) GCLeaderboardLocalizationsOption {
	return func(q *gcLeaderboardLocalizationsQuery) {
		if limit > 0 {
			q.limit = limit
		}
	}
}

// WithGCLeaderboardLocalizationsNextURL uses a next page URL directly.
func WithGCLeaderboardLocalizationsNextURL(next string) GCLeaderboardLocalizationsOption {
	return func(q *gcLeaderboardLocalizationsQuery) {
		if strings.TrimSpace(next) != "" {
			q.nextURL = strings.TrimSpace(next)
		}
	}
}

func buildGCLeaderboardLocalizationsQuery(query *gcLeaderboardLocalizationsQuery) string {
	values := url.Values{}
	addLimit(values, query.limit)
	return values.Encode()
}

// GameCenterAchievementLocalizationAttributes represents a Game Center achievement localization.
type GameCenterAchievementLocalizationAttributes struct {
	Locale                  string `json:"locale"`
	Name                    string `json:"name,omitempty"`
	BeforeEarnedDescription string `json:"beforeEarnedDescription,omitempty"`
	AfterEarnedDescription  string `json:"afterEarnedDescription,omitempty"`
}

// GameCenterAchievementLocalizationCreateAttributes describes attributes for creating a localization.
type GameCenterAchievementLocalizationCreateAttributes struct {
	Locale                  string `json:"locale"`
	Name                    string `json:"name,omitempty"`
	BeforeEarnedDescription string `json:"beforeEarnedDescription,omitempty"`
	AfterEarnedDescription  string `json:"afterEarnedDescription,omitempty"`
}

// GameCenterAchievementLocalizationUpdateAttributes describes attributes for updating a localization.
type GameCenterAchievementLocalizationUpdateAttributes struct {
	Name                    *string `json:"name,omitempty"`
	BeforeEarnedDescription *string `json:"beforeEarnedDescription,omitempty"`
	AfterEarnedDescription  *string `json:"afterEarnedDescription,omitempty"`
}

// GameCenterAchievementLocalizationRelationships describes relationships for achievement localizations.
type GameCenterAchievementLocalizationRelationships struct {
	GameCenterAchievement *Relationship `json:"gameCenterAchievement"`
}

// GameCenterAchievementLocalizationCreateData is the data portion of a localization create request.
type GameCenterAchievementLocalizationCreateData struct {
	Type          ResourceType                                      `json:"type"`
	Attributes    GameCenterAchievementLocalizationCreateAttributes `json:"attributes"`
	Relationships *GameCenterAchievementLocalizationRelationships   `json:"relationships,omitempty"`
}

// GameCenterAchievementLocalizationCreateRequest is a request to create a localization.
type GameCenterAchievementLocalizationCreateRequest struct {
	Data GameCenterAchievementLocalizationCreateData `json:"data"`
}

// GameCenterAchievementLocalizationUpdateData is the data portion of a localization update request.
type GameCenterAchievementLocalizationUpdateData struct {
	Type       ResourceType                                       `json:"type"`
	ID         string                                             `json:"id"`
	Attributes *GameCenterAchievementLocalizationUpdateAttributes `json:"attributes,omitempty"`
}

// GameCenterAchievementLocalizationUpdateRequest is a request to update a localization.
type GameCenterAchievementLocalizationUpdateRequest struct {
	Data GameCenterAchievementLocalizationUpdateData `json:"data"`
}

// GameCenterAchievementLocalizationsResponse is the response from achievement localization list endpoints.
type GameCenterAchievementLocalizationsResponse = Response[GameCenterAchievementLocalizationAttributes]

// GameCenterAchievementLocalizationResponse is the response from achievement localization detail endpoints.
type GameCenterAchievementLocalizationResponse = SingleResponse[GameCenterAchievementLocalizationAttributes]

// GameCenterAchievementLocalizationDeleteResult represents CLI output for localization deletions.
type GameCenterAchievementLocalizationDeleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// GCAchievementLocalizationsOption is a functional option for GetGameCenterAchievementLocalizations.
type GCAchievementLocalizationsOption func(*gcAchievementLocalizationsQuery)

type gcAchievementLocalizationsQuery struct {
	listQuery
}

// WithGCAchievementLocalizationsLimit sets the max number of localizations to return.
func WithGCAchievementLocalizationsLimit(limit int) GCAchievementLocalizationsOption {
	return func(q *gcAchievementLocalizationsQuery) {
		if limit > 0 {
			q.limit = limit
		}
	}
}

// WithGCAchievementLocalizationsNextURL uses a next page URL directly.
func WithGCAchievementLocalizationsNextURL(next string) GCAchievementLocalizationsOption {
	return func(q *gcAchievementLocalizationsQuery) {
		if strings.TrimSpace(next) != "" {
			q.nextURL = strings.TrimSpace(next)
		}
	}
}

func buildGCAchievementLocalizationsQuery(query *gcAchievementLocalizationsQuery) string {
	values := url.Values{}
	addLimit(values, query.limit)
	return values.Encode()
}

// GameCenterLeaderboardReleaseAttributes represents a Game Center leaderboard release resource.
type GameCenterLeaderboardReleaseAttributes struct {
	Live bool `json:"live"`
}

// GameCenterLeaderboardReleaseRelationships describes relationships for leaderboard releases.
type GameCenterLeaderboardReleaseRelationships struct {
	GameCenterDetail      *Relationship `json:"gameCenterDetail"`
	GameCenterLeaderboard *Relationship `json:"gameCenterLeaderboard"`
}

// GameCenterLeaderboardReleaseCreateData is the data portion of a leaderboard release create request.
type GameCenterLeaderboardReleaseCreateData struct {
	Type          ResourceType                               `json:"type"`
	Relationships *GameCenterLeaderboardReleaseRelationships `json:"relationships"`
}

// GameCenterLeaderboardReleaseCreateRequest is a request to create a leaderboard release.
type GameCenterLeaderboardReleaseCreateRequest struct {
	Data GameCenterLeaderboardReleaseCreateData `json:"data"`
}

// GameCenterLeaderboardReleasesResponse is the response from leaderboard release list endpoints.
type GameCenterLeaderboardReleasesResponse = Response[GameCenterLeaderboardReleaseAttributes]

// GameCenterLeaderboardReleaseResponse is the response from leaderboard release detail endpoints.
type GameCenterLeaderboardReleaseResponse = SingleResponse[GameCenterLeaderboardReleaseAttributes]

// GameCenterLeaderboardReleaseDeleteResult represents CLI output for leaderboard release deletions.
type GameCenterLeaderboardReleaseDeleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// GCLeaderboardReleasesOption is a functional option for GetGameCenterLeaderboardReleases.
type GCLeaderboardReleasesOption func(*gcLeaderboardReleasesQuery)

type gcLeaderboardReleasesQuery struct {
	listQuery
}

// WithGCLeaderboardReleasesLimit sets the max number of leaderboard releases to return.
func WithGCLeaderboardReleasesLimit(limit int) GCLeaderboardReleasesOption {
	return func(q *gcLeaderboardReleasesQuery) {
		if limit > 0 {
			q.limit = limit
		}
	}
}

// WithGCLeaderboardReleasesNextURL uses a next page URL directly.
func WithGCLeaderboardReleasesNextURL(next string) GCLeaderboardReleasesOption {
	return func(q *gcLeaderboardReleasesQuery) {
		if strings.TrimSpace(next) != "" {
			q.nextURL = strings.TrimSpace(next)
		}
	}
}

func buildGCLeaderboardReleasesQuery(query *gcLeaderboardReleasesQuery) string {
	values := url.Values{}
	addLimit(values, query.limit)
	return values.Encode()
}

// GameCenterAchievementReleaseAttributes represents a Game Center achievement release resource.
type GameCenterAchievementReleaseAttributes struct {
	Live bool `json:"live"`
}

// GameCenterAchievementReleaseRelationships describes relationships for achievement releases.
type GameCenterAchievementReleaseRelationships struct {
	GameCenterDetail      *Relationship `json:"gameCenterDetail"`
	GameCenterAchievement *Relationship `json:"gameCenterAchievement"`
}

// GameCenterAchievementReleaseCreateData is the data portion of an achievement release create request.
type GameCenterAchievementReleaseCreateData struct {
	Type          ResourceType                               `json:"type"`
	Relationships *GameCenterAchievementReleaseRelationships `json:"relationships"`
}

// GameCenterAchievementReleaseCreateRequest is a request to create an achievement release.
type GameCenterAchievementReleaseCreateRequest struct {
	Data GameCenterAchievementReleaseCreateData `json:"data"`
}

// GameCenterAchievementReleasesResponse is the response from achievement release list endpoints.
type GameCenterAchievementReleasesResponse = Response[GameCenterAchievementReleaseAttributes]

// GameCenterAchievementReleaseResponse is the response from achievement release detail endpoints.
type GameCenterAchievementReleaseResponse = SingleResponse[GameCenterAchievementReleaseAttributes]

// GameCenterAchievementReleaseDeleteResult represents CLI output for achievement release deletions.
type GameCenterAchievementReleaseDeleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// GCAchievementReleasesOption is a functional option for GetGameCenterAchievementReleases.
type GCAchievementReleasesOption func(*gcAchievementReleasesQuery)

type gcAchievementReleasesQuery struct {
	listQuery
}

// WithGCAchievementReleasesLimit sets the max number of achievement releases to return.
func WithGCAchievementReleasesLimit(limit int) GCAchievementReleasesOption {
	return func(q *gcAchievementReleasesQuery) {
		if limit > 0 {
			q.limit = limit
		}
	}
}

// WithGCAchievementReleasesNextURL uses a next page URL directly.
func WithGCAchievementReleasesNextURL(next string) GCAchievementReleasesOption {
	return func(q *gcAchievementReleasesQuery) {
		if strings.TrimSpace(next) != "" {
			q.nextURL = strings.TrimSpace(next)
		}
	}
}

func buildGCAchievementReleasesQuery(query *gcAchievementReleasesQuery) string {
	values := url.Values{}
	addLimit(values, query.limit)
	return values.Encode()
}

// GameCenterLeaderboardSetReleaseAttributes represents a Game Center leaderboard set release resource.
type GameCenterLeaderboardSetReleaseAttributes struct {
	Live bool `json:"live"`
}

// GameCenterLeaderboardSetReleaseRelationships describes relationships for leaderboard set releases.
type GameCenterLeaderboardSetReleaseRelationships struct {
	GameCenterDetail         *Relationship `json:"gameCenterDetail"`
	GameCenterLeaderboardSet *Relationship `json:"gameCenterLeaderboardSet"`
}

// GameCenterLeaderboardSetReleaseCreateData is the data portion of a leaderboard set release create request.
type GameCenterLeaderboardSetReleaseCreateData struct {
	Type          ResourceType                                  `json:"type"`
	Relationships *GameCenterLeaderboardSetReleaseRelationships `json:"relationships"`
}

// GameCenterLeaderboardSetReleaseCreateRequest is a request to create a leaderboard set release.
type GameCenterLeaderboardSetReleaseCreateRequest struct {
	Data GameCenterLeaderboardSetReleaseCreateData `json:"data"`
}

// GameCenterLeaderboardSetReleasesResponse is the response from leaderboard set release list endpoints.
type GameCenterLeaderboardSetReleasesResponse = Response[GameCenterLeaderboardSetReleaseAttributes]

// GameCenterLeaderboardSetReleaseResponse is the response from leaderboard set release detail endpoints.
type GameCenterLeaderboardSetReleaseResponse = SingleResponse[GameCenterLeaderboardSetReleaseAttributes]

// GameCenterLeaderboardSetReleaseDeleteResult represents CLI output for leaderboard set release deletions.
type GameCenterLeaderboardSetReleaseDeleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// GCLeaderboardSetReleasesOption is a functional option for GetGameCenterLeaderboardSetReleases.
type GCLeaderboardSetReleasesOption func(*gcLeaderboardSetReleasesQuery)

type gcLeaderboardSetReleasesQuery struct {
	listQuery
}

// WithGCLeaderboardSetReleasesLimit sets the max number of leaderboard set releases to return.
func WithGCLeaderboardSetReleasesLimit(limit int) GCLeaderboardSetReleasesOption {
	return func(q *gcLeaderboardSetReleasesQuery) {
		if limit > 0 {
			q.limit = limit
		}
	}
}

// WithGCLeaderboardSetReleasesNextURL uses a next page URL directly.
func WithGCLeaderboardSetReleasesNextURL(next string) GCLeaderboardSetReleasesOption {
	return func(q *gcLeaderboardSetReleasesQuery) {
		if strings.TrimSpace(next) != "" {
			q.nextURL = strings.TrimSpace(next)
		}
	}
}

func buildGCLeaderboardSetReleasesQuery(query *gcLeaderboardSetReleasesQuery) string {
	values := url.Values{}
	addLimit(values, query.limit)
	return values.Encode()
}

// GameCenterLeaderboardSetMembersUpdateRequest is a request to replace leaderboard set members.
type GameCenterLeaderboardSetMembersUpdateRequest struct {
	Data []RelationshipData `json:"data"`
}

// GameCenterLeaderboardSetMembersUpdateResult represents CLI output for member updates.
type GameCenterLeaderboardSetMembersUpdateResult struct {
	SetID       string   `json:"setId"`
	MemberCount int      `json:"memberCount"`
	MemberIDs   []string `json:"memberIds"`
	Updated     bool     `json:"updated"`
}

// GCLeaderboardSetMembersOption is a functional option for GetGameCenterLeaderboardSetMembers.
type GCLeaderboardSetMembersOption func(*gcLeaderboardSetMembersQuery)

type gcLeaderboardSetMembersQuery struct {
	listQuery
}

// WithGCLeaderboardSetMembersLimit sets the max number of members to return.
func WithGCLeaderboardSetMembersLimit(limit int) GCLeaderboardSetMembersOption {
	return func(q *gcLeaderboardSetMembersQuery) {
		if limit > 0 {
			q.limit = limit
		}
	}
}

// WithGCLeaderboardSetMembersNextURL uses a next page URL directly.
func WithGCLeaderboardSetMembersNextURL(next string) GCLeaderboardSetMembersOption {
	return func(q *gcLeaderboardSetMembersQuery) {
		if strings.TrimSpace(next) != "" {
			q.nextURL = strings.TrimSpace(next)
		}
	}
}

func buildGCLeaderboardSetMembersQuery(query *gcLeaderboardSetMembersQuery) string {
	values := url.Values{}
	addLimit(values, query.limit)
	return values.Encode()
}

// GameCenterLeaderboardImageAttributes represents a Game Center leaderboard image resource.
type GameCenterLeaderboardImageAttributes struct {
	FileSize           int64               `json:"fileSize"`
	FileName           string              `json:"fileName"`
	ImageAsset         *ImageAsset         `json:"imageAsset,omitempty"`
	UploadOperations   []UploadOperation   `json:"uploadOperations,omitempty"`
	AssetDeliveryState *AssetDeliveryState `json:"assetDeliveryState,omitempty"`
}

// GameCenterLeaderboardImageCreateAttributes describes attributes for reserving an image upload.
type GameCenterLeaderboardImageCreateAttributes struct {
	FileSize int64  `json:"fileSize"`
	FileName string `json:"fileName"`
}

// GameCenterLeaderboardImageUpdateAttributes describes attributes for committing an image upload.
type GameCenterLeaderboardImageUpdateAttributes struct {
	Uploaded *bool `json:"uploaded,omitempty"`
}

// GameCenterLeaderboardImageRelationships describes relationships for leaderboard images.
type GameCenterLeaderboardImageRelationships struct {
	GameCenterLeaderboardLocalization *Relationship `json:"gameCenterLeaderboardLocalization"`
}

// GameCenterLeaderboardImageCreateData is the data portion of an image create (reserve) request.
type GameCenterLeaderboardImageCreateData struct {
	Type          ResourceType                               `json:"type"`
	Attributes    GameCenterLeaderboardImageCreateAttributes `json:"attributes"`
	Relationships *GameCenterLeaderboardImageRelationships   `json:"relationships"`
}

// GameCenterLeaderboardImageCreateRequest is a request to reserve an image upload.
type GameCenterLeaderboardImageCreateRequest struct {
	Data GameCenterLeaderboardImageCreateData `json:"data"`
}

// GameCenterLeaderboardImageUpdateData is the data portion of an image update (commit) request.
type GameCenterLeaderboardImageUpdateData struct {
	Type       ResourceType                                `json:"type"`
	ID         string                                      `json:"id"`
	Attributes *GameCenterLeaderboardImageUpdateAttributes `json:"attributes,omitempty"`
}

// GameCenterLeaderboardImageUpdateRequest is a request to commit an image upload.
type GameCenterLeaderboardImageUpdateRequest struct {
	Data GameCenterLeaderboardImageUpdateData `json:"data"`
}

// GameCenterLeaderboardImagesResponse is the response from leaderboard image list endpoints.
type GameCenterLeaderboardImagesResponse = Response[GameCenterLeaderboardImageAttributes]

// GameCenterLeaderboardImageResponse is the response from leaderboard image detail endpoints.
type GameCenterLeaderboardImageResponse = SingleResponse[GameCenterLeaderboardImageAttributes]

// GameCenterLeaderboardImageDeleteResult represents CLI output for image deletions.
type GameCenterLeaderboardImageDeleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// GameCenterLeaderboardImageUploadResult represents CLI output for image uploads.
type GameCenterLeaderboardImageUploadResult struct {
	ID                 string `json:"id"`
	LocalizationID     string `json:"localizationId"`
	FileName           string `json:"fileName"`
	FileSize           int64  `json:"fileSize"`
	AssetDeliveryState string `json:"assetDeliveryState,omitempty"`
	Uploaded           bool   `json:"uploaded"`
}

// GameCenterAchievementImageAttributes represents a Game Center achievement image resource.
type GameCenterAchievementImageAttributes struct {
	FileSize           int64               `json:"fileSize"`
	FileName           string              `json:"fileName"`
	ImageAsset         *ImageAsset         `json:"imageAsset,omitempty"`
	UploadOperations   []UploadOperation   `json:"uploadOperations,omitempty"`
	AssetDeliveryState *AssetDeliveryState `json:"assetDeliveryState,omitempty"`
}

// GameCenterAchievementImageCreateAttributes describes attributes for reserving an image upload.
type GameCenterAchievementImageCreateAttributes struct {
	FileSize int64  `json:"fileSize"`
	FileName string `json:"fileName"`
}

// GameCenterAchievementImageUpdateAttributes describes attributes for committing an image upload.
type GameCenterAchievementImageUpdateAttributes struct {
	Uploaded *bool `json:"uploaded,omitempty"`
}

// GameCenterAchievementImageRelationships describes relationships for achievement images.
type GameCenterAchievementImageRelationships struct {
	GameCenterAchievementLocalization *Relationship `json:"gameCenterAchievementLocalization"`
}

// GameCenterAchievementImageCreateData is the data portion of an image create (reserve) request.
type GameCenterAchievementImageCreateData struct {
	Type          ResourceType                               `json:"type"`
	Attributes    GameCenterAchievementImageCreateAttributes `json:"attributes"`
	Relationships *GameCenterAchievementImageRelationships   `json:"relationships"`
}

// GameCenterAchievementImageCreateRequest is a request to reserve an image upload.
type GameCenterAchievementImageCreateRequest struct {
	Data GameCenterAchievementImageCreateData `json:"data"`
}

// GameCenterAchievementImageUpdateData is the data portion of an image update (commit) request.
type GameCenterAchievementImageUpdateData struct {
	Type       ResourceType                                `json:"type"`
	ID         string                                      `json:"id"`
	Attributes *GameCenterAchievementImageUpdateAttributes `json:"attributes,omitempty"`
}

// GameCenterAchievementImageUpdateRequest is a request to commit an image upload.
type GameCenterAchievementImageUpdateRequest struct {
	Data GameCenterAchievementImageUpdateData `json:"data"`
}

// GameCenterAchievementImagesResponse is the response from achievement image list endpoints.
type GameCenterAchievementImagesResponse = Response[GameCenterAchievementImageAttributes]

// GameCenterAchievementImageResponse is the response from achievement image detail endpoints.
type GameCenterAchievementImageResponse = SingleResponse[GameCenterAchievementImageAttributes]

// GameCenterAchievementImageDeleteResult represents CLI output for image deletions.
type GameCenterAchievementImageDeleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// GameCenterAchievementImageUploadResult represents CLI output for image uploads.
type GameCenterAchievementImageUploadResult struct {
	ID                 string `json:"id"`
	LocalizationID     string `json:"localizationId"`
	FileName           string `json:"fileName"`
	FileSize           int64  `json:"fileSize"`
	AssetDeliveryState string `json:"assetDeliveryState,omitempty"`
	Uploaded           bool   `json:"uploaded"`
}

// GameCenterLeaderboardSetImageAttributes represents a Game Center leaderboard set image resource.
type GameCenterLeaderboardSetImageAttributes struct {
	FileSize           int64               `json:"fileSize"`
	FileName           string              `json:"fileName"`
	ImageAsset         *ImageAsset         `json:"imageAsset,omitempty"`
	UploadOperations   []UploadOperation   `json:"uploadOperations,omitempty"`
	AssetDeliveryState *AssetDeliveryState `json:"assetDeliveryState,omitempty"`
}

// GameCenterLeaderboardSetImageCreateAttributes describes attributes for reserving an image upload.
type GameCenterLeaderboardSetImageCreateAttributes struct {
	FileSize int64  `json:"fileSize"`
	FileName string `json:"fileName"`
}

// GameCenterLeaderboardSetImageUpdateAttributes describes attributes for committing an image upload.
type GameCenterLeaderboardSetImageUpdateAttributes struct {
	Uploaded *bool `json:"uploaded,omitempty"`
}

// GameCenterLeaderboardSetImageRelationships describes relationships for leaderboard set images.
type GameCenterLeaderboardSetImageRelationships struct {
	GameCenterLeaderboardSetLocalization *Relationship `json:"gameCenterLeaderboardSetLocalization"`
}

// GameCenterLeaderboardSetImageCreateData is the data portion of an image create (reserve) request.
type GameCenterLeaderboardSetImageCreateData struct {
	Type          ResourceType                                  `json:"type"`
	Attributes    GameCenterLeaderboardSetImageCreateAttributes `json:"attributes"`
	Relationships *GameCenterLeaderboardSetImageRelationships   `json:"relationships"`
}

// GameCenterLeaderboardSetImageCreateRequest is a request to reserve an image upload.
type GameCenterLeaderboardSetImageCreateRequest struct {
	Data GameCenterLeaderboardSetImageCreateData `json:"data"`
}

// GameCenterLeaderboardSetImageUpdateData is the data portion of an image update (commit) request.
type GameCenterLeaderboardSetImageUpdateData struct {
	Type       ResourceType                                   `json:"type"`
	ID         string                                         `json:"id"`
	Attributes *GameCenterLeaderboardSetImageUpdateAttributes `json:"attributes,omitempty"`
}

// GameCenterLeaderboardSetImageUpdateRequest is a request to update a leaderboard set image.
type GameCenterLeaderboardSetImageUpdateRequest struct {
	Data GameCenterLeaderboardSetImageUpdateData `json:"data"`
}

// GameCenterLeaderboardSetImageResponse is the response from leaderboard set image detail endpoints.
type GameCenterLeaderboardSetImageResponse = SingleResponse[GameCenterLeaderboardSetImageAttributes]

// GameCenterLeaderboardSetImageDeleteResult represents CLI output for leaderboard set image deletions.
type GameCenterLeaderboardSetImageDeleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// GameCenterLeaderboardSetImageUploadResult represents CLI output for leaderboard set image uploads.
type GameCenterLeaderboardSetImageUploadResult struct {
	ID                 string `json:"id"`
	LocalizationID     string `json:"localizationId"`
	FileName           string `json:"fileName"`
	FileSize           int64  `json:"fileSize"`
	AssetDeliveryState string `json:"assetDeliveryState,omitempty"`
	Uploaded           bool   `json:"uploaded"`
}
