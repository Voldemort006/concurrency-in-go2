package the_context_package

import (
	"context"
	"fmt"
)

func main() {
	ProcessRequest("siddharth", "siddharth06")
}

func ProcessRequest(name, userName string) {
	ctx := context.WithValue(context.Background(), "name", name)
	ctx = context.WithValue(ctx, "userName", userName)

	HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
	fmt.Printf("handling response for name = %v and username = %v", ctx.Value("name"), ctx.Value("userName"))
}
