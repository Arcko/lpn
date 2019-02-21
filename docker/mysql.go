package docker

// MySQLDefaultTag Default tag for MySQL
const MySQLDefaultTag = "5.7"

// MySQLRepository Namespace for the official Docker releases for MySQL
const MySQLRepository = "mdelapenya/mysql-utf8"

// MySQL represents a MySQL image
type MySQL struct {
	LpnType string
	Tag     string
}

// GetContainerName returns the name of the container generated by this type of image
func (m MySQL) GetContainerName() string {
	return "db-" + m.GetLpnType()
}

// GetDataFolder returns the data folder for the database
func (m MySQL) GetDataFolder() string {
	return "/var/lib/mysql"
}

// GetDockerHubTagsURL returns the URL of the available tags on Docker Hub
func (m MySQL) GetDockerHubTagsURL() string {
	return "mysql"
}

// GetEnvVariables returns the specific environment variables to configure the docker image
func (m MySQL) GetEnvVariables() EnvVariables {
	return EnvVariables{
		Database: "MYSQL_DATABASE=" + DBName,
		Password: "MYSQL_ROOT_PASSWORD=" + DBPassword,
	}
}

// GetJDBCConnection returns the JDBC connection
func (m MySQL) GetJDBCConnection() JDBCConnection {
	return JDBCConnection{
		DriverClassName: "com.mysql.jdbc.Driver",
		Password:        DBPassword,
		URL:             "jdbc:mysql://" + GetAlias() + "/lportal?characterEncoding=UTF-8&dontTrackOpenResources=true&holdResultsOpenOverStatementClose=true&useFastDateParsing=false&useUnicode=true",
		User:            "root",
	}
}

// GetFullyQualifiedName returns the fully qualified name of the image
func (m MySQL) GetFullyQualifiedName() string {
	return m.GetRepository() + ":" + m.GetTag()
}

// GetLpnType returns the type of the lpn image
func (m MySQL) GetLpnType() string {
	return m.LpnType
}

// GetPort returns the bind port of the service
func (m MySQL) GetPort() int {
	return 3301
}

// GetRepository returns the repository for MySQL
func (m MySQL) GetRepository() string {
	return MySQLRepository
}

// GetTag returns the tag of the image
func (m MySQL) GetTag() string {
	if m.Tag == "" {
		return MySQLDefaultTag
	}

	return m.Tag
}

// GetType returns the type of the image
func (m MySQL) GetType() string {
	return "mysql"
}
