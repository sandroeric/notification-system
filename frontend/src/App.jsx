import React, { useState } from 'react';
import SubmissionForm from './components/SubmissionForm';
import LogHistory from './components/LogHistory';
import './index.css';

function App() {
  const [refreshHash, setRefreshHash] = useState(Date.now());

  const handleNotificationSent = () => {
    // changing hash forces LogHistory to refresh
    setRefreshHash(Date.now());
  };

  return (
    <div className="container">
      <header className="header">
        <h1>Command Center</h1>
        <p>Notification Strategy Multi-Channel Dispatcher</p>
      </header>
      
      <main className="grid">
        <section>
          <SubmissionForm onSent={handleNotificationSent} />
        </section>
        
        <section>
          <LogHistory refreshHash={refreshHash} />
        </section>
      </main>
    </div>
  );
}

export default App;
