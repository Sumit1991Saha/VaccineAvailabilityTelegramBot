# VaccineAvailabilityTelegramBot

There are 2 ways to use this project:-
1. Using continuous polling from Telegram
2. Using webhook

comment either of the option in main.go
Also generate a bot token using bot father in telegram, since this service primarily interact with telegram

If using option 2 ie using webhook,
we would need a 3rd party service which can forward public requests to our local server, for that this project is using ngrok. 
Its a very light weight service.
Download it and start it using "./ngrok http 3000" ie ngrok forwards the incoming requests to localhost:3000 since we are starting our local http server at 3000

In the file webhook.go, there is a method SetWebhook(), which fetches the tunnels created by ngrok and registers a  webhook with telegram using your bot token.
