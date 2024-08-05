package modelsv1

type RequestStruct struct {
	Name       string            `yaml:"name"`
	Endpoint   string            `yaml:"endpoint"`
	Method     string            `yaml:"method"`
	Headers    map[string]string `yaml:"headers"`
	Parameters map[string]string `yaml:"parameters"`
}
