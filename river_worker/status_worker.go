package river_worker

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jackc/pgx/v5"
    "github.com/riverqueue/river"

    "net/http"
)
type ServiceStatus struct {
    Status string `json:"status"`
}

func (ServiceStatus) Kind() string { return "indicate_status" }

type StatusWorker struct {
    river.WorkerDefaults[ServiceStatus]
}

func (w *StatusWorker) Work(ctx context.Context, job *river.Job[ServiceStatus]) error {
    fmt.Println("status: " + job.Args.Status)
    return nil
}

func StartStatusJob(riverClient *river.Client[pgx.Tx], ctx context.Context, dbPool *pgxpool.Pool) http.HandlerFunc{
  return func(w http.ResponseWriter, r *http.Request) {
    tx, err := dbPool.Begin(ctx)
    if err != nil {
        panic(err)
    }
    defer tx.Rollback(ctx)
    str := r.URL.Query().Get("status")

    _, err = riverClient.InsertTx(ctx, tx, ServiceStatus{ Status: str }, nil)
    if err != nil {
        panic(err)
    }

    if err := tx.Commit(ctx); err != nil {
        panic(err)
    }
    w.Write([]byte("ok"))
  }
}
