import { Flex, useColorModeValue } from "@chakra-ui/react"
import React from "react"
import Navbar from "../components/Navbar"
import { useVerifyLogin } from "../hooks/useVerifyLogin"
import { colors } from "../theme"

export default function HomePage() {
  useVerifyLogin()

  return (
    <>
      <Navbar />
      HOMEPAGE
    </>
  )
}
