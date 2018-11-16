package core

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestMergeResults(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	options := &EnumerationOptions{
		Sources:   []Source{&FakeSource1{}, &FakeSource2{}},
		Recursive: true,
	}

	counter := 0

	apple := EnumerateSubdomains(ctx, "apple.com", options)
	google := EnumerateSubdomains(ctx, "google.com", options)

	merged := MergeResults(apple, google)

	for result := range merged {
		counter++
		fmt.Println(result)
	}

	fmt.Println(counter, ctx.Err())

}
