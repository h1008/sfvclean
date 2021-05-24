# sfvclean

A simple tool to remove outdated files from SyncThing's 
[Simple File Versioning](https://docs.syncthing.net/users/versioning.html#simple-file-versioning) 
history.

## Usage Example

```
# Find all files in all stversion folders (within /home/user) that are older than 30 days
sfvclean analyze -keep=30D /home/user

# Delete the identified files
sfvclean clean -force
```