package smirror

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"smirror/base"
	"smirror/gs"
	"time"
)

//StorageMirror cloud function entry point
func StorageMirror(ctx context.Context, event gs.Event) (err error) {
	start := time.Now()
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	response, err := storageMirror(ctx, event)
	elapsed := time.Since(start)
	if err != nil {
		err = errors.Wrap(err, "failed to mirror "+event.URL())
		return err
	}

	if base.IsLoggingEnabled() {
		fmt.Printf("mirrored %v: %v in %v", response.Status, event.URL(), elapsed)
	}
	return err
}

func storageMirror(ctx context.Context, event gs.Event) (*Response, error) {
	service, err := NewFromEnv(ctx, ConfigEnvKey)
	if err != nil {
		return nil, err
	}
	if base.IsLoggingEnabled() {
		fmt.Printf("uses service %T(%p), err: %v\n", service, service, err)
	}
	response := service.Mirror(ctx, NewRequest(event.URL()))
	if base.IsLoggingEnabled() {
		if data, err := json.Marshal(response); err == nil {
			fmt.Printf("%v\n", string(data))
		}
	}
	if response.Error != "" {
		return nil, fmt.Errorf(response.Error)
	}
	return response, nil
}
