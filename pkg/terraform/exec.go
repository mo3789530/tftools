package terraform

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
)

const (
	outPaht = "/tmp/terraform.out"
)

type Exec interface {
	Init(workingDir string) (*tfexec.Terraform, error)
	Plan(tf *tfexec.Terraform) (bool, error)
	Show(tf *tfexec.Terraform, isJson bool) (string, error)
	Apply(tf *tfexec.Terraform) error
	State(tf *tfexec.Terraform, isRemote bool) (string, error)
}

type exec struct {
	Version  string
	IsOutput bool
	IsJson   bool
}

func NewExec(version string, isOutput bool, isJson bool) Exec {
	return &exec{
		Version:  version,
		IsOutput: isOutput,
		IsJson:   isJson,
	}
}

func (e *exec) Init(workingDir string) (*tfexec.Terraform, error) {
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion(e.Version)),
	}
	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Printf("error installing terraform: %s", err)
	}

	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		log.Printf("error new terraform: %s", err)
	}

	return tf, nil
}

func (e *exec) Plan(tf *tfexec.Terraform) (bool, error) {
	err := tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Printf("error init terraform: %s", err)
		return false, err
	}

	var planOptions []tfexec.PlanOption
	if e.IsOutput {
		planOptions = []tfexec.PlanOption{
			tfexec.Out(outPaht),
		}
	}

	result, err := tf.Plan(context.Background(), planOptions...)
	if err != nil {
		log.Printf("error terraform plan: %s", err)
		return result, err
	}
	return result, nil
}

func (e *exec) Show(tf *tfexec.Terraform) (string, error) {
	if e.IsJson {
		result, err := e.showJson(tf)
		if err != err {
			return "", err
		}
		return result, nil
	}

	result, err := e.showPlanText(tf)
	if err != err {
		return "", err
	}
	return result, nil
}

func (e *exec) showPlanText(tf *tfexec.Terraform) (string, error) {
	// log.Print(e.IsOutput)
	if !e.IsOutput {
		_, err := tf.Show(context.Background())
		if err != nil {
			log.Printf("error terrform show: %s", err)
			return "", err
		}

		return "", nil
	} else {
		show, err := tf.ShowPlanFileRaw(context.Background(), outPaht)
		if err != nil {
			log.Printf("error terraform show: %s", err)
			return "", err
		}
		return show, nil
	}
}

func (e *exec) showJson(tf *tfexec.Terraform) (string, error) {
	if !e.IsOutput {
		show, err := tf.Show(context.Background())
		if err != nil {
			log.Printf("error terrform show: %s", err)
			return "", err
		}
		jsonData, err := json.Marshal(show)
		if err != nil {
			log.Printf("error convert json %s", err)
			return "", err
		}
		return string(jsonData), nil

	} else {
		show, err := tf.ShowPlanFile(context.Background(), outPaht)
		if err != nil {
			log.Printf("error terraform show: %s", err)
		}
		jsonData, err := json.Marshal(show)
		if err != nil {
			log.Printf("error convert json %s", err)
			return "", err
		}
		return string(jsonData), nil
	}
}

func (e *exec) Apply(tf *tfexec.Terraform) error {
	err := tf.Apply(context.Background())

	if err != nil {
		log.Printf("error terraform apply: %s", err)
		return err
	}
	return nil
}

func (e *exec) State(tf *tfexec.Terraform, isRemote bool) (string, error) {
	log.Println("terraform state")
	if isRemote {
		var pullOptions tfexec.StatePullOption
		_, err := tf.StatePull(context.Background(), pullOptions)
		if err != nil {
			log.Printf("error terraform pull %s", err)
			return "", err
		}
	}

	state, err := tf.ShowStateFile(context.Background(), "terraform.tfstate")
	if err != nil {
		log.Printf("error terraform state file %s", err)
		return "", err
	}

	if e.IsJson {
		return ShowStateFileJson(state), nil
	} else {
		return ShowStateFileRaw(state), nil
	}

	return "", nil

}
