package pwsh

import (
	"fmt"
	"os"

	"get.porter.sh/porter/pkg/exec/builder"
)

var _ builder.ExecutableAction = Action{}
var _ builder.BuildableAction = Action{}

type Action struct {
	Name  string
	Steps []Step // using UnmarshalYAML so that we don't need a custom type per action
}

// MarshalYAML converts the action back to a YAML representation
// install:
//
//	pwsh:
//	  ...
func (a Action) MarshalYAML() (interface{}, error) {
	return map[string]interface{}{a.Name: a.Steps}, nil
}

// MakeSteps builds a slice of Step for data to be unmarshaled into.
func (a Action) MakeSteps() interface{} {
	return &[]Step{}
}

// UnmarshalYAML takes any yaml in this form
// ACTION:
// - pwsh: ...
// and puts the steps into the Action.Steps field
func (a *Action) UnmarshalYAML(unmarshal func(interface{}) error) error {
	results, err := builder.UnmarshalAction(unmarshal, a)
	if err != nil {
		return err
	}

	for actionName, action := range results {
		a.Name = actionName
		for _, result := range action {
			step := result.(*[]Step)
			a.Steps = append(a.Steps, *step...)
		}
		break // There is only 1 action
	}
	return nil
}

func (a Action) GetSteps() []builder.ExecutableStep {
	// Go doesn't have generics, nothing to see here...
	steps := make([]builder.ExecutableStep, len(a.Steps))
	for i := range a.Steps {
		steps[i] = a.Steps[i]
	}

	return steps
}

type Step struct {
	Instruction `yaml:"pwsh"`
}

// Actions is a set of actions, and the steps, passed from Porter.
type Actions []Action

// UnmarshalYAML takes chunks of a porter.yaml file associated with this mixin
// and populates it on the current action set.
// install:
//
//	pwsh:
//	  ...
//	pwsh:
//	  ...
//
// upgrade:
//
//	pwsh:
//	  ...
func (a *Actions) UnmarshalYAML(unmarshal func(interface{}) error) error {
	results, err := builder.UnmarshalAction(unmarshal, Action{})
	if err != nil {
		return err
	}

	for actionName, action := range results {
		for _, result := range action {
			s := result.(*[]Step)
			*a = append(*a, Action{
				Name:  actionName,
				Steps: *s,
			})
		}
	}
	return nil
}

var _ builder.HasOrderedArguments = Instruction{}
var _ builder.ExecutableStep = Instruction{}
var _ builder.StepWithOutputs = Instruction{}

type Instruction struct {
	Name           string   `yaml:"name"`
	Description    string   `yaml:"description"`
	WorkingDir     string   `yaml:"workingDir,omitempty"`
	InlineScript   string   `yaml:"inlineScript,omitempty"`
	File           string   `yaml:"file,omitempty"`
	Arguments      []string `yaml:"arguments,omitempty"`
	Outputs        []Output `yaml:"outputs,omitempty"`
	SuppressOutput bool     `yaml:"suppress-output,omitempty"`

	builder.IgnoreErrorHandler `yaml:"ignoreError,omitempty"`
}

func (s Instruction) GetCommand() string {
	return "pwsh"
}

func (s Instruction) GetWorkingDir() string {
	return s.WorkingDir
}

func (s Instruction) GetArguments() []string {
	return s.getArguments()
}

func (s Instruction) GetSuffixArguments() []string {
	return nil
}

func (s Instruction) GetFlags() builder.Flags {
	return nil
}

func (s Instruction) SuppressesOutput() bool {
	return s.SuppressOutput
}

func (s *Instruction) getArguments() []string {
	args := make([]string, 0, len(s.Arguments)+2)
	args = append(args, "-NonInteractive")
	if s.InlineScript != "" {
		args = append(args, "-Command", s.InlineScript)
	} else if s.File != "" {
		args = append(args, "-File", s.File)
	} else {
		fmt.Fprintf(os.Stderr, "inlineScript or file needs to be specified")
		os.Exit(1)
	}

	if len(s.Arguments) > 0 {
		args = append(args, s.Arguments...)
	}

	return args
}

func (s Instruction) GetOutputs() []builder.Output {
	// Go doesn't have generics, nothing to see here...
	outputs := make([]builder.Output, len(s.Outputs))
	for i := range s.Outputs {
		outputs[i] = s.Outputs[i]
	}
	return outputs
}

var _ builder.OutputJsonPath = Output{}
var _ builder.OutputFile = Output{}
var _ builder.OutputRegex = Output{}

type Output struct {
	Name string `yaml:"name"`

	JsonPath string `yaml:"jsonPath,omitempty"`
	FilePath string `yaml:"path,omitempty"`
	Regex    string `yaml:"regex,omitempty"`
}

func (o Output) GetName() string {
	return o.Name
}

func (o Output) GetJsonPath() string {
	return o.JsonPath
}

func (o Output) GetFilePath() string {
	return o.FilePath
}

func (o Output) GetRegex() string {
	return o.Regex
}
