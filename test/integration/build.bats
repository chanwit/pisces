#!/usr/bin/env bats

load helpers

function teardown() {
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
	run pisces build web
	[[ ${status} -eq 0 ]]

	echo $output

	[[ "${lines[0]}" == "build testdata_web" ]]
}
