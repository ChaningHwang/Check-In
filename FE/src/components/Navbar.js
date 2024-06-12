// src/components/Navbar.js
import Link from 'next/link';

export default function Navbar() {
  return (
    <nav className="bg-blue-600 p-4">
      <ul className="flex space-x-4">
        <li>
          <Link href="/register">
            <span className="text-white cursor-pointer">Register</span>
          </Link>
        </li>
        {/* <li>
          <Link href="/checkin">
            <span className="text-white cursor-pointer">Check-In</span>
          </Link>
        </li>
        <li>
          <Link href="/familyInfo">
            <span className="text-white cursor-pointer">Get Family Info</span>
          </Link>
        </li>
        <li>
          <Link href="/selectEvent">
            <span className="text-white cursor-pointer">Select Event</span>
          </Link>
        </li> */}
        <li>
          <Link href="/eventInfo">
            <span className="text-white cursor-pointer">Event Info</span>
          </Link>
        </li>
      </ul>
    </nav>
  );
}

