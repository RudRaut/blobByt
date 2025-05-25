# BlobByt

BlobByt is a BlobByt decentralized file storage system designed to securely upload, encrypt, store, and retrieve files using client-side encryption and decentralized blob storage. It ensures privacy, integrity, and future extensibility with blockchain-based access control and cryptographic key sharing.

## Tech Stack

- Frontend: React, TypeScript, Tailwind CSS, Vite
- Backend: Go (Golang)
- Database: MongoDB
- Storage: Walrus (decentralized blob storage)
- Encryption: AES-256-GCM (client-side)

## Features

- AES-256-GCM encryption performed entirely in the browser before upload.
- Encrypted files are stored via Walrus and only retrievable with proper decryption.
- Metadata (blob ID, file name, size, type, encryption key, upload time, epochs, description) stored in MongoDB.
- File listing interface with metadata visualization.
- Secure file download and client-side decryption.
- Proper error handling, upload/download feedback, and status indication.

## Getting Started

### Prerequisites

- Node.js & npm
- Go
- MongoDB (running locally or remotely)

### Backend Setup

1. Clone the repository.
2. Navigate to the backend directory.
3. Update MongoDB URI in your environment or code.
4. Run the backend:
   ```bash
   go run main.go
   ```

### Frontend Setup

1. Navigate to the frontend directory.
2. Install dependencies:
   ```bash
   npm install
   ```
3. Start the development server:
   ```bash
   npm run dev
   ```

## File Flow

1. User selects a file in the React interface.
2. File is encrypted in the browser using AES-256-GCM.
3. The encrypted blob is uploaded to Walrus.
4. Metadata including blob ID and encryption key is stored in MongoDB via Go backend.
5. User can view uploaded files and securely download and decrypt them.

## Security

- Files are never uploaded in plaintext.
- Encryption keys are generated per file and stored only with the backend.
- Decryption happens in the browser upon download.

## Future Development

BlobByt will evolve into a fully decentralized access-controlled file system by integrating the Sui blockchain and Move smart contracts. Access control and permissions will be managed on-chain, enabling verifiable ownership and shared access rights. We will also implement X25519-based key exchange for secure recipient-side decryption, and Ed25519 digital signatures for authenticating file authorship. Logging mechanisms will be added to record file access and download events transparently. These additions aim to create a fully decentralized, privacy-preserving file-sharing platform suitable for real-world applications.

## License

BlobByt is open-source and available under the MIT License.

Built with security, privacy, and decentralization in mind.
