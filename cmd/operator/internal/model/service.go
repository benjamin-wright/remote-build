package model

type Service struct {
	Name      string
	Namespace string
}

func (s Service) Equal(old Service) bool {
	return s == old
}
