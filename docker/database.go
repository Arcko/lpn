package docker

// DBName name of the default database
const DBName = "lportal"

// DBPassword default credentials for the database
const DBPassword = "my-secret-pw"

// DatabaseImage interface defining the contract for database docker images
type DatabaseImage interface {
	GetContainerName() string
	GetDataFolder() string
	GetDockerHubTagsURL() string
	GetEnvVariables() EnvVariables
	GetJDBCConnection() JDBCConnection
	GetFullyQualifiedName() string
	GetLpnType() string
	GetPort() int
	GetRepository() string
	GetTag() string
	GetType() string
}

// GetAlias returns the alias used to link containers
func GetAlias() string {
	return "db"
}

// EnvVariables defines how to configure the internal variables for the database
type EnvVariables struct {
	Password string
	Database string
}

// JDBCConnection defines the JDBC connection to the database
type JDBCConnection struct {
	DriverClassName string
	Password        string
	User            string
	URL             string
}
