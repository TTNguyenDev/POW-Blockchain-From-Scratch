#!/bin/bash

# Add execute permissions 
chmod +x ./scripts/build.sh

./scripts/build.sh

WALLET_1=$(./main createwallet)
echo "Wallet 1 address is: $WALLET_1"
./main sendcoins -from $CENTRAL_NODE -to $WALLET_1 -amount 10

WALLET_2=$(./main createwallet)
echo "Wallet 2 address is: $WALLET_2"
./main sendcoins -from $CENTRAL_NODE -to $WALLET_2 -amount 10

./main startnode

BALANCE_1=$(./main getbalance $WALLET_1)
echo "Wallet 1 balance is: $BALANCE_1"
BALANCE_2=$(./main getbalance $WALLET_2)
echo "Wallet 2 balance is: $BALANCE_2"
