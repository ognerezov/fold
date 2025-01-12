package threads

import (
	"fmt"
	"fold/console"
	"fold/db"
	"fold/interfaces"
	"time"
)

type AsyncCall[A any, T any] interface {
	Call(args A) (Message[T], ErrorMessage)
}

type CallBackReceiver[T any] interface {
	OnError(err error)
	OnResult(result T)
}

var (
	// FlushInterval this should go to config
	FlushInterval time.Duration = 5
	Ticker                      = time.NewTicker(FlushInterval * time.Second)
	Errors                      = make(chan ErrorMessage)
	Results                     = make(chan Message[string])
	Quit                        = make(chan struct{})
	providers     *interfaces.Providers
)

func wrap[A any](caller AsyncCall[A, string]) func(A) {
	return func(a A) {
		res, err := caller.Call(a)
		if err.e != nil {
			Errors <- err
		}
		if res.hasResult {
			Results <- res
		}
	}
}

func Async[A any](caller AsyncCall[A, string], arg A) {
	wrapped := wrap[A](caller)
	go wrapped(arg)
}

func readChannels() {
	for {
		select {
		case msg := <-Results:
			if msg.receivers != nil {
				for _, receiver := range msg.receivers {
					receiver.OnResult(msg.t)
				}
			}
			console.GreenPrintln(fmt.Sprintf("Process %s returns result %v", msg.process, msg.t))

		case msg := <-Errors:
			if msg.receivers != nil {
				for _, receiver := range msg.receivers {
					receiver.OnError(msg.e)
				}
			}
			console.RedPrintln(msg.e.Error())
		case <-Quit:
			return
		}
	}
}

func flushDb() {
	for {
		select {
		case <-Ticker.C:
			tables := db.Db()
			for file, table := range *tables {
				fmt.Println("write")
				fmt.Println(table)
				WriteCsvAsync(file, table)
			}
			updates := db.Pending()

			for file, n := range *updates {
				WriteNosqlAsync(file, n)
			}
			driveUpdates := db.DrivePending()
			for _, update := range *driveUpdates {
				WriteDriveAsync(update)
			}

		case <-Quit:
			Ticker.Stop()
			return
		}
	}
}

func Start(p *interfaces.Providers) {
	providers = p
	go readChannels()
	go flushDb()
}
