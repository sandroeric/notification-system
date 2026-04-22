import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import LogHistory from './LogHistory';

global.fetch = vi.fn();

describe('LogHistory', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows loading state initially', () => {
    fetch.mockImplementation(() => new Promise(() => {})); // Never resolves
    render(<LogHistory refreshHash={1} />);
    expect(screen.getByText('Synchronizing logs...')).toBeInTheDocument();
  });

  it('renders "No notifications" state if empty array is returned', async () => {
    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => []
    });

    render(<LogHistory refreshHash={1} />);
    
    expect(await screen.findByText('No notifications have been dispatched yet.')).toBeInTheDocument();
  });

  it('renders a table of logs successfully', async () => {
    const mockLogs = [
      {
        id: 1,
        timestamp: '2023-10-01T12:00:00Z',
        user_email: 'test@example.com',
        user_phone: '123-456',
        category: 'Sports',
        channel: 'E-Mail',
        message: 'Goal!',
        delivery_status: 'Success',
        retry_count: 0,
        error_message: ''
      },
      {
        id: 2,
        timestamp: '2023-10-01T12:05:00Z',
        user_email: 'fail@example.com',
        user_phone: '000-000',
        category: 'Finance',
        channel: 'SMS',
        message: 'Crash',
        delivery_status: 'Failed',
        retry_count: 3,
        error_message: 'Network error'
      }
    ];

    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => mockLogs
    });

    render(<LogHistory refreshHash={1} />);

    // Wait for the rows to render
    expect(await screen.findByText('test@example.com')).toBeInTheDocument();
    expect(screen.getByText('Sports')).toBeInTheDocument();
    expect(screen.getByText('E-Mail')).toBeInTheDocument();
    expect(screen.getByText('Goal!')).toBeInTheDocument();

    // Check failed row rendering
    expect(screen.getByText('fail@example.com')).toBeInTheDocument();
    expect(screen.getByText('Failed')).toBeInTheDocument();
    expect(screen.getByText('Retries: 3')).toBeInTheDocument();
  });

  it('renders error alert when fetching fails', async () => {
    fetch.mockRejectedValueOnce(new Error('Network error'));

    render(<LogHistory refreshHash={1} />);

    expect(await screen.findByText(/Error retrieving logs: Network error/i)).toBeInTheDocument();
  });
});
