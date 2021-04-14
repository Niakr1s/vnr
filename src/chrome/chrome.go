package chrome

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

type Chrome struct {
	allocCtx context.Context
	ctx      context.Context

	options ChromeOptions
}

type ChromeOptions struct {
	Headless bool
	Timeout  time.Duration
}

func NewChrome(ctx context.Context, options ChromeOptions) (*Chrome, error) {
	chrome := &Chrome{options: options}

	dir, err := ioutil.TempDir("", "vnr")
	if err != nil {
		return nil, err
	}

	allocatorOptions := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.UserDataDir(dir),

		// without this flag headless mode is very slow
		chromedp.Flag("proxy-server", ""),
		chromedp.Flag("headless", options.Headless),
	)

	chrome.allocCtx, _ = chromedp.NewExecAllocator(ctx, allocatorOptions...)

	ctx, _ = chromedp.NewContext(chrome.allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	chrome.ctx = ctx

	if err := chromedp.Run(ctx); err != nil {
		return nil, err
	}

	return chrome, nil
}

func (chrome *Chrome) getCtx() (context.Context, context.CancelFunc) {
	return chromedp.NewContext(chrome.ctx)
}

func (chrome *Chrome) Run(actions ...chromedp.Action) error {
	taskCtx, cancel := chrome.getCtx()
	defer cancel()
	taskCtx, cancel = context.WithTimeout(taskCtx, chrome.options.Timeout)
	defer cancel()

	err := chromedp.Run(taskCtx, actions...)
	return err
}
