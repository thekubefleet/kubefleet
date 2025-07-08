import React, { useState, useEffect } from 'react';
import {
    Card,
    CardContent,
    Typography,
    Box,
    ToggleButton,
    ToggleButtonGroup,
} from '@mui/material';
import {
    LineChart,
    Line,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    Legend,
    ResponsiveContainer,
    BarChart,
    Bar,
} from 'recharts';

interface MetricData {
    name: string;
    cpu: number;
    memory: number;
    timestamp: number;
}

interface MetricsChartProps {
    selected: { type: 'pod' | 'deployment'; namespace: string; name: string } | null;
}

const MetricsChart: React.FC<MetricsChartProps> = ({ selected }) => {
    const [metrics, setMetrics] = useState<MetricData[]>([]);
    const [chartType, setChartType] = useState<'line' | 'bar'>('line');
    const [metricType, setMetricType] = useState<'cpu' | 'memory'>('cpu');

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch('http://localhost:3000/api/data');
                const result = await response.json();

                if (result.data) {
                    const metricData: MetricData[] = [];

                    result.data.forEach((dataPoint: any) => {
                        if (dataPoint.metrics) {
                            dataPoint.metrics.forEach((metric: any) => {
                                metricData.push({
                                    name: `${metric.namespace}/${metric.name}`,
                                    cpu: metric.cpu || 0,
                                    memory: metric.memory || 0,
                                    timestamp: dataPoint.timestamp,
                                });
                            });
                        }
                    });

                    setMetrics(metricData);
                }
            } catch (error) {
                console.error('Error fetching metrics data:', error);
            }
        };

        fetchData();
        const interval = setInterval(fetchData, 30000);

        return () => clearInterval(interval);
    }, []);

    const handleChartTypeChange = (
        event: React.MouseEvent<HTMLElement>,
        newChartType: 'line' | 'bar' | null,
    ) => {
        if (newChartType !== null) {
            setChartType(newChartType);
        }
    };

    const handleMetricTypeChange = (
        event: React.MouseEvent<HTMLElement>,
        newMetricType: 'cpu' | 'memory' | null,
    ) => {
        if (newMetricType !== null) {
            setMetricType(newMetricType);
        }
    };

    const filteredMetrics = selected
        ? metrics.filter(m => {
            const [ns, name] = m.name.split('/');
            return ns === selected.namespace && name === selected.name;
        })
        : metrics;

    return (
        <Card>
            <CardContent>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                    <Typography variant="h5" component="h2">
                        Performance Metrics
                    </Typography>

                    <Box sx={{ display: 'flex', gap: 1 }}>
                        <ToggleButtonGroup
                            value={chartType}
                            exclusive
                            onChange={handleChartTypeChange}
                            size="small"
                        >
                            <ToggleButton value="line">Line</ToggleButton>
                            <ToggleButton value="bar">Bar</ToggleButton>
                        </ToggleButtonGroup>

                        <ToggleButtonGroup
                            value={metricType}
                            exclusive
                            onChange={handleMetricTypeChange}
                            size="small"
                        >
                            <ToggleButton value="cpu">CPU</ToggleButton>
                            <ToggleButton value="memory">Memory</ToggleButton>
                        </ToggleButtonGroup>
                    </Box>
                </Box>

                {filteredMetrics.length === 0 ? (
                    <Typography color="text.secondary" align="center">
                        No metrics data available
                    </Typography>
                ) : (
                    <ResponsiveContainer width="100%" height={300}>
                        {chartType === 'line' ? (
                            <LineChart data={filteredMetrics}>
                                <CartesianGrid strokeDasharray="3 3" />
                                <XAxis dataKey="name" />
                                <YAxis />
                                <Tooltip />
                                <Legend />
                                <Line
                                    type="monotone"
                                    dataKey={metricType}
                                    stroke="#8884d8"
                                    strokeWidth={2}
                                />
                            </LineChart>
                        ) : (
                            <BarChart data={filteredMetrics}>
                                <CartesianGrid strokeDasharray="3 3" />
                                <XAxis dataKey="name" />
                                <YAxis />
                                <Tooltip />
                                <Legend />
                                <Bar dataKey={metricType} fill="#8884d8" />
                            </BarChart>
                        )}
                    </ResponsiveContainer>
                )}

                <Box sx={{ mt: 2 }}>
                    <Typography variant="body2" color="text.secondary">
                        Showing {filteredMetrics.length} metrics • Last updated: {new Date().toLocaleTimeString()}
                    </Typography>
                </Box>
            </CardContent>
        </Card>
    );
};

export default MetricsChart; 