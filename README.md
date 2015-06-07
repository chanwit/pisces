Pisces
======

`Pisces` is a Fig clone that understands Docker Swarm.

Problem
-------

The Docker Swarm project was started several months after Fig.
Their design rationale do not match.
Fig is built for the stand-alone Docker, while Swarm is a clustering system.
This causes incompatibilty between them at the level that
users cannot build, scale or link containers on their cluster.
Making Fig to be fully compatible with Swarm takes time, but our system needs it today.

Composition Redefined
---------------------

`Pisces` is carefully designed and test against Swarm since the beginning.
It consumes the same `yml` format of Fig.
`Pisces` shares the same integration tests, taken from Swarm.
It is to ensure that `Pisces` will not break good things done by Swarm.

Unlike Fig, `Pisces` does not and will not run against the stand-alone Docker.
`Pisces` works only with a cluster.
This means that `Pisces` will not replace Fig.
It's a clustering counter-part of Fig from the beginning.
Developers, who are familiar with Fig, can just use `Pisces` along side Fig.

Requirements
------------

`Pisces` requires a proper cluster setup done using `docker-machine`.
It also requires the `docker` client executable.
Connection must be done via TLS.

FAQ
---

  * Why don't you contribute to Fig or Docker Compose rather than implement a new tool?

  I'd love to. But I just do not write Python. That's all.
