#!/usr/bin/env bats

load helpers

@test "swarm version" {

	run swarm -v
	echo $output

	[ "$status" -eq 0 ]
	[[ ${lines[0]} =~ version\ [0-9]+\.[0-9]+\.[0-9]+ ]]

}
