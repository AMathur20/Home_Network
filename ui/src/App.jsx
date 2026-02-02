import React, { useState } from 'react';
import { Activity, LayoutGrid, Map as MapIcon, Settings } from 'lucide-react';
import NetworkMap from './components/NetworkMap';
import PriorityList from './components/PriorityList';

function App() {
  const [view, setView] = useState('map');

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
          <StatCard title="Devices" value="3" icon={<Activity className="text-noc-cyan w-4 h-4" />} />
          <StatCard title="Links" value="12" icon={<Activity className="text-noc-emerald w-4 h-4" />} />
          <StatCard title="Speed" value="450 Mbps" icon={<Activity className="text-noc-yellow w-4 h-4" />} />
          <StatCard title="Uptime" value="14d" icon={<Activity className="text-noc-forest w-4 h-4" />} />
        </div>

        <div className="flex flex-col lg:flex-row gap-6">
          <div className="flex-1">
            {view === 'map' ? (
              <NetworkMap />
            ) : (
              <div className="bg-noc-card rounded-xl border border-white/5 aspect-video flex items-center justify-center text-white/20">
                Grid View Placeholder
              </div>
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
