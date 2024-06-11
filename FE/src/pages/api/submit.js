// pages/api/submit.js
export default function handler(req, res) {
    if (req.method === 'POST') {
      const { name, email } = req.body;
      // Simulate a call to a backend service and return the data
      res.status(200).json({ name, email, receivedAt: new Date().toISOString() });
    } else {
      res.status(405).json({ message: 'Method not allowed' });
    }
  }
  