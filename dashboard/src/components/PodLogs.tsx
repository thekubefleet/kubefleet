import React, { useState, useEffect, useRef } from 'react';
import {
    Box,
    Paper,
    Typography,
    List,
    ListItem,
    ListItemText,
    Chip,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    TextField,
    IconButton,
    Tooltip,
    Divider,
    Alert,
} from '@mui/material';
import {
    Refresh as RefreshIcon,
    FilterList as FilterIcon,
    Clear as ClearIcon,
    Download as DownloadIcon,
} from '@mui/icons-material';

interface PodLog {
    namespace: string;
    pod_name: string;
    container_name: string;
    log_line: string;
    timestamp: number;
    level: string;
}

interface PodLogsProps {
    namespace: string;
    podName: string;
    onClose?: () => void;
}

const PodLogs: React.FC<PodLogsProps> = ({ namespace, podName, onClose }) => {
    const [logs, setLogs] = useState<PodLog[]>([]);
    const [containers, setContainers] = useState<string[]>([]);
    const [selectedContainer, setSelectedContainer] = useState<string>('');
    const [filterLevel, setFilterLevel] = useState<string>('ALL');
    const [searchTerm, setSearchTerm] = useState<string>('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string>('');
    const [autoRefresh, setAutoRefresh] = useState(true);
    const logsEndRef = useRef<HTMLDivElement>(null);

    const logLevels = ['ALL', 'ERROR', 'WARN', 'INFO', 'DEBUG'];

    const scrollToBottom = () => {
        logsEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    };

    useEffect(() => {
        scrollToBottom();
    }, [logs]);

    useEffect(() => {
        if (namespace && podName) {
            fetchLogs();
            if (autoRefresh) {
                const interval = setInterval(fetchLogs, 5000); // Refresh every 5 seconds
                return () => clearInterval(interval);
            }
        }
    }, [namespace, podName, selectedContainer, autoRefresh]);

    const fetchLogs = async () => {
        setLoading(true);
        setError('');

        try {
            const url = selectedContainer
                ? `/api/logs/${namespace}/${podName}/${selectedContainer}`
                : `/api/logs/${namespace}/${podName}`;

            const response = await fetch(url);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();
            setLogs(data.logs || []);

            // Extract unique containers from logs
            const containerNames = data.logs?.map((log: PodLog) => log.container_name) || [];
            const uniqueContainers = containerNames.filter((container: string, index: number) =>
                containerNames.indexOf(container) === index
            );
            setContainers(uniqueContainers);

            if (uniqueContainers.length > 0 && !selectedContainer) {
                setSelectedContainer(uniqueContainers[0]);
            }
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to fetch logs');
        } finally {
            setLoading(false);
        }
    };

    const getLevelColor = (level: string) => {
        switch (level.toUpperCase()) {
            case 'ERROR':
                return 'error';
            case 'WARN':
                return 'warning';
            case 'DEBUG':
                return 'info';
            default:
                return 'default';
        }
    };

    const filteredLogs = logs.filter(log => {
        const matchesLevel = filterLevel === 'ALL' || log.level.toUpperCase() === filterLevel;
        const matchesSearch = searchTerm === '' ||
            log.log_line.toLowerCase().includes(searchTerm.toLowerCase()) ||
            log.container_name.toLowerCase().includes(searchTerm.toLowerCase());
        return matchesLevel && matchesSearch;
    });

    const formatTimestamp = (timestamp: number) => {
        return new Date(timestamp * 1000).toLocaleString();
    };

    const downloadLogs = () => {
        const logText = filteredLogs
            .map(log => `[${formatTimestamp(log.timestamp)}] [${log.level}] [${log.container_name}] ${log.log_line}`)
            .join('\n');

        const blob = new Blob([logText], { type: 'text/plain' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `${namespace}-${podName}-logs.txt`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    };

    const clearLogs = () => {
        setLogs([]);
    };

    return (
        <Paper sx={{ p: 2, height: '100%', display: 'flex', flexDirection: 'column' }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                <Typography variant="h6" component="h2">
                    Pod Logs: {podName}
                </Typography>
                <Box sx={{ display: 'flex', gap: 1 }}>
                    <Tooltip title="Refresh logs">
                        <IconButton onClick={fetchLogs} disabled={loading}>
                            <RefreshIcon />
                        </IconButton>
                    </Tooltip>
                    <Tooltip title="Download logs">
                        <IconButton onClick={downloadLogs} disabled={logs.length === 0}>
                            <DownloadIcon />
                        </IconButton>
                    </Tooltip>
                    <Tooltip title="Clear logs">
                        <IconButton onClick={clearLogs} disabled={logs.length === 0}>
                            <ClearIcon />
                        </IconButton>
                    </Tooltip>
                </Box>
            </Box>

            {error && (
                <Alert severity="error" sx={{ mb: 2 }}>
                    {error}
                </Alert>
            )}

            <Box sx={{ display: 'flex', gap: 2, mb: 2 }}>
                <FormControl size="small" sx={{ minWidth: 150 }}>
                    <InputLabel>Container</InputLabel>
                    <Select
                        value={selectedContainer}
                        label="Container"
                        onChange={(e) => setSelectedContainer(e.target.value)}
                    >
                        <MenuItem value="">All Containers</MenuItem>
                        {containers.map((container) => (
                            <MenuItem key={container} value={container}>
                                {container}
                            </MenuItem>
                        ))}
                    </Select>
                </FormControl>

                <FormControl size="small" sx={{ minWidth: 120 }}>
                    <InputLabel>Level</InputLabel>
                    <Select
                        value={filterLevel}
                        label="Level"
                        onChange={(e) => setFilterLevel(e.target.value)}
                    >
                        {logLevels.map((level) => (
                            <MenuItem key={level} value={level}>
                                {level}
                            </MenuItem>
                        ))}
                    </Select>
                </FormControl>

                <TextField
                    size="small"
                    label="Search"
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    sx={{ minWidth: 200 }}
                />
            </Box>

            <Divider sx={{ mb: 2 }} />

            <Box sx={{ flex: 1, overflow: 'auto', bgcolor: 'grey.900', borderRadius: 1 }}>
                <List dense sx={{ p: 0 }}>
                    {filteredLogs.length === 0 ? (
                        <ListItem>
                            <ListItemText
                                primary={
                                    <Typography variant="body2" color="text.secondary">
                                        {loading ? 'Loading logs...' : 'No logs found'}
                                    </Typography>
                                }
                            />
                        </ListItem>
                    ) : (
                        filteredLogs.map((log, index) => (
                            <ListItem key={index} sx={{ borderBottom: '1px solid rgba(255, 255, 255, 0.1)' }}>
                                <ListItemText
                                    primary={
                                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 0.5 }}>
                                            <Typography variant="caption" color="text.secondary">
                                                {formatTimestamp(log.timestamp)}
                                            </Typography>
                                            <Chip
                                                label={log.level}
                                                size="small"
                                                color={getLevelColor(log.level) as any}
                                                variant="outlined"
                                            />
                                            <Typography variant="caption" color="text.secondary">
                                                {log.container_name}
                                            </Typography>
                                        </Box>
                                    }
                                    secondary={
                                        <Typography
                                            variant="body2"
                                            component="pre"
                                            sx={{
                                                fontFamily: 'monospace',
                                                fontSize: '0.875rem',
                                                whiteSpace: 'pre-wrap',
                                                wordBreak: 'break-word',
                                                m: 0,
                                                color: 'text.primary',
                                            }}
                                        >
                                            {log.log_line}
                                        </Typography>
                                    }
                                />
                            </ListItem>
                        ))
                    )}
                    <div ref={logsEndRef} />
                </List>
            </Box>

            <Box sx={{ mt: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Typography variant="caption" color="text.secondary">
                    {filteredLogs.length} of {logs.length} log entries
                </Typography>
                <Chip
                    label={autoRefresh ? 'Auto-refresh ON' : 'Auto-refresh OFF'}
                    size="small"
                    color={autoRefresh ? 'success' : 'default'}
                    onClick={() => setAutoRefresh(!autoRefresh)}
                    clickable
                />
            </Box>
        </Paper>
    );
};

export default PodLogs; 