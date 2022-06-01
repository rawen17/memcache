package api

type Person struct {
	Name string
	Age  int32
}

type SetPersonRequest struct {
	Person Person
	TTL    int64
}

type SetPersonResponse struct {
	ID string
}

type GetPersonRequest struct {
	ID string
}

type GetPersonResponse struct {
	Person Person
}

type DeletePersonRequest struct {
	ID string
}

type DeletePersonResponse struct {
}
