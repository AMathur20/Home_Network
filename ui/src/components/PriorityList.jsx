import React, { useEffect, useState } from 'react';
import { ArrowUpRight, ArrowDownRight, AlertCircle } from 'lucide-react';
import axios from 'axios';

const PriorityList = () => {
    const [metrics, setMetrics] = useState([]);

    useEffect(() => {
        const fetchMetrics = async () => {
            try {
                const response = await axios.get('/api/metrics/live');
                // For demonstration, since the backend API is still partial
                // we will mock a few "hot" links if the response is just a status ok
                if (response.data.status === 'ok') {
                    setMetrics([
                        { id: 1, device: 'mikrotik-core', interface: 'sfp-sfpplus1', speed: '450 Mbps', direction: 'up', status: 'up' },
                        { id: 2, device: 'mikrotik-core', interface: 'ether2', speed: '12 Mbps', direction: 'down', status: 'up' },
                        { id: 3, device: 'living-room-sw', interface: 'port5', speed: '0 bps', direction: 'none', status: 'down' },
                    ]);
                } else {
                    setMetrics(response.data);
                }
            } catch (err) {
                console.error("Failed to fetch metrics:", err);
            }
        };

        fetchMetrics();
        const interval = setInterval(fetchMetrics, 5000);
        return () => clearInterval(interval);
    }, []);

    return (
        <div className="flex flex-col gap-3">
            <h2 className="text-xs font-bold text-white/40 uppercase tracking-widest px-1">Active Links</h2>
            <div className="flex flex-col gap-2">
                {metrics.map((m) => (
                    <div key={m.id} className={`bg-noc-card border border-white/5 p-3 rounded-xl flex items-center justify-between transition-all hover:border-white/10 ${m.status === 'down' ? 'border-red-500/20 bg-red-500/5' : ''}`}>
                        <div className="flex items-center gap-3">
                            <div className={`p-2 rounded-lg ${m.status === 'down' ? 'bg-red-500/20 text-red-500' : 'bg-white/5 text-white/60'}`}>
                                {m.status === 'down' ? <AlertCircle className="w-4 h-4" /> : m.direction === 'up' ? <ArrowUpRight className="w-4 h-4 text-noc-cyan" /> : <ArrowDownRight className="w-4 h-4 text-noc-emerald" />}
                            </div>
                            <div>
                                <p className="text-sm font-bold">{m.device}</p>
                                <p className="text-[10px] text-white/40 font-medium">{m.interface}</p>
                            </div>
                        </div>
                        <div className="text-right">
                            <p className={`text-sm font-mono font-bold ${m.status === 'down' ? 'text-red-500' : 'text-noc-cyan'}`}>
                                {m.status === 'down' ? 'OFFLINE' : m.speed}
                            </p>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default PriorityList;
