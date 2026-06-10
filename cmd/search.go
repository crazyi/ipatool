package cmd

import (
	"github.com/majd/ipatool/v2/pkg/appstore"
	"github.com/spf13/cobra"
)

// nolint:wrapcheck
func searchCmd() *cobra.Command {
	var (
		limit         int64
		countryCode   string
		platformValue string
	)

	cmd := &cobra.Command{
		Use:   "search <term>",
		Short: "Search for iOS and tvOS apps available on the App Store",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			infoResult, err := dependencies.AppStore.AccountInfo()
			if err != nil {
				return err
			}

			platform, err := appstore.ParsePlatform(platformValue)
			if err != nil {
				return err
			}

			output, err := dependencies.AppStore.Search(appstore.SearchInput{
				Account:     infoResult.Account,
				Term:        args[0],
				Limit:       limit,
				CountryCode: countryCode,
				Platform:    platform,
			})
			if err != nil {
				return err
			}

			dependencies.Logger.Log().
				Int("count", output.Count).
				Array("apps", appstore.Apps(output.Results)).
				Send()

			return nil
		},
	}

	cmd.Flags().Int64VarP(&limit, "limit", "l", 5, "maximum amount of search results to retrieve")
	cmd.Flags().StringVarP(&countryCode, "country", "c", "", "The two-letter (ISO 3166-1) country code for the iTunes Store")
	cmd.Flags().StringVar(&platformValue, "platform", "", "Platform to search: iphone, ipad, or appletv")

	return cmd
}
