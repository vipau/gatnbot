# Gatnbot

Ver 0.0.1 turbo alpha

## Files needed for deployment
- client_secret.json
    - Credentials to your Oauth service. To be placed in the binary folder.
- gatnbot
    - compiled binary
- settings.hcl
    - rename settings.hcl.example to settings.hcl and place it in the binary folder.

## Google Authentication

- Get a GCP account and enable the Gmail API on your project.
- Create a new Oauth client with the scope to read Gmail messages.
- Follow the Go Quickstart seen (here)[https://developers.google.com/drive/api/v3/quickstart/go] to get the oauth json.
- Rename it to client_secret.json and put it in the folder where the Go binary is.
- On first run, the Go program will use Oauth to authenticate your Google account. You need to check the "read Gmail emails" permission manually.
- Once you paste the oauth code back in the terminal window, you are now authenticated to Gmail.
- Your credentials will be in ~/.credentials/gmail-go-quickstart.json

## Telegram authentication

- Get a Bot ID and token from @BotFather on Telegram and insert them into settings.hcl "bottoken"

## Other options
- timezone
    - Set to a string with ISO timezone such as "Europe/Rome".
  Used for task scheduling.

- apiurl
    - Custom Telegram API URL, if empty defaults to 'https://api.telegram.org'

- chatid
    - An array containing list of chats to send emails to, and also chats where the bot will reply to commands.
    NOTE: Please only use one chat ID for now.

- adminid
    - An array contining chat IDs of admins. The bot will treat these as additional chats to handle messages in, but it will not forward emails.

- ouremail
    - Our Gmail email, used for mail querying. we don't want to display mails not addressed directly to our address.

- linksmsg
    - The message containing the response of the /links command in Markdown format.

## How do I do \<thing\>?

This is still a huge WIP! Don't expect to kill dragons with this code, but issues and PRs are accepted.  
As for now, here's what to check if you're curious:  

- `settings` loads the HCL config file into memory
- `crontasks` handles the scheduled actions, including refreshing the Markov model
- `commands` handles the /commands that the bot replies to.
- `sendemail` handles checking for unread email, forwarding it, and marking it as read
- `fakernews-mod` modified standalone version of [fakernews](https://github.com/mb-14/gomarkov/blob/master/examples/fakernews/fakernews.go) generates fake Hacker News stories
