package constants

type User struct {
	Id    uint64
	Email string
	Name  string
}

type List struct {
	Id       uint64
	Name     string
	ShareURL string
}

type Item struct {
	Id       uint64
	List     List
	Name     string
	URL      string
	Notes    string
	Priority int
}
