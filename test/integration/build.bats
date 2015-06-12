#!/usr/bin/env bats

load helpers

function teardown() {
	swarm_manage_cleanup
	stop_docker
}

@test "pisces build" {
	start_docker_with_busybox 2
	swarm_manage

	run pisces build web
	[[ ${status} -eq 0 ]]

	[[ "${lines[0]}" == "build web" ]]
}

@test "pisces build no DOCKER_HOST" {
	start_docker_with_busybox 1
	swarm_manage

	run "${PISCES_BINARY}" build web

	# no DOCKER_HOST defined, error should be 1
	[[ ${status} -eq 1 ]]
	# should have some error message
	[[ ${output} == *"Environment variable \"DOCKER_HOST\" is required."* ]]
}
