package RepoPath

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type EzGitConfig struct {
	Name  string
	Email string
}

const defaultIgnoreTemplate = `# EzGit default ignore template
.git/
*.log
*.tmp
node_modules/
__pycache__/
*.pyc
.DS_Store
dist/
build/
*.o
*.exe
.env
`

// ConfigPath returns the platform-correct path for the config file
// Linux:   ~/.config/ezgit/config
// Windows: %APPDATA%\ezgit\config
func ConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		// fallback if os.UserConfigDir() fails
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine config directory: %w", err)
		}
		return filepath.Join(home, ".ezgit", "config"), nil
	}
	return filepath.Join(configDir, "ezgit", "config"), nil
}

// Load reads the ezgit config file
func Load() (*EzGitConfig, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no ezgit config found — run: ezgit config init")
		}
		return nil, fmt.Errorf("failed to open config: %w", err)
	}
	defer f.Close()

	cfg := &EzGitConfig{}
	inUserSection := false

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") {
			section := strings.ToLower(strings.Trim(line, "[]"))
			inUserSection = strings.TrimSpace(section) == "user"
			continue
		}

		if !inUserSection {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(strings.ToLower(parts[0]))
		val := strings.TrimSpace(parts[1])

		switch key {
		case "name":
			cfg.Name = val
		case "email":
			cfg.Email = val
		}
	}

	if cfg.Name == "" || cfg.Email == "" {
		return nil, fmt.Errorf("config incomplete - set user.name and user.email")
	}

	return cfg, scanner.Err()
}

// Save writes the config file to the correct platform path
func Save(cfg *EzGitConfig) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	// create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	content := fmt.Sprintf("[user]\n\tname = %s\n\temail = %s\n", cfg.Name, cfg.Email)
	return os.WriteFile(path, []byte(content), 0644)
}

// Format returns git-compatible author string
func (c *EzGitConfig) Format() string {
	return fmt.Sprintf("%s <%s>", c.Name, c.Email)
}

func IgnoreDefaultsPath() (string, error) {
	path, err := ConfigPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(path), "ignore_defaults"), nil
}

// EnsureIgnoreDefaults creates the default ignore template config on first run.
func EnsureIgnoreDefaults() (string, error) {
	path, err := IgnoreDefaultsPath()
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(path); err == nil {
		return path, nil // already exists
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to check ignore defaults: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", fmt.Errorf("failed to create config dir: %w", err)
	}
	if err := os.WriteFile(path, []byte(defaultIgnoreTemplate), 0644); err != nil {
		return "", fmt.Errorf("failed to write ignore defaults: %w", err)
	}
	return path, nil
}

func LoadIgnoreDefaults() ([]string, error) {
	path, err := EnsureIgnoreDefaults()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rules []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		rules = append(rules, line)
	}
	return rules, nil
}
