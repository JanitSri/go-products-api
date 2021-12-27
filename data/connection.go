package data

type dataStore interface {
	connect()
	disconnect()
	ping()
}

func InitializeDBConnection(d dataStore) {
	d.connect()
}

func HealthCheck(d dataStore) {
	d.ping()
}

func CloseDBConnection(d dataStore) {
	d.disconnect()
}
