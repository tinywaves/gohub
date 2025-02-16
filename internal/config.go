package internal

const (
	Port                       = 8080
	SessionDataKey             = "gohub-user"
	SessionLastRefreshKey      = "gohub-last-refresh"
	SessionLastRefreshInterval = 60
	MysqlDsn                   = "root:root@tcp(localhost:13306)/gohub"
	DevUrl                     = "http://localhost"
	ProdUrl                    = "https://gohub.com"
	AuthenticationKey          = "f2N07VpRqXpOdzWySWdz9LRlmQlqLMkv"
	EncryptionKey              = "C8fB0Ryk2BGQcJ08TqtKObw7AEhVIjB1"
	SessionStoreName           = "gohub-session"
	RedisSize                  = 16
	RedisNetwork               = "tcp"
	RedisAddress               = "localhost:16379"
	RedisPassword              = ""
)
