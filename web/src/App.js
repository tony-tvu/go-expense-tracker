import { ChakraProvider, theme } from '@chakra-ui/react'
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import OverviewPage from './pages/OverviewPage'
import LinkedAccountsPage from './pages/LinkedAccountsPage'
import AdminPage from './pages/AdminPage'
import LoginPage from './pages/LoginPage'
import Protected from './components/Protected'
import ExpensesPage from './pages/ExpensesPage'

export default function App() {
  return (
    <ChakraProvider theme={theme}>
      <Router>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route
            path="/"
            element={
              <Protected>
                <OverviewPage />
              </Protected>
            }
          />
          <Route
            path="/accounts"
            element={
              <Protected>
                <LinkedAccountsPage />
              </Protected>
            }
          />
           <Route
            path="/expenses"
            element={
              <Protected>
                <ExpensesPage />
              </Protected>
            }
          />
          <Route
            path="/admin"
            element={
              <Protected adminOnly={true}>
                <AdminPage />
              </Protected>
            }
          />

          <Route path="*" element={<LoginPage />} />
        </Routes>
      </Router>
    </ChakraProvider>
  )
}
