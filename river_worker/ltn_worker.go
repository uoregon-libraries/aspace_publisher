package river_worker

import (
    "context"
    "net/url"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jackc/pgx/v5"
    "github.com/riverqueue/river"

    "aspace_publisher/alma" //packages worker needs
    "net/http"
)
type LinkToNetwork struct {
    MmsID string `json:"mms_id"`
    Filename string `json:"filename"`
}

func (LinkToNetwork) Kind() string { return "link_to_network" }

type LTNWorker struct {
    river.WorkerDefaults[LinkToNetwork]
}

func (w *LTNWorker) Work(ctx context.Context, job *river.Job[LinkToNetwork]) error {
    alma.LinkToNetwork([]string{job.Args.MmsID}, job.Args.Filename)
    return nil
}

func StartLTNJob(riverClient *river.Client[pgx.Tx], ctx context.Context, dbPool *pgxpool.Pool) http.HandlerFunc{
  return func(w http.ResponseWriter, r *http.Request) {
    tx, err := dbPool.Begin(ctx)
    if err != nil {
        panic(err)
    }
    defer tx.Rollback(ctx)
    id := r.URL.Query().Get("id")
    filename := r.URL.Query().Get("filename")
    _, err = riverClient.InsertTx(ctx, tx, LinkToNetwork{ MmsID: id, Filename: filename }, nil)
    if err != nil {
        panic(err)
    }

    if err := tx.Commit(ctx); err != nil {
        panic(err)
    }
    w.Write([]byte("ok"))
  }
}
