import React, { useState } from 'react';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { Box, AppBar, Toolbar, Typography, Container, Dialog, DialogContent } from '@mui/material';
import DashboardIcon from '@mui/icons-material/Dashboard';
import ClusterOverview from './components/ClusterOverview';
import NamespaceList from './components/NamespaceList';
import MetricsChart from './components/MetricsChart';
import PodLogs from './components/PodLogs';

const theme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#2196f3',
    },
    secondary: {
      main: '#f50057',
    },
    background: {
      default: '#0a1929',
      paper: '#132f4c',
    },
  },
});

function App() {
  const [selected, setSelected] = useState<{ type: 'pod' | 'deployment'; namespace: string; name: string } | null>(null);
  const [logViewer, setLogViewer] = useState<{ namespace: string; podName: string } | null>(null);

  const handleViewLogs = (namespace: string, podName: string) => {
    setLogViewer({ namespace, podName });
  };

  const handleCloseLogs = () => {
    setLogViewer(null);
  };

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Box sx={{ flexGrow: 1 }}>
        <AppBar position="static">
          <Toolbar>
            <DashboardIcon sx={{ mr: 2 }} />
            <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
              KubeFleet Dashboard
            </Typography>
          </Toolbar>
        </AppBar>

        <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
          <ClusterOverview />

          <Box sx={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 3, mt: 3 }}>
            <NamespaceList
              selected={selected}
              setSelected={setSelected}
              onViewLogs={handleViewLogs}
            />
            <MetricsChart selected={selected} />
          </Box>
        </Container>

        {/* Log Viewer Dialog */}
        <Dialog
          open={logViewer !== null}
          onClose={handleCloseLogs}
          maxWidth="xl"
          fullWidth
          PaperProps={{
            sx: {
              height: '80vh',
              maxHeight: '80vh',
            },
          }}
        >
          <DialogContent sx={{ p: 0 }}>
            {logViewer && (
              <PodLogs
                namespace={logViewer.namespace}
                podName={logViewer.podName}
                onClose={handleCloseLogs}
              />
            )}
          </DialogContent>
        </Dialog>
      </Box>
    </ThemeProvider>
  );
}

export default App;
