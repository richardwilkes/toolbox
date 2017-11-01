package cmdline

//go:generate go run cmd/genvalues/main.go

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/richardwilkes/toolbox/atexit"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/i18n"
	"github.com/richardwilkes/toolbox/xio/term"
)

// CmdLine holds information about the command line.
type CmdLine struct {
	// UsageSuffix, if set, will be appended to the 'Usage: ' line of the
	// output.
	UsageSuffix string
	// Description, if set, will be inserted after the program identity
	// section, before the usage.
	Description     string
	options         Options
	cmds            map[string]Cmd
	parent          *CmdLine
	cmd             Cmd
	out             *term.ANSI
	showHelp        bool
	showVersion     bool
	showLongVersion bool
}

// New creates a new CmdLine. If 'includeDefaultOptions' is true, help (-h,
// --help) and version (-v, --version, along with hidden -V, --Version for
// long variants) options will be added, otherwise, only the help options will
// be added, although they will be hidden.
func New(includeDefaultOptions bool) *CmdLine {
	cl := &CmdLine{cmds: make(map[string]Cmd), out: term.NewANSI(os.Stderr)}
	help := cl.NewBoolOption(&cl.showHelp).SetSingle('h').SetName("help")
	if includeDefaultOptions {
		help.SetUsage(i18n.Text("Display this help information and exit"))
		cl.NewBoolOption(&cl.showVersion).SetSingle('v').SetName("version").SetUsage(i18n.Text("Display short version information and exit"))
		cl.NewBoolOption(&cl.showLongVersion).SetSingle('V').SetName("Version").SetUsage(i18n.Text("Display the full version information and exit"))
	}
	return cl
}

// NewOption creates a new Option and attaches it to this CmdLine.
func (cl *CmdLine) NewOption(value Value) *Option {
	option := new(Option)
	option.value = value
	if option.isString() {
		option.def = fmt.Sprintf(`"%s"`, value.String())
	} else {
		option.def = value.String()
	}
	option.arg = i18n.Text("value")
	cl.options = append(cl.options, option)
	return option
}

// Parse the 'args', filling in any options. Returns the remaining arguments
// that weren't used for option content.
func (cl *CmdLine) Parse(args []string) []string {
	const (
		LookForOptionState = iota
		SetOptionValueState
		CollectRemainingState
	)
	var current *Option
	var currentArg string
	state := LookForOptionState
	var remainingArgs []string
	options := cl.availableOptions()
	max := len(args)
	seen := make(map[string]bool)
	for i := 0; i < max; i++ {
		arg := args[i]
		switch state {
		case LookForOptionState:
			if strings.HasPrefix(arg, "@") {
				path := arg[1:]
				if seen[path] {
					cl.FatalMsg(fmt.Sprintf(i18n.Text("Recursive loading of arguments from a file is not permitted: %s"), path))
				}
				seen[path] = true
				insert, err := cl.loadArgsFromFile(path)
				cl.FatalIfError(err)
				args = append(args[:i], append(insert, args[i+1:]...)...)
				max = len(args)
				i--
				continue
			}
			if arg == "--" {
				state = CollectRemainingState
			} else if strings.HasPrefix(arg, "--") {
				var value string
				arg = arg[2:]
				sep := strings.Index(arg, "=")
				if sep != -1 {
					value = arg[sep+1:]
					arg = arg[:sep]
				}
				option := options[arg]
				if option == nil {
					cl.FatalMsg(fmt.Sprintf(i18n.Text("Invalid option: --%s"), arg))
				} else if option.isBool() {
					if sep != -1 {
						cl.FatalMsg(fmt.Sprintf(i18n.Text("Option --%[1]s does not allow an argument: %[2]s"), arg, value))
					} else {
						cl.setOrFail(option, "--"+arg, "true")
					}
				} else if sep != -1 {
					cl.setOrFail(option, "--"+arg, value)
				} else {
					state = SetOptionValueState
					current = option
					currentArg = "--" + arg
				}
			} else if strings.HasPrefix(arg, "-") {
				arg = arg[1:]
				for j, ch := range arg {
					if option := options[string(ch)]; option != nil {
						if option.isBool() {
							cl.setOrFail(option, "-"+arg, "true")
						} else if j == len(arg)-1 {
							state = SetOptionValueState
							current = option
							currentArg = "-" + arg[j:j+1]
						} else if arg[j+1:j+2] == "=" {
							cl.setOrFail(option, "-"+arg, arg[j+2:])
							break
						} else {
							cl.setOrFail(option, "-"+arg, arg[j+1:])
							break
						}
					} else {
						cl.FatalMsg(fmt.Sprintf(i18n.Text("Invalid option: -%s"), arg[j:]))
						break
					}
				}
			} else {
				remainingArgs = append(remainingArgs, arg)
				state = CollectRemainingState
			}
		case SetOptionValueState:
			cl.setOrFail(current, currentArg, arg)
			state = LookForOptionState
		case CollectRemainingState:
			remainingArgs = append(remainingArgs, arg)
		}
	}
	if state == SetOptionValueState {
		cl.FatalMsg(fmt.Sprintf(i18n.Text("Option %s requires an argument"), currentArg))
	}
	if cl.showHelp {
		cl.DisplayUsage()
		atexit.Exit(1)
	}
	if cl.showVersion || cl.showLongVersion {
		fmt.Fprintln(cl, NewVersionFromString(AppVersion).Format(false, cl.showLongVersion))
		atexit.Exit(0)
	}
	return remainingArgs
}

func (cl *CmdLine) setOrFail(op *Option, arg, value string) {
	if err := op.value.Set(value); err != nil {
		cl.FatalMsg(fmt.Sprintf(i18n.Text("Unable to set option %s to %s\n%v"), arg, value, err))
	}
}

// FatalMsg emits an error message and causes the program to exit.
func (cl *CmdLine) FatalMsg(msg string) {
	cl.out.Bell()
	cl.out.Foreground(term.Red, term.Normal)
	fmt.Fprint(cl, msg)
	cl.out.Reset()
	fmt.Fprintln(cl)
	atexit.Exit(1)
}

// FatalError emits an error message and causes the program to exit.
func (cl *CmdLine) FatalError(err error) {
	cl.FatalMsg(err.Error())
}

// FatalIfError emits an error message and causes the program to exit if
// err != nil.
func (cl *CmdLine) FatalIfError(err error) {
	if err != nil {
		cl.FatalError(err)
	}
}

func (cl *CmdLine) availableOptions() (available map[string]*Option) {
	available = make(map[string]*Option)
	for _, option := range cl.options {
		if ok, err := option.isValid(); !ok {
			cl.FatalMsg(fmt.Sprintf(i18n.Text("Invalid option specification: %v"), err))
		} else {
			if option.single != 0 {
				name := string(option.single)
				if available[name] != nil {
					cl.FatalMsg(fmt.Sprintf(i18n.Text("Option specification -%s already exists"), name))
				} else {
					available[name] = option
				}
			}
			if option.name != "" {
				if available[option.name] != nil {
					cl.FatalMsg(fmt.Sprintf(i18n.Text("Option specification --%s already exists"), option.name))
				} else {
					available[option.name] = option
				}
			}
		}
	}
	return available
}

func (cl *CmdLine) loadArgsFromFile(path string) (args []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errs.NewWithCause(fmt.Sprintf(i18n.Text("Unable to open: %s"), path), err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()
	args = make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		args = append(args, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		return nil, errs.Wrap(err)
	}
	return args, nil
}

// SetWriter sets the io.Writer to use for output. By default, a new CmdLine
// uses os.Stderr.
func (cl *CmdLine) SetWriter(w io.Writer) {
	cl.out = term.NewANSI(w)
}

// Write implements the io.Writer interface.
func (cl *CmdLine) Write(p []byte) (n int, err error) {
	return cl.out.Write(p)
}
