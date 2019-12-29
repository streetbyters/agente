package model

// Model structure interface
type Model interface {
	ToJSON() string
	Validate() (map[string]string, error)
}
