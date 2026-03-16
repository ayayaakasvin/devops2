#!/bin/bash
number=0
while true; do
 curl -s -o /dev/null -w "%{http_code} %{time_total} $number\n" https://lab7.local/api/health
 sleep 0.5
 ((number++))
done