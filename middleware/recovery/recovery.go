package recovery

import (
	"bytes"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"io/ioutil"
	"log"
	"runtime"
)

type jsonPanicError struct {
	Code   int           `json:",omitempty"` // the http response code
	Short  string        `json:",omitempty"` // a short explanation of the response (usually one or two words). for internal use only
	Errors []interface{} `json:",omitempty"` // any errors that may have occured with the request and should be displayed to the user
	From   string
}

func (je jsonPanicError) Error() string {
	if len(je.Errors) == 0 {
		return "Unkown error"
	} else if len(je.Errors) == 1 {
		return fmt.Sprintf("%s", je.Errors[0])
	} else {
		return fmt.Sprintf("Multiple errors: %s", je.Errors)
	}
}

func newJSONPanicError(err interface{}) jsonPanicError {
	_, file, line, _ := runtime.Caller(3)
	return jsonPanicError{
		Code:   500,
		Short:  "internalError",
		Errors: []interface{}{err},
		From:   fmt.Sprintf("%s:%d", file, line),
	}
}

func JSONRecovery() martini.Handler {
	return func(c martini.Context, log *log.Logger, r render.Render) {
		defer func() {
			if err := recover(); err != nil {
				stack := stack(3)
				log.Printf("PANIC: %s\n%s", err, stack)

				// respond with panic message while in development mode
				if martini.Env == martini.Dev {
					if e, ok := err.(error); ok {
						// If err is an error, get the string message
						r.JSON(500, newJSONPanicError(e.Error()))
					} else {
						r.JSON(500, newJSONPanicError(err))
					}
				} else {
					r.JSON(500, newJSONPanicError("Sorry, something went wrong."))
				}
			}
		}()

		c.Next()
	}
}

// stack returns a nicely formated stack frame, skipping skip frames
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
