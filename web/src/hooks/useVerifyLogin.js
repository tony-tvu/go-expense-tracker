import { useEffect } from "react"
import { useNavigate } from "react-router-dom"
import { useQuery, gql } from "@apollo/client"

const query = gql`
  query {
    isLoggedIn
  }
`

export function useVerifyLogin() {
  const navigate = useNavigate()

  const { data } = useQuery(query, {
    fetchPolicy: "no-cache",
  })

  useEffect(() => {
    if (data && !data.isLoggedIn) {
      navigate("/login")
    }
  }, [data, navigate])
}
