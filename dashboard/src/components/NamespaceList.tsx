import React, { useState, useEffect } from 'react';
import {
    Card,
    CardContent,
    Typography,
    List,
    ListItem,
    ListItemText,
    ListItemSecondaryAction,
    Chip,
    Box,
    Accordion,
    AccordionSummary,
    AccordionDetails,
} from '@mui/material';
import { ExpandMore, Storage, Memory } from '@mui/icons-material';

interface NamespaceData {
    namespace: string;
    pods: string[];
    deployments: string[];
}

interface NamespaceListProps {
    selected: { type: 'pod' | 'deployment'; namespace: string; name: string } | null;
    setSelected: (sel: { type: 'pod' | 'deployment'; namespace: string; name: string } | null) => void;
}

const NamespaceList: React.FC<NamespaceListProps> = ({ selected, setSelected }) => {
    const [namespaces, setNamespaces] = useState<NamespaceData[]>([]);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch('http://localhost:3000/api/data/latest');
                const result = await response.json();

                if (result.data?.resources) {
                    const namespaceData = result.data.resources.map((resource: any) => ({
                        namespace: resource.namespace,
                        pods: resource.pods || [],
                        deployments: resource.deployments || [],
                    }));
                    setNamespaces(namespaceData);
                }
            } catch (error) {
                console.error('Error fetching namespace data:', error);
            }
        };

        fetchData();
        const interval = setInterval(fetchData, 30000);

        return () => clearInterval(interval);
    }, []);

    return (
        <Card>
            <CardContent>
                <Typography variant="h5" component="h2" gutterBottom>
                    Namespaces ({namespaces.length})
                </Typography>

                {namespaces.length === 0 ? (
                    <Typography color="text.secondary" align="center">
                        No namespaces found
                    </Typography>
                ) : (
                    <Box sx={{ maxHeight: 400, overflow: 'auto' }}>
                        {namespaces.map((ns: NamespaceData) => (
                            <Accordion key={ns.namespace} sx={{ mb: 1 }}>
                                <AccordionSummary expandIcon={<ExpandMore />}>
                                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, width: '100%' }}>
                                        <Typography variant="h6">{ns.namespace}</Typography>
                                        <Chip
                                            icon={<Memory />}
                                            label={`${ns.pods.length} pods`}
                                            size="small"
                                            color="primary"
                                        />
                                        <Chip
                                            icon={<Storage />}
                                            label={`${ns.deployments.length} deployments`}
                                            size="small"
                                            color="secondary"
                                        />
                                    </Box>
                                </AccordionSummary>
                                <AccordionDetails>
                                    <Box sx={{ width: '100%' }}>
                                        {ns.pods.length > 0 && (
                                            <Box sx={{ mb: 2 }}>
                                                <Typography variant="subtitle2" gutterBottom>
                                                    Pods ({ns.pods.length})
                                                </Typography>
                                                <List dense>
                                                    {ns.pods.map((pod: string) => (
                                                        <ListItem
                                                            key={pod}
                                                            {...({ button: true, selected: selected?.type === 'pod' && selected.namespace === ns.namespace && selected.name === pod, onClick: () => setSelected({ type: 'pod', namespace: ns.namespace, name: pod }) } as any)}
                                                        >
                                                            <ListItemText primary={pod} />
                                                        </ListItem>
                                                    ))}
                                                </List>
                                            </Box>
                                        )}

                                        {ns.deployments.length > 0 && (
                                            <Box>
                                                <Typography variant="subtitle2" gutterBottom>
                                                    Deployments ({ns.deployments.length})
                                                </Typography>
                                                <List dense>
                                                    {ns.deployments.map((deployment: string) => (
                                                        <ListItem
                                                            key={deployment}
                                                            {...({ button: true, selected: selected?.type === 'deployment' && selected.namespace === ns.namespace && selected.name === deployment, onClick: () => setSelected({ type: 'deployment', namespace: ns.namespace, name: deployment }) } as any)}
                                                        >
                                                            <ListItemText primary={deployment} />
                                                        </ListItem>
                                                    ))}
                                                </List>
                                            </Box>
                                        )}
                                    </Box>
                                </AccordionDetails>
                            </Accordion>
                        ))}
                    </Box>
                )}
            </CardContent>
        </Card>
    );
};

export default NamespaceList; 