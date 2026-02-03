package river_worker

import (
  "context"
  "net/url"
  "github.com/jackc/pgx/v5/pgxpool"
  "github.com/jackc/pgx/v5"
  "github.com/riverqueue/river"

  "aspace_publisher/alma"
  "aspace_publisher/db"
  "net/http"
  "time"
  "fmt"
  "errors"
)
type LinkToNetwork struct {
  MmsID string `json:"mms_id"`
  Filename string `json:"filename"`
  NumTries int `json:"numtries"`
}

func (LinkToNetwork) Kind() string { return "link_to_network" }

type LTNWorker struct {
  river.WorkerDefaults[LinkToNetwork]
}

func (w *LTNWorker) Work(ctx context.Context, job *river.Job[LinkToNetwork]) error {
  var dbcred db.DBPool
  dbcred.Init()
  dbPool, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s/%s", dbcred.Pguser, dbcred.Pgpass, dbcred.Pgaddress, dbcred.Pgdb))
  if err != nil { return err }
  defer dbPool.Close()

  // check to see if it's safe to run job
  result, t, err := AddTime(dbPool, ctx)
  if err != nil { return err }
  //If AddTime fails, reschedule the job
  if !result {
    if job.Args.NumTries > 10 { return errors.New("Max attempts for job exceeded") }
    riverClient := river.ClientFromContext[pgx.Tx](ctx)
    newtime, err := getNewTime(t)
    if err != nil { return err }
    job.Args.NumTries += 1
    opts := &river.InsertOpts{ ScheduledAt: newtime }
    InsertWorker(riverClient, ctx, dbPool, job.Args, opts)
  } else {
    alma.DummyLinkToNetwork([]string{job.Args.MmsID}, job.Args.Filename)
  }
  return nil
}

func StartLTNJob(riverClient *river.Client[pgx.Tx], ctx context.Context, dbPool *pgxpool.Pool) http.HandlerFunc{
  return func(w http.ResponseWriter, r *http.Request) {
    _unescaped, err := url.QueryUnescape(r.URL.String())
    if err != nil { panic(err) }
    _url,_ := url.Parse(_unescaped)
    q := _url.Query()
    args := LinkToNetwork{ MmsID: q.Get("id"), Filename: q.Get("filename"), NumTries: 0 }
    InsertWorker(riverClient, ctx, dbPool, args, nil)
    w.Write([]byte("ok"))
  }
}

func getNewTime(timestr string)(time.Time, error){
  t, err := time.Parse(time.RFC3339, timestr)
  return t, err
}
