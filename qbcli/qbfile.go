package qbcli

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/QuickBase/quickbase-cli/qbclient"
	"github.com/QuickBase/quickbase-cli/qberrors"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

// QuickbaseFile models a quickbase.yml file.
type QuickbaseFile struct {
	Test   *QuickbaseFileTest   `yaml:"test"`
	Deploy *QuickbaseFileDeploy `yaml:"deploy"`
}

type QuickbaseFileTest struct {
	Formulas []*QuickbaseFileTestFormula `yaml:"formulas"`
}

type QuickbaseFileTestFormula struct {
	File     string `validate:"required" yaml:"file"`
	TableID  string `validate:"required" yaml:"table_id"`
	RecordID int    `validate:"required" yaml:"record_id"`
	Expected string `validate:"required" yaml:"expected"`
}

type QuickbaseFileDeploy struct {
	Formulas []*QuickbaseFileDeployFormula `yaml:"formulas"`
}

type QuickbaseFileDeployFormula struct {
	File    string `validate:"required" yaml:"file"`
	TableID string `validate:"required" yaml:"table_id"`
	FieldID int    `validate:"required" yaml:"field_id"`
}

func ParseQuickbaseFile(file string) (f *QuickbaseFile, err error) {
	f = &QuickbaseFile{}

	// Read the quickbase.yml file.
	b, err := ioutil.ReadFile(file)
	if err != nil {
		err = qberrors.Client(nil).Safef(qberrors.InvalidInput, "error reading quickbase file: %w", err)
		return
	}

	// Parse the yaml file using strict settings.
	dec := yaml.NewDecoder(bytes.NewBuffer(b))
	dec.KnownFields(true)
	derr := dec.Decode(f)
	if derr != nil {
		err = qberrors.Client(nil).Safef(qberrors.InvalidSyntax, "yaml not valid: %w", derr)
	}

	// Validate the decoded file.
	if verr := validator.New().Struct(f); verr != nil {
		err = qberrors.HandleErrorValidation(verr)
	}

	return
}

type TestFormulaInput struct {
	File string `cliutil:"option=file default=quickbase.yml"`
}

type TestFormulaOutput struct {
	Passed []int          `json:"passed"`
	Failed map[int]string `json:"failed"`
}

func TestFormula(qb *qbclient.Client, in *TestFormulaInput) (out *TestFormulaOutput, err error) {
	out = &TestFormulaOutput{
		Passed: []int{},
		Failed: map[int]string{},
	}

	var file *QuickbaseFile
	file, err = ParseQuickbaseFile(in.File)
	if err != nil {
		return
	}

	for idx, f := range file.Test.Formulas {

		b, ferr := ioutil.ReadFile(f.File)
		if ferr != nil {
			err = qberrors.Client(nil).Safef(qberrors.InvalidInput, "error reading formula file: %w", ferr)
			return
		}

		rfi := &qbclient.RunFormulaInput{
			From:     f.TableID,
			RecordID: f.RecordID,
			Formula:  string(b),
		}

		rfo, rerr := qb.RunFormula(rfi)
		if rerr == nil {
			if f.Expected == rfo.Result {
				out.Passed = append(out.Passed, idx)
			} else {
				out.Failed[idx] = fmt.Sprintf("expected %q, got %q", f.Expected, rfo.Result)
			}
		} else {
			out.Failed[idx] = rerr.Error()
		}
	}

	return
}

type DeployFormulaInput struct {
	File string `cliutil:"option=file default=quickbase.yml"`
}

type DeployFormulaOutput struct {
	Deployed map[string][]int          `json:"deployed"`
	Errors   map[string]map[int]string `json:"errors"`
}

func DeployFormula(qb *qbclient.Client, in *DeployFormulaInput) (out *DeployFormulaOutput, err error) {
	out = &DeployFormulaOutput{
		Deployed: map[string][]int{},
		Errors:   map[string]map[int]string{},
	}

	var file *QuickbaseFile
	file, err = ParseQuickbaseFile(in.File)
	if err != nil {
		return
	}

	for _, f := range file.Deploy.Formulas {

		b, ferr := ioutil.ReadFile(f.File)
		if ferr != nil {
			err = qberrors.Client(nil).Safef(qberrors.InvalidInput, "error reading formula file: %w", ferr)
			return
		}

		ufi := &qbclient.UpdateFieldInput{
			TableID:    f.TableID,
			FieldID:    f.FieldID,
			Properties: &qbclient.UpdateFieldInputProperties{},
		}
		ufi.Properties.Formula = string(b)

		// TODO See https://github.com/QuickBase/quickbase-cli/issues/36
		ufi.AddToNewReports = true
		ufi.Searchable = true

		_, uerr := qb.UpdateField(ufi)
		if uerr == nil {
			if _, ok := out.Deployed[f.TableID]; !ok {
				out.Deployed[f.TableID] = []int{}
			}
			out.Deployed[f.TableID] = append(out.Deployed[f.TableID], f.FieldID)
		} else {
			if _, ok := out.Errors[f.TableID]; !ok {
				out.Errors[f.TableID] = make(map[int]string)
			}
			out.Errors[f.TableID][f.FieldID] = uerr.Error()
		}
	}

	return
}
