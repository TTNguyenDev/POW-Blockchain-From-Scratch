
#!/bin/bash

# Add execute permissions 
chmod +x ./scripts/clean.sh
chmod +x ./scripts/build.sh

# Build
./scripts/clean.sh
./scripts/build.sh

# Run 
CENTRAL_NODE=$(./main createwallet)
echo "Your wallet is: $CENTRAL_NODE"

## New blockchain
./main newblockchain -address $CENTRAL_NODE
## Check balance 
BALANCE=$(./main getbalance -address $CENTRAL_NODE)

# Check if the balance is 0
if [ "$BALANCE" -eq 0 ]
then
    echo "Error: Balance is zero"
else
    echo "Balance of $CENTRAL_NODE are $BALANCE"

    # # Copy the entire folders with the genesis block to start node_3001
    #
    # WALLET_1=$(./main createwallet)
    # echo "Wallet 1 address is: $WALLET_1"
    # ./main sendcoins -from $CENTRAL_NODE -to $WALLET_1 -amount 10
    #
    # WALLET_2=$(./main createwallet)
    # echo "Wallet 2 address is: $WALLET_2"
    # ./main sendcoins -from $CENTRAL_NODE -to $WALLET_2 -amount 10
    #
    # ./main startnode
fi
