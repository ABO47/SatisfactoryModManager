package installfinders

import (
	"log/slog"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/settings"

	_ "github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/launchers/all" // register all launchers
)

func FindInstallations() ([]*common.Installation, []error) {
	registrations := launchers.GetInstallFinders()

	slog.Debug("finding installations", slog.String("launchers", strings.Join(maps.Keys(registrations), ",")))

	installs, errors := common.FindAll(maps.Values(registrations)...)

	// Add custom game paths from settings
	customPaths := settings.Settings.GetCustomGamePaths()
	customInstalls, customErrors := findCustomInstallations(customPaths)
	
	installs = append(installs, customInstalls...)
	errors = append(errors, customErrors...)

	return installs, errors
}

func findCustomInstallations(paths []string) ([]*common.Installation, []error) {
	var installations []*common.Installation
	var errors []error

	platform := common.NativePlatform()

	for _, path := range paths {
		installType, version, savedPath, err := common.GetGameInfo(path, platform)
		if err != nil {
			slog.Debug("custom path is not a valid installation", slog.String("path", path), slog.Any("error", err))
			errors = append(errors, err)
			continue
		}

		install := &common.Installation{
			Path:     path,
			Version:  version,
			Type:     installType,
			Location: common.LocationTypeLocal,
			Branch:   common.GameBranchStable,
			Launcher: "Custom",
			SavedPath: savedPath,
		}

		installations = append(installations, install)
	}

	return installations, errors
}
