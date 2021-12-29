package data

type dataStore interface {
	connect()
	disconnect()
	ping()
	read(productId int) Products
	search(searchTerm string) Products
	create(p Product) string
	delete(productId int) int
	update(p Product) int
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

func readData(d dataStore) Products {
	return readDataById(d, -1)
}

func readDataById(d dataStore, productId int) Products {
	return d.read(productId)
}

func insertData(d dataStore, p Product) string {
	return d.create(p)
}

func deleteData(d dataStore, productId int) int {
	return d.delete(productId)
}

func updateData(d dataStore, p Product) int {
	return d.update(p)
}

func searchData(d dataStore, searchTerm string) Products {
	return d.search(searchTerm)
}
