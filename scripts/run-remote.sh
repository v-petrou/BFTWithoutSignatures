#!/bin/bash

MACHINE=$1

N=10
CLIENTS=25
SCEN=0

if [ $MACHINE -eq 2 ]
then
	go install BFTWithoutSignatures
	go install BFTWithoutSignatures_Client
	go run BFTWithoutSignatures generate_keys $N
fi

ID=$MACHINE
echo "STARTED $ID"
go run BFTWithoutSignatures $ID $N $CLIENTS $SCEN 1 &

if [ $MACHINE -eq 2 ] || [ $MACHINE -eq 3 ]
then
	TEMP=(($MACHINE%2))
	for (( ID=(($TEMP+8)); ID<$N; ((ID+=2)) ))
	do
		echo "STARTED $ID"
		go run BFTWithoutSignatures $ID $N $CLIENTS $SCEN 1 &
	done
fi