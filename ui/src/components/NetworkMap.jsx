import React, { useEffect, useState, useRef } from 'react';
import ForceGraph2D from 'react-force-graph-2d';
import axios from 'axios';

const NetworkMap = () => {
    const [data, setData] = useState({ nodes: [], links: [] });
    const fgRef = useRef();

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await axios.get('/api/topology');
                const topo = response.data;

                const nodes = [];
                const links = [];
                const deviceSet = new Set();

                topo.links.forEach(l => {
                    deviceSet.add(l.source_device);
                    deviceSet.add(l.target_device);
                    links.push({
                        source: l.source_device,
                        target: l.target_device,
                        type: l.type,
                        source_int: l.source_interface,
                        target_int: l.target_interface
                    });
                });

                deviceSet.forEach(d => {
                    nodes.push({ id: d, name: d });
                });

                setData({ nodes, links });
            } catch (err) {
                console.error("Failed to fetch topology:", err);
            }
        };

        fetchData();
        const interval = setInterval(fetchData, 30000); // Refresh every 30s
        return () => clearInterval(interval);
    }, []);

    return (
        <div className="w-full h-[calc(100vh-80px)] bg-noc-bg rounded-xl overflow-hidden border border-white/5 relative">
            <ForceGraph2D
                ref={fgRef}
                graphData={data}
                nodeLabel="name"
                nodeColor={() => '#ffffff'}
                nodeRelSize={6}
                linkColor={(link) => {
                    if (link.type === '10g') return '#00f5ff';
                    if (link.type === '1g') return '#10b981'; // SFP Emerald
                    if (link.type === 'ethernet') return '#065f46'; // Forest Green
                    if (link.type === 'wireless') return '#facc15'; // Yellow
                    return '#666';
                }}
                linkWidth={(link) => (link.type === '10g' ? 3 : 1.5)}
                linkLineDash={(link) => (link.type === 'wireless' ? [2, 2] : null)}
                linkDirectionalParticles={(link) => (link.type === '10g' || link.type === '1g' ? 2 : 0)}
                linkDirectionalParticleSpeed={(link) => (link.type === '10g' ? 0.01 : 0.005)}
                linkDirectionalParticleWidth={2}
                backgroundColor="rgba(0,0,0,0)"
            />
            <div className="absolute top-4 left-4 flex flex-col gap-2 bg-black/40 backdrop-blur-sm p-3 rounded-lg border border-white/5">
                <div className="flex items-center gap-2 text-[10px] uppercase font-bold tracking-widest text-white/60">
                    <div className="w-2 h-2 bg-noc-cyan rounded-full pulse-glow" /> 10G SFP+
                </div>
                <div className="flex items-center gap-2 text-[10px] uppercase font-bold tracking-widest text-white/60">
                    <div className="w-2 h-2 bg-noc-emerald rounded-full" /> 1G SFP
                </div>
                <div className="flex items-center gap-2 text-[10px] uppercase font-bold tracking-widest text-white/60">
                    <div className="w-2 h-2 bg-[#065f46] rounded-full" /> 1G Ethernet
                </div>
                <div className="flex items-center gap-2 text-[10px] uppercase font-bold tracking-widest text-white/60">
                    <div className="w-2 h-2 border border-dashed border-noc-yellow rounded-full" /> Wireless
                </div>
            </div>
        </div>
    );
};

export default NetworkMap;
