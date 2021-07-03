export function formatTime(time) {
  let t = new Date(time);
  let minutes = t.getMinutes() < 10 ? '0' + t.getMinutes() : t.getMinutes()
  let hours = t.getHours() < 10 ? '0' + t.getHours() : t.getHours()
  return hours + ":" + minutes
}

export function formatDateToYYYYMMDD(date) {
  let d = new Date(date)
  return d.toISOString().split('T')[0]
}