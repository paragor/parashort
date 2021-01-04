package parashort

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/paragor/parashort/pkg/domain/storage"
	"github.com/paragor/parashort/pkg/domain/user"
)

type ParashortApp struct {
	requestTimeProcessing time.Duration
	storage               storage.StorageEngine
}

func NewParashortApp(requestTimeProcessing time.Duration, storage storage.StorageEngine) *ParashortApp {
	return &ParashortApp{requestTimeProcessing: requestTimeProcessing, storage: storage}
}

func (app *ParashortApp) Delete(key string) error {
	ctx, _ := context.WithTimeout(context.Background(), app.requestTimeProcessing)

	if len(key) == 0 {
		return fmt.Errorf("invalid key")
	}

	return app.storage.Delete(ctx, key)
}
func (app *ParashortApp) List() ([]string, error) {
	ctx, _ := context.WithTimeout(context.Background(), app.requestTimeProcessing)

	return app.storage.ListKeys(ctx)
}

func (app *ParashortApp) SaveText(item SaveItem) (*SaveResult, error) {
	ctx, _ := context.WithTimeout(context.Background(), app.requestTimeProcessing)

	key := ""
	strict := false
	if item.RequiredKey != "" {
		strict = true
		key = item.RequiredKey
	} else {
		strict = false
		key = app.generateKey(app.itemToPayload(item))
	}
	if err := app.validate(key); err != nil {
		return nil, fmt.Errorf("key (%s) is invalid: %w", key, err)
	}
	key, err := app.reserveKey(ctx, key, strict)
	if err != nil {
		return nil, fmt.Errorf("cant reserve key: %w", err)
	}

	if err := app.storage.Save(ctx, key, item.Text, item.TTL); err != nil {
		return nil, fmt.Errorf("cant save: %w", err)
	}

	return &SaveResult{Key: key}, nil
}

func (app *ParashortApp) LoadText(key string) (string, error) {
	ctx, _ := context.WithTimeout(context.Background(), app.requestTimeProcessing)

	text, err := app.storage.Get(ctx, strings.ToLower(key))
	if err != nil {
		return "", err
	}
	return text, nil
}

var validationRe = regexp.MustCompile("^[a-z0-9]{5,100}$")

func (app *ParashortApp) validate(key string) error {
	if !validationRe.MatchString(key) {
		return fmt.Errorf("key should conains only [a-z1-9], 5 <= len <= 100")
	}
	return nil
}

func (app *ParashortApp) reserveKey(ctx context.Context, key string, strict bool) (reservedKey string, err error) {
	for {
		err = app.storage.ReserveKey(ctx, key, app.requestTimeProcessing)
		if err == nil {
			return key, nil
		}
		if errors.Is(err, storage.ErrKeyAlreadyReversed) {
			if strict {
				return "", err
			} else {
				key = app.generateKey(key + strconv.Itoa(rand.Int()) + time.Now().String())
				continue
			}
		}
		return "", fmt.Errorf("cant reserve key: %w", err)
	}
}
func (app *ParashortApp) itemToPayload(item SaveItem) string {
	payload := item.Text
	if item.User.LastIp != nil {
		payload += item.User.LastIp.String()
	}
	payload += time.Now().String()
	return payload
}

func (app *ParashortApp) generateKey(payload string) string {
	return strings.ToLower(fmt.Sprintf("%x", md5.Sum([]byte(payload))))[:5]
}

type SaveItem struct {
	RequiredKey string
	Text        string
	User        *user.UserInfo
	TTL         time.Duration
}
type SaveResult struct {
	Key string
}
