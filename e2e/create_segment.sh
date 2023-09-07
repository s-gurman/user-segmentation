#! /bin/bash

export $(grep -v '^#' ./config/.env | xargs -d '\n')

# prepare: deleting existing segments
for i in {1..8}; do
    curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{"name": "a'$i'"}' > /dev/null
done



echo; echo "=================== success test cases ==================="

# test: creating segment without autoadd_percent option
for i in {1..5}; do
    curl -s -X 'POST' 'http://localhost:8081/api/segment' -d '{"name": "a'$i'"}'
    echo
done

# prepare: adding segments to user -> trigger creates new user in `users` table
for i in {1001..1100}; do
    curl -s -X 'POST' 'http://localhost:8081/api/experiments/user/'$i -d '{"add_segments": ["a1"]}' > /dev/null
done
echo

docker exec $APP psql postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@\
$POSTGRES_ADDR/$POSTGRES_DB?sslmode=$POSTGRES_SSLMODE \
-c 'SELECT segments.slug, count(*) AS users_count FROM experiments JOIN segments '\
'ON segment_id = segments.id GROUP BY segments.slug;'

# test: creating segment with user autoadd_percent option (20%, 50%, 80%)
perc=20

for slug in a6 a7 a8; do
curl -s -X 'POST' 'http://localhost:8081/api/segment' -d '{
    "name": "'$slug'",
    "options": {
        "autoadd_percent": '$perc'
    }
}'; echo
    (( perc += 30))
done
echo

docker exec -it $APP psql postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@\
$POSTGRES_ADDR/$POSTGRES_DB?sslmode=$POSTGRES_SSLMODE \
-c 'SELECT segments.slug, count(*) AS users_count FROM experiments JOIN segments '\
'ON segment_id = segments.id GROUP BY segments.slug ORDER BY segments.slug;'



echo; echo "=================== failed test cases ==================="

curl -s -X 'POST' 'http://localhost:8081/api/segment' -d '{"name": "a1"}'  ; echo
curl -s -X 'POST' 'http://localhost:8081/api/segment' -d '{"name": ""}'    ; echo
curl -s -X 'POST' 'http://localhost:8081/api/segment' -d '{"seg_name": ""}'; echo
curl -s -X 'POST' 'http://localhost:8081/api/segment' -d '{'               ; echo

curl -s -X 'POST' 'http://localhost:8081/api/segment' -d '{
    "name":"a6",
    "options": {
        "autoadd_percent": 120
    }
}'; echo

curl -s -X 'POST' 'http://localhost:8081/api/segment' '{
    "name":"a6",
    "options": {
        "autoadd_percent": 120
    }
}'; echo

curl -s -X 'POST' 'http://localhost:8081/api/segment' '{
    "name":"a6",
    "options": {
        "autoadd_percent": "10"
    }
}'; echo



# clear after test
for i in {1..8}; do
    curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{"name": "a'$i'"}' > /dev/null
done
