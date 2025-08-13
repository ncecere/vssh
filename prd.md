# vssh

A golang cli tool that signs ssh keys with hashicorp vault and uses that signed cert with ssh to login to servers.

it should support signing in multiple sign in methods.

it should support passing in ssh arguments.

Over all it should function like ssh expect you will pass in the command vssh.

there should be a config file where you can configure the required vault server parameters.  like server address, role, etc.

it should support signing ssh keys for multiple users... 

the config file should support specifying what private key to use.

it should support creating the signed ssh public key for example

vssh user1@server.com will sign the key for user1 in vault and create the public key with the of vault_signed_user1.pub

and vssh user2@server.com will sign the key for user2 in vault and create the public key with the of vault_signed_user2.pub

if key are still valid it will use them but if they ar eno longer valid it will recreate them.

It should look for the vault token before prompting the user to sign in.

I want to use golang, cobra, and viper.  if there are any other tool you think i need bring them up.