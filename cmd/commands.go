package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Commands = []*cli.Command{
	{
		Name: "crawl",
		Subcommands: []*cli.Command{
			{
				Name: "groupware",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "job"},
					&cli.StringFlag{Name: "masters"},
				},
				Action: CrawlGroupWareDeclinedPayments,
			},
			{
				Name: "hackernews",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "channel"},
					&cli.IntFlag{Name: "point_threshold"},
				},
				Action: CrawlHackerNews,
			},
			{
				Name: "quasarzone",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "channel"},
				},
				Action: CrawlQuasarZoneSales,
			},
			{
				Name: "book",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "channel"},
				},
				Action: CrawlBook,
			},
			{
				Name: "gitpublic",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "channel"},
				},
				Action: CrawlGitPublic,
			},
			{
				Name: "wanted",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "channel"},
					&cli.StringFlag{Name: "query"},
				},
				Action: CrawlWanted,
			},
			{
				Name: "eomisae",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "channel"},
					&cli.StringFlag{Name: "target"},
				},
				Action: CrawlEomisae,
			},
			{
				Name: "ipo",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "channel"},
				},
				Action: CrawlIPO,
			},
		},
	},
	{
		Name: "restriction",
		Subcommands: []*cli.Command{
			{
				Name: "add",
				Flags: []cli.Flag{
					&cli.TimestampFlag{Name: "start_date"},
					&cli.TimestampFlag{Name: "end_date"},
					&cli.TimestampFlag{Name: "hour_from"},
					&cli.TimestampFlag{Name: "hour_to"},
				},
				Action: AddRestriction,
			},
		},
	},
	{
		Name:        "test",
		Subcommands: []*cli.Command{
			// {Name: "slack", Action: TestSlack},
			// {Name: "chrome", Action: TestChrome},
		},
	},
}

// AddRestriction adds a restriction
func AddRestriction(ctx *cli.Context) error {
	return nil
}

func getChromeURL(logger *zap.Logger, chromeHost string) (string, error) {
	endpoint := fmt.Sprintf("http://%s/json/version", chromeHost)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}
	req.Host = "localhost"

	// request to chrome
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		logger.Error("get", zap.Error(err))
		return "", err
	}
	defer res.Body.Close()

	// read buffer
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("read", zap.Error(err))
		return "", err
	}

	var m map[string]string
	err = json.Unmarshal(body, &m)
	if err != nil {
		return "", err
	}

	wsURL, ok := m["webSocketDebuggerUrl"]
	if !ok {
		return "", errors.New("webSocketDebuggerUrl is not found")
	}

	u, err := url.Parse(wsURL)
	if err != nil {
		return "", err
	}
	u.Host = chromeHost // replace to chrome host
	return u.String(), nil
}

func zapLogger() *zap.Logger {
	// Create logger configuration
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create logger with configurations
	zapLogger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		zap.NewAtomicLevelAt(zapcore.InfoLevel),
	))

	return zapLogger
}
