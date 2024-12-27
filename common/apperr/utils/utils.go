package utils

func IsPsqlDuplicateKeyErrorMessage(errorString string) bool {
	if len(errorString) > 52 && errorString[:53] == "ERROR: duplicate key value violates unique constraint" {
		return true
	}
	return false
}
