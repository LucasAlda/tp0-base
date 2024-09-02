docker run --network tp0_testing_net -it alpine sh -c '
 RESPONSE=$(echo "ping" | nc server:12345 -w 3) 
 if [ "$RESPONSE" == "ping" ]; then
   echo "action: test_echo_server | result: success"
 else
   echo "action: test_echo_server | result: fail"
 fi
'