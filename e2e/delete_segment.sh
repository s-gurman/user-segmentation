#! /bin/bash

# prepare: creating new segments
for i in {1..5}; do
    curl -s -X 'POST' 'http://localhost:8081/api/segment' -d '{"name": "a'$i'"}' > /dev/null
done



echo; echo "=================== success test cases ==================="

# test: deleting segment
for i in {1..5}; do
    curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{"name": "a'$i'"}'
    echo
done



echo; echo "=================== failed test cases ==================="

curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{"name": "a1"}'  ; echo
curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{"name": "a99"}' ; echo
curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{"name": ""}'    ; echo
curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{"seg_name": ""}'; echo
curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{'               ; echo
