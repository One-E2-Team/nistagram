export let server = 'localhost:7005'

export let protocol = 'https'

export function setJWTToken(jwt) {
  let new_roles = [];
  for(let item of jwt.roles){
    new_roles.push(item.name);
  }
  jwt.roles = new_roles;
  sessionStorage.setItem("JWT", JSON.stringify(jwt));
}

export function hasRole(role) {
  let jwt = JSON.parse(sessionStorage.getItem("JWT"));
  if (jwt == undefined || jwt == null || jwt == {}){
    return false;
  }
  return jwt.roles.includes(role);
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

export function isUserLogged() {
  return getLoggedUserID() != 0;
}

export function logOut() {
  sessionStorage.removeItem("JWT");
  console.log(sessionStorage.getItem("JWT"));
}

export function getUrlVars() {
  var vars = {};
  window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, function(m, key, value) {
    vars[key] = value;
  });
  return vars;
}