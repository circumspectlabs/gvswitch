package control

import "context"

func Serve(ctx context.Context, addr string) error {
	<-ctx.Done()

	return nil
}
