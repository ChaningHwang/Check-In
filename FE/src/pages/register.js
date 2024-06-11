// src/pages/register.js
// src/pages/register.js
import { useState } from 'react';

export default function RegisterPage() {
  const [members, setMembers] = useState([{ firstName: '', lastName: '', schoolID: '', schoolName: '' }]);

  const handleChange = (index, e) => {
    const { name, value } = e.target;
    const newMembers = [...members];
    newMembers[index][name] = value;
    setMembers(newMembers);
  };

  const handleAddMember = () => {
    setMembers([...members, { firstName: '', lastName: '', schoolID: '', schoolName: '' }]);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    const postData = { members };

    // const res = await fetch('http://localhost:8080/families', {
    //   method: 'POST',
    //   body: JSON.stringify(postData),
    // });

    // if (res.ok) {
    //   alert('Registration successful!');
    //   setMembers([{ firstName: '', lastName: '', schoolID: '', schoolName: '' }]);
    // } else {
    //   alert('Registration failed. Please try again.');
    // }
    try {
      const res = await fetch('http://localhost:8080/families', {
        method: 'POST',
        mode: 'no-cors',
        body: JSON.stringify(postData),
      });

      console.log(res.headers);

      // if (!res.ok) {
      //   throw new Error('Network response was not ok');
      // }

      alert('Registration successful!');
      setMembers([{ firstName: '', lastName: '', schoolID: '', schoolName: '' }]);
    } catch (error) {
      console.error('Failed to fetch:', error);
      alert('Registration failed. Please try again.');
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <form className="bg-white p-6 rounded shadow-md" onSubmit={handleSubmit}>
        <h1 className="text-2xl mb-4">Register</h1>
        {members.map((member, index) => (
          <div key={index} className="mb-4">
            <label className="block mb-2">
              First Name:
              <input
                type="text"
                name="firstName"
                value={member.firstName}
                onChange={(e) => handleChange(index, e)}
                className="border border-gray-300 p-2 w-full"
              />
            </label>
            <label className="block mb-2">
              Last Name:
              <input
                type="text"
                name="lastName"
                value={member.lastName}
                onChange={(e) => handleChange(index, e)}
                className="border border-gray-300 p-2 w-full"
              />
            </label>
            <label className="block mb-2">
              School ID:
              <input
                type="text"
                name="schoolID"
                value={member.schoolID}
                onChange={(e) => handleChange(index, e)}
                className="border border-gray-300 p-2 w-full"
              />
            </label>
            <label className="block mb-2">
              School Name:
              <input
                type="text"
                name="schoolName"
                value={member.schoolName}
                onChange={(e) => handleChange(index, e)}
                className="border border-gray-300 p-2 w-full"
              />
            </label>
          </div>
        ))}
        <button
          type="button"
          onClick={handleAddMember}
          className="bg-green-500 text-white px-4 py-2 rounded mb-4"
        >
          + Add Member
        </button>
        <button type="submit" className="bg-blue-500 text-white px-4 py-2 rounded">
          Register
        </button>
      </form>
    </div>
  );
}
