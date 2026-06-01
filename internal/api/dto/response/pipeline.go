package response

type Pipeline struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	WorkDirectory string `json:"work_dir"`
	Steps         []Step `json:"steps"`
}

type Step struct {
	Name string `json:"name"`
	Run  string `json:"run"`
}
