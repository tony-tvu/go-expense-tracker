import { ChakraProvider, theme } from "@chakra-ui/react"
import { BrowserRouter as Router, Routes, Route } from "react-router-dom"
import HomePage from "./pages/HomePage"
import AccountsPage from "./pages/AccountsPage"
import PageNotFound from "./pages/PageNotFound"
import AdminPage from "./pages/AdminPage"
import LoginPage from "./pages/LoginPage"

export default function App() {
  return (
    <ChakraProvider theme={theme}>
      <Router>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/accounts" element={<AccountsPage />} />
          <Route path="/admin" element={<AdminPage />} />
          <Route path="/not-found" element={<PageNotFound />} />
          <Route path="*" element={<PageNotFound />} />
        </Routes>
      </Router>
    </ChakraProvider>
  )
}
