package v6

import (
	"code.cloudfoundry.org/cli/actor/sharedaction"
	"code.cloudfoundry.org/cli/actor/v3action"
	"code.cloudfoundry.org/cli/command"
	"code.cloudfoundry.org/cli/command/flag"
	"code.cloudfoundry.org/cli/command/v6/shared"
)

//go:generate counterfeiter . EnableOrgIsolationActor

type EnableOrgIsolationActor interface {
	EntitleIsolationSegmentToOrganizationByName(isolationSegmentName string, orgName string) (v3action.Warnings, error)
}

type EnableOrgIsolationCommand struct {
	RequiredArgs    flag.OrgIsolationArgs `positional-args:"yes"`
	usage           interface{}           `usage:"CF_NAME enable-org-isolation ORG_NAME SEGMENT_NAME"`
	relatedCommands interface{}           `related_commands:"create-isolation-segment, isolation-segments, set-org-default-isolation-segment, set-space-isolation-segment"`

	UI          command.UI
	Config      command.Config
	SharedActor command.SharedActor
	Actor       EnableOrgIsolationActor
}

func (cmd *EnableOrgIsolationCommand) Setup(config command.Config, ui command.UI) error {
	cmd.UI = ui
	cmd.Config = config
	cmd.SharedActor = sharedaction.NewActor(config)

	client, _, err := shared.NewV3BasedClients(config, ui, true, "")
	if err != nil {
		return err
	}
	cmd.Actor = v3action.NewActor(client, config, nil, nil)

	return nil
}

func (cmd EnableOrgIsolationCommand) Execute(args []string) error {
	err := cmd.SharedActor.CheckTarget(false, false)
	if err != nil {
		return err
	}

	user, err := cmd.Config.CurrentUser()
	if err != nil {
		return err
	}

	cmd.UI.DisplayTextWithFlavor("Enabling isolation segment {{.SegmentName}} for org {{.OrgName}} as {{.CurrentUser}}...", map[string]interface{}{
		"SegmentName": cmd.RequiredArgs.IsolationSegmentName,
		"OrgName":     cmd.RequiredArgs.OrganizationName,
		"CurrentUser": user.Name,
	})

	warnings, err := cmd.Actor.EntitleIsolationSegmentToOrganizationByName(cmd.RequiredArgs.IsolationSegmentName, cmd.RequiredArgs.OrganizationName)
	cmd.UI.DisplayWarnings(warnings)
	if err != nil {
		return err
	}

	cmd.UI.DisplayOK()

	return nil
}
