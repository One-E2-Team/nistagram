export let rules = {
  required: value => !!value || 'Required.',
  min: v => v.length >= 8 || 'Min 8 characters',
  email: v => /.+@.+\..+/.test(v) || 'E-mail must be valid',
  name: v => (v && v.length <= 10) || 'Name must be less than 10 characters',
  emailMatch: () => (`The email and password you entered don't match`),
}