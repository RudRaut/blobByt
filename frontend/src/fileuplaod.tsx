import { FormEvent, useState } from "react";
import { WALRUS_API, PACKAGE_ID, MODULE_NAME, FILE_REGISTRY_ID } from '../../config';
import { useSignAndExecuteTransaction, useWallets } from '@mysten/dapp-kit';
import { Transaction } from '@mysten/sui/transactions';
import { fromBase64, toHex } from '@mysten/sui/utils';

import toast, { Toaster } from "react-hot-toast";
import ClipLoader from "react-spinners/ClipLoader";

export default function FileUpload() {
    const [file, setFile] = useState<File | null>(null);
    const [loading, setLoading] = useState(false);
    const currentWallet = useWallets();
	const {mutate: signAndExecute } = useSignAndExecuteTransaction();

    const checkServerHealth = async (): Promise<boolean> => {
        try {
            const res = await fetch(`${WALRUS_API}/health`, { method: "GET" });
            return res.ok;
        } catch {
            return false; 
        }
    };

    const handleUpload = async (e: FormEvent) => {
        e.preventDefault();
        if(!file) {
            toast.error("Please select a file.");
            return;
        }

        if(!currentWallet || currentWallet.length === 0) {
            toast.error("Please connect your Wallet");
        }

        setLoading(true);
        toast.loading("Uploading file to Walrus...");

        const isServerUp = await checkServerHealth();

        if(!isServerUp) {
            toast.error("Backend server is unavailable.");
            setLoading(false);
            return;
        }

        toast("Uploading File to Walrus...");

        const formData = new FormData();
        formData.append("file", file);

        const res = await fetch(`${WALRUS_API}/upload`, {
            method: "PUT",
            body: formData,  
        });

        if(!res.ok) {
            toast.error("Upload to Walrus failed");
            return;
        }

        const { blobID, key } = await res.json();
        console.log("Encrypted AES key: ", key);
        toast.success("File uploaded to Walrus")

        const blobBytes = fromBase64(blobID);
        const blobAddress = '0x' + toHex(blobBytes);

        toast("uploading file to Sui");

        const tx = new Transaction();
        tx.moveCall({
            target: `${PACKAGE_ID}::${MODULE_NAME}::upload_file`,
            arguments: [
                tx.object(FILE_REGISTRY_ID), 
                tx.pure.address(blobAddress),
                tx.object('ox6'),
            ]
        });

        signAndExecute(
            {
                transaction: tx,
                chain: 'sui:testnet', // for use with testnet
            },
            {
                onSuccess: (res) => {
					toast.success(`File registered on-chain. Digest: ${res.digest}`);
					setLoading(false);
					setFile(null);                  
                },
                onError: (err) => {
                    console.error(err);
                    toast.error("On chain Transaction failed.");
                },
            }
        );
    };

    return (
        <div className="max-w-xl mx-auto mt-10 p-6 bg-white rounded-xl shadow space-y-4">
        <Toaster position="top-right" />
        <h2 className="text-xl font-semibold">📤 Upload File</h2>
        <form onSubmit={handleUpload} className="space-y-4">
          <input
            type="file"
            onChange={(e) => setFile(e.target.files?.[0] || null)}
            className="block w-full text-sm file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
            disabled={loading}
          />
          <button
            type="submit"
            disabled={loading}
            className="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded flex items-center justify-center gap-2 disabled:opacity-50"
          >
            {loading && <ClipLoader size={18} color="white" />}
            Upload & Register
          </button>
        </form>
      </div>
    );
}