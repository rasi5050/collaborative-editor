import React from 'react';
import './App.css';
import Editor from './components/Editor';

function App() {
  const documentId = 'single-doc';  // Static single document

  return (
    <div className="App">
      <header className="App-header">
        <h1>Collaborative Document Editor</h1>
      </header>
      <main>
        <Editor docId={documentId} />
      </main>
    </div>
  );
}

export default App;
