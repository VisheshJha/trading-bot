#!/bin/bash

RESPONSE=$(curl --write-out '%{http_code}' --silent --output /dev/null http://localhost:8080/health)

if [ "$RESPONSE" -ne 200 ]; then
    echo "Health check failed: HTTP $RESPONSE"
    exit 1
fi

exit 0