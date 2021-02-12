#!/bin/bash

kill $(ps | egrep 'BFTWithoutSigna' | awk '{print $1}')
