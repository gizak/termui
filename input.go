package termui

import (
	"strconv"
	"strings"

	"github.com/nsf/termbox-go"
)

// default mappings between /sys/kbd events and multi-line inputs
var multiLineCharMap = map[string]string{
	"<space>":  " ",
	"<tab>":    "\t",
	"<enter>":  "\n",
	"<escape>": "",
}

// default mappings between /sys/kbd events and single line inputs
var singleLineCharMap = map[string]string{
	"<space>":  " ",
	"<tab>":    "\t",
	"<enter>":  "",
	"<escape>": "",
}

const NEW_LINE = "\n"
const LINE_NO_MIN_SPACE = 1000

// TextRenderer can modify text content before it is added to the buffer for rendering
type TextRenderer func(word string, lineNo int, linePos int) string

// EvtInput  defines the structure for the /input/* events. The event contains the last keystroke, the full text
// for the current line, and the position of the cursor in the current line as well as the index of the current
// line in the full text of the input
type EvtInput struct {
	KeyStr         string
	LineText       string
	CursorPosition int
	LineIndex      int
}

// Input is the main object for a text input. The object exposes the following public properties:
// TextFgColor: color for the text.
// TextBgColor: background color for the text box.
// IsCapturing: true if the input is currently capturing keyboard events, this is controlled by the StartCapture and
//              StopCapture methods.
// IsMultiline: Whether we should accept multiple lines of input or this is a singe line form field.
// TextBuilder: An implementation of the TextBuilder interface to customize the look of the text on the screen.
// SpecialChars: a map[string]string of characters from the /sys/kbd events to actual strings in the content.
// Name: When specified, the Input uses its name to propagate events, for example /input/<name>/kbd.
type Input struct {
	Block
	TextFgColor  Attribute
	TextBgColor  Attribute
	IsCapturing  bool
	IsMultiLine  bool
	TextBuilder  TextBuilder
	SpecialChars map[string]string
	ShowLineNo   bool
	Name         string
	CursorX      int
	CursorY      int
	Formatter    TextRenderer

	//DebugMode				bool
	//debugMessage		string

	// internal vars
	lines           []string
	cursorLineIndex int
	cursorLinePos   int
}

// NewInput returns a new, initialized Input object. The method receives the initial content for the input (if any)
// and whether it should be initialized as a multi-line innput field or not
func NewInput(s string, isMultiLine bool) *Input {
	textArea := &Input{
		Block:       *NewBlock(),
		TextFgColor: ThemeAttr("par.text.fg"),
		TextBgColor: ThemeAttr("par.text.bg"),
		TextBuilder: NewMarkdownTxBuilder(),
		IsMultiLine: isMultiLine,
		ShowLineNo:  false,

		cursorLineIndex: 0,
		cursorLinePos:   0,
	}

	if s != "" {
		textArea.addString(s)
	}

	if isMultiLine {
		textArea.SpecialChars = multiLineCharMap
	} else {
		textArea.SpecialChars = singleLineCharMap
	}

	return textArea
}

// StartCapture begins catching events from the /sys/kbd stream and updates the content of the Input field. While
// capturing events, the Input field also publishes its own event stream under the /input/kbd path.
func (i *Input) StartCapture() {
	i.IsCapturing = true
	Handle("/sys/kbd", func(e Event) {
		if i.IsCapturing {
			key := e.Data.(EvtKbd).KeyStr

			switch key {
			case "<up>":
				i.moveUp()
			case "<down>":
				i.moveDown()
			case "<left>":
				i.moveLeft()
			case "<right>":
				i.moveRight()
			case "C-8":
				i.backspace()
			default:
				// If it's a CTRL something we don't handle then just ignore it
				if strings.HasPrefix(key, "C-") {
					break
				}
				newString := i.getCharString(key)
				i.addString(newString)
			}
			if i.Name == "" {
				SendCustomEvt("/input/kbd", i.getInputEvt(key))
			} else {
				SendCustomEvt("/input/"+i.Name+"/kbd", i.getInputEvt(key))
			}

			Render(i)
		}
	})
}

// StopCapture tells the Input field to stop accepting events from the /sys/kbd stream
func (i *Input) StopCapture() {
	i.IsCapturing = false
}

// Text returns the text of the input field as a string
func (i *Input) Text() string {
	if len(i.lines) == 0 {
		return ""
	}

	if len(i.lines) == 1 {
		return i.lines[0]
	}

	if i.IsMultiLine {
		return strings.Join(i.lines, NEW_LINE)
	} else {
		// we should never get here!
		return i.lines[0]
	}
}

func (i *Input) SetText(text string) {
	i.lines = strings.Split(text, NEW_LINE)
}

// Lines returns the slice of strings with the content of the input field. By default lines are separated by \n
func (i *Input) Lines() []string {
	return i.lines
}

// Private methods for the input field
// TODO: handle delete key

func (i *Input) backspace() {
	// we have no lines yet, nothing to do.
	if len(i.lines) == 0 {
		return
	}

	curLine := i.lines[i.cursorLineIndex]
	// at the beginning of the buffer, nothing to do
	if len(curLine) == 0 && i.cursorLineIndex == 0 {
		return
	}

	// at the beginning of a line somewhere in the buffer
	if i.cursorLinePos == 0 {
		prevLine := i.lines[i.cursorLineIndex-1]
		// remove the newline character from the prevline
		prevLine = prevLine[:len(curLine)-1] + curLine
		i.lines = append(i.lines[:i.cursorLineIndex], i.lines[i.cursorLineIndex+1:]...)
		i.cursorLineIndex--
		i.cursorLinePos = len(prevLine) - 1
		return
	}

	// I'm at the end of a line
	if i.cursorLinePos == len(curLine)-1 {
		i.lines[i.cursorLineIndex] = curLine[:len(curLine)-1]
		i.cursorLinePos--
		return
	}

	// I'm in the middle of a line
	i.lines[i.cursorLineIndex] = curLine[:i.cursorLinePos-1] + curLine[i.cursorLinePos:]
	i.cursorLinePos--
}

func (i *Input) addString(key string) {
	if len(i.lines) > 0 {
		if key == NEW_LINE {
			// special case when we go back to the beginning of a buffer with multiple lines, prepend a new line
			if i.cursorLineIndex == 0 && len(i.lines) > 1 {
				i.lines = append([]string{""}, i.lines...)

				// this case handles newlines at the end of the file or in the middle of the file
			} else {
				newString := ""

				// if we are inserting a newline in a populated line then set what goes into the new line
				// and what stays in the current line
				if i.cursorLinePos < len(i.lines[i.cursorLineIndex]) {
					newString = i.lines[i.cursorLineIndex][i.cursorLinePos:]
					i.lines[i.cursorLineIndex] = i.lines[i.cursorLineIndex][:i.cursorLinePos]
				}

				// append a newline in the current position with the content we computed in the previous if statement
				i.lines = append(
					i.lines[:i.cursorLineIndex+1],
					append(
						[]string{newString},
						i.lines[i.cursorLineIndex+1:]...,
					)...,
				)
			}
			// increment the line index, reset the cursor to the beginning and return to skip the next step
			i.cursorLineIndex++
			i.cursorLinePos = 0
			return
		}

		// cursor is at the end of the line
		if i.cursorLinePos == len(i.lines[i.cursorLineIndex]) {
			//i.debugMessage ="end"
			i.lines[i.cursorLineIndex] += key

			// cursor at the beginning of the line
		} else if i.cursorLinePos == 0 {
			//i.debugMessage = "beginning"
			i.lines[i.cursorLineIndex] = key + i.lines[i.cursorLineIndex]

			// cursor in the middle of the line
		} else {
			//i.debugMessage = "middle"
			before := i.lines[i.cursorLineIndex][:i.cursorLinePos]
			after := i.lines[i.cursorLineIndex][i.cursorLinePos:]
			i.lines[i.cursorLineIndex] = before + key + after

		}
		i.cursorLinePos += len(key)

	} else {
		//i.debugMessage = "newline"
		i.lines = append(i.lines, key)
		i.cursorLinePos += len(key)
	}
}

func (i *Input) moveUp() {
	// if we are already on the first line then just move the cursor to the beginning
	if i.cursorLineIndex == 0 {
		i.cursorLinePos = 0
		return
	}

	// The previous line is just as long, we can move to the same position in the line
	prevLine := i.lines[i.cursorLineIndex-1]
	if len(prevLine) >= i.cursorLinePos {
		i.cursorLineIndex--
	} else {
		// otherwise we move the cursor to the end of the previous line
		i.cursorLineIndex--
		i.cursorLinePos = len(prevLine) - 1
	}
}

func (i *Input) moveDown() {
	// we are already on the last line, we just need to move the position to the end of the line
	if i.cursorLineIndex == len(i.lines)-1 {
		i.cursorLinePos = len(i.lines[i.cursorLineIndex])
		return
	}

	// check if the cursor can move to the same position in the next line, otherwise move it to the end
	nextLine := i.lines[i.cursorLineIndex+1]
	if len(nextLine) >= i.cursorLinePos {
		i.cursorLineIndex++
	} else {
		i.cursorLineIndex++
		i.cursorLinePos = len(nextLine) - 1
	}
}

func (i *Input) moveLeft() {
	// if we are at the beginning of the line move the cursor to the previous line
	if i.cursorLinePos == 0 {
		origLine := i.cursorLineIndex
		i.moveUp()
		if origLine > 0 {
			i.cursorLinePos = len(i.lines[i.cursorLineIndex])
		}
		return
	}

	i.cursorLinePos--
}

func (i *Input) moveRight() {
	// if we are at the end of the line move to the next
	if i.cursorLinePos >= len(i.lines[i.cursorLineIndex]) {
		origLine := i.cursorLineIndex
		i.moveDown()
		if origLine < len(i.lines)-1 {
			i.cursorLinePos = 0
		}
		return
	}

	i.cursorLinePos++
}

// Buffer implements Bufferer interface.
func (i *Input) Buffer() Buffer {
	buf := i.Block.Buffer()

	// offset used to display the line numbers
	textXOffset := 0

	bufferLines := i.lines[:]
	firstLine := 0
	lastLine := i.innerArea.Dy()
	if i.IsMultiLine {
		if i.cursorLineIndex >= lastLine {
			firstLine += i.cursorLineIndex - lastLine + 1
			lastLine += i.cursorLineIndex - lastLine + 1
		}

		if len(i.lines) < lastLine {
			bufferLines = i.lines[firstLine:]
		} else {
			bufferLines = i.lines[firstLine:lastLine]
		}
	}

	if i.ShowLineNo {
		// forcing space for up to 1K
		if lastLine < LINE_NO_MIN_SPACE {
			textXOffset = len(strconv.Itoa(LINE_NO_MIN_SPACE)) + 2
		} else {
			textXOffset = len(strconv.Itoa(lastLine)) + 2 // one space at the beginning and one at the end
		}
	}

	text := strings.Join(bufferLines, NEW_LINE)

	// if the last line is empty then we add a fake space to make sure line numbers are displayed
	if len(bufferLines) > 0 && bufferLines[len(bufferLines)-1] == "" && i.ShowLineNo {
		text += " "
	}
	finalText := ""
	if i.Formatter != nil {
		for _, w := range strings.Split(text, " ") {
			if strings.HasPrefix(w, "\n") {
				finalText += "\n" + i.Formatter(strings.Trim(w, "\n"), 0, 0) + " "
			} else if strings.HasSuffix(w, "\n") {
				finalText += i.Formatter(strings.Trim(w, "\n"), 0, 0) + "\n"
			} else {
				finalText += i.Formatter(strings.Trim(w, "\n"), 0, 0) + " "
			}
		}
	} else {
		finalText = text
	}

	fg, bg := i.TextFgColor, i.TextBgColor
	cs := i.TextBuilder.Build(finalText, fg, bg)
	y, x, n := 0, 0, 0
	lineNoCnt := 1

	for n < len(cs) {
		w := cs[n].Width()

		if x == 0 && i.ShowLineNo {
			curLineNoString := " " + strconv.Itoa(lineNoCnt) +
				strings.Join(make([]string, textXOffset-len(strconv.Itoa(lineNoCnt))-1), " ")
			//i.debugMessage = "Line no: " + curLineNoString
			curLineNoRunes := i.TextBuilder.Build(curLineNoString, fg, bg)
			for lineNo := 0; lineNo < len(curLineNoRunes); lineNo++ {
				buf.Set(i.innerArea.Min.X+x+lineNo, i.innerArea.Min.Y+y, curLineNoRunes[lineNo])
			}
			lineNoCnt++
		}

		if cs[n].Ch == '\n' {
			y++
			n++
			x = 0 // set x = 0
			continue
		}
		buf.Set(i.innerArea.Min.X+x+textXOffset, i.innerArea.Min.Y+y, cs[n])

		n++
		x += w
	}

	cursorXOffset := i.X + textXOffset
	if i.BorderLeft {
		cursorXOffset++
	}

	cursorYOffset := i.Y //   termui.TermHeight() - i.innerArea.Dy()
	if i.BorderTop {
		cursorYOffset++
	}
	if lastLine > i.innerArea.Dy() {
		cursorYOffset += i.innerArea.Dy() - 1
	} else {
		cursorYOffset += i.cursorLineIndex
	}
	if i.IsCapturing {
		i.CursorX = i.cursorLinePos + cursorXOffset
		i.CursorY = cursorYOffset
		termbox.SetCursor(i.cursorLinePos+cursorXOffset, cursorYOffset)
	}

	return buf
}

func (i *Input) getCharString(s string) string {
	if val, ok := i.SpecialChars[s]; ok {
		return val
	} else {
		return s
	}
}

func (i *Input) getInputEvt(key string) EvtInput {
	return EvtInput{
		KeyStr:         key,
		LineText:       i.lines[i.cursorLineIndex],
		CursorPosition: i.cursorLinePos,
		LineIndex:      i.cursorLineIndex,
	}
}
