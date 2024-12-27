package config

type ConfigServer struct {
	ScratchDir      *string          `json:"scratchdir" validate:"required"`
	PgSQLConfig     *ConfigPgSQL     `json:"pgdb" validate:"required"`
	APIServerConfig *ConfigAPIServer `json:"apiserver" validate:"required"`
}

type ConfigPgSQL struct {
	PgURL *string `json:"pgurl" validate:"required"`
}
type ConfigAPIServer struct {
	Port int `json:"port" validate:"required,numeric"`
}
