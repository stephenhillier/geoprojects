# web

This is the Earthworks web frontend.

Vue components are located in the `/src` directory. The entrypoint is `src/main.js`.

Authentication is provided by an external service (auth0) and is triggered by the router (when accessing protected routes, which are currently all routes). See `src/router.js` and `src/components/common/AuthService.js`

## Project setup
```
npm install
```

### Compiles and hot-reloads for development
```
npm run serve
```

### Compiles and minifies for production
```
npm run build
```

### Run your tests
```
npm run test
```

### Lints and fixes files
```
npm run lint
```

### Run your unit tests
```
npm run test:unit
```
