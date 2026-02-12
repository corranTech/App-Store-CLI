package registry

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/accessibility"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/actors"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/agerating"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/agreements"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/alternativedistribution"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/analytics"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/androidiosmapping"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/app_events"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/appclips"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/apps"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/assets"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/auth"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/backgroundassets"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/betaapplocalizations"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/betabuildlocalizations"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/buildbundles"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/buildlocalizations"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/builds"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/bundleids"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/categories"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/certificates"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/crashes"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/devices"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/docs"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/encryption"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/eula"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/feedback"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/finance"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/gamecenter"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/iap"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/initcmd"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/install"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/localizations"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/marketplace"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/merchantids"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/migrate"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/nominations"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/notarization"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/notify"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/offercodes"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/passtypeids"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/performance"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/preorders"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/prerelease"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/pricing"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/productpages"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/profiles"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/promotedpurchases"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/publish"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/reviews"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/routingcoverage"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/sandbox"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/signing"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/submit"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/subscriptions"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/testflight"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/users"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/validate"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/versions"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/webhooks"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/winbackoffers"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/xcodecloud"
)

// CoreCommands returns authentication, initialization, and utility commands.
func CoreCommands() []*ffcli.Command {
	return []*ffcli.Command{
		auth.AuthCommand(),
		install.InstallCommand(),
		initcmd.InitCommand(),
		docs.DocsCommand(),
	}
}

// FeedbackCommands returns commands related to beta tester feedback and crashes.
func FeedbackCommands() []*ffcli.Command {
	return []*ffcli.Command{
		feedback.FeedbackCommand(),
		crashes.CrashesCommand(),
	}
}

// ReviewCommands returns commands for managing App Store reviews.
func ReviewCommands() []*ffcli.Command {
	return []*ffcli.Command{
		reviews.ReviewsCommand(),
		reviews.ReviewCommand(),
	}
}

// AnalyticsCommands returns analytics, performance, and finance reporting commands.
func AnalyticsCommands() []*ffcli.Command {
	return []*ffcli.Command{
		analytics.AnalyticsCommand(),
		performance.PerformanceCommand(),
		finance.FinanceCommand(),
	}
}

// AppCommands returns commands for managing apps and app metadata.
func AppCommands() []*ffcli.Command {
	return []*ffcli.Command{
		apps.AppsCommand(),
		appclips.AppClipsCommand(),
		androidiosmapping.AndroidIosMappingCommand(),
		apps.AppSetupCommand(),
		apps.AppTagsCommand(),
		marketplace.MarketplaceCommand(),
		alternativedistribution.Command(),
	}
}

// VersionCommands returns commands for managing app versions and metadata.
func VersionCommands() []*ffcli.Command {
	return []*ffcli.Command{
		versions.VersionsCommand(),
		apps.AppInfoCommand(),
		apps.AppInfosCommand(),
		eula.EULACommand(),
		agreements.AgreementsCommand(),
		pricing.PricingCommand(),
		preorders.PreOrdersCommand(),
		routingcoverage.RoutingCoverageCommand(),
		productpages.ProductPagesCommand(),
	}
}

// BuildCommands returns commands for managing builds and build bundles.
func BuildCommands() []*ffcli.Command {
	return []*ffcli.Command{
		builds.BuildsCommand(),
		buildbundles.BuildBundlesCommand(),
		publish.PublishCommand(),
		prerelease.PreReleaseVersionsCommand(),
	}
}

// TestFlightCommands returns TestFlight-related commands.
func TestFlightCommands() []*ffcli.Command {
	return []*ffcli.Command{
		testflight.TestFlightCommand(),
		betaapplocalizations.BetaAppLocalizationsCommand(),
		betabuildlocalizations.BetaBuildLocalizationsCommand(),
	}
}

// AssetCommands returns commands for managing assets and localizations.
func AssetCommands() []*ffcli.Command {
	return []*ffcli.Command{
		localizations.LocalizationsCommand(),
		assets.AssetsCommand(),
		backgroundassets.BackgroundAssetsCommand(),
		buildlocalizations.BuildLocalizationsCommand(),
	}
}

// SandboxCommands returns sandbox tester management commands.
func SandboxCommands() []*ffcli.Command {
	return []*ffcli.Command{
		sandbox.SandboxCommand(),
	}
}

// SigningCommands returns code signing and notarization commands.
func SigningCommands() []*ffcli.Command {
	return []*ffcli.Command{
		signing.SigningCommand(),
		notarization.NotarizationCommand(),
		certificates.CertificatesCommand(),
		profiles.ProfilesCommand(),
	}
}

// IAPCommands returns in-app purchase and subscription commands.
func IAPCommands() []*ffcli.Command {
	return []*ffcli.Command{
		iap.IAPCommand(),
		app_events.Command(),
		subscriptions.SubscriptionsCommand(),
		offercodes.OfferCodesCommand(),
		winbackoffers.WinBackOffersCommand(),
	}
}

// UserCommands returns user and device management commands.
func UserCommands() []*ffcli.Command {
	return []*ffcli.Command{
		users.UsersCommand(),
		actors.ActorsCommand(),
		devices.DevicesCommand(),
	}
}

// BundleIDCommands returns bundle ID and capability management commands.
func BundleIDCommands() []*ffcli.Command {
	return []*ffcli.Command{
		bundleids.BundleIDsCommand(),
		merchantids.MerchantIDsCommand(),
		passtypeids.PassTypeIDsCommand(),
	}
}

// WebhookCommands returns webhook and notification commands.
func WebhookCommands() []*ffcli.Command {
	return []*ffcli.Command{
		webhooks.WebhooksCommand(),
		nominations.NominationsCommand(),
		notify.NotifyCommand(),
	}
}

// SubmissionCommands returns app submission and validation commands.
func SubmissionCommands() []*ffcli.Command {
	return []*ffcli.Command{
		submit.SubmitCommand(),
		validate.ValidateCommand(),
		xcodecloud.XcodeCloudCommand(),
	}
}

// MetadataCommands returns metadata management commands.
func MetadataCommands() []*ffcli.Command {
	return []*ffcli.Command{
		migrate.MigrateCommand(),
	}
}

// AppStoreCommands returns App Store-specific metadata commands.
func AppStoreCommands() []*ffcli.Command {
	return []*ffcli.Command{
		categories.CategoriesCommand(),
		agerating.AgeRatingCommand(),
		accessibility.AccessibilityCommand(),
		encryption.EncryptionCommand(),
		promotedpurchases.PromotedPurchasesCommand(),
	}
}

// GameCenterCommands returns Game Center commands.
func GameCenterCommands() []*ffcli.Command {
	return []*ffcli.Command{
		gamecenter.GameCenterCommand(),
	}
}
