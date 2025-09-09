package service

type Calculator struct{}

func (c *Calculator) Add(params map[string]interface{}) float64 {
	a, _ := params["A"].(float64)
	b, _ := params["B"].(float64)
	return a + b
}

func (c *Calculator) Sub(params map[string]interface{}) float64 {
	a, _ := params["A"].(float64)
	b, _ := params["B"].(float64)
	return a - b
}
