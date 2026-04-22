import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import SubmissionForm from './SubmissionForm';

// Mock the global fetch
global.fetch = vi.fn();

describe('SubmissionForm', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders all categories and default message area', () => {
    render(<SubmissionForm />);
    expect(screen.getByText('Send Notification')).toBeInTheDocument();
    
    // Check categories
    expect(screen.getByRole('option', { name: 'Sports' })).toBeInTheDocument();
    expect(screen.getByRole('option', { name: 'Finance' })).toBeInTheDocument();
    expect(screen.getByRole('option', { name: 'Movies' })).toBeInTheDocument();
    
    // Check elements
    expect(screen.getByLabelText(/message/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /dispatch message/i })).toBeInTheDocument();
  });

  it('shows error if message is empty on submit', async () => {
    render(<SubmissionForm />);
    
    // Submit without typing
    fireEvent.click(screen.getByRole('button', { name: /dispatch message/i }));
    
    expect(await screen.findByText('Message cannot be empty.')).toBeInTheDocument();
    expect(fetch).not.toHaveBeenCalled();
  });

  it('handles successful submission, shows success alert, and clears message', async () => {
    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ status: 'success' })
    });

    const mockOnSent = vi.fn();
    render(<SubmissionForm onSent={mockOnSent} />);

    // Type message
    const textarea = screen.getByPlaceholderText(/type your notification/i);
    fireEvent.change(textarea, { target: { value: 'Hello world' } });
    
    // Change category to Finance
    const select = screen.getByLabelText(/category/i);
    fireEvent.change(select, { target: { value: 'Finance' } });

    // Submit
    fireEvent.click(screen.getByRole('button', { name: /dispatch message/i }));

    // Wait for success status
    expect(await screen.findByText('Message dispatched successfully!')).toBeInTheDocument();
    
    // Verify fetch call payload
    expect(fetch).toHaveBeenCalledWith('http://localhost:8080/api/notifications', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ category: 'Finance', message: 'Hello world' })
    });

    // Message reset, category retained
    expect(textarea.value).toBe('');
    expect(select.value).toBe('Finance');
    
    expect(mockOnSent).toHaveBeenCalled();
  });

  it('handles server errors and displays failure alert', async () => {
    fetch.mockResolvedValueOnce({
      ok: false,
      status: 500
    });

    render(<SubmissionForm />);

    const textarea = screen.getByPlaceholderText(/type your notification/i);
    fireEvent.change(textarea, { target: { value: 'This will fail' } });
    fireEvent.click(screen.getByRole('button', { name: /dispatch message/i }));

    expect(await screen.findByText(/Failed to dispatch.*500/i)).toBeInTheDocument();
  });
});
