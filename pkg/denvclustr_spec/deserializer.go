package denvclustr

import (
	"encoding/json"
	"fmt"
	"strings"
)

// DeserializeConfiguration parses the JSON configuration data and returns
// the deserialized configuration, informational messages, warnings, and any errors.
func DeserializeConfiguration(data []byte) (*DenvclustrConfiguration, []string, []string, error) {
	var config DenvclustrConfiguration
	var warnings []string
	var infos []string

	// Deserialize the JSON data
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to parse configuration: %w", err)
	}

	// Process each devcontainer
	for id, devcontainer := range config.Devcontainers {
		// Trim all string fields
		devcontainer.Name = strings.TrimSpace(devcontainer.Name)
		devcontainer.Description = strings.TrimSpace(devcontainer.Description)
		devcontainer.RepositoryURL = strings.TrimSpace(devcontainer.RepositoryURL)
		devcontainer.RepositoryPath = strings.TrimSpace(devcontainer.RepositoryPath)
		devcontainer.Branch = strings.TrimSpace(devcontainer.Branch)

		// Check if both repository URL and path are specified
		if devcontainer.RepositoryURL != "" && devcontainer.RepositoryPath != "" {
			return nil, infos, warnings, fmt.Errorf("devcontainer '%s' has both repository_url and repository_path specified; exactly one is required", id)
		}

		// Check if neither repository URL nor path are specified
		if devcontainer.RepositoryURL == "" && devcontainer.RepositoryPath == "" {
			return nil, infos, warnings, fmt.Errorf("devcontainer '%s' has neither repository_url nor repository_path specified; exactly one is required", id)
		}

		// Set default name if not provided
		if devcontainer.Name == "" {
			devcontainer.Name = id
			infos = append(infos, fmt.Sprintf("devcontainer '%s' has no name specified; using ID as name", id))
		}

		// Set default branch if not provided and using URL
		if devcontainer.Branch == "" && devcontainer.RepositoryURL != "" {
			devcontainer.Branch = "main"
			infos = append(infos, fmt.Sprintf("devcontainer '%s' has no branch specified; using 'main' as default", id))
		}

		// Warn if branch is specified with repository_path
		if devcontainer.Branch != "" && devcontainer.RepositoryPath != "" {
			warnings = append(warnings, fmt.Sprintf("devcontainer '%s' has branch specified with repository_path; branch will be ignored", id))
		}
	}

	return &config, infos, warnings, nil
}
