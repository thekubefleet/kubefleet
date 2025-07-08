import React, { useState, useEffect } from 'react';
import { Card, CardContent, Typography, Box } from '@mui/material';
import { Storage, Memory, Speed, Timeline } from '@mui/icons-material';

interface ClusterData {
    namespaces: number;
    totalPods: number;
    totalDeployments: number;
    lastUpdate: string;
}

const ClusterOverview: React.FC = () => {
    const [clusterData, setClusterData] = useState<ClusterData>({
        namespaces: 0,
        totalPods: 0,
        totalDeployments: 0,
        lastUpdate: 'Never',
    });

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch('/api/data/latest');
                const result = await response.json();

                if (result.data) {
                    const data = result.data;
                    const totalPods = data.resources?.reduce((acc: number, resource: any) =>
                        acc + (resource.pods?.length || 0), 0) || 0;
                    const totalDeployments = data.resources?.reduce((acc: number, resource: any) =>
                        acc + (resource.deployments?.length || 0), 0) || 0;

                    setClusterData({
                        namespaces: data.resources?.length || 0,
                        totalPods,
                        totalDeployments,
                        lastUpdate: new Date(data.timestamp * 1000).toLocaleString(),
                    });
                }
            } catch (error) {
                console.error('Error fetching cluster data:', error);
            }
        };

        fetchData();
        const interval = setInterval(fetchData, 30000); // Update every 30 seconds

        return () => clearInterval(interval);
    }, []);

    const stats = [
        {
            title: 'Namespaces',
            value: clusterData.namespaces,
            icon: <Storage />,
            color: 'primary' as const,
        },
        {
            title: 'Total Pods',
            value: clusterData.totalPods,
            icon: <Memory />,
            color: 'secondary' as const,
        },
        {
            title: 'Deployments',
            value: clusterData.totalDeployments,
            icon: <Speed />,
            color: 'success' as const,
        },
    ];

    return (
        <Card>
            <CardContent>
                <Typography variant="h5" component="h2" gutterBottom>
                    Cluster Overview
                </Typography>

                <Box sx={{ display: 'flex', gap: 3, flexWrap: 'wrap' }}>
                    {stats.map((stat) => (
                        <Box key={stat.title} sx={{ flex: 1, minWidth: 200, textAlign: 'center' }}>
                            <Box sx={{ color: `${stat.color}.main`, mb: 1 }}>
                                {stat.icon}
                            </Box>
                            <Typography variant="h4" component="div">
                                {stat.value}
                            </Typography>
                            <Typography color="text.secondary">
                                {stat.title}
                            </Typography>
                        </Box>
                    ))}
                </Box>

                <Box sx={{ mt: 2, display: 'flex', alignItems: 'center', gap: 1 }}>
                    <Timeline fontSize="small" />
                    <Typography variant="body2" color="text.secondary">
                        Last updated: {clusterData.lastUpdate}
                    </Typography>
                </Box>
            </CardContent>
        </Card>
    );
};

export default ClusterOverview; 