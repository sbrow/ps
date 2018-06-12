package ps

import "fmt"

// ModeEnum determines how aggressively the package will attempt to sync with Photoshop.
// Loading Photoshop files from scratch takes a long time, so the package saves
// the state of the document in a JSON file in the /data folder whenever you call
// Document.Dump(). ModeEnum tells the program how trustworthy that file is.
type ModeEnum int

// Mode holds the current mode.
var Mode ModeEnum

const (
	// Normal Mode only verifies layers as they are operated on. The first time a
	// layer's properties would be checked, it first overwrites the data from the
	// Dump with data pulled directly from Photoshop. This allows you to quickly
	// load documents in their current form.
	Normal ModeEnum = iota

	// Safe Mode always loads the document from scratch, ignoring any dumped data.
	// (Very Slow). If a function panics due to outdated data, often times re-running
	// the function in safe mode is enough to remediate it.
	Safe

	// Fast mode skips all verification. Use Fast mode only when certain that the
	// .psd file hasn't changed since the last time Document.Dump() was called.
	Fast
)

// const Fast ModeEnum = 2

// const Normal ModeEnum = 0

// const Safe ModeEnum = 1

// SaveOption is an enum for options when closing a document.
type SaveOption int

func (p *SaveOption) String() string {
	return fmt.Sprint("", *p) // TODO: Fix
}

// SaveChanges Saves changes before closing documents.
const SaveChanges SaveOption = 1

// DoNotSaveChanges Closes documents without saving.
const DoNotSaveChanges SaveOption = 2

// PromptToSaveChanges prompts the user whether the file
// should be saved before closing.
const PromptToSaveChanges SaveOption = 3
