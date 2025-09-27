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
	return s.Name == old.Name &&
		s.Namespace == old.Namespace &&
		s.CPU == old.CPU &&
		s.Memory == old.Memory &&
		s.Disk == old.Disk &&
		s.Image == old.Image
}

func (s StatefulSet) ID() string {
	return s.Namespace + "/" + s.Name
}
