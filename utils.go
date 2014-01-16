package xmpp

import (
    "time"
    "strconv"
)

func generate_random_id() string {
    nano := time.Now().UnixNano()
    return strconv.FormatInt(nano, 16)
}
