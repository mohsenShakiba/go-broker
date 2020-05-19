package storage

import (
	"bytes"
	"fmt"
	"strconv"
)

type indexRow struct {
	msgId   string
	deleted bool
	start   int
	length  int
}

func deserializeIndex(b []byte) *indexRow {
	parts := bytes.Split(b, []byte(","))

	if len(parts) != 4 {
		return nil
	}

	deleted, _ := atob(string(parts[0]))
	msgId := string(parts[1])
	start, _ := strconv.Atoi(string(parts[2]))
	length, _ := strconv.Atoi(string(parts[3]))

	return &indexRow{
		msgId:   msgId,
		deleted: deleted,
		start:   start,
		length:  length,
	}

}

func (r *indexRow) serializeIndex() []byte {

	return []byte(fmt.Sprintf("%d,%s,%d,%d\n", btoi(r.deleted), r.msgId, r.start, r.length))
}

func btoi(b bool) int {
	deletedInt := 0
	if b {
		deletedInt = 1
	}
	return deletedInt
}

func atob(s string) (bool, error) {
	i, err := strconv.Atoi(s)

	if err != nil {
		return false, err
	}

	return i == 1, nil
}
