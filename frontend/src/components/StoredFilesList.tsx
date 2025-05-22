// import React, { useEffect, useState } from 'react';

// interface FileMetadata {
//   blobID: string;
//   name: string;
//   size: number;
//   fileType: string;
//   encryptionKey: string;
//   epochs: number;
//   description?: string;
//   uploadTime: string;
// }

// const FileList: React.FC = () => {
//   const [files, setFiles] = useState<FileMetadata[]>([]);
//   const [error, setError] = useState<string | null>(null);

//   useEffect(() => {
//     fetch('http://localhost:8080/files')
//       .then((res) => res.json())
//       .then((data) => setFiles(data))
//       .catch(() => setError('Failed to fetch files'));
//   }, []);

//   return (
//     <div className="p-4 border rounded-lg shadow-md mt-6">
//       <h2 className="text-xl font-bold mb-4">Stored Files</h2>
//       {error && <p className="text-red-600">{error}</p>}
//       <div className="overflow-auto max-h-96">
//         <table className="w-full text-sm">
//           <thead>
//             <tr className="bg-gray-100 text-left">
//               <th className="p-2">Blob ID</th>
//               <th className="p-2">Name</th>
//               <th className="p-2">Size (bytes)</th>
//               <th className="p-2">Type</th>
//               <th className="p-2">Epochs</th>
//               <th className="p-2">Uploaded</th>
//             </tr>
//           </thead>
//           <tbody>
//           {files.map((file) => (
//               <tr key={file.blobID}>
//                 <td className="p-2">{file.blobID}</td>
//                 <td className="p-2">{file.name}</td>
//                 <td className="p-2">{file.size}</td>
//                 <td className="p-2">{file.fileType}</td>
//                 <td className="p-2">{file.epochs}</td>
//                 <td className="p-2">{new Date(file.uploadTime).toLocaleString()}</td>
//               </tr>
//             ))}
//           </tbody>
//         </table>
//       </div>
//     </div>
//   );
// };

// export default FileList;
import React, { useEffect, useState } from "react";
import axios from "axios";

interface FileMetadata {
  _id: string,
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
        setError(err?.response?.data?.message || "Failed to fetch files");
      } finally {
        setLoading(false);
      }
    };

    fetchFiles();

  }, []);

  if (loading) return <p className="text-gray-700 p-4">Loading files...</p>;
  if (error)
    return (
      <p className="text-red-600 p-4">
        Error loading files: <strong>{error}</strong>
      </p>
    );

  if (files.length === 0)
    return <p className="text-gray-500 p-4">No files stored yet.</p>;

  return (
    <div className="p-4">
      <h2 className="text-xl font-semibold mb-4">Stored Files</h2>
      <table className="min-w-full border border-gray-300 text-sm">
        <thead className="bg-gray-100">
          <tr>
            <th className="p-2 text-left">Blob ID</th>
            <th className="p-2 text-left">Name</th>
            <th className="p-2 text-left">Size</th>
            <th className="p-2 text-left">Type</th>
            <th className="p-2 text-left">Epochs</th>
            <th className="p-2 text-left">Uploaded At</th>
          </tr>
        </thead>
        <tbody>
            {files.map((file) => (
                <tr key={file._id} className="border-t">
                    <td className="p-2">{file._id}</td>
                    <td className="p-2">{file.blobID}</td>
                    <td className="p-2">{file.name}</td>
                    <td className="p-2">{file.size} bytes</td>
                    <td className="p-2">{file.fileType}</td>
                    <td className="p-2">{file.epochs}</td>
                    <td className="p-2">{file.description || "-"}</td>
                    <td className="p-2">{new Date(file.uploadTime).toLocaleString()}</td>
               {/* Do NOT render _id separately */}
                </tr>
           ))}

        </tbody>
      </table>
    </div>
  );
};

export default StoredFilesList;
