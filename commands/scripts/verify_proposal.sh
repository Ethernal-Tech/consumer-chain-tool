#!/bin/bash
set -eu

PREPARE_INPUTS_SCRIPT="$0"
TOOL_INPUT="$1"
CONSUMER_CHAIN_ID="$2"
CONSUMER_CHAIN_MULTISIG_ADDRESS="$3"
CONSUMER_CHAIN_BINARY="$4"
WASM_BINARY="$5"
TOOL_OUTPUT="$6"
CREATE_OUTPUT_SUBFOLDER="$7"
PROPOSAL_GENESIS_HASH="$8"
PROPOSAL_BINARY_HASH="$9"
PROPOSAL_SPAWN_TIME="${10}"

# Delete all generated data.
 clean_up () {
   rm -f "$TOOL_OUTPUT"/proposal_info.json
 } 
 trap clean_up EXIT

if [ "$CREATE_OUTPUT_SUBFOLDER" = "true" ]; 
then
  TOOL_OUTPUT="$TOOL_OUTPUT"/$(date +"%Y-%m-%d_%H-%M-%S")
fi
LOG="$TOOL_OUTPUT"/log_file.txt

# Create directories if they don't exist.
mkdir -p "$TOOL_OUTPUT"

echo "Generating files and hashes for validation..."
if ! bash -c "$PREPARE_INPUTS_SCRIPT" "$TOOL_INPUT" "$CONSUMER_CHAIN_ID" $CONSUMER_CHAIN_MULTISIG_ADDRESS $CONSUMER_CHAIN_BINARY $WASM_BINARY "$TOOL_OUTPUT" $PROPOSAL_SPAWN_TIME;
then
  echo "Error while preparing proposal data! Verify proposal failed. Please check the $LOG for more details."
  exit 1
fi

echo "Validating genesis and binary hashes..."
GENESIS_HASH=$(jq -r ".genesis_hash" "$TOOL_OUTPUT"/sha256hashes.json)
BINARY_HASH=$(jq -r ".binary_hash" "$TOOL_OUTPUT"/sha256hashes.json)

valid=true  

if [ "$GENESIS_HASH" != "$PROPOSAL_GENESIS_HASH" ]; then
  echo "Recalculated genesis hash does not match the one from the proposal!"
  valid=false
fi

if [ "$BINARY_HASH" != "$PROPOSAL_BINARY_HASH" ]; then
  echo "Recalculated binary hash does not match the one from the proposal!"
  valid=false
fi

if [ "$valid" = true ]; then
  echo "Genesis and binary hashes are correct! Verify proposal succeded."
else
  echo "Verify proposal failed."
  exit 1
fi