package entity

type Pipeline struct {
	ID            int
	Name          string `yaml:"name"`
	WorkDirectory string `yaml:"work_dir"`
	Steps         []Step `yaml:"steps"`
}

type Step struct {
	Name string `yaml:"name"`
	Run  string `yaml:"run"`
}
