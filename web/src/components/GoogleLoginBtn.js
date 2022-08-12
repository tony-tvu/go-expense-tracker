import React from 'react'
import { GoogleLogin } from 'react-google-login'

export default function GoogleLoginBtn() {
  const responseGoogle = response => {
    console.log(response)
  }

  return (
    <div>
      <GoogleLogin
        clientId={process.env.REACT_APP_GOOGLE_CLIENT_ID}
        buttonText="Login"
        onSuccess={responseGoogle}
        onFailure={responseGoogle}
        cookiePolicy={'single_host_origin'}
        redirectUri={process.env.REACT_APP_REDIRECT_URI}
        uxMode="redirect"
      />
    </div>
  )
}
