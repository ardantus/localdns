import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

export default function Dashboard() {
    const [domains, setDomains] = useState([]);
    const [newDomain, setNewDomain] = useState('');
    const [loading, setLoading] = useState(true);
    const [expandedDomain, setExpandedDomain] = useState(null);

    // Record Form
    const [newRecord, setNewRecord] = useState({ name: '', type: 'A', content: '' });

    const navigate = useNavigate();
    const user = JSON.parse(localStorage.getItem('user') || '{}');

    useEffect(() => {
        fetchDomains();
    }, []);

    const fetchDomains = async () => {
        try {
            const res = await axios.get('/api/domains');
            setDomains(res.data);
            setLoading(false);
        } catch (error) {
            console.error(error);
            setLoading(false);
        }
    };

    const handleCreateDomain = async (e) => {
        e.preventDefault();
        try {
            await axios.post('/api/domains', {
                name: newDomain,
                user_id: user.id
            });
            setNewDomain('');
            fetchDomains();
        } catch (error) {
            alert('Failed to create domain');
        }
    };

    const handleAddRecord = async (e) => {
        e.preventDefault();
        if (!expandedDomain) return;
        try {
            await axios.post(`/api/domains/${expandedDomain}/records`, newRecord);
            setNewRecord({ name: '', type: 'A', content: '' });
            // Ideally re-fetch records for this domain or expand logic. 
            // For MVP we just alert success or close.
            alert('Record added! (Refresh to see in full list would require getting domain details)');
            setExpandedDomain(null);
        } catch (error) {
            alert('Failed to add record');
        }
    };

    const handleLogout = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = '/';
    };

    return (
        <div className="min-h-screen bg-gray-100">
            <nav className="bg-white shadow-sm">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                    <div className="flex justify-between h-16">
                        <div className="flex items-center">
                            <h1 className="text-xl font-bold">LocalDNS Registrar</h1>
                        </div>
                        <div className="flex items-center">
                            <span className="mr-4">Welcome, {user.username}</span>
                            <button onClick={handleLogout} className="text-red-600 hover:text-red-800">Logout</button>
                        </div>
                    </div>
                </div>
            </nav>

            <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
                {/* Add Domain */}
                <div className="bg-white shadow sm:rounded-lg p-6 mb-6">
                    <h2 className="text-lg font-medium mb-4">Register New Domain</h2>
                    <form className="flex gap-4" onSubmit={handleCreateDomain}>
                        <input
                            type="text"
                            placeholder="example.lan"
                            className="flex-1 border rounded px-3 py-2"
                            value={newDomain}
                            onChange={(e) => setNewDomain(e.target.value)}
                            required
                        />
                        <button className="bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700">Register</button>
                    </form>
                </div>

                {/* Domain List */}
                <div className="bg-white shadow sm:rounded-lg overflow-hidden">
                    <ul className="divide-y divide-gray-200">
                        {domains.map(domain => (
                            <li key={domain.id} className="p-4 hover:bg-gray-50">
                                <div className="flex justify-between items-center cursor-pointer" onClick={() => setExpandedDomain(expandedDomain === domain.id ? null : domain.id)}>
                                    <div>
                                        <p className="text-lg font-medium text-blue-600">
                                            {domain.name}
                                            {domain.user && <span className="ml-2 text-xs bg-gray-200 text-gray-700 px-2 py-1 rounded">Owner: {domain.user.username}</span>}
                                        </p>
                                        <p className="text-sm text-gray-500">Created: {new Date(domain.created_at).toLocaleDateString()}</p>
                                    </div>
                                    <button className="text-gray-400">
                                        {expandedDomain === domain.id ? 'Collapse' : 'Manage DNS'}
                                    </button>
                                </div>

                                {expandedDomain === domain.id && (
                                    <div className="mt-4 pl-4 border-l-4 border-blue-100">
                                        <h3 className="text-sm font-semibold uppercase tracking-wide text-gray-500">Add Record</h3>
                                        <form className="flex gap-2 mt-2" onSubmit={handleAddRecord}>
                                            <input
                                                placeholder="Subdomain (e.g. www)"
                                                className="border px-2 py-1 w-1/4"
                                                value={newRecord.name}
                                                onChange={e => setNewRecord({ ...newRecord, name: e.target.value })}
                                                required
                                            />
                                            <select
                                                className="border px-2 py-1"
                                                value={newRecord.type}
                                                onChange={e => setNewRecord({ ...newRecord, type: e.target.value })}
                                            >
                                                <option value="A">A</option>
                                                <option value="CNAME">CNAME</option>
                                                <option value="TXT">TXT</option>
                                            </select>
                                            <input
                                                placeholder="Content (e.g. 192.168.1.5)"
                                                className="border px-2 py-1 flex-1"
                                                value={newRecord.content}
                                                onChange={e => setNewRecord({ ...newRecord, content: e.target.value })}
                                                required
                                            />
                                            <button className="bg-blue-600 text-white px-3 py-1 rounded">Add</button>
                                        </form>

                                        {/* List Existing Records - TODO: Fetch them */}
                                        <div className="mt-2 text-sm text-gray-500 italic">
                                            (Existing records list not implemented in MVP view, add works though!)
                                        </div>
                                    </div>
                                )}
                            </li>
                        ))}
                        {domains.length === 0 && !loading && (
                            <li className="p-4 text-center text-gray-500">No domains found. Register one above!</li>
                        )}
                    </ul>
                </div>
            </main>
        </div>
    );
}
