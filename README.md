# Mattermost BlueJeans Plugin

Start and join voice calls, video calls and use screen sharing with your team members via BlueJeans.

Usage
-----

Once enabled, clicking the video icon in a Mattermost channel invites team members to join a BlueJeans call, hosted using the credentials of the user who initiated the call.

![Screenshot](https://user-images.githubusercontent.com/177788/42196048-af54d2b8-7e30-11e8-80a0-5e160ae06f03.png)

BlueJeans Setup Guide
-----

You will need a BlueJeans Pro or Enterprise account to use the plugin.

1. Go to **System Console > Plugins > BlueJeans** to configure the BlueJeans Plugin.

2. To generate an **API Key** and **API Secret** requires a [Pro or Enterprise BlueJeans plan](https://store.bluejeans.com/). The current BlueJeans process to get OAuth Access in the BlueJeans Admin Console is to contact BlueJeans Support directly

3. Enable settings for [overriding usernames](https://docs.mattermost.com/administration/config-settings.html#enable-integrations-to-override-usernames) and [overriding profile picture icons](https://docs.mattermost.com/administration/config-settings.html#enable-integrations-to-override-profile-picture-icons).

4. Activate the plugin at **System Console > Plugins > Management** by clicking **Activate** for BlueJeans.

Once activated, you will see a video icon in the channel header. Clicking the icon will create a new BlueJeans meeting, and create a post with a link to the meeting. Anyone in the channel can see the post and can join by clicking on the link.

Note
----
   Users will need to sign-up for their own BlueJeans account using the same email address that they use for Mattermost. If the user attempts to start a BlueJeans meeting without a BlueJeans account, they will see the following error message: "We could not verify your Mattermost account in BlueJeans. Please ensure that your Mattermost email address matches your BlueJeans email address."


## Development

This plugin contains both a server and web app portion.

Use `make dist` to build distributions of the plugin that you can upload to a Mattermost server for testing.

Use `make check-style` to check the style for the whole plugin.

### Server

Inside the `/server` directory, you will find the Go files that make up the server-side of the plugin. Within there, build the plugin like you would any other Go application.

### Web App

Inside the `/webapp` directory, you will find the JS and React files that make up the client-side of the plugin. Within there, modify files and components as necessary. Test your syntax by running `npm run build`.
