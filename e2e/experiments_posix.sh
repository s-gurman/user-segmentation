#! /bin/bash

export $(grep -v '^#' ./config/.env | xargs -0)

function update_user_experiments {
    suffix=''
    if [ $# -eq 4 ]; then suffix=$',"options":{"deletion_time":"'$4'"}'; fi
    
    body='{"add_segments":'$2',"delete_segments":'$3$suffix'}'
    echo "request:  POST /api/experiments/user/$1 $body"
    echo "response: $(curl -s -X 'POST' 'http://localhost:8081/api/experiments/user/'$1 -d "$body")"
    echo
}

function get_user_experiments {
    echo "timestamp: $(date '+%F %T')"
    echo "request:  GET /api/experiments/user/$1"
    echo "response: $(curl -s -X 'GET' 'http://localhost:8081/api/experiments/user/'$1)"
    echo
}

function print_experiments {
    docker exec -it $APP \
psql postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_ADDR/$POSTGRES_DB?sslmode=$POSTGRES_SSLMODE -c \
'SELECT user_id, segments.slug, started_at, expired_at FROM experiments '\
'JOIN segments ON segment_id = segments.id '\
'WHERE user_id = '$1' '\
'ORDER BY segments.slug;'
}

# prepare: creating new segments (a1, a2, ..., a10)
for i in {1..10}; do
    curl -s -X 'POST' 'http://localhost:8081/api/segment' -d '{"name": "a'$i'"}' > /dev/null
done



echo -e '\n============================== TESTING ==============================\n'
echo 'endpoint: /api/experiments/user/{user_id=1234} [POST]'



echo -e '\n============================== SUCCESS ==============================\n'

print_experiments 1234

get_user_experiments 1234
update_user_experiments 1234 '["a1","a2","a3","a4","a5"]' '[]'
sleep 2; print_experiments 1234

get_user_experiments 1234
update_user_experiments 1234 '[]' '["a1","a2"]'
sleep 2; print_experiments 1234

get_user_experiments 1234
ts=$(date -v +1y '+%F %T')
update_user_experiments 1234 '["a4","a5","a6","a7"]' '["a3","a4", "a5"]' "$ts"
print_experiments 1234; sleep 2

get_user_experiments 1234
ts=$(date -v +2S '+%F %T')
update_user_experiments 1234 '["a8","a9", "a10"]' '[]' "$ts"
print_experiments 1234; sleep 3

get_user_experiments 1234
ts=$(date -v +1y '+%F %T')
update_user_experiments 1234 '["a8","a9", "a10"]' '[]' "$ts"
print_experiments 1234; sleep 2

get_user_experiments 1234
update_user_experiments 1234 '["a1","a2","a3"]' '["a8","a9","a10"]' "$ts"
print_experiments 1234; sleep 2

get_user_experiments 1234
update_user_experiments 1234 '["a8","a9","a10"]' '["a1","a2","a3","a4","a5"]'
print_experiments 1234; sleep 2

get_user_experiments 1234
curl -s -X 'GET' 'http://localhost:8081/api/experiments/user/1234'; echo



echo -e '\n============================== FAILED ==============================\n'

update_user_experiments 1234 '[]' '[]'
update_user_experiments 1234 '["non-exist segment"]' '[]'
update_user_experiments 1234 '[]' '["non-exist segment"]'
update_user_experiments 1234 '[]' '["a5"]'
update_user_experiments 1234 '["a5","a6"]' '[]'
update_user_experiments 1234 '["a1","a2","a3"]' '[]' "bad time layout"
update_user_experiments 1234 '["a1","a2","a3"]' '[]' "$(date -d '-1 day' '+%F %T')"



#clear after test
for i in {1..10}; do
    curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{"name": "a'$i'"}' > /dev/null
done
