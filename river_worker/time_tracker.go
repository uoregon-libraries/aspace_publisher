package river_worker

import (
    "fmt"
	"time"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jackc/pgx/v5"
    "context"
)

// example of a time string 2009-11-10 23:00:00 +0000 UTC
// returns true if ok to enqueue job
// false if unable to delete old time
// target time (as a string) for next attempt
// error if one occurred
func AddTime(pool *pgxpool.Pool, ctx context.Context) (bool, string, error) {
  tx, err := pool.Begin(ctx)
  defer tx.Rollback(ctx)
  addQuery := fmt.Sprintf("INSERT INTO tracker VALUES('%s')", NowString())
  count, target, err := removeTime(tx, ctx)
  if err != nil { return false, "", err }
  if count < 5 {
    _, err = tx.Exec(ctx, addQuery)
    if err != nil { return false, "", err }
    err = tx.Commit(ctx)
    if err != nil { return false, "", err }
    return true, "", nil
  } else {
    err = tx.Commit(ctx)
    if err != nil { return false, "", err }
    return false, target, nil
  }
}

func NowString()string{
  return time.Now().Round(0).Format(time.RFC3339)
}

func removeTime(tx pgx.Tx, ctx context.Context)(int, string, error){
  count, err := GetCount(tx, ctx)
  if err != nil { return count, "", err }
  if count == 0 { return count, "", nil }

  var timestr string
  var viewed int
  row := tx.QueryRow(ctx, "SELECT time,viewed FROM tracker ORDER BY created_at LIMIT 1")
  err = row.Scan(&timestr, &viewed)
  if err != nil { return count, "", err }
  deleteQuery := fmt.Sprintf("DELETE FROM tracker WHERE time = '%s'", timestr)
  updateQuery := fmt.Sprintf("UPDATE tracker SET viewed = %v WHERE time = '%s'", viewed+1, timestr)
  result, target, err := TimeExpired(timestr, viewed)
  if err != nil { return count, "", err }
  if result {
    _, err = tx.Exec(ctx, deleteQuery)
    if err != nil { return count, "", err }
    return count-1, "", nil 
  } else if count >= 5 {
    _, err = tx.Exec(ctx, updateQuery)
    if err != nil { return count, "", err }
  }
  return count, target.Format(time.RFC3339), nil
}

func TimeExpired(timestr string, viewed int)(bool, time.Time, error){
  _time, err := time.Parse(time.RFC3339, timestr)
  if err != nil {return false, _time, err }
  target := _time.Add(time.Hour + (time.Second * time.Duration(viewed * 3)))
  now := time.Now().Round(0)
  res := now.Compare(target)
  // res == -1 then target is still in the future
  // res == 1, then target is in the past
  if res > -1 { return true, target, nil }
  return false, target, nil
}

func GetCount(tx pgx.Tx, ctx context.Context)(int, error){
  query := "SELECT count(*) FROM tracker"
  var count int
  row := tx.QueryRow(ctx, query)
  //if err != nil { return 0, err }
  err := row.Scan(&count)
  if err != nil { return 0, err }
  return count, nil 
}
