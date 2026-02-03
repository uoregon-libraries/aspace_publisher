package river_worker

import(
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jackc/pgx/v5"
    "github.com/riverqueue/river"
)

func InsertWorker(riverClient *river.Client[pgx.Tx], ctx context.Context, dbPool *pgxpool.Pool, args river.JobArgs, opts *river.InsertOpts){

    tx, err := dbPool.Begin(ctx)
    if err != nil {
        panic(err)
    }
    defer tx.Rollback(ctx)
    _, err = riverClient.InsertTx(ctx, tx, args, opts)
    if err != nil {
        panic(err)
    }

    if err := tx.Commit(ctx); err != nil {
        panic(err)
    }
}

