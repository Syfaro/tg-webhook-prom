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
	pendingUpdates = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "telegram_pending_updates",
		Help: "The number of updates pending for the webhook",
	}, []string{"username"})

	lastErrorDate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "telegram_last_error_date",
		Help: "The unix timestamp of the last webhook error",
	}, []string{"username"})
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.Info("starting exporter")

	token := os.Getenv("TELEGRAM_API_TOKEN")
	if len(token) == 0 {
		log.Fatal("missing TELEGRAM_API_TOKEN")
	}

	bot, err := tgbotapi.NewBotAPI(token)
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
	log.Trace("updating metrics")

	webhookInfo, err := bot.GetWebhookInfo()
	if err != nil {
		log.Error(err)
		return
	}

	log.Debugf("%+v", webhookInfo)

	pendingUpdates.WithLabelValues(bot.Self.UserName).Set(float64(webhookInfo.PendingUpdateCount))
	lastErrorDate.WithLabelValues(bot.Self.UserName).Set(float64(webhookInfo.LastErrorDate))

	log.Info("update complete")
}
