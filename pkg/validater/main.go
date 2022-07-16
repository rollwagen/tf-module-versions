package validater

import (
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/xanzy/go-gitlab"
)

type Module struct {
	Name     string
	Location struct {
		FileName string
		Line     int
	}
	UsedVersion      string
	AvailableVersion string
}

func NewModule(name string, usedVersion, availableVersion, fileName string, line int) (*Module, error) {
	for _, v := range []string{usedVersion, availableVersion} {
		_, err := version.NewVersion(v)
		if err != nil {
			return nil, fmt.Errorf("'%s' is not a valid version string: %w", v, err)
		}
	}

	m := Module{
		Name:             name,
		UsedVersion:      usedVersion,
		AvailableVersion: availableVersion,
	}
	m.Location.FileName = fileName
	m.Location.Line = line

	return &m, nil
}

func (m Module) HasNewerVersion() bool {
	ver := func(s string) *version.Version {
		v, _ := version.NewVersion(s)
		return v
	}
	return ver(m.UsedVersion).LessThan(ver(m.AvailableVersion))
}

func (m Module) HasSameVersion() bool {
	return m.AvailableVersion == m.UsedVersion
}

// Validate inspect and check used terraform gitlab reference versions
func Validate(dir string, quiet bool) []Module {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var validatedModules []Module

	gitlabToken := os.Getenv("GITLAB_TOKEN")
	if gitlabToken == "" {
		_, _ = fmt.Fprintln(os.Stderr, "Please define gitlab auth token as environment variable GITLAB_TOKEN")
		os.Exit(1)
	}

	tfModule, diag := tfconfig.LoadModule(dir)
	if diag.HasErrors() {
		_, _ = fmt.Fprintf(os.Stderr, "The terraform tfModule contains errors: %v\n", diag)
		os.Exit(1)
	}

	for _, moduleCall := range tfModule.ModuleCalls {

		sourceVersion := moduleCall.Version

		// for git references modules will be empty (vs registry referenced modules)
		if sourceVersion != "" {
			continue
		}

		// get version ref from git url "...ref=1.1.1"
		splitSource := strings.Split(moduleCall.Source, "=")
		sourceVersion = splitSource[len(splitSource)-1]

		// get rid of tf specific generic git:: prefix
		u, _ := url.Parse(strings.Replace(moduleCall.Source, "git::", "", 1))
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

		latestVersion, err := retrieveLatestVersion(gitlabClient, gitlabProjectNamespaceName)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error getting latest version from gitlab: %v\n", err)
			os.Exit(1)
		}

		_, err = version.NewVersion(sourceVersion)
		if err != nil {
			log.Warn().
				Str("tfModule", moduleCall.Name).
				Str("file", fmt.Sprintf("%s:%d", moduleCall.Pos.Filename, moduleCall.Pos.Line)).
				Msg(fmt.Sprintf("ref '%s' is not a valid version string", sourceVersion))

			sourceVersion = "0" // set to default '0' as no valid version number referenced
		}

		module, _ := NewModule(moduleCall.Name, sourceVersion, latestVersion, moduleCall.Pos.Filename, moduleCall.Pos.Line)
		validatedModules = append(validatedModules, *module)

		if module.HasNewerVersion() {
			log.Info().
				Str("tfModule", module.Name).
				Str("version_used", module.UsedVersion).
				Str("version_latest", module.AvailableVersion).
				Str("file", fmt.Sprintf("%s:%d", module.Location.FileName, module.Location.Line)).
				Msg(color.New(color.FgRed).Add(color.Bold).Sprint("✖ ···>"))
		} else if module.HasSameVersion() && module.UsedVersion != "0" {
			if !quiet {
				log.Debug().
					Str("tfModule", module.Name).
					Str("file", fmt.Sprintf("%s:%d", module.Location.FileName, module.Location.Line)).
					Msg(color.New(color.FgGreen).Add(color.Bold).Sprint("✔ latest version used"))
			}
		}
	}
	log.Debug().Msg("validation completed")

	return validatedModules
}

func retrieveLatestVersion(gitlabClient *gitlab.Client, moduleNamespaceName string) (string, error) {
	tags, _, err := gitlabClient.Tags.ListTags(moduleNamespaceName, &gitlab.ListTagsOptions{})
	if err != nil {
		return "", fmt.Errorf("error querying gitlab; potential auth issue; check GITLAB_TOKEN: %w", err)
	}

	var availableVersions []string
	for _, t := range tags {
		_, err := version.NewVersion(t.Name)
		if err == nil {
			availableVersions = append(availableVersions, t.Name)
		}
	}
	sort.Strings(availableVersions)
	latestVersion := "0"
	if len(availableVersions) > 0 {
		latestVersion = availableVersions[len(availableVersions)-1]
	}
	return latestVersion, nil
}
