package cdata

import (
	"runtime"
	"sync"
	"testing"
)

func TestRef(t *testing.T) {
	const N = 10
	runtime.LockOSThread()
	wg := sync.WaitGroup{}
	wg.Add(N)
	ch := make(chan uintptr)
	for i := 0; i < N; i++ {
		go func() {
			runtime.LockOSThread()
			wg.Done()
			ch <- Ref()
			wg.Wait()
		}()
	}
	wg.Wait()
	refs := []uintptr{Ref()}
	for i := 0; i < N; i++ {
		chref := <-ch
		for _, ref := range refs {
			if chref == ref {
				t.Fatalf("found duplicated ref: %d == %d", chref, ref)
			}
		}
		refs = append(refs, chref)
	}
}
