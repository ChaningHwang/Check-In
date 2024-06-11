// src/pages/familyInfo.js
import { useState } from 'react';

export default function FamilyInfoPage() {
  const [fetchOption, setFetchOption] = useState('families');
  const [query, setQuery] = useState('');
  const [familyData, setFamilyData] = useState(null);

  const handleChange = (e) => {
    setFetchOption(e.target.value);
    setQuery('');
  };

  const handleQueryChange = (e) => {
    setQuery(e.target.value);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    let url = 'http://localhost:8080/' + fetchOption;
    let headers = new Headers();

    console.log("headers")
    console.log(headers)
    console.log(fetchOption)

    if (fetchOption === 'GetMembersByLastName') {
        console.log("here is reached", query)
        headers.append('lastName', query);
    } else if (fetchOption === 'GetMembersBySchoolID') {
        headers.append('SchoolID', query);
    }

    console.log(headers)

    try {
      const res = await fetch(url, {
        method: 'GET',
        mode: 'cors',
        headers: headers,
      });

    //   if (!res.ok) {
    //     throw new Error('Network response was not ok');
    //   }

      console.log(res)
      const data = await res.json();
      setFamilyData(data);
    } catch (error) {
      console.error('Failed to fetch:', error);
      alert('Failed to fetch family information. Please try again.');
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <div className="bg-white p-6 rounded shadow-md w-full max-w-lg">
        <h1 className="text-2xl mb-4">Get Family Information</h1>
        <form onSubmit={handleSubmit}>
          <label className="block mb-4">
            <span className="text-gray-700">Select option to fetch family information:</span>
            <select
              value={fetchOption}
              onChange={handleChange}
              className="block w-full mt-1 border border-gray-300 p-2 rounded"
            >
              <option value="families">Get All Families</option>
              <option value="GetMembersByLastName">Get Members by Last Name</option>
              <option value="GetMembersBySchoolID">Get Members by School ID</option>
            </select>
          </label>

          {(fetchOption === 'GetMembersByLastName' || fetchOption === 'GetMembersBySchoolID') && (
            <label className="block mb-4">
              <span className="text-gray-700">
                {fetchOption === 'GetMembersByLastName' ? 'Last Name:' : 'School ID:'}
              </span>
              <input
                type="text"
                value={query}
                onChange={handleQueryChange}
                className="block w-full mt-1 border border-gray-300 p-2 rounded"
              />
            </label>
          )}

          <button type="submit" className="bg-blue-500 text-white px-4 py-2 rounded">
            Fetch Data
          </button>
        </form>

        {familyData && (
          <div className="mt-6">
            <h2 className="text-xl mb-4">Family Data:</h2>
            <pre className="bg-gray-200 p-4 rounded">{JSON.stringify(familyData, null, 2)}</pre>
          </div>
        )}
      </div>
    </div>
  );
}

