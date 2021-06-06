export let server = 'localhost:81'

export function getJWTToken() {
  return JSON.parse(sessionStorage.getItem("JWT"))
}

export function logOut() {
  sessionStorage.removeItem("JWT");
  console.log(sessionStorage.getItem("JWT"));
  window.location.href = '/';
}

export function setJWTToken(jwt) {
  sessionStorage.setItem("JWT", JSON.stringify(jwt));
}