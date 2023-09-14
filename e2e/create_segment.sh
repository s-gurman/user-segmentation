#! /bin/bash

export $(grep -v '^#' ./config/.env | xargs -0)

function create_segment {
    suffix=''
    if [ $# -eq 2 ]; then suffix=$',"options":{"autoadd_percent":'$2'}'; fi

    body='{"name":"'$1'"'$suffix'}'
    echo "request:  POST /api/segment $body"
    echo "response: $(curl -s -X 'POST' 'http://localhost:8081/api/segment' -d "$body")"
    echo
}

function print_segments {
    docker exec -it $APP \
psql postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_ADDR/$POSTGRES_DB?sslmode=$POSTGRES_SSLMODE -c \
'SELECT * FROM segments;'
}

# prepare: creating 10'000 users
docker exec -it $APP \
psql postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_ADDR/$POSTGRES_DB?sslmode=$POSTGRES_SSLMODE -c \
'INSERT INTO users (id) VALUES (generate_series(10001, 20000));' > /dev/null

# prepare: deleting existing segments (a1, a2, ..., a9)
for i in {1..9}; do
    curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{"name": "a'$i'"}' > /dev/null
done



echo -e '\n============================== TESTING ==============================\n'
echo 'endpoint: /api/segment [POST]'



echo -e '\n============================== SUCCESS ==============================\n'

print_segments

# test: creating segments (a1, a2, ..., a5) without autoadd_percent option
for i in {1..5}; do
    create_segment 'a'$i
done

print_segments

docker exec -it $APP \
psql postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_ADDR/$POSTGRES_DB?sslmode=$POSTGRES_SSLMODE -c \
'SELECT count(*) AS user_count FROM users;'

# test: creating segments (a6, a7, a8) with user autoadd_percent option (20%, 50%, 80%)
percent=20

for name in a6 a7 a8; do
    create_segment "$name" "$percent"
    (( percent += 30))
done

docker exec -it $APP \
psql postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_ADDR/$POSTGRES_DB?sslmode=$POSTGRES_SSLMODE -c \
'SELECT segments.slug, count(*) AS user_count FROM experiments '\
'JOIN segments ON segment_id = segments.id '\
'GROUP BY segments.slug '\
'ORDER BY segments.slug;'



echo -e '============================== FAILED ==============================\n'

for body in '{' '{"foo": "bar"}' '{"name": 1234}'; do
    echo "request: POST /api/segment $body"
    echo "response: $(curl -s -X 'POST' 'http://localhost:8081/api/segment' -d "$body")"
    echo
done
create_segment ''
create_segment 'a1'
create_segment 'a9' 120
create_segment 'a9' -10.5
create_segment 'a9' -10,5
create_segment 'a9' '"10.5"'



# clear after test
for i in {1..9}; do
    curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{"name": "a'$i'"}' > /dev/null
done

# clear after test
docker exec -it $APP \
psql postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_ADDR/$POSTGRES_DB?sslmode=$POSTGRES_SSLMODE -c \
'DELETE FROM users WHERE id >= 10001 AND id <= 20000;' > /dev/null
