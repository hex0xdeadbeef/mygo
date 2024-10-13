package main

import (
	"bytes"
	"fmt"
)

/*
	bytes.Buffer, I THOUGHT YOU WERE MY FRIEND
1. I want to build a string to use as a key for a cache, something like
	<type>:<client id>:<id>. But I don't want too many allocations.
2. The naive way is like to make many allocations and be slow. The next way
	is to use fmt.Sprintf(...) function.
3. The fmt.Sprintf(...) seems to be better, but there's a better way that
	is defined as bytes.Buffer using.
*/

func NaiveConcat(itemType, clientId, id string) string {
	return itemType + ":" + clientId + ":" + id
}

func FmtSprintfConcat(itemType, clientId, id string) string {
	return fmt.Sprintf("%s : %s : %s", itemType, clientId, id)
}

func BytesBufferConcat(itemType, clientId, id string) string {
	const columnsCount = 2

	var (
		b   = make([]byte, len(itemType)+len(clientId)+len(id)+columnsCount)
		buf = bytes.NewBuffer(b)
	)

	buf.WriteString(itemType)
	buf.WriteByte(':')

	buf.WriteString(clientId)
	buf.WriteByte(':')

	buf.WriteString(id)

	return buf.String()
}


