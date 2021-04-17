# tg-webhook-prom

Export basic information from Telegram's getWebhookInfo Bot API endpoint.

It provides a `telegram_pending_updates` and `telegram_last_error_date` gauge.

## Configuration

The `TELEGRAM_APITOKEN` environment variable must be set to your bot's token.
Data is refreshed every minute.
