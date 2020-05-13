package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const (
	COMMAND_HELP = `* |/bluejeans start| - Start a BlueJeans meeting.`
)

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          "bluejeans",
		DisplayName:      "BlueJeans",
		Description:      "Integration with BlueJeans.",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: start",
		AutoCompleteHint: "[command]",
	}
}

func (p *Plugin) postCommandResponse(args *model.CommandArgs, text string) {
	post := &model.Post{
		UserId:    p.botUserID,
		ChannelId: args.ChannelId,
		Message:   text,
	}
	_ = p.API.SendEphemeralPost(args.UserId, post)
}

func (p *Plugin) executeCommand(c *plugin.Context, args *model.CommandArgs) (string, error) {

	split := strings.Fields(args.Command)
	command := split[0]
	action := ""

	if command != "/bluejeans" {
		return fmt.Sprintf("Command '%s' is not /bluejeans. Please try again.", command), nil
	}

	if len(split) > 1 {
		action = split[1]
	} else {
		return "Please specify an action for /bluejeans command.", nil
	}

	userID := args.UserId
	user, appErr := p.API.GetUser(userID)
	if appErr != nil {
		return fmt.Sprintf("We could not retrieve user (userId: %v)", args.UserId), nil
	}

	if action == "start" {
		if _, appErr = p.API.GetChannelMember(args.ChannelId, userID); appErr != nil {
			return fmt.Sprintf("We could not get channel members (channelId: %v)", args.ChannelId), nil
		}

		recentMeeting, recentMeetingID, creatorName, appErr := p.checkPreviousMessages(args.ChannelId)
		if appErr != nil {
			return fmt.Sprintf("Error checking previous messages"), nil
		}

		if recentMeeting {
			p.postConfirm(recentMeetingID, args.ChannelId, "", userID, creatorName)
			return "", nil
		}

		// create a personal bluejeans meeting
		rpm, clientErr := p.bluejeansClient.GetPersonalMeeting(user.Email)
		if clientErr != nil {
			return "We could not verify your Mattermost account in BlueJeans. Please ensure that your Mattermost email address matches your BlueJeans login email address.", nil
		}
		meetingID := rpm.NumericMeetingID

		_, appErr = p.postMeeting(user.Username, meetingID, args.ChannelId, "")
		if appErr != nil {
			return "Failed to post message. Please try again.", nil
		}
		return "", nil
	}
	return fmt.Sprintf("Unknown action %v", action), nil
}

// ExecuteCommand method
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	msg, err := p.executeCommand(c, args)
	if err != nil {
		p.API.LogWarn("failed to execute command", "error", err.Error())
	}
	if msg != "" {
		p.postCommandResponse(args, msg)
	}
	return &model.CommandResponse{}, nil
}
