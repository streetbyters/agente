package model

type Model interface {
	ToJSON() string
	Validate() (map[string]string, error)
}
