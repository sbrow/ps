package ps

import (
	"fmt"
)

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

const (
	PSSaveChanges         PSSaveOptions = iota + 1 // Saves changes before closing documents.
	PSDoNotSaveChanges                             // Closes documents without saving.
	PSPromptToSaveChanges                          // Prompts whether to save before closing.
)
