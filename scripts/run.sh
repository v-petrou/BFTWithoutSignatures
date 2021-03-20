#!/bin/bash

N=10
CLIENTS=2
SCEN=0
REM=0

go install BFTWithoutSignatures

BFTWithoutSignatures generate_keys $N

for (( ID=0; ID<$N; ID++ ))
do
	BFTWithoutSignatures $ID $N $CLIENTS $SCEN $REM &
done
