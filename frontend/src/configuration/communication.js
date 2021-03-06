export let server = 'localhost:81'

export let protocol = 'https'

export let wsProtocol = 'wss'

export let wsNotificationServer = 'localhost:7008'

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
}

export function getUrlVars() {
  var vars = {};
  window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, function(m, key, value) {
    vars[key] = value;
  });
  return vars;
}

export async function openWebSocketConn(url, handler){
  let ws = new WebSocket(url + "?token=" + getJWTToken().token)
  let reload = function(event) {
    console.log(event)
    window.location.reload()
  }
  ws.onerror = reload
  ws.onclose = reload
  ws.onmessage = function(event) {
    handler(event.data.response, event.data.data)
  }

  await new Promise(r => setTimeout(r, 2000));
  
  return function(request, data){
    let req = {
      jwt: getJWTToken().token, 
      request: request,
      data: JSON.stringify(data)
    }
    ws.send(JSON.stringify(req))
  }
}
