package config

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/catalystsquad/app-utils-go/logging"
	"github.com/spf13/viper"
)

// GetConfigFromViper is used to marshall viper's settings into a given struct. Instantiate your config struct and pass
// in a pointer to it. Any validation errors will be returned. To ensure accurate marshalling, make sure your struct's
// json tags match your cobra flag names / viper settings.
func GetConfigFromViper(config interface{}) error {
	settings := viper.AllSettings()
	settingsJson, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	logging.Log.WithField("settings", string(settingsJson)).Debug("viper settings")
	err = json.Unmarshal(settingsJson, &config)
	if err != nil {
		return err
	}
	_, err = govalidator.ValidateStruct(config)
	return err
}
