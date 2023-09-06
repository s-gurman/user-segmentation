#! /bin/sh

# prepare: deleting existing segments
curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{"name": "a1"}' > /dev/null
curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{"name": "a2"}' > /dev/null
curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{"name": "a3"}' > /dev/null
curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{"name": "a4"}' > /dev/null
curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{"name": "a5"}' > /dev/null

>&2 echo "================ success test cases ================"

curl -s -X 'POST' 'http://localhost:8081/api/segment' -d '{"name": "a1"}'; echo
curl -s -X 'POST' 'http://localhost:8081/api/segment' -d '{"name": "a2"}'; echo
curl -s -X 'POST' 'http://localhost:8081/api/segment' -d '{"name": "a3"}'; echo
curl -s -X 'POST' 'http://localhost:8081/api/segment' -d '{"name": "a4"}'; echo
curl -s -X 'POST' 'http://localhost:8081/api/segment' -d '{"name": "a5"}'; echo
echo

>&2 echo "================ failed test cases ================"

curl -s -X 'POST' 'http://localhost:8081/api/segment' -d '{"name": "a1"}'; echo
curl -s -X 'POST' 'http://localhost:8081/api/segment' -d '{"name": ""}'; echo
curl -s -X 'POST' 'http://localhost:8081/api/segment' -d '{"seg_name": ""}'; echo
curl -s -X 'POST' 'http://localhost:8081/api/segment' -d '{'; echo

# clear after test
curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{"name": "a1"}' > /dev/null
curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{"name": "a2"}' > /dev/null
curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{"name": "a3"}' > /dev/null
curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{"name": "a4"}' > /dev/null
curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{"name": "a5"}' > /dev/null