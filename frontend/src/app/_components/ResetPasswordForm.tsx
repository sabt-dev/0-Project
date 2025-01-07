/* eslint-disable @typescript-eslint/no-unused-vars */
'use client';
import React, { useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

const ResetPasswordForm = () => {
    const [newPassword, setNewPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const searchParams = useSearchParams();
    const token = searchParams.get('token');

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (newPassword !== confirmPassword) {
            setError('Passwords do not match');
            return;
        }

        try {
            const response = await fetch('https://localhost:5000/api/v1/auth/reset-password', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    token,
                    newPassword,
                })
            });

            const data = await response.json();

            if (response.ok) {
                setSuccess(data.message);
                setError('');
            } else {
                setError(data.error || 'Failed to reset password');
                setSuccess('');
            }
        } catch (err) {
            setError('An unexpected error occurred');
            setSuccess('');
        }
    };

    return (
        <div className='hero bg-base-200'>
            <div className='hero-content flex items-center max-w-md justify-center min-h-screen'>
                <div className='card bg-base-100 w-full  shrink-0 shadow-2xl flex items-center justify-center'>
                    <h1>Reset Password</h1>
                    <form onSubmit={handleSubmit} className='card-body'>
                        <div className='mt-4 form-control'>
                            <label className='label'>New Password</label>
                            <input
                                className='input input-bordered input-md ml-2'
                                type="password"
                                value={newPassword}
                                onChange={(e) => setNewPassword(e.target.value)}
                                required
                            />
                        </div>
                        <div className='mt-4 form-control'>
                            <label className='label'>Confirm Password</label>
                            <input
                                className='input input-bordered input-md ml-2 mt-2'
                                type="password"
                                value={confirmPassword}
                                onChange={(e) => setConfirmPassword(e.target.value)}
                                required
                            />
                        </div>
                        {error && <p style={{ color: 'red' }}>{error}</p>}
                        {success && <p style={{ color: 'green' }}>{success}</p>}
                        <button className='btn btn-primary' type="submit">Reset Password</button>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default ResetPasswordForm;