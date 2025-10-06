package confluence

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/graze/go-throttled"
	"github.com/spf13/viper"
	conflu "github.com/virtomize/confluence-go-api"
	"github.com/vogtp/rag/pkg/cfg"
	vecdb "github.com/vogtp/rag/pkg/vecDB"
	"github.com/vogtp/rag/pkg/vecDB/pdf"
	"golang.org/x/time/rate"
)

// GetDocuments retrives confluence spaces and generates vecdb.EmbeddDocuments
func GetDocuments(ctx context.Context, slog *slog.Logger, config cfg.RagConfig, spaces ...string) (chan vecdb.EmbeddDocument, error) {
	baseURL := config.Confluence.BaseURL
	baseURL = strings.TrimRight(baseURL, "/")
	conf := confluence{
		slog:       slog.With("confluence_url", baseURL),
		baseURL:    baseURL,
		out:        make(chan vecdb.EmbeddDocument, 10),
		accessKey:  config.Confluence.Key,
		rateLimit:  rate.Limit(0.4),
		queryLimit: 100,
		spaces:     spaces,
		config:     config,
	}
	if len(spaces) < 1 {
		conf.spaces = config.Confluence.Spaces
	}
	if err := conf.init(); err != nil {
		return nil, err
	}
	go conf.query(ctx)
	return conf.out, nil
}

type confluence struct {
	slog       *slog.Logger
	baseURL    string
	out        chan vecdb.EmbeddDocument
	api        *conflu.API
	accessKey  string
	rateLimit  rate.Limit
	queryLimit int
	spaces     []string
	mu         sync.Mutex
	config     cfg.RagConfig
}

func (c *confluence) init() error {
	url := c.getAPIURL()
	api, err := conflu.NewAPI(url, "", c.accessKey)
	if err != nil {
		return err
	}
	if api == nil {
		return fmt.Errorf("confluence api was not created")
	}
	api.Client.Transport = &uaRT{
		RoundTripper: api.Client.Transport,
	}

	api.Client = throttled.WrapClient(api.Client, rate.NewLimiter(c.rateLimit, 1))
	c.api = api
	c.slog.Info("loaded confluence rest api", "url", url)
	return nil
}

func (c *confluence) getAPIURL() string {
	return c.getURL("/rest/api")
}

func (c *confluence) getURL(ui string) string {
	return fmt.Sprintf("%s%s", c.baseURL, ui)
}

func (c *confluence) query(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()
	defer close(c.out)
	for _, s := range c.spaces {
		if ctx.Err() != nil {
			return
		}
		c.querySpace(ctx, s)
	}
}

func (c *confluence) querySpace(ctx context.Context, spaceKey string) {
	slogSpace := c.slog.With("space", spaceKey)
	slogSpace.Info("starting to query space")
	maxPageAge := viper.GetDuration(cfg.ConfluenceMaxAge)
	if maxPageAge < time.Minute {
		maxPageAge = cfg.DefaultConfluenceMaxAge
	}
	start := 0
	total := 0
	retryCnt := 0
	retryMax := 3
	retryDelay := 500 * time.Millisecond
	for {
		slogPage := slogSpace.With(slog.Group("paging", "start", start, "limit", c.queryLimit))
		if ctx.Err() != nil {
			slogPage.Warn("confluence canceled by context", "err", ctx.Err())
			return
		}
		res, err := c.api.GetContent(conflu.ContentQuery{
			SpaceKey: spaceKey,
			Start:    start,
			Limit:    c.queryLimit,
			Expand:   []string{"space", "body.view", "version", "container", "body.storage", "metadata", "history.lastUpdated"},
		})
		if err != nil {
			// FIXME handle unauthorised errors
			retryCnt++
			if retryCnt > retryMax {
				slogPage.Error("Max Retries reached, cannot get confluence content...", "err", err, "start_index", start, "retryCnt", retryCnt, "retryMax", retryMax, "retryDelay", retryDelay)
				return
			}
			slogPage.Warn("Cannot get confluence content...", "err", err, "start_index", start, "retryCnt", retryCnt, "retryMax", retryMax, "retryDelay", retryDelay)
			time.Sleep(retryDelay * time.Duration(retryCnt))
			continue
		}
		start += res.Limit
		total += res.Size
		retryCnt = 0

		for _, d := range res.Results {
			slog := slogPage.With("title", d.Title, "doc_url", d.Links.WebUI)
			if ctx.Err() != nil {
				slog.Warn("confluence canceled by context", "err", ctx.Err())
				return
			}
			slog.Debug("processing confluence document", "title", d.Title)
			//	txt := html2text.HTML2Text(d.Body.View.Value)
			bodyMD, pdfLinks := parsePage(slog, d.Body.View.Value)
			doc := vecdb.EmbeddDocument{
				Title:       d.Title,
				URL:         c.getURL(d.Links.WebUI),
				Document:    bodyMD,
				IDMetaKey:   vecdb.MetaURL,
				IDMetaValue: d.Links.WebUI,
				MetaData:    make(map[string]any),
			}
			// 2016-05-30T16:14:07.787+02:00
			if t, err := time.Parse(time.RFC3339Nano, d.History.LastUpdated.When); err == nil {
				doc.Modified = t
			} else {
				slog.Error("Cannot parse time of confluence page", "time", t.String(), "err", err, "title", d.Title, "url", d.Links.WebUI)
				continue
			}
			if time.Since(doc.Modified) > maxPageAge {
				slog.Info("Document is to old", "age", time.Since(doc.Modified).String(), "lastModify", doc.Modified.String(), "maxAge", maxPageAge.String())
				continue
			}

			doc.MetaData["confluence_space"] = spaceKey
			c.out <- doc
			for _, pl := range pdfLinks {
				docs, err := pdf.SplitFromLink(ctx, pl)
				if err != nil {
					slog.Warn("Cannot embedd PDF", "pdf.url", pl, "err", err)
					continue
				}
				for _, d := range docs {
					c.out <- d
				}
			}
		}

		slogPage.Info("confluence query batch done", "start", res.Start, "size", res.Size, "result_len", len(res.Results), "limit", res.Limit, "total", total)
		if res.Limit != res.Size {
			slogPage.Info("Indexing space done")
			return
		}
	}
}

// func (c *confluence) getAllSpaces() {
// 	panic("Not implemented")
// 	// spaces, err := api.GetAllSpaces(conflu.AllSpacesQuery{
// 	// 	Type:   "global",
// 	// 	Start:  0,
// 	// 	Limit:  999,
// 	// 	Expand: []string{"space", "body.view", "version", "container"},
// 	// })
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// for _, space := range spaces.Results {
// }
