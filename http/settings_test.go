package main

func NewTestSettingsHttp(rc *RedisContainer, pc *PostgresContainer) (*SettingsHttp, error) {
	redisDsn := ""
	if rc != nil {
		redisDsn = rc.Endpoint
	}
	postgresDsn := ""
	if pc != nil {
		postgresDsn = pc.Endpoint
	}
	return &SettingsHttp{RedisDsn: redisDsn, PostgresDsn: postgresDsn}, nil
}
