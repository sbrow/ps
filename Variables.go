package ps

import (
	"fmt"
)

// Colors enumerates some basic, commonly used colors.
var Colors map[string]Color = map[string]Color{
	"Black": &RGB{0, 0, 0},
	"Gray":  &RGB{128, 128, 128},
	"White": &RGB{255, 255, 255},
}

// ModeEnum determines how aggressively the package will attempt to sync with Photoshop.
// Loading Photoshop files from scratch takes a long time, so the package saves
// the state of the document in a JSON file in the /data folder whenever you call
// Document.Dump(). ModeEnum tells the program how trustworthy that file is.
type ModeEnum int

// Holds the current mode.
var Mode ModeEnum

// Fast mode skips all verification. Use Fast mode only when certain that the
// .psd file hasn't changed since the last time Document.Dump() was called.
const Fast ModeEnum = 2

// Normal Mode only verifies layers as they are operated on. The first time a
// layer's properties would be checked, it first overwrites the data from the
// Dump with data pulled directly from Photoshop. This allows you to quickly
// load documents in their current form.
const Normal ModeEnum = 0

// Safe Mode always loads the document from scratch, ignoring any dumped data.
// (Very Slow). If a function panics due to outdated data, often times re-running
// the function in safe mode is enough to re-mediate it.
const Safe ModeEnum = 1

// PSSaveOptions is an enum for options when closing a document.
type PSSaveOptions int

func (p *PSSaveOptions) String() string {
	return fmt.Sprint("", *p)
}

// PSSaveChanges saves changes before closing documents.
const PSSaveChanges PSSaveOptions = 1

// PSDoNotSaveChanges closes documents without saving.
const PSDoNotSaveChanges PSSaveOptions = 2

// PSPromptToSaveChanges prompts the user whether to save each
// document before closing it.
const PSPromptToSaveChanges PSSaveOptions = 3
