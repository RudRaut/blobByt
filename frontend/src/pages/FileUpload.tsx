import { FormEvent, useState } from "react";
import { WALRUS_API, PACKAGE_ID, MODULE_NAME, FILE_REGISTRY_ID } from "../../constants.ts";
import { useNetworkVariable } from "../networkConfig";
import { useSignAndExecuteTransaction, useSuiClient, useWallets } from '@mysten/dapp-kit';

import toast, { Toaster } from "react-hot-toast";
import { fromBase64, toHex, SUI_CLOCK_OBJECT_ID } from "@mysten/sui/utils";
import { Transaction } from "@mysten/sui/transactions";
import { randomBytes } from '@noble/hashes/utils';
import axios from "axios";

export default function FileUpload() {
    const [status, setStatus] = useState("");
    const [file, setFile] = useState<File | null>(null);
    const [epochs, setEpochs] = useState<number>(1);
    const [loading, setLoading] = useState(false);
    const vaultPackageId = useNetworkVariable("vaultPackageId");
    const fileRegistryId = useNetworkVariable("fileRegistryId");
    const suiClient = useSuiClient();
    const {
        mutate: signAndExecute,
    } = useSignAndExecuteTransaction();

    
    function generateFileId(): string {
        const bytes = randomBytes(32);
        return '0x' + toHex(bytes);
    }

    const checkServerHealth = async (): Promise<boolean> => {
        try {
          const res = await axios.get(`${WALRUS_API}/health`, { timeout: 3000 });
          return res.status === 200;
        } catch (error: any) {
          // Axios in browser throws a generic Network Error on connection refusal
          if (error.message === "Network Error") {
            console.error("Server connection failed: Server is likely down.");
          } else if (error.code === "ECONNABORTED") {
            console.error("Health check request timed out.");
          } else {
            console.error("Unexpected health check error:", error);
          }
          return false;
        }
      };

    const handleUpload = async (e: FormEvent) => {
        e.preventDefault();

        console.log(file?.name, file?.size, file?.type);

        if(!file) {
            alert("Please select a file");
            return;
        }

        setLoading(true);

        const isServerUp = await checkServerHealth();
        if (!isServerUp) {
            toast.error("Backend server is unavailable.");
            setLoading(false);
            return;
        }

        try {

            const formData = new FormData();
            formData.append("file", file);
            formData.append("epochs", epochs.toString());
            toast.success("Uploading File to walrus");

            const response = await axios.put(
                `${WALRUS_API}/upload?epochs=${epochs}`,
                formData,
                {
                    headers: {
                        "Content-Type": "multipart/form-data",
                    },
                }
            );


            const { blobID, key } = response.data;
            if (!blobID || !key) {
                toast.error("Invalid upload response: Missing blobID or key.");
                setLoading(false);
                return;
            }
            console.log(blobID, key);
            setStatus("upload successful");

            let blobAddress = '';
            try {
                // Convert Base64URL to standard Base64
                const standardBase64 = blobID.replace(/-/g, '+').replace(/_/g, '/');
            
                // Decode to bytes
                const blobBytes = fromBase64(standardBase64);
            
                // Ensure it's 32 bytes
                if (blobBytes.length !== 32) {
                    throw new Error("Decoded blobID is not 32 bytes long.");
                }
            
                // Convert to hex with '0x' prefix
                blobAddress = '0x' + toHex(blobBytes);
            } catch (err) {
                console.error("Failed to decode blobID:", err);
                toast.error("Failed to process blobID from upload.");
                setLoading(false);
                return;
            }
            
            setStatus("uploading file to Sui");
                if (!vaultPackageId || !fileRegistryId) {
                    toast.error("Missing network configuration. Check vaultPackageId or fileRegistryId.");
                    setLoading(false);
                    return;
                }

                const fileId = generateFileId();

                const tx = new Transaction();
                tx.moveCall({
                    arguments: [
                        tx.object(fileRegistryId),
                        tx.pure.address(fileId),
                        tx.pure.address(blobAddress),
                        tx.object(SUI_CLOCK_OBJECT_ID),
                    ],
                    target: `${vaultPackageId}::${MODULE_NAME}::upload_file`,
                });
            
                signAndExecute(
                    {
                        transaction: tx,
                        chain: 'sui:testnet',
                    },
                    {
                        onSuccess: async ({ digest }) => {
                            try {
                                const { effects } = await suiClient.waitForTransaction({
                                    digest,
                                    options: { showEffects: true },
                                });
                                if (!effects) {
                                    toast.error("No transaction effects returned.");
                                    setLoading(false);
                                    return;
                                }
                                toast.success(`File registered on-chain. Digest: ${digest}`);
                                setFile(null);
                            } catch (error) {
                                toast.error(
                                    `Transaction failed: ${error instanceof Error ? error.message : String(error)}`
                                );
                            } finally {
                                setLoading(false);
                            }
                        }
                    }
                );
        }

        


        catch (error: any) {
            if (error.response) {
              // Server responded but with a status code ≠ 2xx
              console.log("Status:", error.response.status);
              console.log("Response data:", error.response.data);
            } else if (error.request) {
              // Request was sent, but no response
              console.log("No response received:", error.request);
            } else {
              // Something went wrong setting up the request
              console.log("Axios setup error:", error.message);
            }
          }
    };
    return (
        <div className="max-w-xl mx-auto mt-10 p-6 bg-white rounded-xl shadow space-y-4">
            <h2 className="text-xl font-semibold">📤 Upload File</h2>
            <form onSubmit={handleUpload} className="space-y-4">
                <input
                    type="file"
                    onChange={(e) => setFile(e.target.files?.[0] || null)}
                    className="block w-full text-sm file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
                />
                
                {/* 🔹 Added: input field for epochs */}
                <input
                    type="number"
                    value={epochs}
                    onChange={(e) => setEpochs(Number(e.target.value))}
                    placeholder="Enter number of epochs"
                    className="block w-full text-sm border rounded px-3 py-2"
                    min={1}
                />

                <button
                    type="submit"
                    className="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded"
                >
                    Upload & Register
                </button>
            </form>
            {status && <p className="text-sm text-gray-700">{status}</p>}
        </div>
    );
}