import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

// Configure axios to send JWT token with every request
const api = axios.create();
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// DNS Record Types with examples
const RECORD_TYPES = [
    { value: 'A', label: 'A (IPv4)', placeholder: '192.168.1.100' },
    { value: 'AAAA', label: 'AAAA (IPv6)', placeholder: '2001:db8::1' },
    { value: 'CNAME', label: 'CNAME (Alias)', placeholder: 'www.example.lan' },
    { value: 'MX', label: 'MX (Mail)', placeholder: 'mail.example.lan' },
    { value: 'NS', label: 'NS (Nameserver)', placeholder: 'ns1.example.lan' },
    { value: 'TXT', label: 'TXT (Text)', placeholder: 'v=spf1 include:example.lan ~all' },
    { value: 'SRV', label: 'SRV (Service)', placeholder: '0 5 5060 sipserver.example.lan' },
    { value: 'PTR', label: 'PTR (Reverse)', placeholder: 'host.example.lan' },
];

const DOMAIN_EXAMPLES = ['myserver.lan', 'homelab.local', 'nas.home', 'printer.internal'];

export default function Dashboard() {
    const [domains, setDomains] = useState([]);
    const [newDomain, setNewDomain] = useState('');
    const [loading, setLoading] = useState(true);
    const [expandedDomain, setExpandedDomain] = useState(null);
    const [domainRecords, setDomainRecords] = useState({});
    const [activeTab, setActiveTab] = useState('domains');
    const [users, setUsers] = useState([]);
    const [registrarConfig, setRegistrarConfig] = useState({});

    // Forms
    const [newRecord, setNewRecord] = useState({ name: '', type: 'A', content: '', ttl: 3600, prio: 0 });
    const [editingRecord, setEditingRecord] = useState(null);
    const [editingUser, setEditingUser] = useState(null);
    const [editingRegistrant, setEditingRegistrant] = useState(null);
    const [editingConfig, setEditingConfig] = useState(false);

    const navigate = useNavigate();
    const user = JSON.parse(localStorage.getItem('user') || '{}');
    const isAdmin = user.role === 'admin';

    useEffect(() => {
        fetchDomains();
        if (isAdmin) {
            fetchUsers();
            fetchConfig();
        }
    }, []);

    const handleTabSwitch = (tab) => {
        setActiveTab(tab);
        if (tab === 'domains') fetchDomains();
        else if (tab === 'users') fetchUsers();
        else if (tab === 'config') fetchConfig();
    };

    const fetchDomains = async () => {
        try {
            const res = await api.get('/api/domains');
            setDomains(res.data || []);
            setLoading(false);
        } catch (error) {
            console.error(error);
            setLoading(false);
        }
    };

    const fetchUsers = async () => {
        try {
            const res = await api.get('/api/users');
            setUsers(res.data || []);
        } catch (error) {
            console.error(error);
        }
    };

    const fetchConfig = async () => {
        try {
            const res = await api.get('/api/config');
            setRegistrarConfig(res.data || {});
        } catch (error) {
            console.error(error);
        }
    };

    const fetchRecords = async (domainId) => {
        try {
            const res = await api.get(`/api/domains/${domainId}/records`);
            setDomainRecords(prev => ({ ...prev, [domainId]: res.data || [] }));
        } catch (error) {
            console.error(error);
        }
    };

    const handleCreateDomain = async (e) => {
        e.preventDefault();
        try {
            await api.post('/api/domains', { name: newDomain });
            setNewDomain('');
            fetchDomains();
        } catch (error) {
            alert('Failed to create domain: ' + (error.response?.data?.error || error.message));
        }
    };

    const handleDeleteDomain = async (domainId) => {
        if (!confirm('Delete this domain and all its records?')) return;
        try {
            await api.delete(`/api/domains/${domainId}`);
            fetchDomains();
        } catch (error) {
            alert('Failed to delete domain');
        }
    };

    const handleAddRecord = async (e, domainId) => {
        e.preventDefault();
        try {
            await api.post(`/api/domains/${domainId}/records`, newRecord);
            setNewRecord({ name: '', type: 'A', content: '', ttl: 3600, prio: 0 });
            fetchRecords(domainId);
        } catch (error) {
            alert('Failed to add record');
        }
    };

    const handleUpdateRecord = async (e) => {
        e.preventDefault();
        if (!editingRecord) return;
        try {
            await api.put(`/api/records/${editingRecord.id}`, editingRecord);
            fetchRecords(editingRecord.domain_id);
            setEditingRecord(null);
        } catch (error) {
            alert('Failed to update record');
        }
    };

    const handleDeleteRecord = async (recordId, domainId) => {
        if (!confirm('Delete this record?')) return;
        try {
            await api.delete(`/api/records/${recordId}`);
            fetchRecords(domainId);
        } catch (error) {
            alert('Failed to delete record');
        }
    };

    const handleUpdateRegistrant = async (e) => {
        e.preventDefault();
        if (!editingRegistrant) return;
        try {
            await api.put(`/api/domains/${editingRegistrant.id}/registrant`, editingRegistrant);
            fetchDomains();
            setEditingRegistrant(null);
        } catch (error) {
            alert('Failed to update registrant: ' + (error.response?.data?.error || error.message));
        }
    };

    const handleUpdateUser = async (e) => {
        e.preventDefault();
        if (!editingUser) return;
        try {
            if (editingUser.id) {
                await api.put(`/api/users/${editingUser.id}`, editingUser);
            } else {
                await api.post('/api/users', editingUser);
            }
            fetchUsers();
            setEditingUser(null);
        } catch (error) {
            alert('Failed to save user: ' + (error.response?.data?.error || error.message));
        }
    };

    const handleDeleteUser = async (userId) => {
        if (!confirm('Delete this user?')) return;
        try {
            await api.delete(`/api/users/${userId}`);
            fetchUsers();
        } catch (error) {
            alert('Failed to delete user');
        }
    };

    const handleUpdateConfig = async (e) => {
        e.preventDefault();
        try {
            await api.put('/api/config', registrarConfig);
            setEditingConfig(false);
            alert('Configuration saved!');
        } catch (error) {
            alert('Failed to update config: ' + (error.response?.data?.error || error.message));
        }
    };

    const toggleDomain = (domainId) => {
        if (expandedDomain === domainId) {
            setExpandedDomain(null);
        } else {
            setExpandedDomain(domainId);
            fetchRecords(domainId);
        }
    };

    const handleLogout = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = '/';
    };

    const selectedType = RECORD_TYPES.find(t => t.value === newRecord.type) || RECORD_TYPES[0];

    return (
        <div className="min-h-screen bg-gray-100">
            <nav className="bg-white shadow-sm">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                    <div className="flex justify-between h-16">
                        <div className="flex items-center">
                            <h1 className="text-xl font-bold text-blue-600">🌐 LocalDNS Registrar</h1>
                        </div>
                        <div className="flex items-center gap-4">
                            <span className="text-sm">
                                Welcome, <strong>{user.username}</strong>
                                {isAdmin && <span className="ml-1 text-xs bg-purple-100 text-purple-800 px-2 py-0.5 rounded">Admin</span>}
                            </span>
                            <button onClick={handleLogout} className="text-red-600 hover:text-red-800 text-sm">Logout</button>
                        </div>
                    </div>
                </div>
            </nav>

            {/* Tabs */}
            {isAdmin && (
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mt-4">
                    <div className="flex gap-2">
                        {['domains', 'users', 'config'].map(tab => (
                            <button
                                key={tab}
                                onClick={() => handleTabSwitch(tab)}
                                className={`px-4 py-2 rounded-t capitalize ${activeTab === tab ? 'bg-white font-medium' : 'bg-gray-200'}`}
                            >
                                {tab === 'config' ? 'Registrar Config' : tab}
                            </button>
                        ))}
                    </div>
                </div>
            )}

            <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
                {/* DOMAINS TAB */}
                {activeTab === 'domains' && (
                    <>
                        <div className="bg-white shadow sm:rounded-lg p-6 mb-6">
                            <h2 className="text-lg font-medium mb-2">Register New Domain</h2>
                            <p className="text-sm text-gray-500 mb-4">
                                Use local TLDs: <code className="bg-gray-100 px-1">.lan</code>, <code className="bg-gray-100 px-1">.local</code>, <code className="bg-gray-100 px-1">.home</code>, <code className="bg-gray-100 px-1">.internal</code>
                            </p>
                            <form className="flex gap-4" onSubmit={handleCreateDomain}>
                                <input
                                    type="text"
                                    placeholder={DOMAIN_EXAMPLES[Math.floor(Math.random() * DOMAIN_EXAMPLES.length)]}
                                    className="flex-1 border rounded px-3 py-2"
                                    value={newDomain}
                                    onChange={(e) => setNewDomain(e.target.value)}
                                    required
                                />
                                <button className="bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700">Register</button>
                            </form>
                        </div>

                        <div className="bg-white shadow sm:rounded-lg overflow-hidden">
                            <div className="px-4 py-3 bg-gray-50 border-b flex justify-between items-center">
                                <h3 className="font-medium">{isAdmin ? 'All Domains' : 'Your Domains'} ({domains.length})</h3>
                                <span className="text-sm text-gray-500">Test WHOIS: <code>whois -h localhost domain.lan</code></span>
                            </div>
                            <ul className="divide-y divide-gray-200">
                                {domains.map(domain => (
                                    <li key={domain.id} className="p-4">
                                        <div className="flex justify-between items-center">
                                            <div className="cursor-pointer flex-1" onClick={() => toggleDomain(domain.id)}>
                                                <p className="text-lg font-medium text-blue-600">
                                                    {domain.name}
                                                    {domain.user && <span className="ml-2 text-xs bg-gray-200 text-gray-700 px-2 py-1 rounded">Owner: {domain.user.username}</span>}
                                                    <span className={`ml-2 text-xs px-2 py-1 rounded ${domain.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>{domain.status || 'active'}</span>
                                                </p>
                                                <p className="text-sm text-gray-500">Created: {new Date(domain.created_at).toLocaleDateString()}</p>
                                            </div>
                                            <div className="flex gap-2">
                                                <button onClick={() => setEditingRegistrant({ ...domain })} className="text-purple-600 hover:text-purple-800 text-sm">WHOIS Info</button>
                                                <button onClick={() => toggleDomain(domain.id)} className="text-blue-600 hover:text-blue-800 text-sm">
                                                    {expandedDomain === domain.id ? '▲ Collapse' : '▼ DNS'}
                                                </button>
                                                <button onClick={() => handleDeleteDomain(domain.id)} className="text-red-600 hover:text-red-800 text-sm">Delete</button>
                                            </div>
                                        </div>

                                        {expandedDomain === domain.id && (
                                            <div className="mt-4 pl-4 border-l-4 border-blue-200">
                                                <h3 className="text-sm font-semibold uppercase text-gray-600 mb-2">Add DNS Record</h3>
                                                <form className="grid grid-cols-12 gap-2 mb-4" onSubmit={(e) => handleAddRecord(e, domain.id)}>
                                                    <input placeholder="@ or subdomain" className="col-span-2 border px-2 py-1 rounded text-sm" value={newRecord.name} onChange={e => setNewRecord({ ...newRecord, name: e.target.value })} required />
                                                    <select className="col-span-2 border px-2 py-1 rounded text-sm" value={newRecord.type} onChange={e => setNewRecord({ ...newRecord, type: e.target.value })}>
                                                        {RECORD_TYPES.map(type => <option key={type.value} value={type.value}>{type.label}</option>)}
                                                    </select>
                                                    <input placeholder={selectedType.placeholder} className="col-span-4 border px-2 py-1 rounded text-sm" value={newRecord.content} onChange={e => setNewRecord({ ...newRecord, content: e.target.value })} required />
                                                    <input type="number" placeholder="TTL" className="col-span-1 border px-2 py-1 rounded text-sm" value={newRecord.ttl} onChange={e => setNewRecord({ ...newRecord, ttl: parseInt(e.target.value) || 3600 })} />
                                                    <button className="col-span-3 bg-blue-600 text-white px-3 py-1 rounded text-sm hover:bg-blue-700">Add Record</button>
                                                </form>

                                                <h3 className="text-sm font-semibold uppercase text-gray-600 mb-2">DNS Records</h3>
                                                {domainRecords[domain.id]?.length > 0 ? (
                                                    <table className="w-full text-sm">
                                                        <thead className="bg-gray-50">
                                                            <tr>
                                                                <th className="text-left px-2 py-1">Name</th>
                                                                <th className="text-left px-2 py-1">Type</th>
                                                                <th className="text-left px-2 py-1">Content</th>
                                                                <th className="text-left px-2 py-1">TTL</th>
                                                                <th className="text-left px-2 py-1">Actions</th>
                                                            </tr>
                                                        </thead>
                                                        <tbody>
                                                            {domainRecords[domain.id].map(record => (
                                                                <tr key={record.id} className="border-t">
                                                                    {editingRecord?.id === record.id ? (
                                                                        <>
                                                                            <td className="px-2 py-1"><input className="border px-1 w-full text-sm" value={editingRecord.name} onChange={e => setEditingRecord({ ...editingRecord, name: e.target.value })} /></td>
                                                                            <td className="px-2 py-1"><select className="border px-1 text-sm" value={editingRecord.type} onChange={e => setEditingRecord({ ...editingRecord, type: e.target.value })}>{RECORD_TYPES.map(t => <option key={t.value} value={t.value}>{t.value}</option>)}</select></td>
                                                                            <td className="px-2 py-1"><input className="border px-1 w-full text-sm" value={editingRecord.content} onChange={e => setEditingRecord({ ...editingRecord, content: e.target.value })} /></td>
                                                                            <td className="px-2 py-1"><input type="number" className="border px-1 w-16 text-sm" value={editingRecord.ttl} onChange={e => setEditingRecord({ ...editingRecord, ttl: parseInt(e.target.value) || 3600 })} /></td>
                                                                            <td className="px-2 py-1 flex gap-1"><button onClick={handleUpdateRecord} className="text-green-600">Save</button><button onClick={() => setEditingRecord(null)} className="text-gray-600">Cancel</button></td>
                                                                        </>
                                                                    ) : (
                                                                        <>
                                                                            <td className="px-2 py-1 font-mono">{record.name}</td>
                                                                            <td className="px-2 py-1"><span className="bg-blue-100 text-blue-800 px-1 rounded">{record.type}</span></td>
                                                                            <td className="px-2 py-1 font-mono text-xs">{record.content}</td>
                                                                            <td className="px-2 py-1">{record.ttl}s</td>
                                                                            <td className="px-2 py-1 flex gap-2"><button onClick={() => setEditingRecord({ ...record })} className="text-blue-600">Edit</button><button onClick={() => handleDeleteRecord(record.id, domain.id)} className="text-red-600">Delete</button></td>
                                                                        </>
                                                                    )}
                                                                </tr>
                                                            ))}
                                                        </tbody>
                                                    </table>
                                                ) : (
                                                    <p className="text-gray-500 text-sm italic">No records yet.</p>
                                                )}
                                            </div>
                                        )}
                                    </li>
                                ))}
                                {domains.length === 0 && !loading && <li className="p-4 text-center text-gray-500">No domains found.</li>}
                            </ul>
                        </div>
                    </>
                )}

                {/* USERS TAB */}
                {activeTab === 'users' && isAdmin && (
                    <div className="bg-white shadow sm:rounded-lg overflow-hidden">
                        <div className="px-4 py-3 bg-gray-50 border-b flex justify-between items-center">
                            <div>
                                <h3 className="font-medium">User Management ({users.length})</h3>
                                <p className="text-sm text-gray-500">User contact info is used for domain WHOIS data</p>
                            </div>
                            <button onClick={() => setEditingUser({ role: 'user', username: '', password: '' })} className="bg-blue-600 text-white px-3 py-1 rounded text-sm hover:bg-blue-700">Add User</button>
                        </div>
                        <table className="w-full">
                            <thead className="bg-gray-50">
                                <tr>
                                    <th className="text-left px-4 py-2">ID</th>
                                    <th className="text-left px-4 py-2">Username</th>
                                    <th className="text-left px-4 py-2">Role</th>
                                    <th className="text-left px-4 py-2">Contact Name</th>
                                    <th className="text-left px-4 py-2">Contact Email</th>
                                    <th className="text-left px-4 py-2">Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {users.map(u => (
                                    <tr key={u.id} className="border-t">
                                        <td className="px-4 py-2">{u.id}</td>
                                        <td className="px-4 py-2 font-medium">{u.username}</td>
                                        <td className="px-4 py-2"><span className={`px-2 py-0.5 rounded text-xs ${u.role === 'admin' ? 'bg-purple-100 text-purple-800' : 'bg-gray-100 text-gray-800'}`}>{u.role}</span></td>
                                        <td className="px-4 py-2 text-sm">{u.contact_name || <span className="text-gray-400">Not set</span>}</td>
                                        <td className="px-4 py-2 text-sm">{u.contact_email || <span className="text-gray-400">Not set</span>}</td>
                                        <td className="px-4 py-2 flex gap-2">
                                            <button onClick={() => setEditingUser({ ...u })} className="text-blue-600 text-sm">Edit</button>
                                            {u.id !== user.id && <button onClick={() => handleDeleteUser(u.id)} className="text-red-600 text-sm">Delete</button>}
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}

                {/* CONFIG TAB */}
                {activeTab === 'config' && isAdmin && (
                    <div className="bg-white shadow sm:rounded-lg p-6">
                        <h2 className="text-lg font-medium mb-4">Registrar Configuration</h2>
                        <p className="text-sm text-gray-500 mb-6">These settings appear in WHOIS responses for all domains. Follows RFC 3912 and ICANN formatting standards.</p>
                        <form onSubmit={handleUpdateConfig} className="grid grid-cols-2 gap-4">
                            <div><label className="block text-sm font-medium text-gray-700">Registrar Name</label><input className="mt-1 block w-full border rounded px-3 py-2" value={registrarConfig.registrar_name || ''} onChange={e => setRegistrarConfig({ ...registrarConfig, registrar_name: e.target.value })} /></div>
                            <div><label className="block text-sm font-medium text-gray-700">Registrar URL</label><input className="mt-1 block w-full border rounded px-3 py-2" value={registrarConfig.registrar_url || ''} onChange={e => setRegistrarConfig({ ...registrarConfig, registrar_url: e.target.value })} /></div>
                            <div><label className="block text-sm font-medium text-gray-700">Registrar Email</label><input className="mt-1 block w-full border rounded px-3 py-2" value={registrarConfig.registrar_email || ''} onChange={e => setRegistrarConfig({ ...registrarConfig, registrar_email: e.target.value })} /></div>
                            <div><label className="block text-sm font-medium text-gray-700">Registrar Phone</label><input className="mt-1 block w-full border rounded px-3 py-2" value={registrarConfig.registrar_phone || ''} onChange={e => setRegistrarConfig({ ...registrarConfig, registrar_phone: e.target.value })} /></div>
                            <div><label className="block text-sm font-medium text-gray-700">Registrar IANA ID</label><input className="mt-1 block w-full border rounded px-3 py-2" value={registrarConfig.registrar_iana_id || '9999'} onChange={e => setRegistrarConfig({ ...registrarConfig, registrar_iana_id: e.target.value })} placeholder="9999" /></div>
                            <div><label className="block text-sm font-medium text-gray-700">Abuse Contact Email</label><input className="mt-1 block w-full border rounded px-3 py-2" value={registrarConfig.abuse_contact_email || ''} onChange={e => setRegistrarConfig({ ...registrarConfig, abuse_contact_email: e.target.value })} placeholder="abuse@localdns.local" /></div>
                            <div><label className="block text-sm font-medium text-gray-700">Abuse Contact Phone</label><input className="mt-1 block w-full border rounded px-3 py-2" value={registrarConfig.abuse_contact_phone || ''} onChange={e => setRegistrarConfig({ ...registrarConfig, abuse_contact_phone: e.target.value })} placeholder="+1-555-0199" /></div>
                            <div><label className="block text-sm font-medium text-gray-700">WHOIS Server</label><input className="mt-1 block w-full border rounded px-3 py-2" value={registrarConfig.whois_server || ''} onChange={e => setRegistrarConfig({ ...registrarConfig, whois_server: e.target.value })} /></div>
                            <div><label className="block text-sm font-medium text-gray-700">Nameserver 1</label><input className="mt-1 block w-full border rounded px-3 py-2" value={registrarConfig.nameserver1 || ''} onChange={e => setRegistrarConfig({ ...registrarConfig, nameserver1: e.target.value })} /></div>
                            <div><label className="block text-sm font-medium text-gray-700">Nameserver 2</label><input className="mt-1 block w-full border rounded px-3 py-2" value={registrarConfig.nameserver2 || ''} onChange={e => setRegistrarConfig({ ...registrarConfig, nameserver2: e.target.value })} /></div>
                            <div><label className="block text-sm font-medium text-gray-700">Default TTL (seconds)</label><input type="number" className="mt-1 block w-full border rounded px-3 py-2" value={registrarConfig.default_ttl || 3600} onChange={e => setRegistrarConfig({ ...registrarConfig, default_ttl: parseInt(e.target.value) })} /></div>
                            <div><label className="block text-sm font-medium text-gray-700">Default Expiry (days)</label><input type="number" className="mt-1 block w-full border rounded px-3 py-2" value={registrarConfig.default_expiry_days || 365} onChange={e => setRegistrarConfig({ ...registrarConfig, default_expiry_days: parseInt(e.target.value) })} /></div>
                            <div className="col-span-2"><button type="submit" className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700">Save Configuration</button></div>
                        </form>
                    </div>
                )}
            </main>

            {/* Registrant Modal - READ ONLY (data comes from User) */}
            {editingRegistrant && (
                <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
                    <div className="bg-white rounded-lg p-6 max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
                        <h2 className="text-lg font-medium mb-2">WHOIS Info - {editingRegistrant.name}</h2>
                        <p className="text-sm text-gray-500 mb-4">
                            Contact data is inherited from the domain owner's profile.
                            To edit, go to <strong>Users tab → Edit</strong> the owner's contact info.
                        </p>

                        <div className="grid grid-cols-2 gap-4 mb-4">
                            <div className="col-span-2 bg-gray-50 p-3 rounded">
                                <p className="text-sm"><strong>Owner:</strong> {editingRegistrant.user?.username || `User #${editingRegistrant.user_id}`}</p>
                                <p className="text-sm"><strong>Created:</strong> {new Date(editingRegistrant.created_at).toLocaleString()}</p>
                                <p className="text-sm"><strong>Updated:</strong> {new Date(editingRegistrant.updated_at).toLocaleString()}</p>
                                <p className="text-sm"><strong>Expires:</strong> {
                                    (() => {
                                        if (!editingRegistrant.expires_at) return 'N/A';
                                        const expiresDate = new Date(editingRegistrant.expires_at);
                                        // Check if date is valid and not zero/invalid date
                                        if (isNaN(expiresDate.getTime()) || expiresDate.getFullYear() < 1970) {
                                            // Fallback: calculate from created_at + 1 year
                                            if (editingRegistrant.created_at) {
                                                const created = new Date(editingRegistrant.created_at);
                                                const expires = new Date(created);
                                                expires.setFullYear(expires.getFullYear() + 1);
                                                return expires.toLocaleString();
                                            }
                                            return 'N/A';
                                        }
                                        return expiresDate.toLocaleString();
                                    })()
                                }</p>
                                <p className="text-sm"><strong>Status:</strong> <span className="px-2 py-0.5 rounded text-xs bg-green-100 text-green-800">{editingRegistrant.status || 'active'}</span></p>
                            </div>
                        </div>

                        <h3 className="font-medium text-blue-600 mb-2">Registrant Contact (from owner profile)</h3>
                        <div className="grid grid-cols-2 gap-2 text-sm mb-4 bg-gray-50 p-3 rounded">
                            <p><strong>Name:</strong> {editingRegistrant.registrant_name || editingRegistrant.user?.contact_name || 'Not set'}</p>
                            <p><strong>Org:</strong> {editingRegistrant.registrant_org || editingRegistrant.user?.contact_org || 'Not set'}</p>
                            <p><strong>Email:</strong> {editingRegistrant.registrant_email || editingRegistrant.user?.contact_email || 'Not set'}</p>
                            <p><strong>Phone:</strong> {editingRegistrant.registrant_phone || editingRegistrant.user?.contact_phone || 'Not set'}</p>
                            <p className="col-span-2"><strong>Address:</strong> {
                                editingRegistrant.registrant_address || editingRegistrant.user?.contact_address || 'Not set'
                            }, {
                                editingRegistrant.registrant_city || editingRegistrant.user?.contact_city || '-'
                            }, {
                                editingRegistrant.registrant_state || editingRegistrant.user?.contact_state || '-'
                            } {
                                editingRegistrant.registrant_zip || editingRegistrant.user?.contact_zip || '-'
                            }, {
                                editingRegistrant.registrant_country || editingRegistrant.user?.contact_country || '-'
                            }</p>
                        </div>

                        <div className="flex gap-4 mt-6">
                            <button type="button" onClick={() => setEditingRegistrant(null)} className="bg-gray-200 text-gray-800 px-6 py-2 rounded hover:bg-gray-300">Close</button>
                        </div>
                    </div>
                </div>
            )}

            {/* User Edit Modal */}
            {editingUser && (
                <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
                    <div className="bg-white rounded-lg p-6 max-w-3xl w-full mx-4 max-h-[90vh] overflow-y-auto">
                        <h2 className="text-lg font-medium mb-2">{editingUser.id ? 'Edit User' : 'Create New User'} - {editingUser.username}</h2>
                        <p className="text-sm text-gray-500 mb-4">Contact info will be used for WHOIS data.</p>

                        <form onSubmit={handleUpdateUser}>
                            {/* User Basic Info */}
                            <h3 className="font-medium text-blue-600 mb-2">Account Info</h3>
                            <div className="grid grid-cols-2 gap-4 mb-4">
                                <div><label className="block text-sm font-medium text-gray-700">Username</label><input className="mt-1 block w-full border rounded px-3 py-2" value={editingUser.username || ''} onChange={e => setEditingUser({ ...editingUser, username: e.target.value })} required /></div>
                                <div><label className="block text-sm font-medium text-gray-700">Role</label><select className="mt-1 block w-full border rounded px-3 py-2" value={editingUser.role || 'user'} onChange={e => setEditingUser({ ...editingUser, role: e.target.value })}><option value="user">user</option><option value="admin">admin</option></select></div>
                                {!editingUser.id && (
                                    <div className="col-span-2"><label className="block text-sm font-medium text-gray-700">Password</label><input type="password" className="mt-1 block w-full border rounded px-3 py-2" value={editingUser.password || ''} onChange={e => setEditingUser({ ...editingUser, password: e.target.value })} required /></div>
                                )}
                            </div>

                            {/* Contact Info */}
                            <h3 className="font-medium text-green-600 mb-2 mt-4">Contact Info (for WHOIS)</h3>
                            <div className="grid grid-cols-2 gap-4 mb-4">
                                <div><label className="block text-sm font-medium text-gray-700">Full Name</label><input className="mt-1 block w-full border rounded px-3 py-2" value={editingUser.contact_name || ''} onChange={e => setEditingUser({ ...editingUser, contact_name: e.target.value })} placeholder="John Doe" /></div>
                                <div><label className="block text-sm font-medium text-gray-700">Organization</label><input className="mt-1 block w-full border rounded px-3 py-2" value={editingUser.contact_org || ''} onChange={e => setEditingUser({ ...editingUser, contact_org: e.target.value })} placeholder="My Homelab Inc." /></div>
                                <div><label className="block text-sm font-medium text-gray-700">Email</label><input type="email" className="mt-1 block w-full border rounded px-3 py-2" value={editingUser.contact_email || ''} onChange={e => setEditingUser({ ...editingUser, contact_email: e.target.value })} placeholder="admin@example.lan" /></div>
                                <div><label className="block text-sm font-medium text-gray-700">Phone</label><input className="mt-1 block w-full border rounded px-3 py-2" value={editingUser.contact_phone || ''} onChange={e => setEditingUser({ ...editingUser, contact_phone: e.target.value })} placeholder="+1-555-0100" /></div>
                                <div className="col-span-2"><label className="block text-sm font-medium text-gray-700">Street Address</label><input className="mt-1 block w-full border rounded px-3 py-2" value={editingUser.contact_address || ''} onChange={e => setEditingUser({ ...editingUser, contact_address: e.target.value })} placeholder="123 Lab Street" /></div>
                                <div><label className="block text-sm font-medium text-gray-700">City</label><input className="mt-1 block w-full border rounded px-3 py-2" value={editingUser.contact_city || ''} onChange={e => setEditingUser({ ...editingUser, contact_city: e.target.value })} placeholder="Tech City" /></div>
                                <div><label className="block text-sm font-medium text-gray-700">State/Province</label><input className="mt-1 block w-full border rounded px-3 py-2" value={editingUser.contact_state || ''} onChange={e => setEditingUser({ ...editingUser, contact_state: e.target.value })} placeholder="CA" /></div>
                                <div><label className="block text-sm font-medium text-gray-700">Postal Code</label><input className="mt-1 block w-full border rounded px-3 py-2" value={editingUser.contact_zip || ''} onChange={e => setEditingUser({ ...editingUser, contact_zip: e.target.value })} placeholder="90210" /></div>
                                <div><label className="block text-sm font-medium text-gray-700">Country</label><input className="mt-1 block w-full border rounded px-3 py-2" value={editingUser.contact_country || ''} onChange={e => setEditingUser({ ...editingUser, contact_country: e.target.value })} placeholder="US" /></div>
                            </div>

                            <div className="flex gap-4 mt-6">
                                <button type="submit" className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700">Save User</button>
                                <button type="button" onClick={() => setEditingUser(null)} className="bg-gray-200 text-gray-800 px-6 py-2 rounded hover:bg-gray-300">Cancel</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
}

