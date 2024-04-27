package constants

type UserRead struct {
	Id    uint64
	Email string
	Name  string
	Admin bool
}

type UserRegister struct {
	Email    string
	Name     string
	Password string
	Admin    bool
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
