import React from 'react'
import { ChakraProvider, theme } from '@chakra-ui/react'
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import LinkedAccounts from './pages/LinkedAccounts'
import Login from './pages/Login'
import Protected from './nav/Protected'
import Transactions from './pages/Transactions'
import Rules from './pages/Rules'
import AppStateProvider from './hooks/AppStateProvider'
import RegisterUser from './pages/RegisterUser'

export default function App() {
  return (
    <AppStateProvider>
      <ChakraProvider theme={theme}>
        <Router>
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<RegisterUser />} />
            <Route
              path="/"
              element={
                <Protected current="transactions">
                  <Transactions />
                </Protected>
              }
            />
            <Route
              path="/rules"
              element={
                <Protected current="rules">
                  <Rules />
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
            <Route path="*" element={<Login />} />
          </Routes>
        </Router>
      </ChakraProvider>
    </AppStateProvider>
  )
}
