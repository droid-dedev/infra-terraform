package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/zdungey/terraform-state-local/cmd/terraform-state-local/config"
	"github.com/zdungey/terraform-state-local/cmd/terraform-state-local/state"
)

func GenerateTerraformConfigFile(tslConfig *config.Config, rootDir string) error {
	// Create the root directory if it doesn't exist
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		err := os.MkdirAll(rootDir, 0755)
		if err != nil {
			return err
		}
	}

	// Write the terraform configuration file
	terraformConfigFile := filepath.Join(rootDir, "main.tf")
	terraformConfig, err := json.Marshal(tslConfig)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(terraformConfigFile, terraformConfig, 0644)
	if err != nil {
		return err
	}

	return nil
}

func GenerateTerraformStateFile(tslConfig *config.Config, rootDir string) error {
	// Create the root directory if it doesn't exist
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		err := os.MkdirAll(rootDir, 0755)
		if err != nil {
			return err
		}
	}

	// Write the terraform state file
	stateFilePath := filepath.Join(rootDir, "terraform.tfstate")
	if _, err := os.Stat(stateFilePath); err == nil {
		return fmt.Errorf("terraform state file already exists at %s", stateFilePath)
	}

	// Generate a state file
	stateFile, err := state.GenerateStateFile(tslConfig)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(stateFilePath, stateFile, 0644)
	if err != nil {
		return err
	}

	return nil
}

func GenerateRandomName(prefix string) string {
	id := uuid.New()
	name := fmt.Sprintf("%s-%s", prefix, id)
	return strings.ToLower(name)
}

func LoadConfigFile(filePath string) (*config.Config, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config config.Config
	err = mapstructure.Decode(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadTerraformConfigFile(rootDir string) (*config.Config, error) {
	terraformConfigFilePath := filepath.Join(rootDir, "main.tf")
	return LoadConfigFile(terraformConfigFilePath)
}

func LoadTerraformStateFile(rootDir string) (*state.TerraformState, error) {
	stateFilePath := filepath.Join(rootDir, "terraform.tfstate")
	return state.LoadStateFile(stateFilePath)
}