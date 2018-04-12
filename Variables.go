package ps

import (
	"fmt"
)

var Colors map[string]Color = map[string]Color{
	"Black": &RGB{0, 0, 0},
	"Gray":  &RGB{128, 128, 128},
	"White": &RGB{255, 255, 255},
}

// ModeEnum determines how aggressively the package will attempt to sync with Photoshop.
type ModeEnum int

// Holds the current mode.
var Mode ModeEnum

// Fast mode never checks layers before returning.
const Fast ModeEnum = 2

// Normal Mode Always checks to see if layers are up to date
// before returning them.
const Normal ModeEnum = 0

// Safe Mode Always loads the document from scratch. (Very Slow)
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
