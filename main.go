package main

import (
	"cmp"
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/chromedp/chromedp"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

func main() {
	allocatorContext, cancel := chromedp.NewRemoteAllocator(
		context.Background(),
		os.Getenv("BROWSER_WS_ENDPOINT"),
		chromedp.NoModifyURL,
	)

	defer cancel()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := chromedp.NewContext(allocatorContext)

		defer cancel()

		// Set viewport size to 1920x1080
		if err := chromedp.Run(ctx, chromedp.EmulateViewport(1920, 1080)); err != nil {
			logger.Error("error emulating viewport", "error", err)
			http.Error(w, "error emulating viewport", http.StatusInternalServerError)
			return
		}

		var buf []byte

		if err := chromedp.Run(ctx, chromedp.Tasks{
			chromedp.Navigate(`https://example.com`),
			chromedp.FullScreenshot(&buf, 100),
		}); err != nil {
			logger.Error("error taking screenshot", "error", err)
			http.Error(w, "error taking screenshot", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")

		w.Write(buf)
	})

	port := cmp.Or(os.Getenv("PORT"), "8080")

	logger.Info("starting server on port", "port", port)

	server := &http.Server{
		Addr:    (":" + port),
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Error("error starting server", "error", err)
		os.Exit(1)
	}
}
