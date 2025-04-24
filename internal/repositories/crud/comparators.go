package crud

import "fmt"

func Equals(field string) string {
	return fmt.Sprintf("`%s` = ?", field)
}

func NotEquals(field string) string {
	return fmt.Sprintf("`%s` != ?", field)
}

func Like(field string) string {
	return fmt.Sprintf("`%s` LIKE ?", field)
}

func NotLike(field string) string {
	return fmt.Sprintf("`%s` NOT LIKE ?", field)
}

func Greater(field string) string {
	return fmt.Sprintf("`%s` > ?", field)
}

func GreaterEq(field string) string {
	return fmt.Sprintf("`%s` >= ?", field)
}

func Less(field string) string {
	return fmt.Sprintf("`%s` < ?", field)
}

func LessEq(field string) string {
	return fmt.Sprintf("`%s` <= ?", field)
}

