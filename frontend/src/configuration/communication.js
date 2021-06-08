export let server = 'localhost:81'

export function getJWTToken() {
  return JSON.parse(sessionStorage.getItem("JWT"));
}

export function getLoggedUserID() {
  if (getJWTToken() != null) {
    return getJWTToken().profileId;
  }
  return 0;
}

export function logOut() {
  sessionStorage.removeItem("JWT");
  console.log(sessionStorage.getItem("JWT"));
  window.location.href = '/';
}

export function setJWTToken(jwt) {
  sessionStorage.setItem("JWT", JSON.stringify(jwt));
}

export function getUrlVars() {
  var vars = {};
  window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, function(m, key, value) {
    vars[key] = value;
  });
  return vars;
}

export function getHeader() {
  if (getJWTToken() != null) {
    return {
      Authorization: "Bearer " + getJWTToken().token
    };
  }
  return {
    Authorization: "Bearer "
  };
}