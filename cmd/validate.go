package cmd

import (
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/xanzy/go-gitlab"

	"github.com/fatih/color"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "validate",
	Short: "Print module version validation on stdout as logs",
	Run: func(cmd *cobra.Command, args []string) {
		validate()
	},
}

func validate() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	gitlabToken := os.Getenv("GITLAB_TOKEN")
	if gitlabToken == "" {
		_, _ = fmt.Fprintln(os.Stderr, "Please define gitlab auth token as environment variable GITLAB_TOKEN")
		os.Exit(1)
	}

	module, _ := tfconfig.LoadModule(".")

	for _, m := range module.ModuleCalls {

		latestVersion := "0"
		sourceVersion := m.Version
		zeroVersion, _ := version.NewVersion("0")

		// for git references modules will be empty (vs registry referenced modules)
		if sourceVersion != "" {
			continue
		}

		// get version ref from git url "...ref=1.1.1"
		splitSource := strings.Split(m.Source, "=")
		sourceVersion = splitSource[len(splitSource)-1]

		// get rid of tf specific generic git:: prefix, for details see
		// https://www.terraform.io/language/modules/sources#generic-git-repository
		u, _ := url.Parse(strings.Replace(m.Source, "git::", "", 1))
		if u.Host == "" { // path source have no version and also no host e.g. source = "./submodule"
			continue
		}

		pathNoSuffix := strings.Replace(u.Path, ".git", "", 1)
		gitlabProjectNamespaceName := strings.Replace(pathNoSuffix, "/", "", 1)
		baseURL := fmt.Sprintf("https://%s/api/v4", u.Host)

		gitlabClient, err := gitlab.NewClient(gitlabToken, gitlab.WithBaseURL(baseURL))
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error creating gitlab client: %s\n", err)
			os.Exit(1)
		}

		tags, _, err := gitlabClient.Tags.ListTags(gitlabProjectNamespaceName, &gitlab.ListTagsOptions{})
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error querying host '%s'; potential auth issue; check GITLAB_TOKEN.\nDetails: %v\n", u.Host, err)
			os.Exit(1)
		}

		var availableVersions []string
		for _, t := range tags {
			_, err := version.NewVersion(t.Name)
			if err == nil {
				availableVersions = append(availableVersions, t.Name)
			}
		}

		sort.Strings(availableVersions)
		if len(availableVersions) > 0 {
			latestVersion = availableVersions[len(availableVersions)-1]
		}

		vLatest, _ := version.NewVersion(latestVersion)
		vUsed, err := version.NewVersion(sourceVersion)
		if err != nil {
			log.Warn().
				Str("module", m.Name).
				Str("file", fmt.Sprintf("%s:%d", m.Pos.Filename, m.Pos.Line)).
				Msg(fmt.Sprintf("ref '%s' is not a valid version string", sourceVersion))
			vUsed, _ = version.NewVersion("0")
		}

		if vUsed.LessThan(vLatest) {
			log.Info().
				Str("module", m.Name).
				Str("version_used", sourceVersion).
				Str("version_latest", latestVersion).
				Str("file", fmt.Sprintf("%s:%d", m.Pos.Filename, m.Pos.Line)).
				Msg(color.New(color.FgRed).Add(color.Bold).Sprint("✖ ···>"))
		} else if vUsed.Equal(vLatest) && !vUsed.Equal(zeroVersion) {
			if !Quiet {
				log.Debug().
					Str("module", m.Name).
					Str("file", fmt.Sprintf("%s:%d", m.Pos.Filename, m.Pos.Line)).
					Msg(color.New(color.FgGreen).Add(color.Bold).Sprint("✔ latest version used"))
			}
		}
	}
	log.Debug().Msg("validation completed")
}
