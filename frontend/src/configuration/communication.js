export let server = 'localhost:81'

export let protocol = 'https'

export function setJWTToken(jwt) {
  sessionStorage.setItem("JWT", JSON.stringify(jwt));
}

export function getJWTToken() {
  return JSON.parse(sessionStorage.getItem("JWT"));
}

export function getHeader() {
  if (getJWTToken()) {
    return {
      Authorization: "Bearer " + getJWTToken().token
    };
  }
  return {
    Authorization: "Bearer "
  };
}

export function getLoggedUserID() {
  if (getJWTToken()) {
    return getJWTToken().profileId;
  }
  return 0;
}

export function getLoggedUserUsername() {
  if (getJWTToken()) {
    return getJWTToken().username;
  }
  return null;
}

export function setLoggedUserUsername(u) {
  let jwt = getJWTToken();
  jwt.username = u;
  setJWTToken(jwt);
}

export function logOut() {
  sessionStorage.removeItem("JWT");
  console.log(sessionStorage.getItem("JWT"));
  window.location.href = '/';
}

export function getUrlVars() {
  var vars = {};
  window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, function(m, key, value) {
    vars[key] = value;
  });
  return vars;
}