package service

type ConfigSorter struct {
	MiddlewareConfigs []*MiddlewareConfig
}

func (o ConfigSorter) Len() int {
	return len(o.MiddlewareConfigs)
}

func (o ConfigSorter) Less(i, j int) bool {
	return o.MiddlewareConfigs[i].Priority < o.MiddlewareConfigs[j].Priority && o.MiddlewareConfigs[i].Priority != 0
}

func (o ConfigSorter) Swap(i, j int) {
	o.MiddlewareConfigs[i], o.MiddlewareConfigs[j] = o.MiddlewareConfigs[j], o.MiddlewareConfigs[i]
}

//NewPollSorter custom sort to sort Polls by creted date
func NewMiddlewareConfigSorter(configs []*MiddlewareConfig) ConfigSorter {
	return ConfigSorter{MiddlewareConfigs: configs}
}
