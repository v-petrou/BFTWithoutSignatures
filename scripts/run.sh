#!/bin/bash

N=10
T=0
CLIENTS=2
SCEN=0

go install BFTWithoutSignatures

BFTWithoutSignatures generate_keys $N

for (( ID=0; ID<$N; ID++ ))
do
	BFTWithoutSignatures $ID $N $T $CLIENTS $SCEN &
done
