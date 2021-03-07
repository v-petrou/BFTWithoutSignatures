#!/bin/bash

TEST=TestABroadcast

N=10
T=0
CLIENTS=1
SCEN=0

go install BFTWithoutSignatures

BFTWithoutSignatures generate_keys $N

for (( ID=0; ID<$N; ID++ ))
do
	go test -v -run $TEST /home/vasilis/go/src/BFTWithoutSignatures/tests -args $ID $N $T $CLIENTS $SCEN &
done
