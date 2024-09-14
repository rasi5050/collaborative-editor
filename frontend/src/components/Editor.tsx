import React, { useEffect, useState } from 'react';

interface EditorProps {
  docId: string;
}

const Editor: React.FC<EditorProps> = ({ docId }) => {
  const [content, setContent] = useState<string>('');
  const [lastEdited, setLastEdited] = useState<string>('');
  const [userCount, setUserCount] = useState<number>(1);
  const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    // Establish WebSocket connection
    // const socket = new WebSocket(`ws://localhost:8082/ws?doc_id=${docId}`);
    const socket = new WebSocket(`ws://192.168.229.31:8082/ws?doc_id=${docId}`);
    setWs(socket);

    // Handle incoming messages
    socket.onmessage = (event) => {
      const update = JSON.parse(event.data);
      if (update.content !== undefined) {
        setContent(update.content);
      }
      if (update.lastEdited !== undefined) {
        setLastEdited(new Date(update.lastEdited).toLocaleString()); // Format the timestamp
      }
      if (update.userCount !== undefined) {
        setUserCount(update.userCount);
      }
    };

    // Clean up on component unmount
    return () => {
      socket.close();
    };
  }, [docId]);

  // Handle text area changes
  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newContent = e.target.value;
    setContent(newContent);

    // Send the updated content to the server
    if (ws) {
      ws.send(JSON.stringify({ content: newContent }));
    }
  };

  return (
    <div>
      <p>Number of users online: {userCount}</p> {/* Display the user count */}
      <textarea
        value={content}
        onChange={handleChange}
        style={{ width: '100%', height: '70vh', fontSize: '16px' }}
      />
      <p>Last edited: {lastEdited}</p> {/* Display the last edited time */}
    </div>
  );
};

export default Editor;
