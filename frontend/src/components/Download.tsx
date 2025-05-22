import React, { useState } from 'react';

const Download: React.FC = () => {
  const [blobID, setBlobID] = useState('');
  const [status, setStatus] = useState('');

  const handleDownload = async () => {
    setStatus('Downloading...');
    try {
      const res = await fetch(`http://localhost:8080/download?blobID=${blobID}`);
      if (!res.ok) throw new Error('File not found');

      const blob = await res.blob();
      const fileName = res.headers.get('Content-Disposition')?.split('filename=')[1] || 'file';

      const link = document.createElement('a');
      link.href = URL.createObjectURL(blob);
      link.download = fileName;
      link.click();

      setStatus('Download successful');
    } catch (err) {
      setStatus('Download failed');
    }
  };

  return (
    <div className="p-4 border rounded-lg shadow-md mt-6">
      <h2 className="text-xl font-bold mb-2">Download File</h2>
      <input
        type="text"
        placeholder="Enter blobID"
        value={blobID}
        onChange={(e) => setBlobID(e.target.value)}
        className="mb-2 block border p-1 w-full"
      />
      <button onClick={handleDownload} className="bg-green-600 text-white px-4 py-2 rounded">
        Download
      </button>
      <p className="mt-2">{status}</p>
    </div>
  );
};

export default Download;
