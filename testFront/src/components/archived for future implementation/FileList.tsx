import React, { useEffect, useState } from 'react';

interface FileListProps {
  onSelectFile: (fileName: string) => void;
}

const FileList: React.FC<FileListProps> = ({ onSelectFile }) => {
  const [files, setFiles] = useState<string[]>([]);

  useEffect(() => {
    fetch('http://localhost:8082/files')
      .then(response => response.json())
      .then(data => setFiles(data));
  }, []);

  return (
    <div>
      <h2>Select a Document to Edit</h2>
      <ul>
        {files.map(file => (
          <li key={file} onClick={() => onSelectFile(file)} style={{ cursor: 'pointer', color: 'blue' }}>
            {file}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default FileList;
