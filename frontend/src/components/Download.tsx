import React, { useState } from 'react';
import toast from 'react-hot-toast';
import { FiDownload, FiHash } from 'react-icons/fi';

const Download: React.FC = () => {
  const [blobID, setBlobID] = useState('');
  const [isDownloading, setIsDownloading] = useState(false);

  const handleDownload = async () => {
    if (!blobID.trim()) {
      toast.error('Please enter a Blob ID');
      return;
    }

    setIsDownloading(true);
    try {
      const res = await fetch(`http://localhost:8080/download?blobID=${blobID}`);
      if (!res.ok) throw new Error('File not found');

      const blob = await res.blob();
      const fileName = res.headers.get('Content-Disposition')?.split('filename=')[1] || 'file';

      const link = document.createElement('a');
      link.href = URL.createObjectURL(blob);
      link.download = fileName;
      link.click();

      toast.success('File downloaded successfully!');
      setBlobID('');
    } catch (err) {
      toast.error('Download failed. Please check the Blob ID and try again.');
    } finally {
      setIsDownloading(false);
    }
  };

  return (
    <div className="h-full">
      <h2 className="text-xl sm:text-2xl font-bold mb-4 text-[#37454d] flex items-center gap-2">
        <FiDownload className="text-[#76a0bd]" /> Download File
      </h2>
      <div className="space-y-4">
        <div className="flex flex-col">
          <label className="text-sm sm:text-base text-[#37454d] mb-2 flex items-center gap-2">
            <FiHash className="text-[#76a0bd]" /> Blob ID
          </label>
          <input
            type="text"
            placeholder="Enter blobID"
            value={blobID}
            onChange={(e) => setBlobID(e.target.value)}
            className="p-2 text-sm sm:text-base border border-[#999898] rounded focus:border-[#76a0bd] focus:outline-none"
          />
        </div>
        <button 
          onClick={handleDownload}
          disabled={isDownloading}
          className="bg-[#76a0bd] text-white px-4 sm:px-6 py-2 rounded hover:bg-[#37454d] transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 w-full text-sm sm:text-base"
        >
          <FiDownload className={isDownloading ? 'animate-bounce' : ''} />
          {isDownloading ? 'Downloading...' : 'Download'}
        </button>
      </div>
    </div>
  );
};

export default Download;
