package testing

type Book struct {
	id   int
	name string
}

func getTrue() bool {
	return true
}

type Library struct {
	books []Book
	addr  string
}
