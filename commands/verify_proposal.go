package commands

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func NewVerifyProposalCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     getVerifyCommandUsage(),
		Example: getVerifyCommandExample(),
		Short:   VerifyProposalShortDesc,
		Long:    getVerifyProposalLongDesc(),
		Args:    cobra.ExactArgs(VerifyProposalCmdParamsCount),
		RunE: func(cmd *cobra.Command, args []string) error {
			inpus, err := NewVerifyProposalArgs(args)
			if err != nil {
				return err
			}

			bashCmd := exec.Command("/bin/bash", "verify_proposal.sh",
				inpus.smartContractsLocation, inpus.consumerChainId, inpus.multisigAddress,
				ConsumerBinary, CosmWasmBinary, inpus.toolOutputLocation, "true", // true for create output subdirectory
				inpus.proposalId, inpus.providerNodeId, ProviderBinary)

			RunCmdAndPrintOutput(bashCmd)

			return nil
		},
	}

	return cmd
}

func getVerifyCommandUsage() string {
	return fmt.Sprintf("%s [%s] [%s] [%s] [%s] [%s] [%s]",
		VerifyProposalCmdName, SmartContractsLocation, ConsumerChainId,
		MultisigAddress, ToolOutputLocation, ProposalId, ProviderNodeId)
}

func getVerifyCommandExample() string {
	return fmt.Sprintf("%s %s %s %s %s %s %s %s",
		ToolName, VerifyProposalCmdName, "$HOME/wasm_contracts", "wasm", "wasm1243cuuy98lxaf7ufgav0w76xt5es93afr8a3ya",
		"$HOME/tool_output_step2", "1", "tcp://localhost:26657")
}

func getVerifyProposalLongDesc() string {
	return fmt.Sprintf(VerifyProposalLongDesc, SmartContractsLocation, ConsumerChainId,
		MultisigAddress, ToolOutputLocation, ProposalId, ProviderNodeId)
}

type VerifyProposalArgs struct {
	smartContractsLocation string
	consumerChainId        string
	multisigAddress        string
	toolOutputLocation     string
	proposalId             string
	providerNodeId         string
}

func NewVerifyProposalArgs(args []string) (*VerifyProposalArgs, error) {
	if len(args) != VerifyProposalCmdParamsCount {
		return nil, fmt.Errorf("Unexpected number of arguments. Expected: %d, received: %d.", VerifyProposalCmdParamsCount, len(args))
	}

	commandArgs := new(VerifyProposalArgs)
	var errors []string

	smartContractsLocation := strings.TrimSpace(args[0])
	if IsValidInputPath(smartContractsLocation) {
		commandArgs.smartContractsLocation = smartContractsLocation
	} else {
		errors = append(errors, fmt.Sprintf("Provided input path '%s' is not a valid directory.", smartContractsLocation))
	}

	consumerChainId := strings.TrimSpace(args[1])
	if IsValidString(consumerChainId) {
		commandArgs.consumerChainId = consumerChainId
	} else {
		errors = append(errors, fmt.Sprintf("Provided chain-id '%s' is not valid.", consumerChainId))
	}

	multisigAddress := strings.TrimSpace(args[2])
	if IsValidString(multisigAddress) {
		commandArgs.multisigAddress = multisigAddress
	} else {
		errors = append(errors, fmt.Sprintf("Provided multisig address '%s' is not valid.", multisigAddress))
	}

	toolOutputLocation := strings.TrimSpace(args[3])
	if IsValidOutputPath(toolOutputLocation) {
		commandArgs.toolOutputLocation = toolOutputLocation
	} else {
		errors = append(errors, fmt.Sprintf("Provided output path '%s' is not a valid directory.", toolOutputLocation))
	}

	proposalId := strings.TrimSpace(args[4])
	if IsValidProposalId(proposalId) {
		commandArgs.proposalId = proposalId
	} else {
		errors = append(errors, fmt.Sprintf("Provided proposal id '%s' is not valid.", proposalId))
	}

	// TODO: not sure if we should validate node id with regex
	providerNodeId := strings.TrimSpace(args[5])
	if IsValidString(providerNodeId) {
		commandArgs.providerNodeId = providerNodeId
	} else {
		errors = append(errors, fmt.Sprintf("Provided provider node id '%s' is not valid.", providerNodeId))
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf(strings.Join(errors, "\n"))
	}

	return commandArgs, nil
}

const (
	VerifyProposalShortDesc = "Verify that genesis and binary hashes created from the provided inputs match the hashes from the 'create consumer chain' proposal with the given proposal ID"
	VerifyProposalLongDesc  = `This command takes the same inputs and goes through the same process as 'prepare-proposal' command to create the genesis.json file and calculate its hash.
It then queries the 'create consumer chain' proposal from the provider chain to obtain the hashes. If the hashes from the proposal match the recalculated ones, then the resulting genesis.json file contains the smart contracts provided to the input of this command.
Command arguments:
    %s - The location of the directory that contains CosmWasm smart contracts source code. TODO: add details about subdirectories structure and other things (Cargo.toml etc.)?
    %s - The chain ID of the consumer chain.
    %s - The multi-signature address that will have the permission to instantiate contracts from the set of predeployed codes.
    %s - The location of the directory where the resulting genesis.json and sha256hashes.json files will be saved.
    %s - The ID of the 'create consumer chain' proposal submitted to the provider chain, whose data will be used to verify if the inputs of this command match the ones from the proposal.
    %s - The address of the provider chain node in the following format: tcp://IP_ADDRESS:PORT_NUMBER. This address is used to query the provider chain to obtain the proposal information.`
)