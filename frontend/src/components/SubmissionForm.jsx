import React, { useState } from 'react';

const API_URL = 'http://localhost:8080/api/notifications';
const CATEGORIES = ["Sports", "Finance", "Movies"];

export default function SubmissionForm({ onSent }) {
  const [category, setCategory] = useState(CATEGORIES[0]);
  const [message, setMessage] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [status, setStatus] = useState({ type: '', text: '' });

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    // Client-side validation required by challenge
    if (!message.trim()) {
      setStatus({ type: 'error', text: 'Message cannot be empty.' });
      return;
    }

    setIsSubmitting(true);
    setStatus({ type: '', text: '' });

    try {
      const response = await fetch(API_URL, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ category, message })
      });

      if (!response.ok) {
        throw new Error(`Server returned status ${response.status}`);
      }

      setStatus({ type: 'success', text: 'Message dispatched successfully!' });
      setMessage(''); // Reset ONLY message, retain category for ease of use
      
      if (onSent) onSent();
      
    } catch (err) {
      setStatus({ type: 'error', text: `Failed to dispatch: ${err.message}` });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="glass-panel">
      <h2 style={{ marginBottom: '1.5rem', fontSize: '1.25rem' }}>Send Notification</h2>
      
      {status.text && (
        <div className={`alert ${status.type}`}>
          {status.text}
        </div>
      )}

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="category">Category</label>
          <select 
            id="category" 
            value={category} 
            onChange={(e) => setCategory(e.target.value)}
          >
            {CATEGORIES.map(cat => (
              <option key={cat} value={cat}>{cat}</option>
            ))}
          </select>
        </div>

        <div className="form-group">
          <label htmlFor="message">Message</label>
          <textarea 
            id="message" 
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            placeholder="Type your notification message here..."
          />
        </div>

        <button type="submit" className="btn" disabled={isSubmitting}>
          {isSubmitting ? 'Dispatching...' : 'Dispatch Message'}
        </button>
      </form>
    </div>
  );
}
