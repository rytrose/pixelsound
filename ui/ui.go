package ui

// UI defines a set of functionality for different UI
// implementations, namely a MacOS thick client implementation
// versus a browser implementation.
type UI interface {
	// Run starts the UI. Must be blocking.
	Run()
}
