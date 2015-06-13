#!/usr/bin/env bats

load ../helpers

function go() {
	export GOPATH=/home/chanwit/projects/pisces
	export GOROOT=/opt/go
	export DOCKER_HOST=$SWARM_HOST
	"${GOROOT}/bin/go" "$@"
}

function teardown() {
	swarm_manage_cleanup
	stop_docker
}

@test "go test: conf" {
	run go test github.com/chanwit/pisces/conf
	echo $output

	[ "$status" -eq 0 ]
	[[ "${lines[0]}" == "ok"* ]]
}


@test "go test: swarm" {
	start_docker 1
	swarm_manage

	run go test github.com/chanwit/pisces/swarm
	echo $output

	[ "$status" -eq 0 ]
	[[ "${lines[0]}" == "ok"* ]]
}
