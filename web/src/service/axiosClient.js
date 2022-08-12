import axios from 'axios'

const client = axios

client.interceptors.response.use(function (response) {
  if (
    response.status !== 403 &&
    response.headers['access-token-renewed'] === 'true'
  ) {
    localStorage.setItem(
      'user-access-token',
      response.headers['user-access-token']
    )
  }
  return response
})

export default client
