/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hirosassa/zerodriver"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	telebot "gopkg.in/telebot.v3"
)

var (
	//Teletoken bot
	Teletoken   = os.Getenv("TELE_TOKEN")
	MetricsHost = os.Getenv("METRICS_HOST")
)

func initMetrics(ctx context.Context) {

	exporter, _ := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(MetricsHost),
		otlpmetricgrpc.WithInsecure(),
	)

	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(fmt.Sprintf("kbot_%s", appVersion)),
	)

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(10*time.Second)),
		),
	)

	otel.SetMeterProvider(mp)

}

func pmetrics(ctx context.Context, payload string) {
	// Get the global MeterProvider and create a new Meter with the name "kbot_light_signal_counter"
	meter := otel.GetMeterProvider().Meter("kbot_light_signal_counter")

	// Get or create an Int64Counter instrument with the name "kbot_light_signal_<payload>"
	counter, _ := meter.Int64Counter(fmt.Sprintf("kbot_light_signal_%s", payload))

	// Add a value of 1 to the Int64Counter
	counter.Add(ctx, 1)
}

// kbotCmd represents the kbot command
var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := zerodriver.NewProductionLogger()
		fmt.Printf("kbot is started", appVersion)
		kbot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  Teletoken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			logger.Fatal().Str("Error", err.Error()).Msg("Please check TELE_TOKEN")
			return
		} else {
			logger.Info().Str("Version", appVersion).Msg("kbot started")

		}

		kbot.Handle(telebot.OnText, func(m telebot.Context) error {

			log.Print(m.Message().Payload, m.Text())
			payload := m.Message().Payload

			switch payload {
			case "version":
				err = m.Send(fmt.Sprintf(appVersion))

			case "hello":
				err = m.Send(fmt.Sprintf("Hello!"))

			case "music":
				err = m.Send(fmt.Sprintf("https://www.youtube.com/watch?v=CTAud5O7Qqk"))
			}

			return err
		})

		kbot.Start()

	},
}

func init() {
	rootCmd.AddCommand(kbotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// kbotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// kbotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
