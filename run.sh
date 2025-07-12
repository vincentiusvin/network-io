#! /bin/sh

if [ "$1" != "blocking" ] && [ "$1" != "epoll" ] && [ "$1" != "uring" ]; then
  echo "Usage: ./run.sh <variant>"
  echo "<variant> can be blocking, epoll, or uring"
  exit 1
fi

VARIANT=$1
OUT_BIN=bin/$VARIANT
IN_FILE=$VARIANT/main.go

go build -o $OUT_BIN $IN_FILE
./$OUT_BIN&
PID=$!
(cd e2e && go test)
kill $PID
