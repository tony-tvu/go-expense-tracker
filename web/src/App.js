import React from 'react'
import { ChakraProvider, theme } from '@chakra-ui/react'
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import LinkedAccounts from './pages/LinkedAccounts'
import Admin from './pages/Admin'
import Login from './pages/Login'
import Protected from './nav/Protected'
import Transactions from './pages/Transactions'
import Rules from './pages/Rules'
import Analytics from './pages/Analytics'
import AppStateProvider from './hooks/AppStateProvider'

export default function App() {
  return (
    <AppStateProvider>
      <ChakraProvider theme={theme}>
        <Router>
          <Routes>
            <Route path="/login" element={<Login />} />
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
            <Route
              path="/analytics"
              element={
                <Protected current="analytics">
                  <Analytics />
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
    </AppStateProvider>
  )
}
