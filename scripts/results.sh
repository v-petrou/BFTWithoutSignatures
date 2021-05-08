#!/bin/bash

INFO=$( cat $(ls ~/logs/out/0_*) | grep 'ID:0' | cut -d" " -f4- | cut -d"|" -f2- )

TEMP=$( cd ~/logs/client && cat $(ls) | grep 'Operation Latency:' | cut -d" " -f4 )
OP="${TEMP// /$'\n'}"
AVG_OP=$( echo "$OP" | awk -vx=0 '{x += $1} END {print x/NR}' )

TEMP=$( cd ~/logs/out && cat $(ls) | grep 'Message Complexity:' | cut -d" " -f3 )
COM="${TEMP// /$'\n'}"
AVG_COM=$( echo "$COM" | awk -vx=0 '{x += $1} END {print x/NR}' )

TEMP=$( cd ~/logs/out && cat $(ls) | grep 'Message Size:' | cut -d" " -f3 )
SIZE="${TEMP// /$'\n'}"
AVG_SIZE=$( echo "$SIZE" | awk -vx=0 '{x += $1} END {print x/NR}' )

cd ~

TEMP=$( echo "$INFO" | cut -d":" -f2 | cut -d"|" -f1 )
FILE=$( echo "${TEMP}_results.txt" | tr -d " " )

echo -e "------------------------------------------------------------------" >> $FILE
echo "$INFO" >> $FILE

echo -e "\nOperation Latency\n-----------------" >> $FILE
echo "$OP" >> $FILE
echo -e "Average Operation Latency: $AVG_OP s" >> $FILE

echo -e "\nMessage Complexity\n------------------" >> $FILE
echo "$COM" >> $FILE
echo -e "Average Message Complexity: $AVG_COM msgs" >> $FILE

echo -e "\nMessage Size\n------------" >> $FILE
echo "$SIZE" >> $FILE
echo -e "Average Message Size: $AVG_SIZE MB" >> $FILE

echo -e "------------------------------------------------------------------\n" >> $FILE
