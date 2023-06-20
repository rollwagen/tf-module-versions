package validater

import (
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/rollwagen/tf-module-versions/internal/tf"
	"github.com/rollwagen/tf-module-versions/pkg/printer"

	"github.com/fatih/color"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/xanzy/go-gitlab"
)

// Validate inspect and check used terraform gitlab reference versions
func Validate(dir string, outputFormat string, verbose bool) []tf.Module {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	if verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	var validatedModules []tf.Module

	gitlabToken := os.Getenv("GITLAB_TOKEN")
	if gitlabToken == "" {
		_, _ = fmt.Fprintln(os.Stderr, "Please define gitlab auth token with environment variable GITLAB_TOKEN")
		os.Exit(1)
	}

	log.Debug().Msg(fmt.Sprintf("loading module from dir '%s'", dir))
	tfModule, diag := tfconfig.LoadModule(dir)
	if diag.HasErrors() {
		_, _ = fmt.Fprintf(os.Stderr, "The terraform tfModule contains errors: %v\n", diag)
		os.Exit(1)
	}

	for _, moduleCall := range tfModule.ModuleCalls {

		sourceVersion := moduleCall.Version

		// for git references modules will be empty (vs registry referenced modules)
		if sourceVersion != "" {
			log.Debug().Msg("terraform 'version=' used. skipping")
			continue
		}

		// get version ref from git url "...ref=1.1.1"
		splitSource := strings.Split(moduleCall.Source, "=")
		sourceVersion = splitSource[len(splitSource)-1]
		log.Debug().Str("tfModule", moduleCall.Name).Msg(fmt.Sprintf("sourceVersion='%s'", sourceVersion))

		// get rid of tf specific generic git:: prefix
		sourceWithoutPrefix := strings.Replace(moduleCall.Source, "git::", "", 1)
		// get rid of tf specific generic git@ prefix
		sourceWithoutUser := strings.Replace(sourceWithoutPrefix, "git@", "", 1)
		log.Debug().Str("tfModule", moduleCall.Name).Msg(fmt.Sprintf("parsing module source url '%s'", sourceWithoutUser))
		u, err := url.Parse(sourceWithoutUser)
		if err != nil {
			panic(err)
		}
		if u.Host == "" { // path source have no version and also no host e.g. source = "./submodule"
			log.Debug().Str("tfModule", moduleCall.Name).Msg("path source has no host")
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

		gitRef := sourceVersion

		_, err = version.NewVersion(sourceVersion)
		const versionNil = "nil"
		if err != nil {
			log.Warn().
				Str("tfModule", moduleCall.Name).
				Str("file", fmt.Sprintf("%s:%d", moduleCall.Pos.Filename, moduleCall.Pos.Line)).
				Msg(fmt.Sprintf("ref '%s' is not a valid version string", sourceVersion))

			sourceVersion = versionNil // set to default 'nil' as no valid version number referenced
		}

		module, err := tf.NewModule(moduleCall.Name, sourceVersion, latestVersion, gitRef, moduleCall.Pos.Filename, moduleCall.Pos.Line)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("could not create new module %s: %v", moduleCall.Name, err))
			os.Exit(1)
		}
		validatedModules = append(validatedModules, *module)

		if module.HasNewerVersion() {
			log.Info().
				Str("tfModule", module.Name).
				Str("version_used", module.UsedVersion).
				Str("version_latest", module.AvailableVersion).
				Str("file", fmt.Sprintf("%s:%d", module.Location.FileName, module.Location.Line)).
				Msg(color.New(color.FgRed).Add(color.Bold).Sprint("✖ ···>"))
		} else if module.HasSameVersion() && module.UsedVersion != versionNil {
			if verbose {
				log.Debug().
					Str("tfModule", module.Name).
					Str("file", fmt.Sprintf("%s:%d", module.Location.FileName, module.Location.Line)).
					Msg(color.New(color.FgGreen).Add(color.Bold).Sprint("✔ latest version used"))
			}
		}
	}
	log.Debug().Msg("validation completed")

	var p printer.ModuleVersionPrinter
	switch outputFormat {
	case "table":
		p = printer.TextPrinter{}
	case "json":
		p = printer.JSONPrinter{}
	case "noout":
		// don't define any printer
	}

	if p != nil {
		_ = p.PrintReport(validatedModules, os.Stdout)
	}

	return validatedModules
}

func retrieveLatestVersion(gitlabClient *gitlab.Client, moduleNamespaceName string) (string, error) {
	tags, _, err := gitlabClient.Tags.ListTags(moduleNamespaceName, &gitlab.ListTagsOptions{})
	if err != nil {
		return "", fmt.Errorf("error querying gitlab; potential auth issue; check GITLAB_TOKEN: %w", err)
	}

	var availableVersionsRaw []string
	for _, t := range tags {
		_, err := version.NewVersion(t.Name)
		if err == nil {
			availableVersionsRaw = append(availableVersionsRaw, t.Name)
		}
	}

	sort.Strings(availableVersionsRaw)
	log.Debug().Msg(fmt.Sprintf("available versions for %s raw: %v", moduleNamespaceName, availableVersionsRaw))
	return determineLatestVersionString(availableVersionsRaw)
}

func determineLatestVersionString(availableVersionsRaw []string) (string, error) {
	var versions []*version.Version

	versionToRaw := make(map[*version.Version]string)

	for _, raw := range availableVersionsRaw {
		v, _ := version.NewVersion(raw)
		if v != nil {
			versions = append(versions, v)
			versionToRaw[v] = raw
		} else {
			log.Debug().Msg(fmt.Sprintf("could not determine latest version for %s", raw))
		}
	}
	sort.Sort(version.Collection(versions))
	log.Debug().Msg(fmt.Sprintf("sorted versions = %v", versions))

	var latestVersion *version.Version
	if len(versions) > 0 {
		latestVersion = versions[len(versions)-1]
	}
	if latestVersion == nil {
		return "nil", nil
	}

	return versionToRaw[latestVersion], nil
}
