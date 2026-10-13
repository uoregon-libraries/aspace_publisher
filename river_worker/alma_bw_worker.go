package river_worker

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jackc/pgx/v5"
    "github.com/riverqueue/river"

    "aspace_publisher/handlers" //packages worker needs
    "net/http"
)
type AlmaBW struct {
    TopCon string `json:"id"`
    Filename string `json:"filename"`
    Session string `json:"session"`
}

func (AlmaBW) Kind() string { return "alma_boundwith" }

type AlmaBWWorker struct {
    river.WorkerDefaults[AlmaBW]
}

func (w *AlmaBWWorker) Work(ctx context.Context, job *river.Job[AlmaBW]) error {
    handlers.AlmaBW(job.Args.TopCon, job.Args.Filename, job.Args.Session)
    return nil
}

func StartAlmaBWJob(riverClient *river.Client[pgx.Tx], ctx context.Context, dbPool *pgxpool.Pool) http.HandlerFunc{
  return func(w http.ResponseWriter, r *http.Request) {
    defer HandlerPanic()
    tx, err := dbPool.Begin(ctx)
    if err != nil { panic(err) }
    defer tx.Rollback(ctx)
    id := r.URL.Query().Get("id")
    filename := r.URL.Query().Get("filename")
    session := r.URL.Query().Get("session")
    _, err = riverClient.InsertTx(ctx, tx, AlmaBW{ TopCon: id, Filename: filename, Session: session}, nil)
    if err != nil { panic(err) }

    if err := tx.Commit(ctx); err != nil { panic(err) }
    w.Write([]byte("ok"))
  }
}
