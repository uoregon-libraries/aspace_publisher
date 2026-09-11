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
    "github.com/DeRuina/timberjack"
    "riverqueue.com/riverui"
    "aspace_publisher/river_worker"
    "net/http"
    "io"
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

    logmode := os.Getenv("LOGMODE")
    path := os.Getenv("HOME_DIR")
    var logr io.Writer
    if logmode == "file" {
      logr = &timberjack.Logger{
      Filename:   path + "logs/river.log", // path of log file
      MaxSize:    50, // file size in MB
      MaxBackups: 7, // number of files to retain
      MaxAge:     8, // how long (in days) to retain files
      Compression: "gzip", // archive files?
      LocalTime:  true, // re timestamps
      RotateAt: []string{"00:00"},
      }
    } else { logr = os.Stdout }

    workers := river.NewWorkers()
    // add each type of workers here
    river.AddWorker(workers, &river_worker.LTNWorker{})
    river.AddWorker(workers, &river_worker.StatusWorker{})
    river.AddWorker(workers, &river_worker.AlmaCrupWorker{})
    logger := slog.New(slog.NewTextHandler(logr, &slog.HandlerOptions{Level: slog.LevelInfo, ReplaceAttr: slogutil.NoLevelTime}))
    slog.SetDefault(logger)
    riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
        Logger: logger,
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
    http.HandleFunc("/startAlmaCrupJob", river_worker.StartAlmaCrupJob(riverClient, ctx, dbPool))
    http.Handle("/riverui/", handler)
    http.ListenAndServe(":3200", nil)
}
