import React, { useState } from 'react';
import { Activity, LayoutGrid, Map as MapIcon, Settings } from 'lucide-react';
import NetworkMap from './components/NetworkMap';
import PriorityList from './components/PriorityList';

const App = () => {
  const [view, setView] = useState('map');
  const [stats, setStats] = useState({ devices: 0, links: 0, speed: '0 Mbps', offline: 0 });
  const [metrics, setMetrics] = useState([]);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const [topoRes, metricsRes] = await Promise.all([
          axios.get('/api/topology'),
          axios.get('/api/metrics/live')
        ]);
        
        const mData = metricsRes.data || [];
        const tData = topoRes.data || { links: [] };
        
        const devices = new Set();
        mData.forEach(m => devices.add(m.DeviceName));
        
        let totalBps = 0;
        let offline = 0;
        mData.forEach(m => {
          if (m.Status === 'down') offline++;
          totalBps += m.InSpeed + m.OutSpeed;
        });

        setStats({
          devices: devices.size || tData.links.length ? 3 : 0, // Fallback if no metrics yet
          links: tData.links.length,
          speed: `${(totalBps / 1000000).toFixed(1)} Mbps`,
          offline
        });
        setMetrics(mData);
      } catch (err) {
        console.error("Failed to fetch stats:", err);
      }
    };

    fetchStats();
    const interval = setInterval(fetchStats, 5000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="min-h-screen bg-noc-bg text-white flex flex-col">
      {/* Header */}
      <header className="h-16 border-b border-white/5 flex items-center justify-between px-6 bg-noc-bg/80 backdrop-blur-md sticky top-0 z-50">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 bg-noc-cyan rounded-lg flex items-center justify-center">
            <Activity className="w-5 h-5 text-noc-bg" />
          </div>
          <h1 className="font-bold text-lg tracking-tight">HNM <span className="text-noc-cyan">v2</span></h1>
        </div>

        <nav className="hidden md:flex items-center bg-noc-card rounded-lg p-1 border border-white/5">
          <button
            onClick={() => setView('map')}
            className={`flex items-center gap-2 px-3 py-1.5 rounded-md transition-all ${view === 'map' ? 'bg-white/10 text-white' : 'text-white/40 hover:text-white'}`}
          >
            <MapIcon className="w-4 h-4" />
            <span className="text-sm font-medium">Map</span>
          </button>
          <button
            onClick={() => setView('grid')}
            className={`flex items-center gap-2 px-3 py-1.5 rounded-md transition-all ${view === 'grid' ? 'bg-white/10 text-white' : 'text-white/40 hover:text-white'}`}
          >
            <LayoutGrid className="w-4 h-4" />
            <span className="text-sm font-medium">Grid</span>
          </button>
        </nav>

        <div className="flex items-center gap-4">
          <div className="hidden sm:block px-3 py-1 bg-noc-emerald/10 text-noc-emerald rounded-full text-[10px] font-bold uppercase tracking-wider border border-noc-emerald/20">
            System Online
          </div>
          <button className="p-2 text-white/40 hover:text-white transition-colors">
            <Settings className="w-5 h-5" />
          </button>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 p-4 md:p-6 flex flex-col gap-6">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <StatCard title="Devices" value={stats.devices} icon={<Activity className="text-noc-cyan w-4 h-4" />} />
          <StatCard title="Links" value={stats.links} icon={<Activity className="text-noc-emerald w-4 h-4" />} />
          <StatCard title="Throughput" value={stats.speed} icon={<Activity className="text-noc-yellow w-4 h-4" />} />
          <StatCard title="Status" value={stats.offline > 0 ? `${stats.offline} Down` : 'All Up'} icon={<Activity className="text-noc-forest w-4 h-4" />} />
        </div>

        <div className="flex flex-col lg:flex-row gap-6">
          <div className="flex-1">
            {view === 'map' ? (
              <NetworkMap />
            ) : (
              <GridView metrics={metrics} />
            )}
          </div>

          <div className="lg:w-80">
            <PriorityList />
          </div>
        </div>
      </main>
    </div>
  );
}

function GridView({ metrics }) {
  return (
    <div className="bg-noc-card rounded-xl border border-white/5 overflow-hidden">
      <table className="w-full text-left border-collapse">
        <thead className="bg-white/5 text-[10px] uppercase font-bold tracking-widest text-white/40">
          <tr>
            <th className="p-4">Device</th>
            <th className="p-4">Interface</th>
            <th className="p-4">Status</th>
            <th className="p-4 text-right">In Speed</th>
            <th className="p-4 text-right">Out Speed</th>
          </tr>
        </thead>
        <tbody className="text-sm">
          {metrics.map((m, i) => (
            <tr key={i} className="border-t border-white/5 hover:bg-white/5 transition-colors">
              <td className="p-4 font-medium">{m.DeviceName}</td>
              <td className="p-4 text-white/60">{m.InterfaceName}</td>
              <td className="p-4">
                <span className={`px-2 py-0.5 rounded-full text-[10px] font-bold uppercase ${m.Status === 'up' ? 'bg-noc-emerald/10 text-noc-emerald' : 'bg-red-500/10 text-red-500'}`}>
                  {m.Status}
                </span>
              </td>
              <td className="p-4 text-right font-mono">{(m.InSpeed / 1000000).toFixed(2)} Mbps</td>
              <td className="p-4 text-right font-mono">{(m.OutSpeed / 1000000).toFixed(2)} Mbps</td>
            </tr>
          ))}
          {metrics.length === 0 && (
            <tr>
              <td colSpan="5" className="p-12 text-center text-white/20">No active ports detected</td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}

function StatCard({ title, value, icon }) {
  return (
    <div className="bg-noc-card border border-white/5 p-4 rounded-xl flex items-center justify-between">
      <div>
        <p className="text-white/40 text-xs font-medium uppercase tracking-wider">{title}</p>
        <p className="text-xl font-bold mt-1">{value}</p>
      </div>
      <div className="p-2 bg-white/5 rounded-lg">
        {icon}
      </div>
    </div>
  );
}

export default App;
