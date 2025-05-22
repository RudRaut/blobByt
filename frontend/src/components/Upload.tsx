import React, { useState, useCallback } from 'react';
import toast from 'react-hot-toast';
import { FiUpload, FiFile, FiType, FiInbox } from 'react-icons/fi';

const Upload: React.FC = () => {
  const [file, setFile] = useState<File | null>(null);
  const [description, setDescription] = useState('');
  const [isUploading, setIsUploading] = useState(false);
  const [isDragging, setIsDragging] = useState(false);
  const [response, setResponse] = useState('');

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFile(e.target.files?.[0] || null);
  };

  const handleDragOver = useCallback((e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(false);
  }, []);

  const handleDrop = useCallback((e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(false);
    
    const droppedFile = e.dataTransfer.files[0];
    if (droppedFile) {
      setFile(droppedFile);
      toast.success('File selected: ' + droppedFile.name);
    }
  }, []);

  const handleUpload = async () => {
    if (!file) {
      toast.error('Please select a file first');
      return;
    }

    setIsUploading(true);
    const formData = new FormData();
    formData.append('file', file);
    formData.append('description', description);

    try {
      const res = await fetch('http://localhost:8080/upload', {
        method: 'POST',
        body: formData,
      });

      const result = await res.json();
      toast.success('File uploaded successfully!');
      setDescription('');
      setFile(null);
    } catch (err) {
      toast.error('Upload failed. Please try again.');
    } finally {
      setIsUploading(false);
    }
  };

  return (
    <div className="h-full">
      <h2 className="text-xl sm:text-2xl font-bold mb-4 text-[#37454d] flex items-center gap-2">
        <FiUpload className="text-[#76a0bd]" /> Upload File
      </h2>
      <div className="space-y-4">
        <div 
          className={`border-2 border-dashed rounded-lg p-2 sm:p-3 text-center transition-colors ${
            isDragging 
              ? 'border-[#76a0bd] bg-[#f1ede7]' 
              : 'border-[#999898] hover:border-[#76a0bd]'
          }`}
          style={{ minHeight: '80px' }}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
        >
          <FiInbox className="mx-auto text-3xl sm:text-4xl text-[#76a0bd] mb-2" />
          <p className="text-sm sm:text-base text-[#37454d] mb-2">Drag and drop your file here</p>
          <p className="text-xs sm:text-sm text-[#999898] mb-4">or</p>
          <input 
            type="file" 
            onChange={handleFileChange} 
            className="hidden"
            id="file-input"
          />
          <label 
            htmlFor="file-input"
            className="bg-[#76a0bd] text-white px-4 sm:px-6 py-2 rounded hover:bg-[#37454d] transition-colors cursor-pointer inline-block text-sm sm:text-base"
          >
            Browse Files
          </label>
          {file && (
            <p className="mt-4 text-xs sm:text-sm text-[#37454d] flex items-center justify-center gap-2">
              <FiFile className="text-[#76a0bd]" /> Selected: {file.name}
            </p>
          )}
        </div>
        <div className="flex flex-col">
          <label className="text-sm sm:text-base text-[#37454d] mb-2 text-left">
            Description
          </label>
          <input
            type="text"
            placeholder="Enter file description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            className="p-2 text-sm sm:text-base border border-[#999898] rounded focus:border-[#76a0bd] focus:outline-none"
          />
        </div>
        <button 
          onClick={handleUpload} 
          disabled={isUploading}
          className="bg-[#76a0bd] text-white px-4 sm:px-6 py-2 rounded hover:bg-[#37454d] transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 w-full text-sm sm:text-base"
        >
          <FiUpload className={isUploading ? 'animate-bounce' : ''} />
          {isUploading ? 'Uploading...' : 'Upload'}
        </button>
      </div>
    </div>
  );
};

export default Upload;
