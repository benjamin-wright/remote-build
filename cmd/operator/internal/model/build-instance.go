package model

type BuildInstance struct {
	Name      string
	Namespace string
	CPU       string
	Memory    string
	Disk      string
	Image     string
	State     string
	Active    bool
}

func (b BuildInstance) Equal(old BuildInstance) bool {
	return b == old
}

func (b BuildInstance) ID() string {
	return b.Namespace + "/" + b.Name
}
