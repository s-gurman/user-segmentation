#! /bin/bash

export $(grep -v '^#' ./config/.env | xargs -d '\n')

# prepare: creating new segments
for i in {1..10}; do
    curl -s -X 'POST' 'http://localhost:8081/api/segment' -d '{"name": "a'$i'"}' > /dev/null
done



echo; echo "================ success test cases ================"

curl -s -X 'POST' 'http://localhost:8081/api/experiments/user/1234' -d '{
  "add_segments": ["a1","a2","a3","a4","a5"],
  "delete_segments": []
}'; echo

curl -s -X 'POST' 'http://localhost:8081/api/experiments/user/1234' -d '{
  "add_segments": [],
  "delete_segments": ["a1","a2","a3"]
}'; echo

ts=$(date -d '+1 day' '+%F %T')
curl -s -X 'POST' 'http://localhost:8081/api/experiments/user/1234' -d '{
  "add_segments": ["a4","a5","a6","a7"],
  "delete_segments": ["a4","a5"],
  "options":{"deletion_time": "'"$ts"'"}
}'; echo

ts=$(date -d '+2 second' '+%F %T')
curl -s -X 'POST' 'http://localhost:8081/api/experiments/user/1234' -d '{
  "add_segments": ["a8","a9", "a10"],
  "delete_segments": [],
  "options":{"deletion_time": "'"$ts"'"}
}'; echo
sleep 2

ts=$(date -d '+1 day' '+%F %T')
curl -s -X 'POST' 'http://localhost:8081/api/experiments/user/1234' -d '{
  "add_segments": ["a8","a9", "a10"],
  "delete_segments": [],
  "options":{"deletion_time": "'"$ts"'"}
}'; echo

curl -s -X 'POST' 'http://localhost:8081/api/experiments/user/1234' -d '{
  "add_segments": ["a1","a2","a3"],
  "delete_segments": ["a8","a9","a10"],
  "options":{"deletion_time": "'"$ts"'"}
}'; echo
sleep 2

curl -s -X 'POST' 'http://localhost:8081/api/experiments/user/1234' -d '{
  "add_segments": ["a8","a9","a10"],
  "delete_segments": ["a1","a2","a3","a4","a5"]
}'; echo
echo

docker exec -it $APP psql postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@\
$POSTGRES_ADDR/$POSTGRES_DB?sslmode=$POSTGRES_SSLMODE \
-c 'SELECT user_id, segments.slug, started_at, expired_at FROM experiments '\
'JOIN segments ON segment_id = segments.id;'
sleep 2

echo timestamp: $(date '+%F %T')
curl -s -X 'GET' 'http://localhost:8081/api/experiments/user/1234'; echo



echo; echo "================ failed test cases ================"

curl -s -X 'POST' 'http://localhost:8081/api/experiments/user/1234' -d '{
  "add_segments": [],
  "delete_segments": []
}'; echo

curl -s -X 'POST' 'http://localhost:8081/api/experiments/user/1234' -d '{
  "add_segments": ["a1","a2","a3"],
  "delete_segments": []
  "options":{"deletion_time": "bad time layout"}
}'; echo

ts=$(date -d '-1 day' '+%F %T')
curl -s -X 'POST' 'http://localhost:8081/api/experiments/user/1234' -d '{
  "add_segments": ["a1","a2", "a3"],
  "delete_segments": [],
  "options":{"deletion_time": "'"$ts"'"}
}'; echo

curl -s -X 'POST' 'http://localhost:8081/api/experiments/user/1234' -d '{
  "add_segments": ["non-exist segment"],
  "delete_segments": []
}'; echo

curl -s -X 'POST' 'http://localhost:8081/api/experiments/user/1234' -d '{
  "add_segments": [],
  "delete_segments": ["non-exist segment"]
}'; echo

curl -s -X 'POST' 'http://localhost:8081/api/experiments/user/1234' -d '{
  "add_segments": [],
  "delete_segments": ["a5"]
}'; echo

curl -s -X 'POST' 'http://localhost:8081/api/experiments/user/1234' -d '{
  "add_segments": ["a5","a6"],
  "delete_segments": []
}'; echo



#clear after test
for i in {1..10}; do
    curl -s -X 'DELETE' 'http://localhost:8081/api/segment' -d '{"name": "a'$i'"}' > /dev/null
done
