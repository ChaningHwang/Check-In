// src/pages/index.js
import Navbar from '../components/Navbar';

export default function Home() {
  return (
    <div className="min-h-screen bg-gray-100">
      <Navbar />
      <div className="flex items-center justify-center p-6">
        <div className="bg-white p-6 rounded shadow-md">
          <h1 className="text-2xl mb-4">Welcome to My Next.js App</h1>
          <p>Please use the navigation bar to register, check-in, or get family information.</p>
        </div>
      </div>
    </div>
  );
}
