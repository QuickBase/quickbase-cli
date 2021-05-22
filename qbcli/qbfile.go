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

type QuickbaseFile struct {
	Deploy *QuickbaseFileDeploy `yaml:"deploy"`
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

	b, err := ioutil.ReadFile(file)
	if err != nil {
		err = fmt.Errorf("error reading quickbase file: %w", err)
		return
	}

	dec := yaml.NewDecoder(bytes.NewBuffer(b))
	dec.KnownFields(true)
	err = dec.Decode(f)

	if verr := validator.New().Struct(f); err != nil {
		err = qberrors.HandleErrorValidation(verr)
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
			err = fmt.Errorf("error reading formula file: %w", ferr)
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
