#!/bin/sh

# build the run-params depending on the environment variables
# (don't change the port, it should be changed via docker port mappings)
if [ $BULB_NAME ]; then
    PARAMS="$PARAMS --bulbName=$BULB_NAME"
fi
if [ $JOBS_FILE_PATH ]; then
    PARAMS="$PARAMS --jobsFilePath=$JOBS_FILE_PATH"
fi
if [ $API_KEY ]; then
    PARAMS="$PARAMS --apiKey=$API_KEY"
fi

# Start the go-wrapper
go-wrapper run $PARAMS

