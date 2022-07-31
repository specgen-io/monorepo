package console

import (
	"fmt"
)

const VERBOSITY_QUIET = 0
const VERBOSITY_NORMAL = 1
const VERBOSITY_VERBOSE = 2

var Verbosity = VERBOSITY_NORMAL

func Print(args ...interface{}) {
	log := fmt.Sprint(args...)
	outNormal(VERBOSITY_QUIET, log)
}

func PrintLn(args ...interface{}) {
	log := fmt.Sprintln(args...)
	outNormal(VERBOSITY_QUIET, log)
}

func PrintF(format string, args ...interface{}) {
	log := fmt.Sprintf(format, args...)
	outNormal(VERBOSITY_QUIET, log)
}

func PrintLnF(format string, args ...interface{}) {
	log := fmt.Sprintln(fmt.Sprintf(format, args...))
	outNormal(VERBOSITY_QUIET, log)
}

func Verbose(args ...interface{}) {
	log := fmt.Sprint(args...)
	outNormal(VERBOSITY_QUIET, log)
}

func VerboseLn(args ...interface{}) {
	log := fmt.Sprintln(args...)
	outNormal(VERBOSITY_QUIET, log)
}

func VerboseF(format string, args ...interface{}) {
	log := fmt.Sprintf(format, args...)
	outNormal(VERBOSITY_QUIET, log)
}

func VerboseLnF(format string, args ...interface{}) {
	log := fmt.Sprintln(fmt.Sprintf(format, args...))
	outNormal(VERBOSITY_QUIET, log)
}

func Problem(args ...interface{}) {
	log := fmt.Sprint(args...)
	outProblem(VERBOSITY_QUIET, log)
}

func ProblemLn(args ...interface{}) {
	log := fmt.Sprintln(args...)
	outProblem(VERBOSITY_QUIET, log)
}

func ProblemF(format string, args ...interface{}) {
	log := fmt.Sprintf(format, args...)
	outProblem(VERBOSITY_QUIET, log)
}

func ProblemLnF(format string, args ...interface{}) {
	log := fmt.Sprintln(fmt.Sprintf(format, args...))
	outProblem(VERBOSITY_QUIET, log)
}

func SuccessLn(args ...interface{}) {
	log := fmt.Sprintln(args...)
	outSuccess(VERBOSITY_QUIET, log)
}

func SuccessLnF(format string, args ...interface{}) {
	log := fmt.Sprintln(fmt.Sprintf(format, args...))
	outSuccess(VERBOSITY_QUIET, log)
}

func out(minVerbosity int, content interface{}) {
	if Verbosity >= minVerbosity {
		print(content)
	}
}

func outNormal(minVerbosity int, content interface{}) {
	out(minVerbosity, content)
}

func outProblem(minVerbosity int, content interface{}) {
	out(minVerbosity, content)
}

func outSuccess(minVerbosity int, content interface{}) {
	out(minVerbosity, content)
}
