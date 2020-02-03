package engine

// Situation represents all the properties of an object in space
// type naming...
type Situation struct {
	Position  Vector2
	Velocity  Vector2
	HorizFlip bool
	VertFlip  bool
	DiagFlip  bool
}
