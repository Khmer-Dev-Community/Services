package config

type GitConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientScrete string `yaml:"client_secret"`
	ClientRedir  string `yaml:"client_redir"`
	ClientEnd    string `yaml:"client_endpoint"`
}
