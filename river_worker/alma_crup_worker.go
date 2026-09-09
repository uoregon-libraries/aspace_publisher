package river_worker

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jackc/pgx/v5"
    "github.com/riverqueue/river"

    "aspace_publisher/handlers" //packages worker needs
    "net/http"
)
type AlmaCrup struct {
    Resource string `json:"id"`
    Filename string `json:"filename"`
    Session string `json:"session"`
    Token string `json:"token"`
}

func (AlmaCrup) Kind() string { return "alma_crup" }

type AlmaCrupWorker struct {
    river.WorkerDefaults[AlmaCrup]
}

func (w *AlmaCrupWorker) Work(ctx context.Context, job *river.Job[AlmaCrup]) error {
    handlers.AlmaCrup(job.Args.Resource, job.Args.Filename, job.Args.Session, job.Args.Token)
    return nil
}

func StartAlmaCrupJob(riverClient *river.Client[pgx.Tx], ctx context.Context, dbPool *pgxpool.Pool) http.HandlerFunc{
  return func(w http.ResponseWriter, r *http.Request) {
    defer HandlerPanic()
    tx, err := dbPool.Begin(ctx)
    if err != nil { panic(err) }
    defer tx.Rollback(ctx)
    id := r.URL.Query().Get("id")
    filename := r.URL.Query().Get("filename")
    session := r.URL.Query().Get("session")
    token := r.URL.Query().Get("token")
    _, err = riverClient.InsertTx(ctx, tx, AlmaCrup{ Resource: id, Filename: filename, Session: session, Token: token}, nil)
    if err != nil { panic(err) }

    if err := tx.Commit(ctx); err != nil { panic(err) }
    w.Write([]byte("ok"))
  }
}
