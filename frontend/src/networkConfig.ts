import { getFullnodeUrl } from "@mysten/sui/client";
import { createNetworkConfig } from "@mysten/dapp-kit";
import {  PACKAGE_ID, FILE_REGISTRY_ID } from "../constants.ts";


const { networkConfig, useNetworkVariable, useNetworkVariables } =
  createNetworkConfig({
    devnet: {
      url: getFullnodeUrl("devnet"),
      variables: {
        vaultPackageId: "Todo",
        fileRegistryId: "todo",
      }
    },
    testnet: {
      url: getFullnodeUrl("testnet"),
      variables: {
        vaultPackageId: PACKAGE_ID,
        fileRegistryId: FILE_REGISTRY_ID,
      }
    },
    mainnet: {
      url: getFullnodeUrl("mainnet"),
      variables: {
        vaultPackageId: "todo",
        fileRegistryId: "todo", //Need to be defined or causes error
      }
    },
  });

export { useNetworkVariable, useNetworkVariables, networkConfig };
