package types

import "time"

type Status struct {
	CN string
	UN string
	UD string
	A bool
	S string
	T time.Time
}
