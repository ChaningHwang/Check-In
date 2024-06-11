// src/pages/eventInfo.js
import { useState, useEffect } from 'react';

export default function EventInfoPage() {
  const [events, setEvents] = useState([]);
  const [selectedEvent, setSelectedEvent] = useState('');
  const [fetchOption, setFetchOption] = useState('');
  const [query, setQuery] = useState('');
  const [familyData, setFamilyData] = useState(null);
  const [selectedMembers, setSelectedMembers] = useState([]);
  const [isFetching, setIsFetching] = useState(false);

  useEffect(() => {
    // Fetch events when the component mounts
    const fetchEvents = async () => {
      console.log('Fetching events...');
      try {
        const res = await fetch('http://localhost:8080/events');
        if (!res.ok) {
          throw new Error(`Network response was not ok: ${res.statusText}`);
        }
        const data = await res.json();
        console.log('Events fetched:', data);
        setEvents(data);
      } catch (error) {
        console.error('Failed to fetch events:', error);
        alert(`Failed to fetch events: ${error.message}`);
      }
    };

    fetchEvents();
  }, []);

  const handleEventChange = (e) => {
    setSelectedEvent(e.target.value);
  };

  const handleFetchOptionChange = (e) => {
    setFetchOption(e.target.value);
    setQuery('');
  };

  const handleQueryChange = (e) => {
    setQuery(e.target.value);
  };

  const handleMemberSelection = (member) => {
    setSelectedMembers((prevSelected) => {
      if (prevSelected.includes(member)) {
        return prevSelected.filter((m) => m !== member);
      } else {
        return [...prevSelected, member];
      }
    });
  };

  const handleFetchSubmit = async (e) => {
    e.preventDefault();

    if (!selectedEvent) {
      alert('Please select an event first.');
      return;
    }

    let url = `http://localhost:8080/${fetchOption}`;
    let headers = new Headers();
    headers.append('Content-Type', 'application/json');

    // if (fetchOption === 'GetMembersByLastName' || fetchOption === 'GetMembersBySchoolID' || fetchOption === 'GetFamilyByID') {
    //   headers.append('Query', query);
    // }
    if (fetchOption === 'GetMembersByLastName') {
        console.log("here is reached", query)
        headers.append('lastName', query);
    } else if (fetchOption === 'GetMembersBySchoolID') {
        headers.append('SchoolID', query);
    }

    setIsFetching(true);
    try {
      const res = await fetch(url, {
        method: 'GET',
        headers: headers,
        mode: 'cors',
      });

      if (!res.ok) {
        throw new Error(`Network response was not ok: ${res.statusText}`);
      }

      const data = await res.json();
      setFamilyData(data);
    } catch (error) {
      console.error('Failed to fetch:', error);
      alert(`Failed to fetch family information: ${error.message}`);
    } finally {
      setIsFetching(false);
    }
  };

  const handleCheckIn = async () => {
    if (!selectedEvent) {
      alert('Please select an event first.');
      return;
    }

    const url = 'http://localhost:8080/checkin';
    const headers = new Headers({
      'Content-Type': 'application/json',
    });


    const temp = selectedMembers.map(obj => obj.MemberId)

    console.log("here is the temp", temp)
    console.log("here is selectedMembers", selectedMembers)

    const body = JSON.stringify({
      EventID: "event_"+selectedEvent,
      MemberIDs: temp,
    });

    console.log("here is hte selectedEvent" , selectedEvent)
    console.log("here is the selectedMembers", selectedMembers)

    try {
      const res = await fetch(url, {
        method: 'POST',
        headers: headers,
        body: body,
        mode: 'cors'
      });

      if (!res.ok) {
        throw new Error(`Network response was not ok: ${res.statusText}`);
      }

      alert('Check-in successful!');
      setSelectedMembers([]);
    } catch (error) {
      console.error('Failed to check in:', error);
      alert(`Failed to check in: ${error.message}`);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <div className="bg-white p-6 rounded shadow-md w-full max-w-lg">
        <h1 className="text-2xl mb-4">Select Event</h1>
        <div className="mb-4">
          <label className="block mb-2 text-gray-700">Select Event:</label>
          <select
            value={selectedEvent}
            onChange={handleEventChange}
            className="block w-full mt-1 border border-gray-300 p-2 rounded"
          >
            <option value="">Select an event</option>
            {events.map((event) => (
              <option key={event.EventID} value={event.EventID}>
                eventID: {event.EventID}, eventName: {event.EventName}
              </option>
            ))}
          </select>
        </div>
        {selectedEvent && (
          <>
            <h1 className="text-2xl mb-4">Get Family Information</h1>
            <form onSubmit={handleFetchSubmit}>
              <label className="block mb-4">
                <span className="text-gray-700">Select option to fetch family information:</span>
                <select
                  value={fetchOption}
                  onChange={handleFetchOptionChange}
                  className="block w-full mt-1 border border-gray-300 p-2 rounded"
                >
                  <option value="">Select an option</option>
                  <option value="GetFamilyByID">Get Family by Family ID</option>
                  <option value="GetMembersByLastName">Get Members by Last Name</option>
                  <option value="GetMembersBySchoolID">Get Members by School ID</option>
                </select>
              </label>

              {(fetchOption === 'GetFamilyByID' || fetchOption === 'GetMembersByLastName' || fetchOption === 'GetMembersBySchoolID') && (
                <label className="block mb-4">
                  <span className="text-gray-700">
                    {fetchOption === 'GetFamilyByID' ? 'Family ID:' : fetchOption === 'GetMembersByLastName' ? 'Last Name:' : 'School ID:'}
                  </span>
                  <input
                    type="text"
                    value={query}
                    onChange={handleQueryChange}
                    className="block w-full mt-1 border border-gray-300 p-2 rounded"
                  />
                </label>
              )}

              <button
                type="submit"
                className="bg-blue-500 text-white px-4 py-2 rounded"
                disabled={isFetching}
              >
                {isFetching ? 'Fetching...' : 'Fetch Data'}
              </button>
            </form>

            {familyData && (
              <div className="mt-6">
                <h2 className="text-xl mb-4">Family Data:</h2>
                <pre className="bg-gray-200 p-4 rounded">{JSON.stringify(familyData, null, 2)}</pre>
                <h2 className="text-xl mt-6 mb-4">Select Members to Check-In</h2>
                {familyData && familyData.length > 0 ? (
                  <div className="mb-4">
                    {familyData.map((member, index) => (
                      <div key={index} className="flex items-center mb-2">
                        <input
                          type="checkbox"
                          id={`member-${index}`}
                          className="mr-2"
                          onChange={() => handleMemberSelection(member)}
                          checked={selectedMembers.includes(member)}
                        />
                        <label htmlFor={`member-${index}`} className="text-gray-700">
                          {member.FirstName} {member.LastName}
                        </label>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p>No members found.</p>
                )}
                <button
                  onClick={handleCheckIn}
                  className="bg-green-500 text-white px-4 py-2 rounded"
                >
                  Check-In Selected Members
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}
