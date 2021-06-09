export let rules = {
  required: value => !!value || 'Required.',
  min: v => v.length >= 8 || 'Min 8 characters',
  email: v => /.+@.+\..+/.test(v) || 'E-mail must be valid',
  name: v => (v && v.length <= 10) || 'Name must be less than 10 characters',
  emailMatch: () => (`The email and password you entered don't match`),
  password: v => new RegExp('^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[*.!@#$%^&(){}\\[\\]:;<>,.?~_+-=|\\/])[A-Za-z0-9*.!@#$%^&(){}\\[\\]:;<>,.?~_+-=|\\/]{8,}$').test(v) || 'Password must contain at least one lower, one capital letter, one number and one special character! Password must have at least 8 characters!'
}