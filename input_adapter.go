package eltee

type InputAdapter interface {
	UpdateControlPoints()  // Update any control points with values from this input adapter
	ObserveControlPoints() // Update this input adapter with any values from the control points
}

type InputAdapterRegistration struct {
	ia InputAdapter

	name string // Doesn't change, and is primary key
}
