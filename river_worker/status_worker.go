package river_worker

import (
    "context"
    "log"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jackc/pgx/v5"
    "github.com/riverqueue/river"

    "net/http"
)
type ServiceStatus struct {
    Status string `json:"status"`
    Other string `json:"other"`
}

func (ServiceStatus) Kind() string { return "indicate_status" }

type StatusWorker struct {
    river.WorkerDefaults[ServiceStatus]
}

func (w *StatusWorker) Work(ctx context.Context, job *river.Job[ServiceStatus]) error {
    log.Println("status: " + job.Args.Status)
    log.Println("other: " + job.Args.Other)
    return nil
}

func StartStatusJob(riverClient *river.Client[pgx.Tx], ctx context.Context, dbPool *pgxpool.Pool) http.HandlerFunc{
  return func(w http.ResponseWriter, r *http.Request) {
    tx, err := dbPool.Begin(ctx)
    if err != nil {
        panic(err)
    }
    defer tx.Rollback(ctx)
    sta := r.URL.Query().Get("status")
    oth := r.URL.Query().Get("other")
    _, err = riverClient.InsertTx(ctx, tx, ServiceStatus{ Status: sta, Other: oth }, nil)
    if err != nil {
        panic(err)
    }

    if err := tx.Commit(ctx); err != nil {
        panic(err)
    }
    w.Write([]byte("ok"))
  }
}
