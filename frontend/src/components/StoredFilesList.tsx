import React, { useEffect, useState } from "react";
import axios from "axios";
import toast from 'react-hot-toast';
import { FiList, FiInfo } from 'react-icons/fi';

interface FileMetadata {
  _id: string;
  blobID: string;
  name: string;
  size: number;
  fileType: string;
  encryptionKey: string;
  epochs: number;
  description?: string;
  uploadTime: string;
}

const StoredFilesList: React.FC = () => {
  const [files, setFiles] = useState<FileMetadata[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchFiles = async () => {
      try {
        const response = await axios.get("http://localhost:8080/files");
        setFiles(response.data);
        setError(null);
      } catch (err: any) {
        const errorMessage = err?.response?.data?.message || "Failed to fetch files";
        setError(errorMessage);
        toast.error(errorMessage);
      } finally {
        setLoading(false);
      }
    };

    fetchFiles();
  }, []);

  if (loading) return (
    <div className="h-full">
      <p className="text-sm sm:text-base text-[#37454d] flex items-center gap-2">
        <FiList className="animate-spin text-[#76a0bd]" /> Loading files...
      </p>
    </div>
  );

  if (error) return (
    <div className="h-full">
      <p className="text-sm sm:text-base text-red-600 flex items-center gap-2">
        <FiInfo className="text-red-600" /> Error loading files: <strong>{error}</strong>
      </p>
    </div>
  );

  if (files.length === 0) return (
    <div className="h-full">
      <p className="text-sm sm:text-base text-[#37454d] flex items-center gap-2">
        <FiList className="text-[#76a0bd]" /> No files stored yet.
      </p>
    </div>
  );

  return (
    <div className="h-full">
      <h2 className="text-xl sm:text-2xl font-bold mb-4 text-[#37454d] flex items-center gap-2">
        <FiList className="text-[#76a0bd]" /> Stored Files
      </h2>
      <div className="overflow-x-auto">
        <table className="min-w-full border border-[#999898] text-xs sm:text-sm">
          <thead className="bg-[#f1ede7]">
            <tr>
              <th className="p-2 sm:p-3 text-left text-[#37454d] border-b border-[#999898]">Blob ID</th>
              <th className="p-2 sm:p-3 text-left text-[#37454d] border-b border-[#999898]">Name</th>
              <th className="p-2 sm:p-3 text-left text-[#37454d] border-b border-[#999898]">Size</th>
              <th className="p-2 sm:p-3 text-left text-[#37454d] border-b border-[#999898]">Type</th>
              <th className="p-2 sm:p-3 text-left text-[#37454d] border-b border-[#999898]">Epochs</th>
              <th className="p-2 sm:p-3 text-left text-[#37454d] border-b border-[#999898]">Description</th>
              <th className="p-2 sm:p-3 text-left text-[#37454d] border-b border-[#999898]">Uploaded At</th>
            </tr>
          </thead>
          <tbody>
            {files.map((file) => (
              <tr key={file._id} className="border-b border-[#999898] hover:bg-[#f1ede7]">
                <td className="p-2 sm:p-3 text-[#37454d]">{file.blobID}</td>
                <td className="p-2 sm:p-3 text-[#37454d]">{file.name}</td>
                <td className="p-2 sm:p-3 text-[#37454d]">{file.size} bytes</td>
                <td className="p-2 sm:p-3 text-[#37454d]">{file.fileType}</td>
                <td className="p-2 sm:p-3 text-[#37454d]">{file.epochs}</td>
                <td className="p-2 sm:p-3 text-[#37454d]">{file.description || "-"}</td>
                <td className="p-2 sm:p-3 text-[#37454d]">{new Date(file.uploadTime).toLocaleString()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default StoredFilesList;
