#!/bin/bash
function monitor() {
    inotifywait --exclude '\.git' -q -m -r -e moved_to $1 |
        while read line; do
            restart
        done
    }

# Terminate and rerun the main Go program
function restart () {
    echo ">> Reloading..."
    if [ "$(pidof $PROCESS_NAME)" ]; then
        killall -q -w -9 $PROCESS_NAME
    fi
    echo ">> Reloading..."
    eval "go run $ARGS &"
}

# Make sure all background processes get terminated
function close () {
    killall -q -w -9 inotifywait
    exit 0
}

trap close INT
echo "== Go-reload"

WATCH_ALL=false
while getopts ":a" opt; do
    case $opt in
        a)
            WATCH_ALL=true
            ;;
        \?)
            echo "Invalid option: -$OPTARG" >&2
            exit 0
            ;;
    esac
done

shift "$((OPTIND - 1))"

FILE_NAME=$(basename $1)
PROCESS_NAME=${FILE_NAME%%.*}

ARGS=$@

# Start the main Go program
echo ">> Watching directories, CTRL+C to stop"
eval "go run $ARGS &"

monitor $PWD $WATCH_ALL

wait
