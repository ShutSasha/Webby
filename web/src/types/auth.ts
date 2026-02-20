type GoogleAuthRes = {
  success: boolean
  message: string
  data: {
    userId: string
    email: string
    username: string
    about: string
    avatarUrl: string
  }
  errors: null
}

export { type GoogleAuthRes }
