package model

type Service struct {
	Name      string
	Namespace string
}

func (s Service) Equal(old Service) bool {
	return s == old
}

func (s Service) ID() string {
	return s.Namespace + "/" + s.Name
}
