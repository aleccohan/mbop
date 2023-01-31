package models

import "encoding/json"

type Status struct {
    ConfiguredModules ConfiguredModules `json:"configured_modules"`
}

type ConfiguredModules struct {
    Users  string `json:"users"`
    Mailer string `json:"mailer"`
    JWT    string `json:"jwt"`
}

func (s *Status) ToJSON() []byte {
    bytes, _ := json.Marshal(s)
    return bytes
}
