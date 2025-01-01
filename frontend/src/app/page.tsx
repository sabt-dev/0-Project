'use client';
import React from 'react';
import Link from "next/link";
import { useRouter } from 'next/navigation';
import Footer from './_components/footer';


export default function Home() {
    const [email, setEmail] = React.useState('');
    const [password, setPassword] = React.useState('');
    const [loading, setLoading] = React.useState(false);
    const router = useRouter();

    const handleLogin = async (e: { preventDefault: () => void; }) => {
        e.preventDefault();
        setLoading(true);
        try {
            const response = await fetch('https://localhost:5000/api/v1/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify({
                    "email": email,
                    "pwd": password
                }),
            });

            const text = await response.text();
            const data = text ? JSON.parse(text) : {};
            //console.log(data);
            setLoading(false);
            if (data.status === "success") {
                router.push(`/profile/${data.uid}`);
            }
        } catch (error) {
            console.warn('An unexpected error happened:', error);
            setLoading(false);
        }
    }


  return (
      <div>
        <div className="hero bg-base-200 min-h-screen">
            <div className="hero-content flex-col lg:flex-row-reverse">
                <div className="text-center lg:text-left">
                    <h1 className="text-5xl font-bold">Login now!</h1>
                    <p className="py-6">
                        Don&#39;t have an account? <Link href="#" className="link link-hover">Sign up here</Link>
                    </p>
                </div>
                <div className="card bg-base-100 w-full max-w-sm shrink-0 shadow-2xl">
                    <form className="card-body" onSubmit={handleLogin}>
                        <div className="form-control">
                            <label className="label">
                                <span className="label-text">Email</span>
                            </label>
                            <input 
                                type="email" 
                                placeholder="email" 
                                className="input input-bordered" 
                                value={email} 
                                onChange={(e) => setEmail(e.target.value)} required/>
                        </div>
                        <div className="form-control">
                            <label className="label">
                                <span className="label-text">Password</span>
                            </label>
                            <input 
                                type="password" 
                                placeholder="password" 
                                className="input input-bordered" 
                                value={password} 
                                onChange={(e) => setPassword(e.target.value)} required/>
                            <label className="label">
                                <a href="#" className="label-text-alt link link-hover">Forgot password?</a>
                            </label>
                        </div>
                        <div className="form-control mt-6">
                            <button className="btn btn-primary" type='submit' disabled={loading}>
                                {loading ? 'Loading...' : 'Login'}
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
        <Footer/>
      </div>
  );
}
