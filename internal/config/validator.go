package config

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type Validator struct {
	EmailValidate    string `mapstructure:"email"`
	PasswordValidate string `mapstructure:"password"`
	UserIDValidate   string `mapstructure:"user_id"`
	AppIDValidate    string `mapstructure:"app_id"`
}

func (v *Validator) mustBeRegex() {
	b, _ := json.Marshal(v)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	for i, v := range m {
		_, err := regexp.Compile(v.(string))
		if err != nil {
			panic(fmt.Sprintf("incorrect %s", i))
		}
	}
}