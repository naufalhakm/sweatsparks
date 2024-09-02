package helpers

import "database/sql"

func CommitOrRollback(tx *sql.Tx) {
	err := recover()
	if err != nil {
		tx.Rollback()
		// PanicIfError(errRollback)
		// panic(err)
	} else {
		tx.Commit()
		// PanicIfError(errCommit)
	}

}
