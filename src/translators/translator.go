package translators

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

type TranslatorInitOptions struct {
	Headless bool
}

type TranslationOptions struct {
	From     string
	To       string
	Sentence string

	Timeont time.Duration
}

func NewTranslationOptions(sentence string) TranslationOptions {
	return TranslationOptions{
		From:     "auto",
		To:       "auto",
		Sentence: sentence,

		Timeont: time.Second * 15,
	}
}

type TranslationResult struct {
	TranslationOptions
	Translation string
}

type Translator interface {
	Init(ctx context.Context, options TranslatorInitOptions) error
	GetTranslation(translationOptions TranslationOptions) (TranslationResult, error)
}

func GetTranslator(translator string) (Translator, error) {
	switch translator {
	case "deepl":
		return NewDeeplTranslator(), nil
	default:
		return nil, fmt.Errorf("invalid translator: got: %s, expected: one of %s", translator, KnownTranslators)
	}
}

type BaseChromeTranslator struct {
	allocCtx context.Context
	ctx      context.Context
}

func (t *BaseChromeTranslator) Init(ctx context.Context, options TranslatorInitOptions) error {
	dir, err := ioutil.TempDir("", "vnr")
	if err != nil {
		return err
	}

	go func() {
		for range ctx.Done() {
			fmt.Printf("DONE")
		}
	}()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.UserDataDir(dir),

		// without this flag headless mode is very slow
		chromedp.Flag("proxy-server", ""),
		chromedp.Flag("headless", options.Headless),
	)

	t.allocCtx, _ = chromedp.NewExecAllocator(ctx, opts...)

	// t.ctx, _ = t.getCtx()
	ctx, _ = chromedp.NewContext(t.allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	t.ctx = ctx
	// defer cancel()
	if err := chromedp.Run(ctx); err != nil {
		return err
	}

	return nil
}

func (t *BaseChromeTranslator) getCtx() (context.Context, context.CancelFunc) {
	return chromedp.NewContext(t.ctx)
}
