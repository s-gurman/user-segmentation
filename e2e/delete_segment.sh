#! /bin/bash

export $(grep -v '^#' ./config/.env | xargs -0)

function delete_segment {
    body='{"name":"'$1'"}'
    echo "request:  DELETE /api/segment $body"
    echo "response: $(curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d "$body")"
    echo
}

function print_segments {
    docker exec -it $APP \
psql postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_ADDR/$POSTGRES_DB?sslmode=$POSTGRES_SSLMODE -c \
'SELECT * FROM segments;'
}

# prepare: creating new segments (a1, a2, ..., a5)
for i in {1..5}; do
    curl -s -X 'POST' 'http://localhost:8081/api/segment' -d '{"name": "a'$i'"}' > /dev/null
done



echo -e '\n============================== TESTING ==============================\n'
echo 'endpoint: /api/segment [DELETE]'

echo -e '\n============================== SUCCESS ==============================\n'

print_segments

# test: deleting segments (a1, a2, ..., a5)
for i in {1..5}; do
    delete_segment 'a'$i
done

print_segments



echo -e '============================== FAILED ==============================\n'

for body in '{' '{"foo": "bar"}' '{"name": 1234}'; do
    echo "request:  $body"
    echo "response: $(curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d "$body")"
    echo
done
delete_segment ''
delete_segment 'a1'
delete_segment 'a99'
