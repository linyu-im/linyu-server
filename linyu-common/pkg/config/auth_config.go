package config

type LdapConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	Host         string `mapstructure:"host"`
	BaseDN       string `mapstructure:"base-dn"`
	BindDN       string `mapstructure:"bind-dn"`
	BindPassword string `mapstructure:"bind-password"`
	UserFilter   string `mapstructure:"user-filter"`
	Unique       struct {
		LdapField  string `mapstructure:"ldap-field"`
		LocalField string `mapstructure:"local-field"`
	} `mapstructure:"unique"`
}

type GiteeConfig struct {
	ClientID     string `mapstructure:"client-id"`
	ClientSecret string `mapstructure:"client-secret"`
	RedirectURL  string `mapstructure:"redirect-url"`
}
