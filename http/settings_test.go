package main

func NewTestSettingsHttp(rc *RedisContainer, pc *PostgresContainer) (SettingsHttp, error) {
	redisDsn := ""
	if rc != nil {
		redisDsn = rc.Endpoint
	}
	postgresDsn := ""
	if pc != nil {
		postgresDsn = pc.Endpoint
	}
	return SettingsHttp{
		RedisDsn: redisDsn,
		PostgresDsn: postgresDsn,
		// ProjectsBaseUrl: "http://projects_base_url",
		ProjectsBaseUrl: "http://localhost:9090",
		ProjectsApiTimeout: 5,
	}, nil
}
