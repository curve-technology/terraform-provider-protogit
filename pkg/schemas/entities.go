package schemas

type Entries []Entry

type Entry struct {
	Topic    string
	Section  Section
	Filepath string
}

type Records []Record

type Record struct {
	Subject    string
	Schema     string
	References []string
}
