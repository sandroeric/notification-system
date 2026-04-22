import React, { useState, useEffect } from 'react';

const API_URL = 'http://localhost:8080/api/notifications/log';

export default function LogHistory({ refreshHash }) {
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchLogs = async () => {
      setLoading(true);
      setError(null);
      try {
        const response = await fetch(API_URL);
        if (!response.ok) {
          throw new Error(`Failed to fetch logs: ${response.status}`);
        }
        const data = await response.json();
        // The backend returns the list natively sorted from newest to oldest via SQL.
        setLogs(data || []);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };
    fetchLogs();
  }, [refreshHash]);

  return (
    <div className="glass-panel" style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <h2 style={{ marginBottom: '1.5rem', fontSize: '1.25rem' }}>Delivery History</h2>
      
      {error && (
        <div className="alert error">
          Error retrieving logs: {error}
        </div>
      )}

      <div style={{ overflowX: 'auto', flexGrow: 1 }}>
        <table style={{ width: '100%', borderCollapse: 'collapse', textAlign: 'left' }}>
          <thead>
            <tr style={{ borderBottom: '1px solid var(--border-color)' }}>
              <th style={{ padding: '1rem', color: 'var(--text-muted)', fontWeight: 600 }}>Time</th>
              <th style={{ padding: '1rem', color: 'var(--text-muted)', fontWeight: 600 }}>User Config</th>
              <th style={{ padding: '1rem', color: 'var(--text-muted)', fontWeight: 600 }}>Category</th>
              <th style={{ padding: '1rem', color: 'var(--text-muted)', fontWeight: 600 }}>Channel</th>
              <th style={{ padding: '1rem', color: 'var(--text-muted)', fontWeight: 600 }}>Message</th>
              <th style={{ padding: '1rem', color: 'var(--text-muted)', fontWeight: 600 }}>Status</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr>
                <td colSpan="6" style={{ padding: '2rem', textAlign: 'center', color: 'var(--text-muted)' }}>
                  Synchronizing logs...
                </td>
              </tr>
            ) : logs.length === 0 ? (
              <tr>
                <td colSpan="6" style={{ padding: '2rem', textAlign: 'center', color: 'var(--text-muted)' }}>
                  No notifications have been dispatched yet.
                </td>
              </tr>
            ) : (
              logs.map((log) => (
                <tr key={log.id} style={{ borderBottom: '1px solid var(--border-color)', transition: 'background 0.2s' }}>
                  <td style={{ padding: '1rem', whiteSpace: 'nowrap', fontSize: '0.85rem' }}>
                    {new Date(log.timestamp).toLocaleString()}
                  </td>
                  <td style={{ padding: '1rem' }}>
                    <div style={{ fontWeight: 500 }}>{log.user_email}</div>
                    <div style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>{log.user_phone}</div>
                  </td>
                  <td style={{ padding: '1rem' }}>
                    <span style={{ 
                      padding: '0.2rem 0.6rem', 
                      background: 'rgba(59, 130, 246, 0.1)', 
                      color: 'var(--accent)', 
                      borderRadius: '12px',
                      fontSize: '0.85rem',
                      fontWeight: 600
                    }}>
                      {log.category}
                    </span>
                  </td>
                  <td style={{ padding: '1rem', fontSize: '0.9rem' }}>{log.channel}</td>
                  <td style={{ padding: '1rem', fontSize: '0.9rem', maxWidth: '200px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }} title={log.message}>
                    {log.message}
                  </td>
                  <td style={{ padding: '1rem' }}>
                    <span style={{ 
                      padding: '0.2rem 0.6rem', 
                      background: log.delivery_status === 'Success' ? 'rgba(16, 185, 129, 0.1)' : 'rgba(239, 68, 68, 0.1)', 
                      color: log.delivery_status === 'Success' ? 'var(--success)' : 'var(--danger)', 
                      borderRadius: '12px',
                      fontSize: '0.85rem',
                      fontWeight: 600
                    }} title={log.error_message || ''}>
                      {log.delivery_status}
                    </span>
                    {log.retry_count > 0 && log.delivery_status === 'Failed' && (
                      <div style={{ fontSize: '0.7rem', color: 'var(--danger)', marginTop: '4px' }}>
                        Retries: {log.retry_count}
                      </div>
                    )}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
