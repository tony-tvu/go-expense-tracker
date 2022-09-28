import { ChakraProvider, theme } from '@chakra-ui/react'
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import Dashboard from './pages/Dashboard'
import LinkedAccounts from './pages/LinkedAccounts'
import Admin from './pages/Admin'
import Login from './pages/Login'
import Protected from './nav/Protected'
import Expenses from './pages/Expenses'

export default function App() {
  return (
    <ChakraProvider theme={theme}>
      <Router>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route
            path="/"
            element={
              <Protected current="dashboard">
                <Dashboard />
              </Protected>
            }
          />
          <Route
            path="/accounts"
            element={
              <Protected current="linked_accounts">
                <LinkedAccounts />
              </Protected>
            }
          />
           <Route
            path="/expenses"
            element={
              <Protected current="expenses">
                <Expenses />
              </Protected>
            }
          />
          <Route
            path="/admin"
            element={
              <Protected adminOnly={true} current="admin">
                <Admin />
              </Protected>
            }
          />

          <Route path="*" element={<Login />} />
        </Routes>
      </Router>
    </ChakraProvider>
  )
}
