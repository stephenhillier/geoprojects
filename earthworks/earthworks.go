package earthworks

// Settings contains global settings that handlers, methods and services can reference
type Settings struct {
	DefaultPageLimit int
	MaxPageLimit     int
}

// ContextKey is used for retrieving data from a context
type ContextKey struct {
	Name string
}
