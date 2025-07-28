# Pre-Requisite

# Starting reports-service locally
    * ci/local/start.sh

# Stoping service
    * you can kill the existing process on the terminal which will issue a 'stop' command and service will cleaup and stop in about 30 seconds
    * Alternatively run 'netstat -vanp tcp | grep 8090' and then run 'kill -9 <PID>' - If you can't wait for 30 seconds

# Notes
    reports-service serves the VSM dashboard with data for reports from Opensearch.

