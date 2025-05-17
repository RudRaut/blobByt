// src/components/WalletUI.tsx
import { ConnectButton, useCurrentAccount } from '@mysten/dapp-kit';
import { useEffect, useState } from 'react';
import { useSuiClient } from '@mysten/dapp-kit';

export default function WalletUI() {
	const account = useCurrentAccount();
	const client = useSuiClient();
	const [balance, setBalance] = useState<string | null>(null);

	useEffect(() => {
		const fetchBalance = async () => {
			if (!account?.address) return;
			const result = await client.getBalance({ owner: account.address });
			setBalance((+result.totalBalance / 1e9).toFixed(3)); // Show SUI in full units
		};
		fetchBalance();
	}, [account?.address]);

	return (
		<div className="flex justify-between items-center mb-4 p-4 bg-white shadow rounded-xl">
			<h1 className="text-xl font-bold">🧱 Sui Vault</h1>
			{account ? (
				<div className="flex items-center gap-4">
					<div className="text-sm">
						<p className="font-medium">{account.address.slice(0, 6)}...{account.address.slice(-4)}</p>
						<p className="text-gray-500 text-xs">{balance ?? 'Loading...'} SUI</p>
					</div>
				</div>
			) : (
				<ConnectButton />
			)}
		</div>
	);
}
