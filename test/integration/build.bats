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
