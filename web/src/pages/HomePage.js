import React from 'react'
import Navbar from '../components/Navbar'
import { useVerifyLogin } from '../hooks/useVerifyLogin'

export default function HomePage() {
  useVerifyLogin()

  return (
    <>
      <Navbar />
      HOMEPAGE
    </>
  )
}
