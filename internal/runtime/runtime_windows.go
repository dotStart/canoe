//go:build windows

/*
 * Copyright 2021 Johannes Donath <johannesd@torchmind.com>
 * and other copyright owners as documented in the project's IP log.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package runtime

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"path"
	"strconv"
)

const rootKey = "SOFTWARE\\JavaSoft\\JDK"
const latestVersionKey = "CurrentVersion"
const javaHomeKey = "JavaHome"

const CliExecutableName = "java.exe"
const GuiExecutableName = "javaw.exe"

// identifies the latest installed version of Java as indicated by the registry.
func getLatestVersion() (uint64, error) {
	root, err := registry.OpenKey(registry.LOCAL_MACHINE, rootKey, registry.QUERY_VALUE)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return 0, ErrNotFound
		}

		return 0, fmt.Errorf("failed to open JavaSoft registry key: %w", err)
	}

	latestVersion, _, err := root.GetStringValue(latestVersionKey)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return 0, ErrNotFound
		}

		return 0, fmt.Errorf("failed to open CurrentVersion registry key: %w", err)
	}

	parsedVersion, err := strconv.ParseUint(latestVersion, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("illegal Java version: \"%s\": %w", latestVersion, err)
	}

	return parsedVersion, nil
}

// locates the installation root for a given version of Java as indicated by the registry.
func findRootForVersion(version uint64) (string, error) {
	p := fmt.Sprintf("%s\\%d", rootKey, version)

	root, err := registry.OpenKey(registry.LOCAL_MACHINE, p, registry.QUERY_VALUE)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return "", ErrNotFound
		}

		return "", fmt.Errorf("failed to open installation registration for v%d: %w", version, err)
	}

	home, _, err := root.GetStringValue(javaHomeKey)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return "", ErrNotFound
		}

		return "", fmt.Errorf("%w for v%d: %s", ErrInvalidInstallation, version, err)
	}

	return home, nil
}

// locates a given minimum version of Java within the current execution environment
func Find(minimumVersion uint64, maximumVersion uint64) (string, error) {
	latestVersion, err := getLatestVersion()
	if err != nil {
		return "", err
	}

	if latestVersion < minimumVersion {
		return "", fmt.Errorf("%w: %d required (%d found)", ErrUnsupported, minimumVersion, latestVersion)
	}
	if maximumVersion != 0 && latestVersion > maximumVersion {
		return "", fmt.Errorf("%w: %d and newer are unsupported (%d found)", ErrUnsupported, maximumVersion, latestVersion)
	}

	root, err := findRootForVersion(latestVersion)
	if err != nil {
		return "", err
	}

	return path.Join(root, "bin"), nil
}
