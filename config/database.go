package config

import "time"

var (
	// Return connection string from database section in toml file
	GetDatabaseConnectionString = getDatabaseConnectionString

	// Return max idle connections allowed from database section in toml file
	GetDatabaseMaxIdleConnections = getDatabaseMaxIdleConnections

	// Return max open connections allowed from database section in toml file
	GetDatabaseMaxOpenConnections = getDatabaseMaxOpenConnections

	// Return connection lifetime from database section in toml file
	GetDatabaseConnectonLifetime = getDatabaseConnectonLifetime
)

func getDatabaseConnectionString() string {
	return getConfigString("database.connection_string")
}

func getDatabaseMaxIdleConnections() int {
	return getConfigInt("database.max_idle_connections")
}

func getDatabaseMaxOpenConnections() int {
	return getConfigInt("database.max_open_connections")
}

func getDatabaseConnectonLifetime() time.Duration {
	return getConfigDuration("database.connection_max_lifetime")
}
