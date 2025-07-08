import React, { useState } from 'react';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { Box, AppBar, Toolbar, Typography, Container } from '@mui/material';
import DashboardIcon from '@mui/icons-material/Dashboard';
import ClusterOverview from './components/ClusterOverview';
import NamespaceList from './components/NamespaceList';
import MetricsChart from './components/MetricsChart';

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
            <NamespaceList selected={selected} setSelected={setSelected} />
            <MetricsChart selected={selected} />
          </Box>
        </Container>
      </Box>
    </ThemeProvider>
  );
}

export default App;
