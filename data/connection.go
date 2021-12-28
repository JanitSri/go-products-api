package data

type dataStore interface {
	connect()
	disconnect()
	ping()
	read(fiter interface{}) Products
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

func readData(d dataStore, filter interface{}) Products {
	return d.read(filter)
}
