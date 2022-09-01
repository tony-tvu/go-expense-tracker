import { ChakraProvider, theme } from "@chakra-ui/react"
import { BrowserRouter as Router, Routes, Route } from "react-router-dom"
import Navbar from "./components/Navbar"
import ConnectAccount from "./pages/ConnectAccount"
import PageNotFound from "./pages/PageNotFound"
import AdminPage from "./pages/AdminPage"
import Login from "./pages/Login"
import { useEffect } from "react"

export default function App() {

  return (
    <ChakraProvider theme={theme}>
      <Router>
        {/* <Navbar /> */}
        <Routes>
          <Route path="/" element={<ConnectAccount />} />
          <Route path="/login" element={<Login />} />
          <Route path="/connect" element={<ConnectAccount />} />
          <Route path="/admin" element={<AdminPage />} />
          <Route path="/not-found" element={<PageNotFound />} />
          <Route path="*" element={<PageNotFound />} />
        </Routes>
      </Router>
    </ChakraProvider>
  )
}
