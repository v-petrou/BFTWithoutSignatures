#!/bin/bash

kill $(ps | egrep 'tests.test|go' | awk '{print $1}')
