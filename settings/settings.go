package settings

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/sirupsen/logrus"
)

const (
	settingsFilePath = "./settings.json"
)

//Config truct for settings.json
type Config struct {
	FTP       *FTP
	Env       *Env
	IgnoreStr []string
	Ignore    []*regexp.Regexp
}

//Load loads the configuration from flags or settings file
func Load() (config *Config) {
	var err error
	err = validateFlags()

	if err != nil && !reqVersion {
		config, err = readSettings()

		if err != nil {
			logrus.WithError(err).Warn("loading settings from file")
		}

		if config == nil {
			config = &Config{
				FTP:    NewFTP("", ""),
				Env:    NewEnv(),
				Ignore: []*regexp.Regexp{},
			}
		}
	} else {
		config = &Config{
			FTP: &FTP{
				Username: username,
				Password: password,
				Host:     host,
				Port:     port,
				RootPath: startPath,
				DestPath: destPath,
			},
			Env: &Env{
				NeedWait:   needWait,
				ReqVersion: reqVersion,
				Store:      store,
			},
		}

		if config.Env.Store {
			go config.Save()
		}
	}

	return
}

//readSettings reads data from settings.json
func readSettings() (*Config, error) {
	data, err := readSettingsFile(settingsFilePath)

	if err != nil && data == nil {
		return nil, err
	}

	if data == nil {
		data = &Config{
			FTP: NewFTP("", ""),
			Env: NewEnv(),
		}
	}
	return data, err
}

//readSettingsFile reads data from the specified file
func readSettingsFile(filePath string) (*Config, error) {
	//checks if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, nil
	}

	file, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	data := new(Config)
	err = json.Unmarshal([]byte(file), data)

	if err != nil {
		return nil, err
	}

	for _, str := range data.IgnoreStr {
		rg, err := regexp.Compile(str)
		if err != nil {
			continue
		}

		data.Ignore = append(data.Ignore, rg)
	}

	return data, nil
}

//Validate validates the configuration
func (c *Config) Validate() error {
	// if c.FTP.Username == "" {
	// 	flag.PrintDefaults()
	// 	return fmt.Errorf("username is required")
	// }

	if c.FTP.Host == "" {
		flag.PrintDefaults()
		return fmt.Errorf("host is required")
	}

	return nil
}

//Save saves the config information into settings
func (c *Config) Save() error {
	file, err := json.MarshalIndent(*c, "", "\t")

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(settingsFilePath, file, 0644)
	if err != nil {
		return err
	}
	return nil
}

//GetURL returns a url, joining host and port
func (c *Config) GetURL() string {
	return fmt.Sprintf("%s:%d", c.FTP.Host, c.FTP.Port)
}

//IgnoreFile returns true if the file path match with the regex
func (c *Config) IgnoreFile(filePath string) bool {
	for _, rg := range c.Ignore {
		if rg.Match([]byte(filePath)) {
			return true
		}
	}

	return false
}
