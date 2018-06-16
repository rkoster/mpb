package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWaitGroup(&wg))
	total := 100
	numBars := 3
	wg.Add(numBars)

	for i := 0; i < numBars; i++ {
		name := fmt.Sprintf("Bar#%d:", i)

		var bOption mpb.BarOption
		if i == 0 {
			bOption = mpb.BarRemoveOnComplete()
		}

		sbEta := make(chan time.Time)
		b := p.AddBar(int64(total), mpb.BarID(i),
			bOption,
			mpb.PrependDecorators(
				decor.Name(name),
				decor.ETA(decor.ET_STYLE_GO, 60, sbEta, decor.WCSyncSpace),
			),
			mpb.AppendDecorators(decor.Percentage()),
		)
		go func() {
			defer wg.Done()
			max := 100 * time.Millisecond
			for i := 0; i < total; i++ {
				sbEta <- time.Now()
				if b.ID() == 2 && i == 42 {
					p.Abort(b)
					return
				}
				time.Sleep(time.Duration(rand.Intn(10)+1) * max / 10)
				b.Increment()
			}
		}()
	}

	p.Wait()
}
