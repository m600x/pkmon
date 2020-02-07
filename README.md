# pkmon

Send Slack message to a channel in random interval for a random time at
specific hours avoiding spamming at 3AM...

Default settings:
- Between 9AM to 7PM
- Rate between 15 and 30 minutes
- TTL between 30 and 60 seconds

You need to input your own auth token (`xoxp-...`) at line 39 from https://api.slack.com/legacy/custom-integrations/legacy-tokens

## Instructions

```
go install pkmon/pkmon
cp pkmon.service /etc/systemd/system/pkmon.service
systemctl enable pkmon
systemctl start pkmon
```