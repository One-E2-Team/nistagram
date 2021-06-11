export let rules = {
  required: value => !!value || 'Required.',
  min: v => v.length >= 8 || 'Min 8 characters',
  email: v => new RegExp('^([a-zA-Z0-9]+.?)*[a-zA-Z0-9]@[a-z0-9]+(.[a-z]{2,3})+$').test(v) || 'E-mail must be valid',
  max: v => (v && v.length <= 255) || 'Max 10 characters',
  password: v => (new RegExp('^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[*.!@#$%^&(){}\\[\\]:;<>,.?~_+\\-=|\\/])[A-Za-z0-9*.!@#$%^&(){}\\[\\]:;<>,.?~_+\\-=|\\/]{8,}$')).test(v) || 'Password must contain at least one lower, one capital letter, one number and one special character! Password must have at least 8 characters!',
  username: v => (new RegExp('^[a-zA-Z0-9_]{3,15}$')).test(v) || 'Only possible characters are letters, numbers and underscore (3-15 characters limit)'
}