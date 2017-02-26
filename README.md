# go-userbase
Go &amp; SQLite version of [Com.Ericmas001.Userbase](https://github.com/Com-Ericmas001/Userbase)

Done
```
    public int IdFromUsername(UserbaseDbContext context, string username)
    public int IdFromEmail(UserbaseDbContext context, string email)
    public ConnectUserResponse ValidateToken(UserbaseDbContext context, string username, Guid token)
    public ConnectUserResponse ValidateCredentials(UserbaseDbContext context, string username, string password)
    public ConnectUserResponse CreateUser(UserbaseDbContext context, CreateUserRequest request)
    public TokenSuccessResponse ModifyCredentials(UserbaseDbContext context, ModifyCredentialsRequest request)
    public TokenSuccessResponse ModifyProfile(UserbaseDbContext context, ModifyProfileRequest request)
    public bool Disconnect(UserbaseDbContext context, string username, Guid token)
```

TODO
```
    public void PurgeUsers(UserbaseDbContext context)
    public void PurgeConnectionTokens(UserbaseDbContext context)
    public void PurgeRecoveryTokens(UserbaseDbContext context)
    public bool Deactivate(UserbaseDbContext context, string username, Guid token)
    public bool SendRecoveryToken(UserbaseDbContext context, string username, IEmailSender smtp)
    public ConnectUserResponse ResetPassword(UserbaseDbContext context, string username, Guid recoveryToken, string newPassword)
    public UserSummaryResponse UserSummary(UserbaseDbContext context, string askingUser, Guid token, string requestedUser)
```