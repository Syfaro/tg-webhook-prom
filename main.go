package main

import (
	"net/http"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	pendingUpdates = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "telegram_pending_updates",
		Help: "The number of updates pending for the webhook",
	})

	lastErrorDate = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "telegram_last_error_date",
		Help: "The unix timestamp of the last webhook error",
	})
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Info("Starting exporter")

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		updateMetrics(bot)
		for range time.Tick(60 * time.Second) {
			updateMetrics(bot)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

func updateMetrics(bot *tgbotapi.BotAPI) {
	log.Trace("Updating metrics")

	webhookInfo, err := bot.GetWebhookInfo()
	if err != nil {
		log.Error(err)
		return
	}

	log.Debugf("%+v", webhookInfo)

	pendingUpdates.Set(float64(webhookInfo.PendingUpdateCount))
	lastErrorDate.Set(float64(webhookInfo.LastErrorDate))

	log.Info("Update complete")
}
