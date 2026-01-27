package main

import (
    "context"
    "fmt"
    "log/slog"
    "os"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/riverqueue/river"
    "github.com/riverqueue/river/riverdriver/riverpgxv5"
    "github.com/riverqueue/river/rivershared/util/slogutil"

    "aspace_publisher/river_worker"
    "net/http"
)

func main() {
    ctx := context.Background()
    fmt.Println("have ctx")
    pguser := os.Getenv("POSTGRES_USER")
    pgpass := os.Getenv("POSTGRES_PASSWORD")
    pgdb := os.Getenv("POSTGRES_DB")
    pgaddress := os.Getenv("DATABASE_URL")

    dbPool, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s/%s", pguser,pgpass,pgaddress,pgdb))
    if err != nil {
        panic(err)
    }
    defer dbPool.Close()
    fmt.Println("have dbpool")
    workers := river.NewWorkers()
    // add each type of workers here
    river.AddWorker(workers, &river_worker.LTNWorker{})
    river.AddWorker(workers, &river_worker.StatusWorker{})

    riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
        Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: slogutil.NoLevelTime})),
        Queues: map[string]river.QueueConfig{
            river.QueueDefault: {MaxWorkers: 100},
        },
        Workers:  workers,
    })
    if err != nil {
        panic(err)
    }
    fmt.Println("have client")

    if err := riverClient.Start(ctx); err != nil {
        panic(err)
    }
    defer riverClient.Stop(ctx)

    // add routes here
    http.HandleFunc("/startLTNJob", river_worker.StartLTNJob(riverClient, ctx, dbPool))
    http.HandleFunc("/startStatusJob", river_worker.StartStatusJob(riverClient, ctx, dbPool))

    http.ListenAndServe(":3200", nil)
}
