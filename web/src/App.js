import { ChakraProvider, theme } from '@chakra-ui/react'
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import HomePage from './pages/HomePage'
import AccountsPage from './pages/AccountsPage'
import AdminPage from './pages/AdminPage'
import LoginPage from './pages/LoginPage'
import Protected from './components/Protected'

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
                <HomePage />
              </Protected>
            }
          />
          <Route
            path="/accounts"
            element={
              <Protected>
                <AccountsPage />
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
