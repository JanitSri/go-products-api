package data

type DataStore interface {
	connect()
	disconnect()
	ping()
	read(productId int) Products
	search(searchTerm string) Products
	create(p Product) string
	delete(productId int) int
	update(p Product) int
}

func InitializeDBConnection(d DataStore) {
	d.connect()
}

func HealthCheck(d DataStore) {
	d.ping()
}

func CloseDBConnection(d DataStore) {
	d.disconnect()
}

func readData(d DataStore) Products {
	return readDataById(d, -1)
}

func readDataById(d DataStore, productId int) Products {
	return d.read(productId)
}

func insertData(d DataStore, p Product) string {
	return d.create(p)
}

func deleteData(d DataStore, productId int) int {
	return d.delete(productId)
}

func updateData(d DataStore, p Product) int {
	return d.update(p)
}

func searchData(d DataStore, searchTerm string) Products {
	return d.search(searchTerm)
}
