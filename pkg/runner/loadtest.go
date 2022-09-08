package runner

import (
	"errors"
	"fmt"
	"io"

	"github.com/u-cto-devops/lguctl/pkg/color"
	"github.com/u-cto-devops/lguctl/pkg/constants"
	"github.com/u-cto-devops/lguctl/pkg/tools"
)

func printStatusWithColor(out io.Writer, status string) {
	switch status {
	case "available":
		color.Green.Fprintf(out, "[ available ]")
	case "stopped":
		color.Red.Fprintf(out, "[ stopped ]")
	default:
		color.Blue.Fprintf(out, "[ %s ]", status)
	}
}

func printCapacityWithColor(out io.Writer, capacity int64) {
	switch capacity {
	case 0:
		color.Red.Fprintf(out, "[ 0 ]")
	default:
		color.Green.Fprintf(out, "[ %d ]", capacity)
	}
}

func (r Runner) GetLoadtestStatus(out io.Writer) error {
	if r.Config == nil {
		return errors.New(constants.ConfigErrorMsg)
	}

	for _, rdsClusterName := range r.Config.Loadtest.RDS {
		status, err := r.AWSClient.GetLoadtestRDSStatus(rdsClusterName)
		fmt.Fprintf(out, "RDS - %s\t", rdsClusterName)
		if err != nil {
			color.Red.Fprintf(out, "[ %v ]", err.Error())
		} else {
			printStatusWithColor(out, *status)
		}
	}

	for _, asgName := range r.Config.Loadtest.ASG {
		exactASGNames, _ := r.AWSClient.GetExactASGNames(asgName)
		for _, exactASGName := range exactASGNames {
			desiredCapacity, err := r.AWSClient.GetLoadtestASGStatus(exactASGName)
			fmt.Fprintf(out, "ASG - %s\t", *exactASGName)
			if err != nil {
				color.Red.Fprintf(out, "[ %v ]", err.Error())
			} else {
				printCapacityWithColor(out, *desiredCapacity)
			}
		}
	}

	return nil
}

func (r Runner) StartLoadtest(out io.Writer) error {
	if r.Config == nil {
		return errors.New(constants.ConfigErrorMsg)
	}

	for _, rdsClusterName := range r.Config.Loadtest.RDS {
		if err := tools.AskContinue("Are you sure to start RDS " + rdsClusterName); err != nil {
			color.Red.Fprintf(out, "Starting RDS %s has been cancelled", rdsClusterName)
			continue
		}
		_, err := r.AWSClient.StartLoadtestRDS(rdsClusterName)
		fmt.Fprintf(out, "RDS - %s\t", rdsClusterName)
		if err != nil {
			color.Red.Fprintf(out, "[ %v ]", err.Error())
		} else {
			color.Blue.Fprintf(out, "[ starting ]")
		}
	}

	for _, asgName := range r.Config.Loadtest.ASG {
		exactASGNames, _ := r.AWSClient.GetExactASGNames(asgName)
		for _, exactASGName := range exactASGNames {
			if err := tools.AskContinue("Are you sure to update ASG " + *exactASGName); err != nil {
				color.Red.Fprintf(out, "Updating ASG %s has been cancelled", *exactASGName)
				continue
			}
			err := r.AWSClient.StartLoadtestASG(exactASGName)
			fmt.Fprintf(out, "ASG - %s\t", *exactASGName)
			if err != nil {
				color.Red.Fprintf(out, "[ %v ]", err.Error())
			} else {
				color.Green.Fprintf(out, "[ 1 ]")
			}
		}
	}

	return nil
}

func (r Runner) StopLoadtest(out io.Writer) error {
	if r.Config == nil {
		return errors.New(constants.ConfigErrorMsg)
	}

	for _, rdsClusterName := range r.Config.Loadtest.RDS {
		if err := tools.AskContinue("Are you sure to stop RDS " + rdsClusterName); err != nil {
			color.Red.Fprintf(out, "Stopping RDS %s has been cancelled", rdsClusterName)
			continue
		}
		_, err := r.AWSClient.StopLoadtestRDS(rdsClusterName)
		fmt.Fprintf(out, "RDS - %s\t", rdsClusterName)
		if err != nil {
			color.Red.Fprintf(out, "[ %v ]", err.Error())
		} else {
			color.Blue.Fprintf(out, "[ stopping ]")
		}
	}

	for _, asgName := range r.Config.Loadtest.ASG {
		exactASGNames, _ := r.AWSClient.GetExactASGNames(asgName)
		for _, exactASGName := range exactASGNames {
			if err := tools.AskContinue("Are you sure to update ASG " + *exactASGName); err != nil {
				color.Red.Fprintf(out, "Updating ASG %s has been cancelled", *exactASGName)
				continue
			}
			err := r.AWSClient.StopLoadtestASG(exactASGName)
			fmt.Fprintf(out, "ASG - %s\t", *exactASGName)
			if err != nil {
				color.Red.Fprintf(out, "[ %v ]", err.Error())
			} else {
				color.Red.Fprintf(out, "[ 0 ]")
			}
		}
	}

	return nil
}
