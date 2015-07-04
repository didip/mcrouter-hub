package models

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

func NewMcRouterConfigManager(mcRouterConfigFile string) (*McRouterConfigManager, error) {
	if mcRouterConfigFile == "" {
		return nil, errors.New("McRouter config file is missing")
	}

	m := &McRouterConfigManager{}
	m.McRouterConfigFile = mcRouterConfigFile

	fileInfo, err := os.Stat(mcRouterConfigFile)
	if err != nil {
		return nil, err
	}
	m.McRouterConfigFileInfo = fileInfo

	return m, nil
}

type McRouterConfigManager struct {
	McRouterConfigFile     string
	McRouterConfigFileInfo os.FileInfo
}

func (m *McRouterConfigManager) ConfigJson() ([]byte, error) {
	return ioutil.ReadFile(m.McRouterConfigFile)
}

func (m *McRouterConfigManager) UpdateConfigJson(mcRouterConfigJson []byte) error {
	return ioutil.WriteFile(m.McRouterConfigFile, mcRouterConfigJson, m.McRouterConfigFileInfo.Mode())
}

func (m *McRouterConfigManager) Config() (map[string]interface{}, error) {
	mcRouterConfigJson, err := m.ConfigJson()
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}

	err = json.Unmarshal(mcRouterConfigJson, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (m *McRouterConfigManager) PoolsJson() ([]byte, error) {
	config, err := m.Config()
	if err != nil {
		return nil, err
	}

	return json.Marshal(config["pools"])
}

func (m *McRouterConfigManager) UpdatePoolsJson(poolsJson []byte) error {
	config, err := m.Config()
	if err != nil {
		return err
	}

	var pools map[string]interface{}

	err = json.Unmarshal(poolsJson, &pools)
	if err != nil {
		return err
	}

	config["pools"] = pools
	configJson, err := json.Marshal(config)

	return ioutil.WriteFile(m.McRouterConfigFile, configJson, m.McRouterConfigFileInfo.Mode())
}
