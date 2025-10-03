package main

import (
	"context"
	"fmt"
	"vercel_deployment_service/aws"
	"vercel_deployment_service/utils"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var subscriber = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func main() {
	godotenv.Load()
	fmt.Println("VERCEL DEPLOYMENT SERVICE STARTED")
	ctx := context.Background()
	for true {

		res, err := subscriber.BRPop(ctx, 0, "build-queue").Result()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res)
		aws.DownloadS3Object(context.Background(), "vercel-arch", res[1], "./output")
		utils.BuildProject(res[1])
		aws.CopyFinalDist(res[1])
	}
}
