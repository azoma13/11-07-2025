package entity

type Task struct {
	Id     int
	Status string
	Files  []File
}

type File struct {
	IdFile  int
	UrlFile string
}
