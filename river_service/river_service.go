package main

import (
    "context"
    "fmt"
    "log"
    "log/slog"
    "os"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/riverqueue/river"
    "github.com/riverqueue/river/riverdriver/riverpgxv5"
    "github.com/riverqueue/river/rivershared/util/slogutil"
    "riverqueue.com/riverui"
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

    //add UI
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn, ReplaceAttr: slogutil.NoLevelTime}))
    endpoints := riverui.NewEndpoints(riverClient, nil)
    opts := &riverui.HandlerOpts{
        Endpoints: endpoints,
        Logger: logger,
        Prefix: "/riverui", // mount the UI and its APIs at /riverui or some path
        // ...
    }
    handler, err := riverui.NewHandler(opts)
    if err != nil {
        log.Fatal(err)
    }
    // Start the handler to initialize background processes for caching and periodic queries:
    handler.Start(ctx)
    // add routes here
    http.HandleFunc("/startLTNJob", river_worker.StartLTNJob(riverClient, ctx, dbPool))
    http.HandleFunc("/startStatusJob", river_worker.StartStatusJob(riverClient, ctx, dbPool))
    http.Handle("/riverui/", handler)
    http.ListenAndServe(":3200", nil)
}
