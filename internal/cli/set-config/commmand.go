package setconfig

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	gitpath "gocmd/testfiles/Gitrepostruct"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage ezgit config",
		RunE:  configRun,
	}

	cmd.Flags().Bool("reset", false, "Reset and re-enter config details")
	return cmd
}

func configRun(cmd *cobra.Command, args []string) error {
	reset, err := cmd.Flags().GetBool("reset")
	if err != nil {
		return err
	}

	cfg, err := gitpath.Load()

	// config exists and no reset flag = just print it
	if err == nil && !reset {
		fmt.Printf("name:  %s\n", cfg.Name)
		fmt.Printf("email: %s\n", cfg.Email)
		fmt.Println("\nto update your config run: ezgit config --reset")
		return nil
	}

	// config missing or reset flag will prompt user
	// to re-enter creds
	newCfg, err := PromptUser()
	if err != nil {
		return err
	}

	if err := gitpath.Save(newCfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	path, _ := gitpath.ConfigPath()
	fmt.Printf("\nconfig saved!\n")
	fmt.Printf("name:  %s\n", newCfg.Name)
	fmt.Printf("email: %s\n", newCfg.Email)
	fmt.Printf("path:  %s\n", path)
	return nil
}

func PromptUser() (*gitpath.EzGitConfig, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("You need an email and name set to use ezgit ")
	fmt.Println("Fill in below to save your versions with Ezgit. TIP: use *ezgit config --reset* to reset these credentials")
	fmt.Print("\n enter your name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read name: %w", err)
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	fmt.Print("enter your email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read email: %w", err)
	}
	email = strings.TrimSpace(email)
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	return &gitpath.EzGitConfig{
		Name:  name,
		Email: email,
	}, nil
}

func GetOrPromptConfig() (*gitpath.EzGitConfig, error) {
	cfg, err := gitpath.Load()
	if err == nil {
		// config exists and is complete — just use it
		return cfg, nil
	}

	// config missing or incomplete — prompt user
	fmt.Println("No ezgit config found. Setting up your email and name to save a version.")
	cfg, err = PromptUser()
	if err != nil {
		return nil, err
	}

	if err := gitpath.Save(cfg); err != nil {
		return nil, fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("Config saved — name: %s, email: %s\n\n", cfg.Name, cfg.Email)
	return cfg, nil
}
