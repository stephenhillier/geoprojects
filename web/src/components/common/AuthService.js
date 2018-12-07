// src/Auth/AuthService.js

import auth0 from 'auth0-js'
import EventEmitter from 'eventemitter3'
import router from '../../router'

const redirectUri = process.env.VUE_APP_AUTH_CALLBACK || 'http://localhost:8080/callback'

class AuthService {
  auth0 = new auth0.WebAuth({
    domain: 'earthworks.auth0.com',
    audience: 'https://earthworks.islandcivil.com',
    clientID: 'BYgv3PTBtCEtS76GFXF7uDv1vf4XT5N7',
    redirectUri: redirectUri,
    responseType: 'token id_token',
    scope: 'openid email profile'
  })

  login () {
    this.auth0.authorize()
  }

  accessToken
  idToken
  expiresAt
  name
  authenticated = this.isAuthenticated()
  authNotifier = new EventEmitter()

  // ...
  handleAuthentication () {
    this.auth0.parseHash((err, authResult) => {
      if (authResult && authResult.accessToken && authResult.idToken) {
        console.log(authResult)
        this.setSession(authResult)
        router.push({ name: 'projects' })
      } else if (err) {
        router.push({ name: 'projects' })
        console.log(err)
        alert(`Error: ${err.error}. Check the console for further details.`)
      }
    })
  }

  setSession (authResult) {
    this.accessToken = authResult.accessToken
    this.idToken = authResult.idToken
    this.expiresAt = authResult.expiresIn * 1000 + new Date().getTime()
    this.name = authResult.idTokenPayload['name']

    this.authNotifier.emit('authChange', { authenticated: true })

    router.app.$http.defaults.headers.common['Authorization'] = `Bearer ${authResult.accessToken}`

    localStorage.setItem('loggedIn', true)
  }

  renewSession () {
    this.auth0.checkSession({}, (err, authResult) => {
      if (authResult && authResult.accessToken && authResult.idToken) {
        this.setSession(authResult)
      } else if (err) {
        console.log(err)
        this.logout()
      }
    })
  }

  logout () {
    // Clear access token and ID token from local storage
    this.accessToken = null
    this.idToken = null
    this.expiresAt = null
    this.name = null
    this.authNotifier.emit('authChange', false)

    localStorage.removeItem('loggedIn')

    delete router.app.$http.defaults.headers.common['Authorization']

    // navigate to the home route
    router.replace('/')
  }

  isAuthenticated () {
    // Check whether the current time is past the
    // access token's expiry time
    return new Date().getTime() < this.expiresAt && localStorage.getItem('loggedIn') === 'true'
  }
}

export default new AuthService()
