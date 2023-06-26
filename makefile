url_shortener:
	@echo "Building url_shortener"
	docker-compose -p url \
	-f docker-compose.yml up -d --build 
	@echo "Done"

url_shortener_reids:
	@echo "Building url_shortener wirh redis"
	docker-compose -p url \
	-f docker-compose.yml \
	-f docker-compose.redis.yml up -d --build 
	@echo "Done"

url_shortener_kgs:
	@echo "Building url_shortener wirh kgs"
	docker-compose -p url \
	-f docker-compose.yml \
	-f docker-compose.kgs.yml up -d --build 
	@echo "Done"

url_shortener_redis_kgs:
	@echo "Building url_shortener wirh redis and kgs"
	docker-compose -p url \
	-f docker-compose.yml \
	-f docker-compose.redis.yml \
	-f docker-compose.kgs.yml up -d --build 
	@echo "Done"
url_shortener_redis_kgs_ratelimit:
	@echo "Building url_shortener wirh redis and kgs"
	docker-compose -p url \
	-f docker-compose.yml \
	-f docker-compose.redis.yml \
	-f docker-compose.kgs.yml \
	-f docker-compose.ratelimit.yml up -d --build 
	@echo "Done"

stop_url_shortener_redis_kgs:
	@echo "Stopping url_shortener wirh redis and kgs"
	docker-compose -p url \
	-f docker-compose.yml \
	-f docker-compose.redis.yml \
	-f docker-compose.kgs.yml down
	@echo "Done"

stop_url_shortener_redis:
	@echo "Stopping url_shortener wirh redis"
	docker-compose -p url \
	-f docker-compose.yml \
	-f docker-compose.redis.yml down
	@echo "Done"

stop_url_shortener_kgs:
	@echo "Stopping url_shortener wirh kgs"
	docker-compose -p url \
	-f docker-compose.yml \
	-f docker-compose.kgs.yml down
	@echo "Done"

stop_url_shortener:
	@echo "Stopping url_shortener"
	docker-compose -p url \
	-f docker-compose.yml down
	@echo "Done"

insert_loading_test:
	k6 run --config test/load_test_config.json test/loading_test.js
get_loading_test:
	k6 run --config test/load_test_config.json test/get_loading_test.js