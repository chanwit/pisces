#!/usr/bin/env bats

load helpers

function teardown() {
	swarm_manage_cleanup
	stop_docker
}

@test "pisces up" {
	start_docker_with_busybox 2
	swarm_manage

	run pisces up -d web
	[[ ${status} -eq 0 ]]


	[[ "${lines[0]}" == "up -d web" ]]
}
