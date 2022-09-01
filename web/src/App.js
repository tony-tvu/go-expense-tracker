import { ChakraProvider } from "@chakra-ui/react"
import { BrowserRouter as Router, Routes, Route } from "react-router-dom"
import ConnectAccount from "./pages/ConnectAccount"
import PageNotFound from "./pages/PageNotFound"
import AdminPage from "./pages/AdminPage"
import Login from "./pages/Login"
import { extendedTheme } from "./theme"

export default function App() {

  return (
    <ChakraProvider theme={extendedTheme}>
      <Router>
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
