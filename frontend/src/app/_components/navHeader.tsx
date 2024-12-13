import Link from 'next/link'
import React from 'react'

export default function navHeader() {
  return (
    <div className="drawer">
    <input id="my-drawer-3" type="checkbox" className="drawer-toggle"/>
    <div className="drawer-content flex flex-col">
        {/* Navbar */}
        <div className="navbar bg-base-300 w-full">
            <div className="flex-none lg:hidden">
                <label htmlFor="my-drawer-3" aria-label="open sidebar" className="btn btn-square btn-ghost">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                        className="inline-block h-6 w-6 stroke-current">
                        <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth="2"
                            d="M4 6h16M4 12h16M4 18h16">
                        </path>
                    </svg>
                </label>
            </div>

            <div className="drawer">
                <input id="my-drawer" type="checkbox" className="drawer-toggle"/>
                <div className="drawer-content">
                    {/* Page content here */}
                    <label htmlFor="my-drawer" className="btn btn-primary drawer-button">Open drawer</label>
                </div>
                <div className="drawer-side">
                    <label htmlFor="my-drawer" aria-label="close sidebar" className="drawer-overlay"></label>
                    <ul className="menu bg-base-200 text-base-content min-h-full w-80 p-4">
                        {/* Sidebar content here */}
                        <li><a>Sidebar Item 1</a></li>
                        <li><a>Sidebar Item 2</a></li>
                    </ul>
                </div>
            </div>

            <div className="hidden flex-none lg:block">
                <ul className="menu menu-horizontal">
                    {/* Navbar menu content here */}
                    <li><Link href='/about'>About</Link></li>
                    <li><Link href='/'>Home</Link></li>
                </ul>
            </div>
        </div>
    </div>
</div>
  )
}
