#!/usr/bin/env bats

load helpers

function teardown() {
	docker_swarm rmi -f testdata_web
	swarm_manage_cleanup
	stop_docker
}

@test "pisces build no DOCKER_HOST" {
	run "${PISCES_BINARY}" build web

	# no DOCKER_HOST defined, error should be 1
	[[ ${status} -eq 1 ]]
	# should have some error message
	[[ ${output} == *"Environment variable \"DOCKER_HOST\" is required."* ]]
}

@test "pisces build" {
	start_docker 1
	swarm_manage

	cd $TESTDATA
	EXPECTED=$(pisces build web)
	[[ ${status} -eq 0 ]]

	restart_swarm_manage

	ACTUAL=$(docker_swarm images -a | grep testdata_web | awk '{print $3}')
	[[ $EXPECTED == $ACTUAL ]]
}

@test "pisces build --no-cache" {
	start_docker 1
	swarm_manage

	cd $TESTDATA
	EXPECTED=$(pisces build --no-cache web)

	restart_swarm_manage

	ACTUAL=$(docker_swarm images -a | grep testdata_web | awk '{print $3}')
	[[ $EXPECTED == $ACTUAL ]]
}

@test "pisces build: many nodes" {
	start_docker 2
	swarm_manage

	# pre-condition, image count must be 0
	IMAGE_COUNT=$(docker_swarm images -a | grep testdata_web | wc -l)
	[[ ${IMAGE_COUNT} -eq 0 ]]

	cd $TESTDATA
	EXPECTED=$(pisces build --no-cache web | sort)

	restart_swarm_manage

	ACTUAL=$(docker_swarm images -a | grep testdata_web | awk '{print $3}' | sort)
	[[ $EXPECTED == $ACTUAL ]]
}

@test "pisces build: many services" {
	start_docker 2
	swarm_manage

	# pre-condition, image count must be 0
	IMAGE_COUNT=$(docker_swarm images -a | grep testdata_ | wc -l)
	[[ ${IMAGE_COUNT} -eq 0 ]]

	cd $TESTDATA
	EXPECTED=$(pisces build web front | sort)
	echo ">> EXPECTED"
	echo "$EXPECTED"

	restart_swarm_manage

	ACTUAL=$(docker_swarm images -a | grep testdata_ | awk '{print $3}' | sort)
	ACTUAL_COUNT=$(echo "$ACTUAL" | wc -l)
	echo ">> ACTUAL"
	echo "$ACTUAL"
	echo ">> ACTUAL_COUNT: $ACTUAL_COUNT"

	[[ $ACTUAL_COUNT -eq 4 ]]
	[[ $EXPECTED == $ACTUAL ]]
}
