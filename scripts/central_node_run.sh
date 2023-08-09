#!/bin/bash

WALLET_1=$(./main createwallet)
echo "Wallet 1 address is: $WALLET_1"
./main sendcoins -from $CENTRAL_NODE -to $WALLET_1 -amount 10

WALLET_2=$(./main createwallet)
echo "Wallet 2 address is: $WALLET_2"
./main sendcoins -from $CENTRAL_NODE -to $WALLET_2 -amount 10

./main startnode
