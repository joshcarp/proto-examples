package main

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	connect "github.com/bufbuild/connect-go"
	"github.com/getsentry/sentry-go"
	"log"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:          "https://7735f2e9fe5f4094a0ebee6cf19a5865@o4505268263911424.ingest.sentry.io/4505268264763392",
		Integrations: integrationsFunc,
	})
	// These are global tags that are applied to all events.
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTags(map[string]string{
			"environment": "production",
		})
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	//err = connect.NewError(connect.CodeInternal, fmt.Errorf("hello world 3"))

	//sentry.CaptureException(err)

	for i := 0; i < 10; i++ {
		captureEvent(NewConnectError())
	}
}

func integrationsFunc(integrations []sentry.Integration) []sentry.Integration {
	var filtered []sentry.Integration
	for _, integration := range integrations {
		// Do not capture sentry integrations because they lead to noisy events. Instead we manually
		// build this up.
		//
		// Integrations available:
		//  - ContextifyFrames
		//  - Environment
		//  - Modules
		//  - IgnoreErrors
		//
		// See docs for more details:
		// https://docs.sentry.io/platforms/go/configuration/options/#removing-default-integrations
		switch strings.ToLower(integration.Name()) {
		case "modules", "environment", "contextifyframes":
			continue
		}
		filtered = append(filtered, integration)
	}
	return filtered
}

func NewConnectError() error {
	//switch rand.Int() % 4 {
	//case 0:
	//	return connect.NewError(connect.CodeInternal, fmt.Errorf("error with something important:"+randomString()))
	//case 1:
	//	return connect.NewError(connect.CodeInternal, fmt.Errorf("error with %s something something", randomString()))
	//case 2:
	//	return connect.NewError(connect.CodeInternal, fmt.Errorf("error with something important:"+randomString()))
	//case 3:
	//	return connect.NewError(connect.CodeInternal, fmt.Errorf("error with something important:"+randomString()))
	//}
	return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to parse package name %s", randomString(1)))
}

//func unwrapError(){
//
//}

//for i := 0; i < 10; i++ {
//	err = errorInternalConnect()
//	eventID := sentry.CurrentHub().Clone().CaptureException(errors.New("hello world"))
//	println(*eventID)
//	time.Sleep(time.Second * 2)
//}

// foo calculates the value of the flux capacitor based on the inputs provided.
// It returns the value of the flux capacitor and an error if the flux capacitor
// could not be calculated.
func errorInternalConnect() error {
	return connect.NewError(connect.CodeInternal, fmt.Errorf("randomString()"))
}

// randomString uses gofakeit to create a random string.
func randomString(l int) string {
	return gofakeit.New(rand.Int63()).FirstName()
}

func captureEvent(err error) {
	event := sentry.NewEvent()
	hub := sentry.CurrentHub()
	event.Level = sentry.LevelError
	event.Message = err.Error()
	event.Contexts = map[string]map[string]interface{}{
		"Go Runtime": {
			"version":      runtime.Version(),
			"numcpu":       runtime.NumCPU(),
			"maxprocs":     runtime.GOMAXPROCS(0),
			"numgoroutine": runtime.NumGoroutine(),
		},
	}
	//var subtitle string
	//if unwrappedErr := errors.Unwrap(err); unwrappedErr != nil {
	//	subtitle = unwrappedErr.Error()
	//}
	//True := true

	for i := 0; i < hub.Client().Options().MaxErrorDepth && err != nil; i++ {
		event.Exception = append(event.Exception, sentry.Exception{
			Value: err.Error(),
			Type:  reflect.TypeOf(err).String(),
		})
		switch previous := err.(type) {
		case interface{ Unwrap() error }:
			err = previous.Unwrap()
		case interface{ Cause() error }:
			err = previous.Cause()
		default:
			err = nil
		}
	}
	//stack := sentry.NewStacktrace()
	//stack.Frames = filterFrames(stack.Frames)
	//event.Exception = []sentry.Exception{{
	//	// Type is used as the main issue title.
	//	Type: err.Error(),
	//	// Value is used as the main issue subtitle.
	//	Value: subtitle,
	//	//Stacktrace: stack,
	//}}

	println(*hub.CaptureEvent(event))
}
