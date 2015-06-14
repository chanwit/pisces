#!/usr/bin/env bats

load ../helpers

# TODO make this configurable
export GOPATH=/home/chanwit/projects/pisces
export GOROOT=/opt/go

function go() {
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

	export DOCKER_NODE_0_ADDR=${HOSTS[0]}
	export TESTDATA

	run go test github.com/chanwit/pisces/swarm
	echo $output

	[ "$status" -eq 0 ]
	[[ "${lines[0]}" == "ok"* ]]
}
