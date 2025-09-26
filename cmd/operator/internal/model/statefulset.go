package model

type StatefulSet struct {
	Name      string
	Namespace string
	CPU       string
	Memory    string
	Disk      string
	Image     string
	Ready     bool
}

func (s StatefulSet) Equal(old StatefulSet) bool {
	return s == old
}
