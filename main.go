package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1" // imports as package "cli"
)

const (
	// VERSION represents the version of the generator tool
	VERSION = "1.2.0"

	// TFModulesVersion represents the version of the Base TF modules
	TFModulesVersion = "1.2.1"

	tfDoNotEditStamp = `// DO NOT EDIT
// This file has been auto-generated by taskhawk-terraform-generator ` + VERSION
)

const (
	// alertingFlag represents the cli flag that indicates if alerting should be generated
	alertingFlag = "alerting"

	// dlqAlertAlarmActionsFlag represents the cli flag for DLQ alert actions on ALARM
	dlqAlertAlarmActionsFlag = "dlq-alert-alarm-actions"

	// dlqAlertOKActionsFlag represents the cli flag for DLQ alert actions on OK
	dlqAlertOKActionsFlag = "dlq-alert-ok-actions"

	// iamFlag represents the cli flag for iam generation
	iamFlag = "iam"

	// moduleFlag represents the cli flag for output module name
	moduleFlag = "module"

	// queueAlertAlarmActionsFlag represents the cli flag for DLQ alert actions on ALARM
	queueAlertAlarmActionsFlag = "queue-alert-alarm-actions"

	// queueAlertOKActionsFlag represents the cli flag for DLQ alert actions on OK
	queueAlertOKActionsFlag = "queue-alert-ok-actions"
)

func validateArgs(c *cli.Context) *cli.ExitError {
	if c.NArg() != 1 {
		return cli.NewExitError("<config-file> is required", 1)
	}

	alertingFlagsOkay := true
	alertingFlags := []string{queueAlertAlarmActionsFlag, queueAlertOKActionsFlag, dlqAlertAlarmActionsFlag,
		dlqAlertOKActionsFlag}
	if c.Bool(alertingFlag) {
		for _, f := range alertingFlags {
			if !c.IsSet(f) {
				alertingFlagsOkay = false
				msg := fmt.Sprintf("--%s is required", f)
				if _, err := fmt.Fprint(cli.ErrWriter, msg); err != nil {
					return cli.NewExitError(msg, 1)
				}
			}
		}
		if !alertingFlagsOkay {
			return cli.NewExitError("missing required flags for --alerting", 1)
		}
	} else {
		for _, f := range alertingFlags {
			if c.IsSet(f) {
				alertingFlagsOkay = false
				msg := fmt.Sprintf("--%s is disallowed", f)
				if _, err := fmt.Fprint(cli.ErrWriter, msg); err != nil {
					return cli.NewExitError(msg, 1)
				}
			}
		}
		if !alertingFlagsOkay {
			return cli.NewExitError("disallowed flags specified with missing --alerting", 1)
		}
	}
	return nil
}

func generateModule(c *cli.Context) error {
	if err := validateArgs(c); err != nil {
		return err
	}

	configFile := c.Args().Get(0)

	config, err := NewConfig(configFile)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	err = writeTerraform(config, c)
	if err != nil {
		return cli.NewExitError(errors.Wrap(err, "failed to create terraform module"), 1)
	}

	fmt.Println("Created Terraform Taskhawk module successfully!")
	return nil
}

func generateConfigFileStructure(c *cli.Context) error {
	structure := Config{
		QueueApps: []*QueueApp{
			{
				Queue: "DEV-MYAPP",
				Tags: map[string]string{
					"App": "myapp",
					"Env": "dev",
				},
				Schedule: []ScheduleItem{
					{
						Name:          "nightly-job (unique for each app)",
						Description:   "{optional description}",
						FormatVersion: "{optional format version}",
						Headers: map[string]string{
							"header": "{optional headers}",
						},
						Task: "tasks.send_email",
						Args: []interface{}{"{optional args}"},
						Kwargs: map[string]interface{}{
							"kwarg1": "{optional keyword args}",
						},
					},
				},
			},
		},
		LambdaApps: []*LambdaApp{
			{
				FunctionARN:  "arn:aws:lambda:us-west-2:12345:function:my_function:deployed",
				FunctionName: "{optional - this value is inferred from FunctionARN if that's not an interpolated value}",
				FunctionQualifier: "{optional - this value is inferred from FunctionARN if that's not an interpolated" +
					" value}",
				Name: "myapp",
				Schedule: []ScheduleItem{
					{
						Name:          "nightly-job (unique for each app)",
						Description:   "{optional description}",
						FormatVersion: "{optional format version}",
						Headers: map[string]string{
							"header": "{optional headers}",
						},
						Task: "tasks.send_email",
						Args: []interface{}{"{optional args}"},
						Kwargs: map[string]interface{}{
							"kwarg1": "{optional keyword args}",
						},
					},
				},
			},
		},
	}
	structureAsJSON, err := json.MarshalIndent(structure, "", "    ")
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	fmt.Println(string(structureAsJSON))
	return nil
}

func runApp(args []string) error {
	cli.VersionFlag = cli.BoolFlag{Name: "version, V"}

	app := cli.NewApp()
	app.Name = "TaskHawk Terraform"
	app.Usage = "Manage Terraform configuration for TaskHawk apps"
	app.Version = VERSION
	app.HelpName = "taskhawk-terraform"
	app.Commands = []cli.Command{
		{
			Name:      "generate",
			Usage:     "Generates Terraform module for TaskHawk apps",
			ArgsUsage: "<config-file>",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  moduleFlag,
					Usage: "Terraform module name to generate",
					Value: "taskhawk",
				},
				cli.BoolFlag{
					Name:  alertingFlag,
					Usage: "Should Cloudwatch alerting be generated?",
				},
				cli.BoolFlag{
					Name:  iamFlag,
					Usage: "Should IAM policies be generated?",
				},
				cli.StringSliceFlag{
					Name:  queueAlertAlarmActionsFlag,
					Usage: "Cloudwatch Action ARNs for high message count in queue when in ALARM",
				},
				cli.StringSliceFlag{
					Name:  queueAlertOKActionsFlag,
					Usage: "Cloudwatch Action ARNs for high message count in queue when OK",
				},
				cli.StringSliceFlag{
					Name:  dlqAlertAlarmActionsFlag,
					Usage: "Cloudwatch Action ARNs for high message count in dead-letter queue when in ALARM",
				},
				cli.StringSliceFlag{
					Name:  dlqAlertOKActionsFlag,
					Usage: "Cloudwatch Action ARNs for high message count in dead-letter queue when OK",
				},
			},
			Action: generateModule,
		},
		{
			Name:   "config-file-structure",
			Usage:  "Outputs the structure for config file required for generate command",
			Action: generateConfigFileStructure,
		},
	}

	return app.Run(args)
}

func main() {
	if err := runApp(os.Args); err != nil {
		log.Fatal(err)
	}
}