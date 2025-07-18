/* eslint-disable @typescript-eslint/no-explicit-any */
'use client';
import React, { useEffect, useState } from 'react';

const UserProfileID = ({ params }: { params: any }) => {
    const [authorized, setAuthorized] = useState<boolean>(false);
    const [userID, setUserID] = useState<string | null>(null);
    const [res, setRes] = useState<any>(null);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchParams = async () => {
            try {
                setLoading(true);
                const resolvedParams = await params;
                const id = resolvedParams.id;
                setUserID(id);

                const response = await fetch(`http://localhost:5000/api/v1/users/me`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    credentials: 'include',
                });

                const text = await response.text();
                const resp = text ? JSON.parse(text) : {};
                setRes(resp);
                console.log(res);
                
                if (resp.data.id !== id) {
                    setAuthorized(false);
                    setError('Unauthorized');
                } else {
                    setAuthorized(true);
                    setError(null);
                }
            } catch (error) {
                setError('An unexpected error happened: ' + error);
            } finally {
                setLoading(false);
            }
        };

        fetchParams();

        return () => {
            setAuthorized(false);
            setUserID(null);
            setRes(null);
            setLoading(false);
        };
    }, [params]);

    if (loading) {
        return <div className='min-h-screen'>Loading...</div>;
    }

    if (error) {
        return <div className='min-h-screen'>Error: {error}</div>;
    }

    if (!authorized) {
        return <div className='min-h-screen'>You are not authorized to view this page</div>;
    }

    return (
        <div className='min-h-screen'>
            <h1>User Profile ID: {userID}</h1>
            <pre>{JSON.stringify(res, null, 2)}</pre>
        </div>
    );
};

export default UserProfileID;